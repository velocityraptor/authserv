package main

import (
	"net"
	"fmt"
	"bufio"
	"os"
)

func main() {
	fmt.Println("Velocityraptor Authserver")
	fmt.Println("Revision 0")
	ln, err := net.Listen("tcp", ":12345")
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	for {
		cn, err := ln.Accept()
		if err != nil {
			continue
		}
		go handler(cn)
	}
}

func handler(cn net.Conn) {
	cnr := bufio.NewReader(cn)
	for {
		str, err := cnr.ReadString('\n')
		str = strings.Trim(str, " \n\r\t")
		s := strings.Split(str, " ")
		switch s[0] {
		case "getuser":
			getuser(s[1], cn)
		case "verify":
			verify(s[1], s[2], cn)
		case "netkey":
			netkey_verify(s[1], s[2], cn)
		default:
			fmt.Fprintf(cn, "Invalid\n\r\n")
			cn.Close()
			return
		}
	}
}
