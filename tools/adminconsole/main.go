package main

import (
	"net"
	"bufio"
	"fmt"
	"os"
	"crypto/md5"
	"strings"
)

func main() {
	ln, err := net.Listen("tcp", ":1223")
	if err != nil {
		fmt.Println(err)
		return
	}
	for {
		cn, err := ln.Accept()
		if err != nil {
			continue
		}
		go adminhandler(cn)
	}
}

func adminhandler(cn net.Conn) {
	fmt.Fprintf(cn, "Authserver Administrative Console\n")
	fmt.Fprintf(cn, "Password: ")
	r := bufio.NewReader(cn)
	pass, err := r.ReadString('\n')
	pass = strings.Trim(pass, " \n\r\t")
	h := md5.New()
	fmt.Fprintf(h, "%s", pass)
	pass = fmt.Sprintf("%x", h.Sum(nil))
	as, err := net.Dial("tcp", "localhost:12345")
	if err != nil {
		fmt.Fprintf(cn, "Error connecting to authserver")
		return
	}
	fmt.Fprintf(as, "verify admin %s", pass)
	str, err := r.ReadString('\n')
	if str == "OK\n" {
		as.Close()
		admin_console(cn)
	} else {
		fmt.Fprintf(cn, "Error Authenticating\n")
		as.Close()
		cn.Close()
		return
	}
}

func admin_console(cn net.Conn) {
	dblocked := 0
	r := bufio.NewReader(cn)
	for {
		fmt.Fprintf(cn, "authserv> ")
		str, err := r.ReadString('\n')
		if err != nil {
			continue
		}
		str = strings.Trim(str, " \n\r\t")
		s := strings.Split(str, " ")
		length := len(s)
		fmt.Fprintf(cn, "len:%v", length)
		switch s[0] {
		case "lockdb":
			dblock()
			dblocked = 1
		case "unlockdb":
			dbunlock()
			dblocked = 0
		case "adduser":
			if dblocked == 1 {
				userf, err := os.Create(s[1])
				if err != nil {
					fmt.Fprintf(cn, "User Exists\n")
					continue
				}
				h := md5.New()
				fmt.Fprintf(h, "%s", s[2])
				fmt.Fprintf(userf, "%x\n", h.Sum(nil))
				fmt.Fprintf(cn, "User Created\n")
				userf.Close()
			} else {
				fmt.Fprintf(cn, "Database is not locked\n")
			}
		case "logout":
			cn.Close()
			return
		}
	}
	dbunlock()
}

func dblock() {
	lockfile, err := os.Create("lockfile")
	for err != nil {
		lockfile, err = os.Create("lockfile")
	}
	lockfile.Close()
}

func dbunlock() {
	os.Remove("lockfile")
}
