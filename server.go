package netkit

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/brotherhoodhk/net-kit/transport"
)

type Server interface {
	Send(Request) error
	Accept() (Response, error)
	ListenAndServe(content chan<- DataContent, errchan chan<- error)
	Close() error
}
type Request struct {
}
type Response struct {
}
type serverconn struct {
	listener     net.Listener
	maxerrortime int32
}

func NewServer(add string, port int) Server {
	con, err := net.Listen("tcp", add+":"+strconv.Itoa(port))
	if err == nil {
		return &serverconn{listener: con}
	}
	return nil
}
func (s *serverconn) Send(request Request) error {
	return nil
}

func (s *serverconn) Accept() (Response, error) {
	return Response{}, nil
}

func Listen[T any](s *serverconn, content chan<- []byte, errs chan<- error) {
	buffer := make([]byte, RDWR_BUFFER_SIZE)
	errtime := 0
	var tiobj transport.Common_Proto
	for true {
		con, err := s.listener.Accept()
		if err == nil {
			lang, err := con.Read(buffer)
			if err == nil {
				err = json.Unmarshal(buffer[:lang], &tiobj)
				if err == nil {
					content <- tiobj.Content
				}
			}
		}
		if err != nil {
			if errtime == Max_Error_Time { //when it touch max error times ,process will exit,-1 is never stop
				break
			}
			errtime++
			errs <- err
		}
	}
}
func (s *serverconn) ListenAndServe(content chan<- DataContent, errchan chan<- error) {
	buffer := make([]byte, RDWR_BUFFER_SIZE)
	var errtime int32 = 0
	var tiobj transport.Common_Proto
	var isbreak bool = true
	for isbreak {
		con, err := s.listener.Accept()
		if err == nil {
			fmt.Println("accept message success")
			conerrtime := Single_Max_Err_Times
			for conerrtime > 0 {
				fmt.Println("err time", errtime)
				_, err := con.Write([]byte(""))
				if err != nil {
					//lost connection with client
					con.Close()
					fmt.Println("close connection")
					conerrtime = -1
					break
				}
				fmt.Println("arrive read")
				lang, err := con.Read(buffer)
				fmt.Println("finish read")
				if err == nil {
					err = json.Unmarshal(buffer[:lang], &tiobj)
					if err == nil {
						fmt.Println("accept message") //debugline
						content <- DataContent{Data: tiobj.Content, Types: strings.Split(tiobj.Content_Type, "/")[1]}
					}
				}
				if err != nil {
					if errtime == s.maxerrortime { //when it touch max error times ,process will exit,-1 is never stop
						isbreak = false
						break
					}
					conerrtime--
					errtime++
					errchan <- err
				}
				fmt.Println("arrive end")
			}
		} else {
			fmt.Println("accept message failed")
			errtime++
			errchan <- err
		}

	}
}
func (s *serverconn) Close() error {
	return s.Close()
}
