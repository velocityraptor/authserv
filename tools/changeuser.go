package main

import (
	"net"
	"bufio"
	"fmt"
	"os"
	"flag"
)

func main() {
	cn, err := net.Dial("tcp","auth.velocityraptor.net:12345")
	dblock()
	defer dbunlock()
	if err != nil {
		fmt.Println(err)
		return
	}
	r := bufio.NewReader(cn)
	
