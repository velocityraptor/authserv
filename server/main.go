package main

import (
	"net"
	"fmt"
	"bufio"
	"os"
	"strings"
	"crypto/md5"
	"time"
)

func main() {
	fmt.Println("Velocityraptor Authserver")
	fmt.Println("Revision 1")
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
	var chall string
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
		case "challenge":
			chall = challenge(s[1],cn)
		case "response":
			response(s[1],s[2],chall,cn)
		default:
			fmt.Fprintf(cn, "Invalid\n")
			cn.Close()
			return
		}
	}
}

func challenge(username string, cn net.Conn) string {
	dblock()
	userf, err := os.Open(username)
	if err != nil {
		fmt.Fprintf(cn, "Error: No user\n")
		dbunlock()
		return ""
	}
	r := bufio.NewReader(userf)
	str,err := r.ReadString('\n')
	str = strings.Trim(str," \n\t\r")
	userf.Close()
	dbunlock()
	h := md5.New()
	t := time.Now()
	fmt.Fprintf(h,"%s\n",str)
	fmt.Fprintf(h,"%s\n",t.Local())
	chal := fmt.Sprintf("%x",h.Sum(nil))
	fmt.Fprintf(cn,"chal:%s\n",chal)
	return chal
}

func response(user string,response string,chal string,cn net.Conn) {
	fmt.Println("getting response")
	dblock()
	userf,err := os.Open(user)
	if err != nil {
		fmt.Fprintf(cn,"Error: No User\n")
		dbunlock()
		return
	}
	r := bufio.NewReader(userf)
	str,err := r.ReadString('\n')
	str = strings.Trim(str," \n\t\r")
	userf.Close()
	dbunlock()
	h := md5.New()
	fmt.Fprintf(h,"%s",chal)
	fmt.Fprintf(h,"%s",str)
	x := fmt.Sprintf("%x",h.Sum(nil))
	fmt.Println("User: ",str)
	fmt.Println("Challenge: ",chal)
	fmt.Println("Expected Response: ", x)
	fmt.Println("Actual Response: ", response)
	if x == response {
		fmt.Fprintf(cn,"OK\n")
	} else {
		fmt.Fprintf(cn,"Invalid\n")
	}
}

func getuser(user string, cn net.Conn) {
	dblock()
	defer dbunlock()
	userf, err := os.Open(user)
	if err != nil {
		fmt.Fprintf(cn, "Error: No user\n")
		return
	}
	r := bufio.NewReader(userf)
	str, err := r.ReadString('\n')
	str = strings.Trim(str, " \n\r\t")
	fmt.Fprintf(cn, "pass:%s\n", str)
	return
}

func verify(user string, challenge string, cn net.Conn) {
	dblock()
	defer dbunlock()
	userf, err := os.Open(user)
	if err != nil {
		fmt.Fprintf(cn, "Error: No user\n")
		return
	}
	r := bufio.NewReader(userf)
	str, err := r.ReadString('\n')
	str = strings.Trim(str, " \n\r\t")
	h := md5.New()
	fmt.Fprintf(h, "%s", challenge)
	challenge = fmt.Sprintf("%x", h.Sum(nil))
	if challenge == str {
		fmt.Fprintf(cn, "OK\n")
	} else {
		fmt.Fprintf(cn, "Fail\n")
	}
	return
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

func adduser(username string, password string, cn net.Conn) {
	dblock()
	defer dbunlock()
	userf, err := os.Create(username)
	if err != nil {
		fmt.Fprintf(cn, "User Exists\n")
		return
	}
	h := md5.New()
	fmt.Fprintf(h, "%s", password)
	fmt.Fprintf(userf, "%x\n", h.Sum(nil))
	fmt.Fprintf(cn, "User Created\n")
	userf.Close()
	return
}
