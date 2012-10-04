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
	defer cn.Close()
	for {
		str, err := cnr.ReadString('\n')
		if err != nil {
			continue
		}
		str = strings.Trim(str, " \n\r\t")
		s := strings.Split(str, " ")
		switch s[0] {
		case "getuser":
			getuser(s[1], cn)
		case "verify":
			verify(s[1], s[2], cn)
		case "adduser":
			adduser(s[1], s[2], cn)
		case "close":
			return
		default:
			fmt.Fprintf(cn, "Invalid\n")
			cn.Close()
			return
		}
	}
}

func getuser(user string, cn net.Conn) {
	dblock()
	defer dbunlock()
	userf,err := os.Open(user)
	if err != nil {
		fmt.Fprintf(cn,"Error: No user\n")
		return
	}
	r := bufio.NewReader(userf)
	str,err := r.ReadString('\n')
	str = strings.Trim(str," \n\r\t")
	fmt.Fprintf(cn,"pass:%s\n",str)
	return
}

func verify(user string, challenge string, cn net.Conn) {
	dblock()
	defer dbunlock()
	userf, err := os.Open(user)
	if err != nil {
		fmt.Fprintf(cn,"Error: No user\n")
		return
	}
	r := bufio.NewReader(userf)
	str,err := r.ReadString('\n')
	str = strings.Trim(str," \n\r\t")
	if challenge == str {
		fmt.Fprintf(cn,"OK\n")
	} else {
		fmt.Fprintf(cn,"Fail\n")
	}
	return
}

func dblock() {
	lockfile,err := os.Create("lockfile")
	for err != nil {
		lockfile,err = os.Create("lockfile")
	}
	lockfile.Close()
}

func dbunlock() {
	os.Remove("lockfile")
}

func adduser(username string, password string, cn net.Conn) {
	dblock()
	defer dbunlock()
	userf,err := os.Create(username)
	if err != nil {
		fmt.Fprintf(cn,"User Exists\n")
		return
	}
	fmt.Fprintf(userf,"%s\n",password)
	fmt.Fprintf(cn,"User Created\n")
	userf.Close()
	return
}

	
