package tests

//go:generate gengo msg std_msgs/String
import (
	"fmt"
	"log"
	"os"
	"rospy_tutorials"
	"std_msgs"
	"strings"
	"testing"
	"time"

	"github.com/fetchrobotics/rosgo/ros"
)

func TestPublisherAndSubscriber(t *testing.T) {
	node, err := ros.NewNode("subpub_test", os.Args)
	if err != nil {
		t.Fatal(err)
	}
	defer node.Shutdown()

	// Create message that we want to publish
	msgWanted := new(std_msgs.String)
	msgWanted.Data = fmt.Sprintf("hello %s", time.Now().String())

	// Create message channel for subscription messages
	msgCh := make(chan *std_msgs.String, 10)
	// Create new subscriber to earlier created publisher
	sub := node.NewSubscriber("/chatter", std_msgs.MsgString,
		func(msg *std_msgs.String) { msgCh <- msg })
	defer sub.Shutdown()

	// Create new publisher
	pub := node.NewPublisher("/chatter", std_msgs.MsgString)
	defer pub.Shutdown()

	// Wait for the publisher to have at least one subscriber
	for pub.GetNumSubscribers() == 0 {
		if !node.OK() {
			t.Fatalf("ROS node shutdown while waiting for subscriber")
		}
		t.Log("Waiting for subscriber")
		time.Sleep(500 * time.Millisecond)
	}

	pub.Publish(msgWanted)
	node.SpinOnce()

	select {
	case msgGot := <-msgCh:
		if msgGot == nil {
			t.Errorf("Nil message received")
		}
		if strings.Compare(msgWanted.Data, msgGot.Data) != 0 {
			t.Errorf("Expected %v Got %v", msgWanted.Data, msgGot.Data)
		}
	case <-time.After(2 * time.Second):
		t.Errorf("Message not received within timeout")
	}
}

func TestPublisherWithCallback(t *testing.T) {
	node, err := ros.NewNode("publisher_test", os.Args)
	if err != nil {
		log.Fatal(err)
	}
	defer node.Shutdown()

	// Create message that we want to publish
	msgWanted := new(std_msgs.String)
	msgWanted.Data = fmt.Sprintf("hello %s", time.Now().String())
	msgCh := make(chan *std_msgs.String, 10)
	pubCh := make(chan ros.SingleSubscriberPublisher, 10)

	pubCallback := func(sp ros.SingleSubscriberPublisher) { pubCh <- sp }
	subCallback := func(msg *std_msgs.String) { msgCh <- msg }

	// Create new publisher with callbacks
	pub := node.NewPublisherWithCallbacks("/chatter", std_msgs.MsgString, pubCallback, pubCallback)
	defer pub.Shutdown()

	subs := pub.GetNumSubscribers()
	if subs != 0 {
		t.Errorf("Expected 0 subscribers got %v", subs)
	}

	sub := node.NewSubscriber("/chatter", std_msgs.MsgString, subCallback)

	// Wait for the publisher to have at least one subscriber
	for pub.GetNumSubscribers() != 0 {
		if !node.OK() {
			t.Fatalf("ROS node shutdown while waiting for subscriber")
		}
		t.Log("Waiting for subscriber")
		time.Sleep(500 * time.Millisecond)
	}

	var ssp ros.SingleSubscriberPublisher
	select {
	case ssp = <-pubCh:
	case <-time.After(2 * time.Second):
		t.Fatal("Expected single publisher when subscriber connected but got nothing")
	}

	if ssp.GetSubscriberName() != "/publisher_test" {
		t.Errorf("Subscriber name does not match: Expected \"/publisher_test\" Got \"%s\"", ssp.GetSubscriberName())
	}

	if ssp.GetTopic() != "/chatter" {
		t.Errorf("Subscriber topic does not match: Expected \"/chatter\" Got \"%s\"", ssp.GetTopic())
	}

	ssp.Publish(msgWanted)
	node.SpinOnce()

	select {
	case msgGot := <-msgCh:
		if msgGot == nil {
			t.Errorf("Nil message received")
		}
		if strings.Compare(msgWanted.Data, msgGot.Data) != 0 {
			t.Errorf("Expected \"%v\" Got \"%v\"", msgWanted.Data, msgGot.Data)
		}
	case <-time.After(2 * time.Second):
		t.Errorf("Message not received within timeout")
	}

	sub.Shutdown()

	select {
	case ssp = <-pubCh:
	case <-time.After(2 * time.Second):
		t.Error("Expected single publisher when subscriber disconnected but got nothing")
	}
}

func TestSubscriberWithCallback(t *testing.T) {
	node, err := ros.NewNode("subscriber_test", os.Args)
	if err != nil {
		log.Fatal(err)
	}
	defer node.Shutdown()

	// Create message that we want to publish
	msgWanted := new(std_msgs.String)
	msgWanted.Data = fmt.Sprintf("hello %s", time.Now().String())

	// Create message channel for subscription messages
	nopCh := make(chan struct{}, 2)
	msgCh := make(chan *std_msgs.String, 2)
	eventCh := make(chan ros.MessageEvent, 2)

	nopCb := func() { nopCh <- struct{}{} }
	msgCb := func(msg *std_msgs.String) { msgCh <- msg }
	eventCb := func(msg *std_msgs.String, event ros.MessageEvent) { eventCh <- event }

	// Create new subscriber for each callback type to earlier created publisher
	sub := node.NewSubscriber("/chatter", std_msgs.MsgString, nopCb)
	node.NewSubscriber("/chatter", std_msgs.MsgString, msgCb)
	node.NewSubscriber("/chatter", std_msgs.MsgString, eventCb)
	defer sub.Shutdown()

	// Create new publisher
	pub := node.NewPublisher("/chatter", std_msgs.MsgString)
	defer pub.Shutdown()

	// Wait for the publisher to have at least one subscriber
	for pub.GetNumSubscribers() == 0 {
		if !node.OK() {
			t.Fatalf("ROS node shutdown while waiting for subscriber")
		}
		t.Log("Waiting for subscriber")
		time.Sleep(500 * time.Millisecond)
	}

	timeSent := time.Now()
	pub.Publish(msgWanted)
	node.SpinOnce()

	select {
	case <-nopCh:
	case <-time.After(2 * time.Second):
		t.Errorf("Notification not received within timeout")
	}

	select {
	case <-msgCh:
	case <-time.After(2 * time.Second):
		t.Errorf("Message not received within timeout")
	}

	select {
	case event := <-eventCh:
		if event.PublisherName != "/subscriber_test" {
			t.Error("Publisher and node name do not match")
		}
		if !event.ReceiptTime.After(timeSent) {
			t.Error("Invalid event receipt time received")
		}
		if _, ok := event.ConnectionHeader["topic"]; !ok {
			t.Error("topic field missing in connection header")
		}
		if _, ok := event.ConnectionHeader["type"]; !ok {
			t.Error("type field missing in connection header")
		}
		if _, ok := event.ConnectionHeader["latching"]; !ok {
			t.Error("latching field missing in connection header")
		}
		if _, ok := event.ConnectionHeader["md5sum"]; !ok {
			t.Error("md5sum field missing in connection header")
		}
		if _, ok := event.ConnectionHeader["callerid"]; !ok {
			t.Error("callerid field missing in connection header")
		}
		if _, ok := event.ConnectionHeader["message_definition"]; !ok {
			t.Error("message_definition field missing in connection header")
		}
	case <-time.After(2 * time.Second):
		t.Errorf("Event not received within timeout")
	}
}

func TestParam(t *testing.T) {
	node, err := ros.NewNode("/test_param", os.Args)
	if err != nil {
		t.Fatal(err)
	}
	defer node.Shutdown()

	t.Run("HasParam", func(t *testing.T) {
		if hasParam, err := node.HasParam("/rosdistro"); err != nil {
			t.Errorf("HasParam failed: %v", err)
		} else {
			if !hasParam {
				t.Error("HasParam() failed.")
			}
		}
	})

	t.Run("SearchParam", func(t *testing.T) {
		if foundKey, err := node.SearchParam("rosdistro"); err != nil {
			t.Errorf("SearchParam failed: %v", err)
		} else {
			if foundKey != "/rosdistro" {
				t.Error("SearchParam() failed.")
			}
		}
	})

	t.Run("GetParam", func(t *testing.T) {
		distro := fmt.Sprintf("%v\n", os.Getenv("ROS_DISTRO"))
		if param, err := node.GetParam("/rosdistro"); err != nil {
			t.Errorf("GetParam: %v", err)
		} else {
			if value, ok := param.(string); !ok {
				t.Error("GetParam() failed.")
			} else {
				if value != distro {
					t.Errorf("Expected '%s' but '%s'", distro, value)
				}
			}
		}
	})

	t.Run("SetAndGetParam", func(t *testing.T) {
		if err := node.SetParam("/test_param", 42); err != nil {
			t.Errorf("SetParam failed: %v", err)
		}

		if param, err := node.GetParam("/test_param"); err != nil {
			t.Errorf("GetParam failed: %v", err)
		} else {
			if value, ok := param.(int32); ok {
				if value != 42 {
					t.Errorf("Expected 42 but %d", value)
				}
			} else {
				t.Error("GetParam('/test_param') failed.")
			}
		}
	})

	t.Run("DeleteParam", func(t *testing.T) {
		if err := node.DeleteParam("/test_param"); err != nil {
			t.Errorf("DeleteParam failed: %v", err)
		}

		if hasParam, err := node.HasParam("/rosdistro"); err != nil {
			t.Errorf("HasParam failed: %v", err)
		} else {
			if !hasParam {
				t.Error("HasParam() failed.")
			}
		}
	})
}

func TestServiceClientAndServer(t *testing.T) {
	node, err := ros.NewNode("clientserver_test", os.Args)
	if err != nil {
		t.Fatal(err)
	}
	defer node.Shutdown()
	go node.Spin()

	sumCb := func(srv *rospy_tutorials.AddTwoInts) error {
		srv.Response.Sum = srv.Request.A + srv.Request.B
		return nil
	}

	server := node.NewServiceServer("/add_two_ints", rospy_tutorials.SrvAddTwoInts, sumCb)
	if server == nil {
		t.Fatal("Failed to initialize '/add_two_ints' service server")
	}
	defer server.Shutdown()
	time.Sleep(1 * time.Second)

	client := node.NewServiceClient("/add_two_ints", rospy_tutorials.SrvAddTwoInts)
	defer client.Shutdown()
	time.Sleep(1 * time.Second)

	var a int64 = 1
	var b int64 = 2
	sum := a + b
	srv := new(rospy_tutorials.AddTwoInts)
	srv.Request.A = a
	srv.Request.B = b

	if err = client.Call(srv); err != nil {
		t.Error(err)
	}

	if srv.Response.Sum != sum {
		t.Errorf("Wrong sum received:  %d + %d != %d",
			srv.Request.A, srv.Request.B, srv.Response.Sum)
	}
}

func TestNodeName(t *testing.T) {
	nodeName := "name_test"
	node1, err := ros.NewNode(nodeName, os.Args)
	if err != nil {
		t.Fatal(err)
	}
	defer node1.Shutdown()

	if node1.Name() != nodeName {
		t.Errorf("Invalid node name: Expected %s Got %s",
			nodeName, node1.Name())
	}

	nodeNameWithSlash := fmt.Sprintf("/%s", nodeName)
	node2, err := ros.NewNode(nodeNameWithSlash, os.Args)
	if err != nil {
		t.Fatal(err)
	}
	defer node2.Shutdown()

	if node2.Name() != nodeName {
		t.Errorf("Invalid node name: Expected %s Got %s",
			nodeName, node2.Name())
	}
}

func TestNodeLogger(t *testing.T) {
	node, err := ros.NewNode("logger_test", os.Args)
	if err != nil {
		t.Fatal(err)
	}
	defer node.Shutdown()

	if node.Logger().Severity() != ros.LogLevelInfo {
		t.Errorf("Default log level should be %v but got %v",
			ros.LogLevelInfo, node.Logger().Severity())
	}

	node.Logger().SetSeverity(ros.LogLevelError)
	if node.Logger().Severity() != ros.LogLevelError {
		t.Errorf("Default log level should be %v but got %v",
			ros.LogLevelError, node.Logger().Severity())
	}
}
