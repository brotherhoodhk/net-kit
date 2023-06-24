package bs

import (
	"encoding/json"
	"fmt"
	"net"
	"regexp"
	"strconv"
	"time"

	"github.com/oswaldoooo/cmirco/kits"
)

const (
	url_reg = "^bs(t|u|)://(.*):([\\d]{1,5})/(.*)"
)

var (
	Max_Cut_Time = 10
)

type BsClient struct {
	Socket_Type string
	Address     string
	Port        int
	Path        string
	Type        string
	net.Conn
	isclose     bool
	Buffer_Size int
	req         Request
	rep         chan response
	msgchan     chan any  //the content need to send
	lastconn    time.Time //the last one connection time
}

// return nil is create client failed,[fullurl,socket_type,address,port,path]
func CreateClient(url string) *BsClient {
	reg := regexp.MustCompile(url_reg)
	params := reg.FindStringSubmatch(url)
	if len(params) > 0 {
		port, _ := strconv.Atoi(params[3])
		sct := ""
		switch params[1] {
		case "t":
			sct = "tcp"
		case "u":
			sct = "udp"
		case "":
			sct = "unix"
		default:
			return nil
		}
		return &BsClient{Socket_Type: sct, Address: params[2], Port: port, Path: params[4], Type: "json", Buffer_Size: 5 << 10, isclose: false, msgchan: make(chan any), rep: make(chan response), req: Request{Type: "json", Path: params[4]}}
	} else {
		return nil
	}
}
func (s *BsClient) Register(con net.Conn) {
	s.Conn = con
}

func (s *BsClient) IsClose() bool {
	return s.isclose
}

func (s *BsClient) NeedWaitReturn() bool {
	return true
}

func (s *BsClient) GetBack() error {
	fmt.Println("get backdo") //debug
	buffer := make([]byte, s.Buffer_Size)
	fmt.Print("read buffer\t")
	lang, err := s.Conn.Read(buffer)
	fmt.Println("read buffer success")
	var rep response
	if err == nil {
		err = json.Unmarshal(buffer[:lang], &rep)
	}
	fmt.Print("send to channel\t")
	s.rep <- rep
	fmt.Print("getback end")
	return err
}
func BackEnd(cl *BsClient) {
	for {
		time.Sleep(time.Duration(5) * time.Second)
		if time.Since(cl.lastconn) > time.Duration(Max_Cut_Time)*time.Second {
			//close connection
			cl.isclose = true
			break
		}
	}
}

func Start(cl *BsClient, errchan chan<- error) {
	var address string
	if cl.Socket_Type == "unix" {
		address = cl.Address
	} else {
		address = cl.Address + ":" + strconv.Itoa(cl.Port)
	}
	go BackEnd(cl)
	go kits.Dial(cl.Socket_Type, address, cl, cl.msgchan, errchan)
}
func Stop(cl *BsClient) {
	cl.isclose = true
}
func (s *BsClient) Do(v any, rep any) error {
	content, err := json.Marshal(v)
	if err == nil {
		s.req.Content = string(content)
	}
	s.msgchan <- s.req
	trep := <-s.rep
	err = json.Unmarshal([]byte(trep.Content), rep)
	s.lastconn = time.Now()
	return err
}
