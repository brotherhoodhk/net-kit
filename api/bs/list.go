package bs

import (
	"fmt"
	"sync"
)

// register user session to global list
var ConList sync.Map

// global register
var Register func(*BSServer) any = defalut_register

func (s *BSServer) register() { //use client ip as special identity
	ConList.LoadOrStore(Register(s), s)
	// debugline
	ConList.Range(func(key, value any) bool {
		conid, ok := key.(string)
		if ok {
			fmt.Print(conid + "\t")
		}
		fmt.Println()
		return ok
	})
}
func (s *BSServer) deregister() {
	ConList.LoadAndDelete(Register(s))
}

// default register
func defalut_register(ser *BSServer) any {
	return ser.ClientIP()
}
