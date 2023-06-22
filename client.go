package netkit

import (
	"encoding/json"
	"net"
	"reflect"
	"strconv"

	"github.com/brotherhoodhk/net-kit/transport"
)

type Client interface {
	Send(val any) error
	Accept(any) error
	Close() error
}
type client struct {
	con net.Conn
}

func NewClient(add string, port int) Client {
	scon, err := net.Dial("tcp", add+":"+strconv.Itoa(port))
	if err == nil {
		return &client{con: scon}
	}
	return nil
}

func (s *client) Send(val any) error {
	var typename string
	types := reflect.TypeOf(val)
	if types.Kind() == reflect.Ptr {
		typename = "any"
	} else {
		typename = types.Name()
		if len(typename) < 1 {
			typename = "any"
		}
	}
	valbin, err := json.Marshal(val)
	if err == nil {
		obj := &transport.Common_Proto{Content_Type: "json/" + typename, Content: valbin}
		content, err := json.Marshal(obj)
		if err == nil {
			_, err = s.con.Write(content)
		}
	}

	return err
}
func (s *client) Accept(origin any) error {
	buffer := make([]byte, Max_Error_Time)
	lang, err := s.con.Read(buffer)
	if err == nil {
		comproto := transport.Common_Proto{}
		err = json.Unmarshal(buffer[:lang], comproto)
		if err == nil {
			origin = comproto.Content
		}
	}
	return err
}
func (s *client) Close() error {
	return s.con.Close()
}
