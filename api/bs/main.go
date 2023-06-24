package bs

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"strings"
)

// signal
const (
	Deregister_Sig = 2
	max_cut_time   = 10
)

type Handler func(*Request, *Response) error

var globalfuncmap = make(map[string]Handler)
var NeedRegister = false

type Request struct {
	Type    string `json:"type"`
	Path    string `json:"path"`
	Content string `json:"content"`
}
type Response struct {
	close *bool
	con   net.Conn
}
type BSServer struct {
	net.Conn
	req          Request
	rep          Response
	funcmap      map[string]Handler
	NeedRegister bool
	close        bool
}

func (s *BSServer) Decode(content []byte) error {
	// debugline
	fmt.Println("accpet from client:",string(content))
	err:=json.Unmarshal(content, &s.req)
	if err!=nil{
		fmt.Printf("[error] %s\n",err.Error())
	}
	return err
}

func (s *BSServer) NeedSave() bool {
	return true
}

func (s *BSServer) Save(data []byte) {
	switch string(data) {
	case strconv.Itoa(Deregister_Sig):
		s.deregister()
		s.close = true
		s.Conn.Close()
	}
}

func (s *BSServer) Do() error {
	if s.NeedRegister {
		s.register()
		s.NeedRegister = false
	}
	if s.req.Path[0] != '/' {
		s.req.Path = "/" + s.req.Path
	}
	if fc, ok := s.funcmap[s.req.Path]; ok {
		fmt.Println("get " + s.req.Path) //debugline
		return fc(&s.req, &s.rep)
	} else {
		fmt.Println("get bad request " + s.req.Path) //debugline
		return s.Bad_Response()
	}
}

func (s *BSServer) IsClose() bool {
	return s.close
}

func (s *BSServer) Response() error {
	return nil
}

func (s *BSServer) Bad_Response() error {
	return s.rep.Write("json", "bad request")
}
func (s *BSServer) ClientIP() string {
	iparr := strings.Split(s.Conn.RemoteAddr().String(), ":")
	return strings.Join(iparr[:len(iparr)-1], ":")
}

func (s *Request) ShouldBind(v any) error {
	return json.Unmarshal([]byte(s.Content), v)
}

type response struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}

func (s *Response) Write(typename string, v any) error {
	content, err := json.MarshalIndent(v, "", "\t")
	if err == nil {
		var finalcontent []byte
		finalcontent, err = json.MarshalIndent(&response{Type: typename, Content: string(content)}, "", "\t")
		if err == nil {
			_, err = s.con.Write(finalcontent)
		}
	}
	return err
}
func (s *Response) Ping() bool {
	_, err := s.con.Write([]byte(""))
	return err == nil
}

// register the function for path
func (s *BSServer) RegisterFunc(path string, regfunc func(*Request, *Response) error) {
	if _, ok := s.funcmap[path]; !ok {
		s.funcmap[path] = regfunc
	}
}

// prepare zone
func NewSession(con net.Conn) *BSServer {
	scon := &BSServer{funcmap: globalfuncmap, Conn: con, close: false, rep: Response{con: con}, NeedRegister: NeedRegister}
	scon.rep.close = &scon.close
	return scon
}
func RegisterFunc(path string, handler Handler) {
	if path[0] != '/' {
		path = "/" + path
	}
	if _, ok := globalfuncmap[path]; !ok {
		globalfuncmap[path] = handler
	}
}
