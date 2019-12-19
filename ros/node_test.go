package ros

import (
	"os"
	"testing"
	"time"
)

func TestLoadJsonFromString(t *testing.T) {
	value, err := loadParamFromString("42")
	if err != nil {
		t.Error(err)
	}
	i, ok := value.(float64)
	if !ok {
		t.Fail()
	}
	if i != 42.0 {
		t.Error(i)
	}
}

func TestNodeAPI(t *testing.T) {
	node, err := newDefaultNode("/test_node", []string{})
	if err != nil {
		t.Fatalf("Error starting new test node: %v", err)
	}
	defer node.Shutdown()

	t.Run("GetBusStats", func(t *testing.T) {
		var statusWanted int32 = errorStatus
		var messageWanted string = "Not implemented"
		var valueWanted int = 0

		result, err := node.getBusStats("test_caller")
		if err != nil {
			t.Fatalf("Error calling getBusStats: %v", err)
		}

		res, ok := result.([]interface{})
		if !ok {
			t.Errorf("Result is malformed: Result should be an array of interface values.")
		}

		if len(res) != 3 {
			t.Errorf("Result is malformed: Length of result array should be 3")
		}

		statusGot, ok := res[0].(int32)
		if !ok {
			t.Errorf("Error getting status code from result. Type assertion failed.")
			return
		}

		messageGot, ok := res[1].(string)
		if !ok {
			t.Errorf("Error getting message from result. Type assertion failed.")
			return
		}
		valueGot, ok := res[2].(int)
		if !ok {
			t.Errorf("Error getting message from result. Type assertion failed.")
			return
		}

		if statusWanted != statusGot {
			t.Errorf("Expected %v but got %v", statusWanted, statusGot)
		}
		if messageWanted != messageGot {
			t.Errorf("Expected `%v` but got `%v`", messageWanted, messageGot)
		}
		if valueWanted != valueGot {
			t.Errorf("Expected %v but got %v", valueWanted, valueGot)
		}
	})

	t.Run("GetBusInfo", func(t *testing.T) {
		var statusWanted int32 = errorStatus
		var messageWanted string = "Not implemented"
		var valueWanted int = 0

		result, err := node.getBusInfo("test_caller")
		if err != nil {
			t.Fatalf("Error calling getBusInfo: %v", err)
		}

		res, ok := result.([]interface{})
		if !ok {
			t.Errorf("Result is malformed: Result should be an array of interface values.")
		}

		if len(res) != 3 {
			t.Errorf("Result is malformed: Length of result array should be 3")
		}

		statusGot, ok := res[0].(int32)
		if !ok {
			t.Errorf("Error getting status code from result. Type assertion failed.")
			return
		}

		messageGot, ok := res[1].(string)
		if !ok {
			t.Errorf("Error getting message from result. Type assertion failed.")
			return
		}
		valueGot, ok := res[2].(int)
		if !ok {
			t.Errorf("Error getting message from result. Type assertion failed.")
			return
		}

		if statusWanted != statusGot {
			t.Errorf("Expected %v but got %v", statusWanted, statusGot)
		}
		if messageWanted != messageGot {
			t.Errorf("Expected `%v` but got `%v`", messageWanted, messageGot)
		}
		if valueWanted != valueGot {
			t.Errorf("Expected %v but got %v", valueWanted, valueGot)
		}
	})

	t.Run("GetMasterURI", func(t *testing.T) {
		var statusWanted int32 = successStatus
		var messageWanted string = "Success"
		var valueWanted string = os.Getenv("ROS_MASTER_URI")

		result, err := node.getMasterURI("test_caller")
		if err != nil {
			t.Fatalf("Error calling getMasterURI: %v", err)
		}

		res, ok := result.([]interface{})
		if !ok {
			t.Errorf("Result is malformed: Result should be an array of interface values.")
		}

		if len(res) != 3 {
			t.Errorf("Result is malformed: Length of result array should be 3")
		}

		statusGot, ok := res[0].(int32)
		if !ok {
			t.Errorf("Error getting status code from result. Type assertion failed.")
			return
		}

		messageGot, ok := res[1].(string)
		if !ok {
			t.Errorf("Error getting message from result. Type assertion failed.")
			return
		}
		valueGot, ok := res[2].(string)
		if !ok {
			t.Errorf("Error getting message from result. Type assertion failed.")
			return
		}

		if statusWanted != statusGot {
			t.Errorf("Expected %v but got %v", statusWanted, statusGot)
		}
		if messageWanted != messageGot {
			t.Errorf("Expected `%v` but got `%v`", messageWanted, messageGot)
		}
		if valueWanted != valueGot {
			t.Errorf("Expected %v but got %v", valueWanted, valueGot)
		}
	})

	t.Run("GetPid", func(t *testing.T) {
		var statusWanted int32 = successStatus
		var messageWanted string = "Success"
		var valueWanted int = os.Getpid()

		result, err := node.getPid("test_caller")
		if err != nil {
			t.Fatalf("Error calling getPid: %v", err)
		}

		res, ok := result.([]interface{})
		if !ok {
			t.Errorf("Result is malformed: Result should be an array of interface values.")
		}

		if len(res) != 3 {
			t.Errorf("Result is malformed: Length of result array should be 3")
		}

		statusGot, ok := res[0].(int32)
		if !ok {
			t.Errorf("Error getting status code from result. Type assertion failed.")
			return
		}

		messageGot, ok := res[1].(string)
		if !ok {
			t.Errorf("Error getting message from result. Type assertion failed.")
			return
		}
		valueGot, ok := res[2].(int)
		if !ok {
			t.Errorf("Error getting message from result. Type assertion failed.")
			return
		}

		if statusWanted != statusGot {
			t.Errorf("Expected %v but got %v", statusWanted, statusGot)
		}
		if messageWanted != messageGot {
			t.Errorf("Expected `%v` but got `%v`", messageWanted, messageGot)
		}
		if valueWanted != valueGot {
			t.Errorf("Expected %v but got %v", valueWanted, valueGot)
		}
	})

	t.Run("GetSubscriptions", func(t *testing.T) {
		t.Run("ZeroSubscriptions", func(t *testing.T) {
			var statusWanted int32 = successStatus
			var messageWanted string = "Success"
			var valueWanted []interface{}

			result, err := node.getSubscriptions("test_caller")
			if err != nil {
				t.Fatalf("Error calling getSubscriptions: %v", err)
			}

			res, ok := result.([]interface{})
			if !ok {
				t.Errorf("Result is malformed: Result should be an array of interface values.")
			}

			if len(res) != 3 {
				t.Errorf("Result is malformed: Length of result array should be 3")
			}

			statusGot, ok := res[0].(int32)
			if !ok {
				t.Errorf("Error getting status code from result. Type assertion failed.")
				return
			}

			messageGot, ok := res[1].(string)
			if !ok {
				t.Errorf("Error getting message from result. Type assertion failed.")
				return
			}
			valueGot, ok := res[2].([]interface{})
			if !ok {
				t.Errorf("Error getting message from result. Type assertion failed.")
				return
			}

			if statusWanted != statusGot {
				t.Errorf("Expected %v but got %v", statusWanted, statusGot)
			}
			if messageWanted != messageGot {
				t.Errorf("Expected `%v` but got `%v`", messageWanted, messageGot)
			}
			if len(valueWanted) != len(valueGot) {
				t.Errorf("Expected %v but got %v", valueWanted, valueGot)
			}
		})

		t.Run("OneSubscription", func(t *testing.T) {
			var statusWanted int32 = successStatus
			var messageWanted string = "Success"
			topic := "/test_topic"
			msgType := &dummyMessage{}
			valueWanted := []interface{}{[]interface{}{topic, msgType.Name()}}

			sub := newDefaultSubscriber(topic, msgType, nil)
			node.subscribers["/test_topic"] = sub

			result, err := node.getSubscriptions("test_caller")
			if err != nil {
				t.Fatalf("Error calling getSubscriptions: %v", err)
			}

			res, ok := result.([]interface{})
			if !ok {
				t.Errorf("Result is malformed: Result should be an array of interface values.")
			}

			if len(res) != 3 {
				t.Errorf("Result is malformed: Length of result array should be 3")
			}

			statusGot, ok := res[0].(int32)
			if !ok {
				t.Errorf("Error getting status code from result. Type assertion failed.")
				return
			}

			messageGot, ok := res[1].(string)
			if !ok {
				t.Errorf("Error getting message from result. Type assertion failed.")
				return
			}
			valueGot, ok := res[2].([]interface{})
			if !ok {
				t.Errorf("Error getting message from result. Type assertion failed.")
				return
			}

			if statusWanted != statusGot {
				t.Errorf("Expected %v but got %v", statusWanted, statusGot)
			}
			if messageWanted != messageGot {
				t.Errorf("Expected `%v` but got `%v`", messageWanted, messageGot)
			}
			if len(valueWanted) != len(valueGot) {
				t.Errorf("Expected %v but got %v", valueWanted, valueGot)
			}
		})
	})

	t.Run("GetPublications", func(t *testing.T) {
		t.Run("ZeroPublishers", func(t *testing.T) {
			var statusWanted int32 = successStatus
			var messageWanted string = "Success"
			var valueWanted []interface{}

			result, err := node.getPublications("test_caller")
			if err != nil {
				t.Fatalf("Error calling getPublications: %v", err)
			}

			res, ok := result.([]interface{})
			if !ok {
				t.Errorf("Result is malformed: Result should be an array of interface values.")
			}

			if len(res) != 3 {
				t.Errorf("Result is malformed: Length of result array should be 3")
			}

			statusGot, ok := res[0].(int32)
			if !ok {
				t.Errorf("Error getting status code from result. Type assertion failed.")
				return
			}

			messageGot, ok := res[1].(string)
			if !ok {
				t.Errorf("Error getting message from result. Type assertion failed.")
				return
			}
			valueGot, ok := res[2].([]interface{})
			if !ok {
				t.Errorf("Error getting message from result. Type assertion failed.")
				return
			}

			if statusWanted != statusGot {
				t.Errorf("Expected %v but got %v", statusWanted, statusGot)
			}
			if messageWanted != messageGot {
				t.Errorf("Expected `%v` but got `%v`", messageWanted, messageGot)
			}
			if len(valueWanted) != len(valueGot) {
				t.Errorf("Expected %v but got %v", valueWanted, valueGot)
			}
		})

		t.Run("OnePublisher", func(t *testing.T) {
			var statusWanted int32 = successStatus
			var messageWanted string = "Success"
			topic := "/test_topic"
			msgType := &dummyMessage{}
			valueWanted := []interface{}{[]interface{}{topic, msgType.Name()}}

			pub := newDefaultPublisher(node, topic, msgType, nil, nil)
			node.publishers["/test_topic"] = pub

			result, err := node.getPublications("test_caller")
			if err != nil {
				t.Fatalf("Error calling getPublications: %v", err)
			}

			res, ok := result.([]interface{})
			if !ok {
				t.Errorf("Result is malformed: Result should be an array of interface values.")
			}

			if len(res) != 3 {
				t.Errorf("Result is malformed: Length of result array should be 3")
			}

			statusGot, ok := res[0].(int32)
			if !ok {
				t.Errorf("Error getting status code from result. Type assertion failed.")
				return
			}

			messageGot, ok := res[1].(string)
			if !ok {
				t.Errorf("Error getting message from result. Type assertion failed.")
				return
			}
			valueGot, ok := res[2].([]interface{})
			if !ok {
				t.Errorf("Error getting message from result. Type assertion failed.")
				return
			}

			if statusWanted != statusGot {
				t.Errorf("Expected %v but got %v", statusWanted, statusGot)
			}
			if messageWanted != messageGot {
				t.Errorf("Expected `%v` but got `%v`", messageWanted, messageGot)
			}
			if len(valueWanted) != len(valueGot) {
				t.Errorf("Expected %v but got %v", valueWanted, valueGot)
			}
		})
	})

	t.Run("ParamUpdate", func(t *testing.T) {
		var statusWanted int32 = errorStatus
		var messageWanted string = "Not implemented"
		var valueWanted int = 0

		result, err := node.paramUpdate("test_caller", "/test_param", 0)
		if err != nil {
			t.Fatalf("Error calling paramUpdate: %v", err)
		}

		res, ok := result.([]interface{})
		if !ok {
			t.Errorf("Result is malformed: Result should be an array of interface values.")
		}

		if len(res) != 3 {
			t.Errorf("Result is malformed: Length of result array should be 3")
		}

		statusGot, ok := res[0].(int32)
		if !ok {
			t.Errorf("Error getting status code from result. Type assertion failed.")
			return
		}

		messageGot, ok := res[1].(string)
		if !ok {
			t.Errorf("Error getting message from result. Type assertion failed.")
			return
		}
		valueGot, ok := res[2].(int)
		if !ok {
			t.Errorf("Error getting message from result. Type assertion failed.")
			return
		}

		if statusWanted != statusGot {
			t.Errorf("Expected %v but got %v", statusWanted, statusGot)
		}
		if messageWanted != messageGot {
			t.Errorf("Expected `%v` but got `%v`", messageWanted, messageGot)
		}
		if valueWanted != valueGot {
			t.Errorf("Expected %v but got %v", valueWanted, valueGot)
		}
	})

	t.Run("PublisherUpdate", func(t *testing.T) {
		pubList := []interface{}{"10.100.1.225:5647"}

		t.Run("UnknownTopic", func(t *testing.T) {
			var statusWanted int32 = failureStatus
			var messageWanted string = "No such topic"
			var valueWanted int = 0

			result, err := node.publisherUpdate("test_caller", "/unknown_topic", pubList)
			if err != nil {
				t.Fatalf("Error calling publisherUpdate: %v", err)
			}

			res, ok := result.([]interface{})
			if !ok {
				t.Errorf("Result is malformed: Result should be an array of interface values.")
			}

			if len(res) != 3 {
				t.Errorf("Result is malformed: Length of result array should be 3")
			}

			statusGot, ok := res[0].(int32)
			if !ok {
				t.Errorf("Error getting status code from result. Type assertion failed.")
				return
			}

			messageGot, ok := res[1].(string)
			if !ok {
				t.Errorf("Error getting message from result. Type assertion failed.")
				return
			}
			valueGot, ok := res[2].(int)
			if !ok {
				t.Errorf("Error getting message from result. Type assertion failed.")
				return
			}

			if statusWanted != statusGot {
				t.Errorf("Expected %v but got %v", statusWanted, statusGot)
			}
			if messageWanted != messageGot {
				t.Errorf("Expected `%v` but got `%v`", messageWanted, messageGot)
			}
			if valueWanted != valueGot {
				t.Errorf("Expected %v but got %v", valueWanted, valueGot)
			}
		})

		t.Run("KnownTopic", func(t *testing.T) {
			var statusWanted int32 = successStatus
			var messageWanted string = "Success"
			var valueWanted int = 0

			result, err := node.publisherUpdate("test_caller", "/test_topic", pubList)
			if err != nil {
				t.Fatalf("Error calling publisherUpdate: %v", err)
			}

			res, ok := result.([]interface{})
			if !ok {
				t.Errorf("Result is malformed: Result should be an array of interface values.")
			}

			if len(res) != 3 {
				t.Errorf("Result is malformed: Length of result array should be 3")
			}

			statusGot, ok := res[0].(int32)
			if !ok {
				t.Errorf("Error getting status code from result. Type assertion failed.")
				return
			}

			messageGot, ok := res[1].(string)
			if !ok {
				t.Errorf("Error getting message from result. Type assertion failed.")
				return
			}
			valueGot, ok := res[2].(int)
			if !ok {
				t.Errorf("Error getting message from result. Type assertion failed.")
				return
			}

			if statusWanted != statusGot {
				t.Errorf("Expected %v but got %v", statusWanted, statusGot)
			}
			if messageWanted != messageGot {
				t.Errorf("Expected `%v` but got `%v`", messageWanted, messageGot)
			}
			if valueWanted != valueGot {
				t.Errorf("Expected %v but got %v", valueWanted, valueGot)
			}

			sub := node.subscribers["/test_topic"]

			select {
			case pubs := <-sub.pubListChan:
				if len(pubs) != len(pubList) {
					t.Errorf("Expected %v but got %v", pubs, pubList)
				}
			case <-time.After(2 * time.Second):
				t.Error("Did not receive publishers list within timeout period")
			}
		})
	})

	t.Run("RequestTopic", func(t *testing.T) {
		t.Run("UnknownTopic", func(t *testing.T) {
			var statusWanted int32 = failureStatus
			var messageWanted string = "No such topic"
			var valueWanted int = 0

			result, err := node.requestTopic("test_caller", "/unknown_topic", []interface{}{"TCPROS"})
			if err != nil {
				t.Fatalf("Error calling publisherUpdate: %v", err)
			}

			res, ok := result.([]interface{})
			if !ok {
				t.Errorf("Result is malformed: Result should be an array of interface values.")
			}

			if len(res) != 3 {
				t.Errorf("Result is malformed: Length of result array should be 3")
			}

			statusGot, ok := res[0].(int32)
			if !ok {
				t.Errorf("Error getting status code from result. Type assertion failed.")
				return
			}

			messageGot, ok := res[1].(string)
			if !ok {
				t.Errorf("Error getting message from result. Type assertion failed.")
				return
			}
			valueGot, ok := res[2].(int)
			if !ok {
				t.Errorf("Error getting message from result. Type assertion failed.")
				return
			}

			if statusWanted != statusGot {
				t.Errorf("Expected %v but got %v", statusWanted, statusGot)
			}
			if messageWanted != messageGot {
				t.Errorf("Expected `%v` but got `%v`", messageWanted, messageGot)
			}
			if valueWanted != valueGot {
				t.Errorf("Expected %v but got %v", valueWanted, valueGot)
			}
		})

		t.Run("KnownTopic", func(t *testing.T) {
			var statusWanted int32 = successStatus
			var messageWanted string = "Success"
			value := []interface{}{[]interface{}{"TCPROS"}}

			result, err := node.requestTopic("test_caller", "/test_topic", value)
			if err != nil {
				t.Fatalf("Error calling publisherUpdate: %v", err)
			}

			res, ok := result.([]interface{})
			if !ok {
				t.Errorf("Result is malformed: Result should be an array of interface values.")
			}

			if len(res) != 3 {
				t.Errorf("Result is malformed: Length of result array should be 3")
			}

			statusGot, ok := res[0].(int32)
			if !ok {
				t.Errorf("Error getting status code from result. Type assertion failed.")
				return
			}

			messageGot, ok := res[1].(string)
			if !ok {
				t.Errorf("Error getting message from result. Type assertion failed.")
				return
			}
			valueGot, ok := res[2].([]interface{})
			if !ok {
				t.Errorf("Error getting message from result. Type assertion failed.")
				return
			}

			if statusWanted != statusGot {
				t.Errorf("Expected %v but got %v", statusWanted, statusGot)
			}
			if messageWanted != messageGot {
				t.Errorf("Expected `%v` but got `%v`", messageWanted, messageGot)
			}
			if len(valueGot) != 3 {
				t.Errorf("Expected Protocol, host and port but got %v", valueGot)
			}
		})
	})

	t.Run("Shutdown", func(t *testing.T) {
		if !node.OK() {
			t.Errorf("Node is already shutdown")
		}

		var statusWanted int32 = successStatus
		var messageWanted string = "Success"
		var valueWanted int = 0

		result, err := node.shutdown("test_caller", "")
		if err != nil {
			t.Fatalf("Error calling getBusStats: %v", err)
		}

		res, ok := result.([]interface{})
		if !ok {
			t.Errorf("Result is malformed: Result should be an array of interface values.")
		}

		if len(res) != 3 {
			t.Errorf("Result is malformed: Length of result array should be 3")
		}

		statusGot, ok := res[0].(int32)
		if !ok {
			t.Errorf("Error getting status code from result. Type assertion failed.")
			return
		}

		messageGot, ok := res[1].(string)
		if !ok {
			t.Errorf("Error getting message from result. Type assertion failed.")
			return
		}
		valueGot, ok := res[2].(int)
		if !ok {
			t.Errorf("Error getting message from result. Type assertion failed.")
			return
		}

		if statusWanted != statusGot {
			t.Errorf("Expected %v but got %v", statusWanted, statusGot)
		}
		if messageWanted != messageGot {
			t.Errorf("Expected `%v` but got `%v`", messageWanted, messageGot)
		}
		if valueWanted != valueGot {
			t.Errorf("Expected %v but got %v", valueWanted, valueGot)
		}
	})
}

type dummyMessage struct {
}

func (m *dummyMessage) Text() string {
	return ""
}

func (m *dummyMessage) MD5Sum() string {
	return "d41d8cd98f00b204e9800998ecf8427e"
}

func (m *dummyMessage) Name() string {
	return "empty_msg"
}

func (m *dummyMessage) NewMessage() Message {
	return nil
}
