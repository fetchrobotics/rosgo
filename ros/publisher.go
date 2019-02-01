package ros

import (
	"bytes"
	"container/list"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"net"
	"sync"
	"time"
)

type remoteSubscriberSessionError struct {
	session *remoteSubscriberSession
	err     error
}

func (e *remoteSubscriberSessionError) Error() string {
	return fmt.Sprintf("remoteSubscriberSession %v error: %v", e.session, e.err)
}

type defaultPublisher struct {
	node               *defaultNode
	topic              string
	msgType            MessageType
	msgChan            chan []byte
	shutdownChan       chan struct{}
	sessions           *list.List
	sessionChan        chan *remoteSubscriberSession
	sessionErrorChan   chan error
	listenerErrorChan  chan error
	listener           net.Listener
	connectCallback    func(SingleSubscriberPublisher)
	disconnectCallback func(SingleSubscriberPublisher)
}

func newDefaultPublisher(node *defaultNode,
	topic string, msgType MessageType,
	connectCallback, disconnectCallback func(SingleSubscriberPublisher)) *defaultPublisher {
	pub := new(defaultPublisher)
	pub.node = node
	pub.topic = topic
	pub.msgType = msgType
	pub.shutdownChan = make(chan struct{}, 10)
	pub.msgChan = make(chan []byte, 10)
	pub.listenerErrorChan = make(chan error, 10)
	pub.sessionChan = make(chan *remoteSubscriberSession, 10)
	pub.sessionErrorChan = make(chan error, 10)
	pub.sessions = list.New()
	pub.sesssionIDCount = 0
	pub.connectCallback = connectCallback
	pub.disconnectCallback = disconnectCallback
	if listener, err := net.Listen("tcp", ":0"); err != nil {
		panic(err)
	} else {
		pub.listener = listener
	}
	return pub
}

func (pub *defaultPublisher) start(wg *sync.WaitGroup) {
	logger := pub.node.logger
	logger.Debugf("Publisher goroutine for %s started.", pub.topic)
	wg.Add(1)
	defer func() {
		logger.Debug("defaultPublisher.start exit")
		wg.Done()
	}()

	go pub.listenRemoteSubscriber()

	for {
		logger.Debug("defaultPublisher.start loop")
		select {
		case msg := <-pub.msgChan:
			logger.Debug("Receive msgChan")
			for e := pub.sessions.Front(); e != nil; e = e.Next() {
				session := e.Value.(*remoteSubscriberSession)
				session.msgChan <- msg
			}
		case err := <-pub.listenerErrorChan:
			logger.Debug("Listener closed unexpectedly: %s", err)
			pub.listener.Close()
			return
		case s := <-pub.sessionChan:
			pub.sessions.PushBack(s)
			go s.start()
		case err := <-pub.sessionErrorChan:
			logger.Error(err)
			if sessionError, ok := err.(*remoteSubscriberSessionError); ok {
				for e := pub.sessions.Front(); e != nil; e = e.Next() {
					if e.Value == sessionError.session {
						pub.sessions.Remove(e)
						break
					}
				}
			}
		case <-pub.shutdownChan:
			logger.Debug("defaultPublisher.start Receive shutdownChan")
			pub.listener.Close()
			logger.Debug("defaultPublisher.start closed listener")
			_, err := callRosApi(pub.node.masterUri, "unregisterPublisher", pub.node.qualifiedName, pub.topic, pub.node.xmlrpcUri)
			if err != nil {
				logger.Warn(err)
			}
			for e := pub.sessions.Front(); e != nil; e = e.Next() {
				session := e.Value.(*remoteSubscriberSession)
				session.quitChan <- struct{}{}
			}
			pub.sessions.Init() // Clear all sessions
			return
		}
	}
}

func (pub *defaultPublisher) listenRemoteSubscriber() {
	logger := pub.node.logger
	logger.Infof("Start listen %s.", pub.listener.Addr().String())
	defer func() {
		logger.Debug("defaultPublisher.listenRemoteSubscriber exit")
	}()

	for {
		logger.Debug("defaultPublisher.listenRemoteSubscriber loop")
		conn, err := pub.listener.Accept()
		if err != nil {
			logger.Debugf("pub.listner.Accept() failed")
			pub.listenerErrorChan <- err
			close(pub.listenerErrorChan)
			logger.Debugf("defaultPublisher.listenRemoteSubscriber loop exit")
			return
		}

		logger.Debugf("Connected %s", conn.RemoteAddr().String())
		session := newRemoteSubscriberSession(pub, conn)
		pub.sessionChan <- session
	}
}

func (pub *defaultPublisher) Publish(msg Message) {
	var buf bytes.Buffer
	_ = msg.Serialize(&buf)
	pub.msgChan <- buf.Bytes()
}

func (pub *defaultPublisher) Shutdown() {
	pub.shutdownChan <- struct{}{}
}

func (pub *defaultPublisher) hostAndPort() (string, string) {
	_, port, err := net.SplitHostPort(pub.listener.Addr().String())
	if err != nil {
		// Not reached
		panic(err)
	}
	return pub.node.hostname, port
}

func (pub *defaultPublisher) getPublisherStats() (uint32, []interface{}) {
	var msgDataSent uint32
	pubStats := []interface{}{}
	for e := pub.sessions.Front(); e != nil; e = e.Next() {
		session := e.Value.(*remoteSubscriberSession)
		stat := []interface{}{
			session.id,
			session.sizeBytesSent + session.msgBytesSent,
			session.numSent,
			true,
		}
		msgDataSent += session.msgBytesSent
		pubStats = append(pubStats, stat)
	}
	return msgDataSent, pubStats
}

func (pub *defaultPublisher) getPublisherInfo() []interface{} {
	pubInfo := []interface{}{}
	for e := pub.sessions.Front(); e != nil; e = e.Next() {
		session := e.Value.(*remoteSubscriberSession)
		stat := []interface{}{
			session.id,
			session.callerId,
			"o",
			"TCPROS",
			session.topic,
			true,
		}
		pubInfo = append(pubInfo, stat)
	}
	return pubInfo
}

type remoteSubscriberSession struct {
	id                 int
	conn               net.Conn
	nodeId             string
	callerId           string
	topic              string
	typeText           string
	md5sum             string
	typeName           string
	sizeBytesSent      uint32
	msgBytesSent       uint32
	numSent            int64
	quitChan           chan struct{}
	msgChan            chan []byte
	errorChan          chan error
	logger             Logger
	connectCallback    func(SingleSubscriberPublisher)
	disconnectCallback func(SingleSubscriberPublisher)
}

func newRemoteSubscriberSession(pub *defaultPublisher, id int, conn net.Conn) *remoteSubscriberSession {
	session := new(remoteSubscriberSession)
	session.id = id
	session.conn = conn
	session.nodeId = pub.node.qualifiedName
	session.topic = pub.topic
	session.typeText = pub.msgType.Text()
	session.md5sum = pub.msgType.MD5Sum()
	session.typeName = pub.msgType.Name()
	session.sizeBytesSent = 0
	session.msgBytesSent = 0
	session.numSent = 0
	session.quitChan = make(chan struct{})
	session.msgChan = make(chan []byte, 10)
	session.errorChan = pub.sessionErrorChan
	session.logger = pub.node.logger
	session.connectCallback = pub.connectCallback
	session.disconnectCallback = pub.disconnectCallback
	return session
}

type singleSubPub struct {
	subName string
	topic   string
	msgChan chan []byte
}

func (ssp *singleSubPub) Publish(msg Message) {
	var buf bytes.Buffer
	_ = msg.Serialize(&buf)
	ssp.msgChan <- buf.Bytes()
}

func (ssp *singleSubPub) GetSubscriberName() string {
	return ssp.subName
}

func (ssp *singleSubPub) GetTopic() string {
	return ssp.topic
}

func (session *remoteSubscriberSession) start() {
	logger := session.logger
	logger.Debug("remoteSubscriberSession.start enter")

	ssp := &singleSubPub{
		topic:   session.topic,
		msgChan: session.msgChan,
		// callerId is filled in after header gets read later in this function.
	}

	defer func() {
		logger.Debug("remoteSubscriberSession.start exit")

		if session.disconnectCallback != nil {
			session.disconnectCallback(ssp)
		}
	}()
	defer func() {
		if err := recover(); err != nil {
			if e, ok := err.(error); ok {
				session.errorChan <- &remoteSubscriberSessionError{session, e}
			} else {
				e = fmt.Errorf("Unkonwn error value")
				session.errorChan <- &remoteSubscriberSessionError{session, e}
			}
		} else {
			e := fmt.Errorf("Normal exit")
			session.errorChan <- &remoteSubscriberSessionError{session, e}
		}
	}()
	// 1. Read connection header
	headers, err := readConnectionHeader(session.conn)
	if err != nil {
		panic(errors.New("Failed to read connection header."))
	}
	logger.Info("TCPROS Connection Header:")
	headerMap := make(map[string]string)
	for _, h := range headers {
		headerMap[h.key] = h.value
		logger.Infof("  `%s` = `%s`", h.key, h.value)
	}
	if headerMap["type"] != session.typeName || headerMap["md5sum"] != session.md5sum {
		panic(errors.New("Incomatible message type!"))
	}
	session.callerId = headerMap["callerid"]
	ssp.subName = headerMap["callerid"]
	if session.connectCallback != nil {
		go session.connectCallback(ssp)
	}

	// 2. Return reponse header
	var resHeaders []header
	resHeaders = append(resHeaders, header{"message_definition", session.typeText})
	resHeaders = append(resHeaders, header{"callerid", session.nodeId})
	resHeaders = append(resHeaders, header{"latching", "0"})
	resHeaders = append(resHeaders, header{"md5sum", session.md5sum})
	resHeaders = append(resHeaders, header{"topic", session.topic})
	resHeaders = append(resHeaders, header{"type", session.typeName})
	logger.Debug("TCPROS Response Header")
	for _, h := range resHeaders {
		logger.Debugf("  `%s` = `%s`", h.key, h.value)
	}
	err = writeConnectionHeader(resHeaders, session.conn)
	if err != nil {
		panic(errors.New("Failed to write response header."))
	}

	// 3. Start sending message
	logger.Debug("Start sending messages...")
	queueMaxSize := 100
	queue := make(chan []byte, queueMaxSize)
	for {
		//logger.Debug("session.remoteSubscriberSession")
		select {
		case msg := <-session.msgChan:
			logger.Debug("Receive msgChan")
			if len(queue) == queueMaxSize {
				<-queue
			}
			queue <- msg

		case <-session.quitChan:
			logger.Debug("Receive quitChan")
			return

		case msg := <-queue:
			logger.Debug("writing")
			logger.Debug(hex.EncodeToString(msg))
			session.conn.SetDeadline(time.Now().Add(10 * time.Millisecond))
			size := uint32(len(msg))
			if err := binary.Write(session.conn, binary.LittleEndian, size); err != nil {
				if neterr, ok := err.(net.Error); ok && neterr.Timeout() {
					logger.Debug("timeout")
					continue
				} else {
					logger.Error(err)
					panic(err)
				}
			}
			logger.Debug(len(msg))
			session.conn.SetDeadline(time.Now().Add(10 * time.Millisecond))
			if _, err := session.conn.Write(msg); err != nil {
				if neterr, ok := err.(net.Error); ok && neterr.Timeout() {
					logger.Debug("timeout")
					continue
				} else {
					logger.Error(err)
					panic(err)
				}
			}
			logger.Debug(hex.EncodeToString(msg))
		}
	}
}
