package ros

import (
	"time"
)

// Node defines the interface that a ros node should implement
type Node interface {
	// NewPublisher creates a publisher which can used to publish ros messages of type MessageType
	// to the specified topic.
	NewPublisher(topic string, msgType MessageType) Publisher

	// NewPublisherWithCallbacks creates a publisher which gives you callbacks when subscribers
	// connect and disconnect.  The callbacks are called in their own goroutines, so they don't
	// need to return immediately to let the connection proceed.
	NewPublisherWithCallbacks(topic string, msgType MessageType, connectCallback, disconnectCallback func(SingleSubscriberPublisher)) Publisher

	// NewSubscriber creates a subscriber to a topic and calls callback on receiving a message.
	// Callback should be a function which takes 0, 1, or 2 arguments.
	//
	// 0-arguments - Callback will simply be called without the message.
	// 1-arguments - Callback argument should be of the generated message type.
	// 2-arguments - Callback first argument should be of the generated message type and
	//               the second argument should be of type MessageEvent.
	NewSubscriber(topic string, msgType MessageType, callback interface{}) Subscriber

	// NewServiceClient creates a service client which can be used to connect to a service server
	// send service requests.
	NewServiceClient(service string, srvType ServiceType) ServiceClient

	// NewServiceServer creates a service server that advertises the server and responds to the
	// service requests from service clients.
	NewServiceServer(service string, srvType ServiceType, callback interface{}) ServiceServer

	// OK represents the status of ros node.
	OK() bool

	// SpinOnce executes the job at the top of the Job queue.
	// Job queue consists of callback functions for subscribers and service server.
	// Blocking callbacks or callbacks that take too long to execute will result in new messages
	// being dropped due to job queue being full.
	SpinOnce()

	// Spin is job executor that continuosly starts executing callback jobs in the job queue.
	// Spin is a blocking call that unblocks only on node shutdown.
	Spin()

	// Shutdown stops the ros node
	Shutdown()

	// GetParam gets parameter value from the parameter server it it exists.
	// Non nil error will be returned if the parameter doesn't exists or the node
	// is unable to send the get param request to the paramter server.
	GetParam(name string) (interface{}, error)

	// SetParam sets parameter value on the parameter server which can be
	// retrieved or updated using parameter server API.
	// Returns a non nil error when unable to set the param on the server.
	SetParam(name string, value interface{}) error

	// HasParam checks if the parameter exists on the parameter server.
	// Return a non nil error when the node is unable to connect to paramter server.
	HasParam(name string) (bool, error)

	// SearchParam searches for the param in caller's namespace and proceeds upwards
	// through parent namespaces until Parameter Server finds a matching key.
	// More about parameter server API: https://wiki.ros.org/ROS/Parameter%20Server%20API
	SearchParam(name string) (string, error)

	// DeleteParam deletes the parameter from the parameter server.
	DeleteParam(name string) error

	// Logger returns the logger being used by the ros node.
	Logger() Logger

	// SetLogger changes the default logger used by ros node to the one specified in the call.
	// Custom loggers that implement the Logger interface can be set as the default logger.
	SetLogger(Logger)

	// NonRosArgs returns an array of all the non ros arguments of the ros node.
	NonRosArgs() []string
	Name() string
}

// NewNode creates and returns a new instance of ros node
// Returns a non nil error when unable connect to ros master
// or create a new node.
func NewNode(name string, args []string) (Node, error) {
	return newDefaultNode(name, args)
}

// Publisher can be used to publish a ros messages to a specific topic.
// Shutdown can be used to shutdown this particular publisher.
// All running subscribers, publishers and services are shutdown on node exit.
type Publisher interface {
	// Publish publishes ros message
	Publish(msg Message)

	// GetNumSubscribers gets the number of subscribers to the publishing topic
	GetNumSubscribers() int

	// Shutdown stops the publisher
	Shutdown()
}

// SingleSubscriberPublisher is publisher which only sends to one specific subscriber.
// This is sent as an argument to the connect and disconnect callback functions passed
// to Node.NewPublisherWithCallbacks().
type SingleSubscriberPublisher interface {
	// Publish publishes ros message
	Publish(msg Message)

	// GetSubscriberName gets the name of the single subscriber
	GetSubscriberName() string

	// GetTopic gets the name of the topic that we are publishing to.
	GetTopic() string
}

// Subscriber can be used to get the number of publishers to topic or shutdown
// the subscriber. Messages are received using the callback passed during subscriber creation.
type Subscriber interface {
	// GetNumPublishers gets the numbers of publishers to the topic we are subscribed to has.
	GetNumPublishers() int

	// Shutdown stop subscriber.
	Shutdown()
}

// MessageEvent defines the event information struct that contains meta information
// for each ros message received from a publisher.
type MessageEvent struct {
	// PublisherName is the name of message publisher.
	PublisherName string

	// ReceiptTime is the timestamp of when the message was received.
	ReceiptTime time.Time

	// ConnectionHeader is the HTTP connection header of the connection between
	// the subscriber node and the publisher node.
	ConnectionHeader map[string]string
}

// ServiceServer can be used to shutdown a running service server
type ServiceServer interface {
	// Shutdown stops the service server.
	Shutdown()
}

// ServiceClient can be used to call a service server with a service request.
// It can also be used to shutdown the client.
type ServiceClient interface {
	// Call calls a service server with a service request.
	Call(srv Service) error

	// Shutdown stops the service client.
	Shutdown()
}
