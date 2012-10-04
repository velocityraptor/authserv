package main

import (
	"net"
	"fmt"
	"bufio"
	"os"
	"strings"
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

func getuser(user string, cn net.Conn) {
	dblock()
	defer dbunlock()
	user,err := os.Open(user, O_RDWR)
	if err != nil {
		fmt.Fprintf(cn,"Error: No user\n\r\n")
		return
	}
	r := bufio.NewReader(user)
	str,err := r.ReadString('\n')
	str = strings.Trim(str," \n\r\t")
	fmt.Fprintf(cn,"pass:%s\n\r\n",str)
	return
}

func verify(user string, challenge string, cn net.Conn) {
	dblock()
	defer dbunlock()
	user, err := os.Open(user, O_RDWR)
	if err != nil {
		fmt.Fprintf(cn,"Error: No user\n\r\n")
		return
	}
	r := bufio.NewReader(user)
	str,err := r.ReadString('\n')
	str = strings.Trim(str," \n\r\t")
	if challenge == str {
		fmt.Fprintf(cn,"OK\n\r\n")
	} else {
		fmt.Fprintf(cn,"Fail\n\r\n")
	}
	return
}

func dblock() {
	
