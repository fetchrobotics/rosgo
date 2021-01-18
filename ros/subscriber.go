package ros

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"reflect"
	"sync"
	"time"
)

type messageEvent struct {
	bytes []byte
	event MessageEvent
}

// The subscription object runs in own goroutine (startSubscription).
// Do not access any properties from other goroutine.
type defaultSubscriber struct {
	topic            string
	msgType          MessageType
	pubList          []string
	pubListChan      chan []string
	msgChan          chan messageEvent
	callbacks        []interface{}
	addCallbackChan  chan interface{}
	shutdownChan     chan struct{}
	connections      map[string]chan struct{}
	disconnectedChan chan string
}

func newDefaultSubscriber(topic string, msgType MessageType, callback interface{}) *defaultSubscriber {
	sub := new(defaultSubscriber)
	sub.topic = topic
	sub.msgType = msgType
	sub.msgChan = make(chan messageEvent, 10)
	sub.pubListChan = make(chan []string, 10)
	sub.addCallbackChan = make(chan interface{}, 10)
	sub.shutdownChan = make(chan struct{}, 10)
	sub.disconnectedChan = make(chan string, 10)
	sub.connections = make(map[string]chan struct{})
	sub.callbacks = []interface{}{callback}
	return sub
}

func (sub *defaultSubscriber) start(wg *sync.WaitGroup, nodeID string, nodeURI string, masterURI string, jobChan chan func(), logger Logger, unregisterFromNode func()) {
	logger.Debugf("Subscriber goroutine for %s started.", sub.topic)
	wg.Add(1)
	defer wg.Done()
	defer func() {
		logger.Debug("defaultSubscriber.start exit")
	}()
	for {
		logger.Debug("Loop")
		select {
		case list := <-sub.pubListChan:
			logger.Debug("Receive pubListChan")
			deadPubs := setDifference(sub.pubList, list)
			newPubs := setDifference(list, sub.pubList)
			sub.pubList = list

			for _, pub := range deadPubs {
				quitChan := sub.connections[pub]
				quitChan <- struct{}{}
				delete(sub.connections, pub)
			}

			for _, pub := range newPubs {
				protocols := []interface{}{[]interface{}{"TCPROS"}}
				result, err := callRosAPI(pub, "requestTopic", nodeID, sub.topic, protocols)
				if err != nil {
					logger.Fatalf("[DefaultSubscriber] %v", err)
					continue
				}

				protocolParams := result.([]interface{})
				for _, x := range protocolParams {
					logger.Debug(x)
				}

				name := protocolParams[0].(string)
				if name == "TCPROS" {
					addr := protocolParams[1].(string)
					port := protocolParams[2].(int32)
					uri := fmt.Sprintf("%s:%d", addr, port)
					quitChan := make(chan struct{}, 10)
					sub.connections[pub] = quitChan
					go startRemotePublisherConn(logger,
						uri, sub.topic,
						sub.msgType.MD5Sum(),
						sub.msgType.Name(), nodeID,
						sub.msgChan,
						quitChan,
						sub.disconnectedChan)
				} else {
					logger.Warnf("rosgo Not support protocol '%s'", name)
				}
			}

		case callback := <-sub.addCallbackChan:
			logger.Debug("Receive addCallbackChan")
			sub.callbacks = append(sub.callbacks, callback)

		case msgEvent := <-sub.msgChan:
			// Pop received message then bind callbacks and enqueue to the job channle.
			logger.Debug("Receive msgChan")
			callbacks := make([]interface{}, len(sub.callbacks))
			copy(callbacks, sub.callbacks)
			jobChan <- func() {
				m := sub.msgType.NewMessage()
				reader := NewReader(msgEvent.bytes)
				if err := m.Deserialize(reader); err != nil {
					logger.Error(err)
				}
				args := []reflect.Value{reflect.ValueOf(m), reflect.ValueOf(msgEvent.event)}
				for _, callback := range callbacks {
					fun := reflect.ValueOf(callback)
					numArgsNeeded := fun.Type().NumIn()
					if numArgsNeeded <= 2 {
						fun.Call(args[0:numArgsNeeded])
					}
				}
			}
			logger.Debug("Callback job enqueued.")

		case pubURI := <-sub.disconnectedChan:
			logger.Debugf("Connection to %s was disconnected.", pubURI)
			delete(sub.connections, pubURI)

		case <-sub.shutdownChan:
			// Shutdown subscription goroutine
			logger.Debug("Receive shutdownChan")
			for _, closeChan := range sub.connections {
				closeChan <- struct{}{}
				close(closeChan)
			}
			_, err := callRosAPI(masterURI, "unregisterSubscriber", nodeID, sub.topic, nodeURI)
			if err != nil {
				logger.Warn(err)
			}

			unregisterFromNode()
			return
		}
	}
}

func startRemotePublisherConn(logger Logger,
	pubURI string, topic string, md5sum string,
	msgType string, nodeID string,
	msgChan chan messageEvent,
	quitChan chan struct{},
	disconnectedChan chan string) {
	logger.Debug("startRemotePublisherConn()")

	defer func() {
		logger.Debug("startRemotePublisherConn() exit")
	}()

	conn, err := net.Dial("tcp", pubURI)
	if err != nil {
		logger.Fatalf("Failed to connect %s!", pubURI)
	}

	// 1. Write connection header
	var headers []header
	headers = append(headers, header{"topic", topic})
	headers = append(headers, header{"md5sum", md5sum})
	headers = append(headers, header{"type", msgType})
	headers = append(headers, header{"callerid", nodeID})
	logger.Debug("TCPROS Connection Header")
	for _, h := range headers {
		logger.Debugf("  `%s` = `%s`", h.key, h.value)
	}
	err = writeConnectionHeader(headers, conn)
	if err != nil {
		logger.Fatal("Failed to write connection header.")
	}

	// 2. Read reponse header
	var resHeaders []header
	resHeaders, err = readConnectionHeader(conn)
	if err != nil {
		logger.Fatal("Failed to read reasponse header.")
	}
	logger.Debug("TCPROS Response Header:")
	resHeaderMap := make(map[string]string)
	for _, h := range resHeaders {
		resHeaderMap[h.key] = h.value
		logger.Debugf("  `%s` = `%s`", h.key, h.value)
	}

	if md5sum != resHeaderMap["md5sum"] && md5sum != "*" {
		panic("incompatible message type: md5sum mismatch")
	}

	logger.Debug("Start receiving messages...")
	event := MessageEvent{ // Event struct to be sent with each message.
		PublisherName:    resHeaderMap["callerid"],
		ConnectionHeader: resHeaderMap,
	}

	// 3. Start reading messages
	readingSize := true
	var msgSize uint32
	var buffer []byte
	for {
		select {
		case <-quitChan:
			return
		default:
			conn.SetDeadline(time.Now().Add(1000 * time.Millisecond))
			if readingSize {
				//logger.Debug("Reading message size...")
				err := binary.Read(conn, binary.LittleEndian, &msgSize)
				if err != nil {
					if err == io.EOF {
						logger.Infof("Publisher %s on topic %s disconnected", pubURI, topic)
						disconnectedChan <- pubURI
						return
					}
					if neterr, ok := err.(net.Error); ok && neterr.Timeout() {
						// Timed out
						//logger.Debug(neterr)
						continue
					} else {
						logger.Error("Failed to read a message size", err)
						disconnectedChan <- pubURI
						return
					}
				}
				logger.Debugf("  %d", msgSize)
				buffer = make([]byte, int(msgSize))
				readingSize = false
			} else {
				//logger.Debug("Reading message body...")
				_, err = io.ReadFull(conn, buffer)
				if err != nil {
					if err == io.EOF {
						logger.Info("Publisher disconnected")
						disconnectedChan <- pubURI
						return
					}
					if neterr, ok := err.(net.Error); ok && neterr.Timeout() {
						// Timed out
						//logger.Debug(neterr)
						continue
					} else {
						logger.Error("Failed to read a message body", err)
						disconnectedChan <- pubURI
						return
					}
				}
				event.ReceiptTime = time.Now()
				msgChan <- messageEvent{bytes: buffer, event: event}
				readingSize = true
			}
		}
	}
}

func (sub *defaultSubscriber) Shutdown() {
	sub.shutdownChan <- struct{}{}
}

func (sub *defaultSubscriber) GetNumPublishers() int {
	return len(sub.pubList)
}
