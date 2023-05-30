package server

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/sirupsen/logrus"
)

const (
	CommandPing = 101
	CommandPong = 102
)

type ServerOptions struct {
	writewait time.Duration
	readwait  time.Duration
}

type Server struct {
	once    sync.Once
	options ServerOptions
	id      string
	address string
	sync.Mutex
	users map[string]net.Conn
}

func NewServer(id, address string) *Server {
	return newServer(id, address)
}

func newServer(id, address string) *Server {
	return &Server{
		id:      id,
		address: address,
		users:   make(map[string]net.Conn, 1000),
		options: ServerOptions{
			writewait: time.Second * 10,
			readwait:  time.Minute * 2,
		},
	}
}

func (s *Server) Start() error {
	mux := http.NewServeMux()
	log := logrus.WithFields(logrus.Fields{
		"module": "server",
		"listen": s.address,
		"id":     s.id,
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// 1. update protocol
		conn, _, _, err := ws.UpgradeHTTP(r, w)
		if err != nil {
			conn.Close()
			return
		}
		log.Infoln("update protocol")

		// 2. get userId
		user := r.URL.Query().Get("user")
		if user == "" {
			conn.Close()
			return
		}
		log.Infoln("get userId")

		// 3. add Chan
		old, ok := s.addUser(user, conn)
		if ok {
			old.Close()
		}
		log.Infof("add user: %s", user)

		// 4. run Chan
		go func(user string, conn net.Conn) {
			defer func() {
				log.Infof("connection of %s closed", conn)
				conn.Close()
				s.delUser(user)
			}()

			err := s.readLoop(user, conn)
			if err != nil {
				log.Error(err)
			}

		}(user, conn)
	})

	return http.ListenAndServe(s.address, mux)
}

func (s *Server) addUser(user string, conn net.Conn) (net.Conn, bool) {
	s.Lock()
	defer s.Unlock()

	old, ok := s.users[user]
	s.users[user] = conn

	return old, ok
}

func (s *Server) delUser(user string) {
	s.Lock()
	defer s.Unlock()
	delete(s.users, user)
}

func (s *Server) Shutdown() {
	s.once.Do(func() {
		s.Lock()
		defer s.Unlock()
		for _, conn := range s.users {
			conn.Close()
		}
	})
}

func (s *Server) readLoop(user string, conn net.Conn) error {
	for {

		_ = conn.SetReadDeadline(time.Now().Add(s.options.readwait))

		frame, err := ws.ReadFrame(conn)
		if err != nil {
			return err
		}

		if frame.Header.OpCode == ws.OpPing {
			_ = wsutil.WriteServerMessage(conn, ws.OpPong, nil)
			logrus.Info("write back a pong...")
			continue
		}

		if frame.Header.OpCode == ws.OpClose {
			return errors.New("user side close the conn")
		}

		if frame.Header.Masked {
			ws.Cipher(frame.Payload, frame.Header.Mask, 0)
		}

		if frame.Header.OpCode == ws.OpText {
			go s.handle(user, string(frame.Payload))
		} else if frame.Header.OpCode == ws.OpBinary {
			go s.handleBinary(user, frame.Payload)
		}
	}
}

func (s *Server) handle(user, msg string) {
	logrus.Infof("recv msg: %s from user: %s", msg, user)
	s.Lock()
	defer s.Unlock()

	broadcast := fmt.Sprintf("%s ----from %s", msg, user)

	for u, conn := range s.users {
		if u == user {
			continue
		}
		logrus.Infof("%s send to %s: %s", user, u, broadcast)
		err := s.writeText(conn, broadcast)
		if err != nil {
			logrus.Errorf("write to %s failed, error: %v", u, err)
		}
	}
}

func (s *Server) handleBinary(user string, msg []byte) {
	logrus.Infof("recv msg %v from %s", msg, user)
	s.Lock()
	defer s.Unlock()

	i := 0
	command := binary.BigEndian.Uint16(msg[i : i+2]) // [0 2)
	i += 2
	payloadLen := binary.BigEndian.Uint32(msg[i : i+4]) // [2 6)
	logrus.Info("command: %v payloadLen: %v", command, payloadLen)

	if command == CommandPing {
		u := s.users[user]
		err := wsutil.WriteServerBinary(u, []byte{0, CommandPong, 0, 0, 0, 0})
		if err != nil {
			logrus.Errorf("write to %s failed, error %v", user, err)
		}
	}
}

func (s *Server) writeText(conn net.Conn, msg string) error {
	return ws.WriteFrame(conn, ws.NewTextFrame([]byte(msg)))
}
