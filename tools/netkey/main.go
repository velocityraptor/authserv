package main

import (
	"net"
	"bufio"
	"fmt"
	"crypto/md5"
	"os"
	"strings"
)

func main() {
	stdin := bufio.NewReader(os.Stdin)
	fmt.Printf("Server: ")	
	str,err := stdin.ReadString('\n')
	str = strings.Trim(str," \n\r\t")
	cn, err := net.Dial("tcp",str)
	if err != nil {
		fmt.Println(err)
		return
	}
	r := bufio.NewReader(cn)
	fmt.Printf("User: ")
	str, err = stdin.ReadString('\n')
	str = strings.Trim(str," \n\t\r")
	user := str
	fmt.Printf("Password: ")
	pass, err := stdin.ReadString('\n')
	pass = strings.Trim(pass," \n\r\t")
	ph := md5.New()
	fmt.Fprintf(cn,"challenge %s\n",str)
	str,err = r.ReadString('\n')
	str = strings.Trim(str," \n\r\t")
	s := strings.Split(str,":")
	if s[0] == "chal" {
		fmt.Fprintf(ph,"%s",s[1])
		fmt.Fprintf(ph,"%s",pass)
		x := fmt.Sprintf("%x",ph.Sum(nil))
		fmt.Fprintf(cn,"response %s %s\n",user,x)
		str,err = r.ReadString('\n')
		fmt.Println("response: ",str)
		if str == "OK\n" {
			fmt.Println("Verified")
		} else {
			fmt.Println("Error: Invalid password")
		}
	} else {
		fmt.Println("Error: Protocol Botch")
	}
}

	
	
