package rediserver

import (
	"fmt"
	"net"
	"strconv"
	"sync"
)

type Conf struct {
	Port int `yaml:"port"`
}

type Handler func(cmd *Command) (resp []byte, err error)

type Server struct {
	*Conf
	ln       net.Listener
	connMap  map[int]*conn
	mutex    sync.Mutex
	handler  Handler
	newestID int
}

func (s *Server) Close() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for _, conn := range s.connMap {
		conn.Close()
	}
}

func (s *Server) AddConn(c net.Conn) *conn {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.newestID++
	conn := NewConn(s.newestID, c)
	s.connMap[s.newestID] = conn
	return conn
}

func (s *Server) RemoveConn(id int) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.connMap[id].Close()
	delete(s.connMap, id)
}

func ListenAndServe(port int, handler Handler) error {
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		return err
	}
	s := Server{
		Conf: &Conf{
			Port: port,
		},
		ln:      listener,
		connMap: make(map[int]*conn),
		handler: handler,
	}
	for {
		c, err := listener.Accept()
		if err != nil {
			s.Close()
		}
		conn := s.AddConn(c)
		go func() {
			defer s.RemoveConn(conn.id)
			if err = conn.Write([]byte("OK")); err != nil {
				fmt.Printf("err = %s\n", err.Error())
				return
			}
			for {
				cmd, err := conn.ReadCommand()
				if err != nil {
					fmt.Printf("err = %s\n", err.Error())
					return
				}
				resp, err := handler(cmd)
				if err != nil {
					fmt.Printf("err = %s\n", err.Error())
					return
				}
				if err := conn.Write(resp); err != nil {
					fmt.Printf("err = %s\n", err.Error())
					return
				}
			}
		}()
	}
}
