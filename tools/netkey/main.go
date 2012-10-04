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
	str = strings.Trim(" \n\r\t")
	cn, err := net.Dial("tcp",str)
	if err != nil {
		fmt.Println(err)
		return
	}
	r := bufio.NewReader(cn)
	fmt.Printf("User: ")
	str, err = stdin.ReadString('\n')
	str = strings.Trim(" \n\t\r")
	
