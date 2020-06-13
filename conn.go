package rediserver

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

type CommandType string

const (
	CommandTypeSet CommandType = "set"
	CommandTypeGet CommandType = "get"

	CommandTypeHset    CommandType = "hset"
	CommandTypeHget    CommandType = "hget"
	CommandTypeHgetAll CommandType = "hgetall"

	CommandTypeDel CommandType = "del"
)

type Command struct {
	Type CommandType
	Args [][]byte
}

func (c *Command) String() string {
	cmdParts := []string{string(c.Type)}
	for _, arg := range c.Args {
		cmdParts = append(cmdParts, string(arg))
	}
	return strings.Join(cmdParts, " ")
}

type conn struct {
	id     int
	c      net.Conn
	reader *bufio.Reader
}

func NewConn(id int, c net.Conn) *conn {
	r := bufio.NewReader(c)
	return &conn{
		id:     id,
		c:      c,
		reader: r,
	}
}
func (c *conn) Close() {
	c.c.Close()
}

func (c *conn) ReadCommand() (cmd *Command, err error) {
	buf := make([]byte, 1024)
	for i := 0; i < 1000000; i++ {
		n, err := c.c.Read(buf)
		if err != nil {
			return nil, err
		}
		if n == 0 {
			time.Sleep(time.Microsecond * 10)
			continue
		}
		return Rsfp2Cmd(buf[:n])
	}
	return nil, errors.New("none command")
}

func (c *conn) Write(resp []byte) error {
	resp = []byte("+" + string(resp) + "\r\n")
	for len(resp) > 0 {
		n, err := c.c.Write(resp)
		if err != nil {
			return err
		}
		resp = resp[n:]
	}
	return nil
}

func getFirstSeperator(data []byte) int {
	if len(data) < 2 {
		return -1
	}
	for i := 0; i < len(data)-1; i++ {
		if data[i] == '\r' && data[i+1] == '\n' {
			// fmt.Println("get seperator, buf = " + string(data))
			// fmt.Println("seperator = " + strconv.Itoa(i))
			return i
		}
	}
	return -1
}

func Rsfp2Cmd(buf []byte) (cmd *Command, err error) {
	// example: hget a b -> *3\r\n$4\r\nhget\r\n$5\r\nfield\r\n$3\r\nval\r\n
	if buf[0] != '*' {
		return nil, errors.New("invalid command: " + string(buf))
	}
	buf = buf[1:] // buf = 3\r\n$4\r\nhget\r\n$5\r\nfield\r\n$3\r\nval\r\n
	argCountSeperatorIdx := getFirstSeperator(buf)
	if argCountSeperatorIdx == -1 {
		return nil, errors.New("invalid command: " + string(buf))
	}
	argCount, err := strconv.Atoi(string(buf[:argCountSeperatorIdx]))
	if err != nil {
		fmt.Println("get arg count: " + string(buf[:argCountSeperatorIdx]))
		return nil, err
	}
	// fmt.Println("argcount = " + string(string(buf[:argCountSeperatorIdx])))

	buf = buf[argCountSeperatorIdx+2:] // buf = $4\r\nhget\r\n$1\r\na\r\n$b\r\n

	args := make([][]byte, 0)
	for i := 0; i < argCount; i++ {
		if buf[0] != '$' {
			return nil, errors.New("invalid command: " + string(buf))
		}
		argLenSeperatorIdx := getFirstSeperator(buf)
		if argLenSeperatorIdx == -1 {
			return nil, errors.New("invalid command: " + string(buf))
		}
		// fmt.Println("arg = " + string(buf[1:argLenSeperatorIdx]))
		argLen, err := strconv.Atoi(string(buf[1:argLenSeperatorIdx]))
		if err != nil {
			return nil, errors.New("invalid command: " + string(buf))
		}
		buf = buf[argLenSeperatorIdx+2:]

		arg := buf[:argLen]
		// fmt.Println("arg = " + string(buf[:argLen]))
		args = append(args, arg)
		buf = buf[argLen+2:]
	}
	cmd = &Command{
		Type: CommandType(strings.ToLower(string(args[0]))),
		Args: args[1:],
	}
	return cmd, nil
}
