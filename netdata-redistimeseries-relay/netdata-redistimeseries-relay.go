package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"net/textproto"

)

type chart_datapoint struct {
	chart_context string      `json:"chart_context"`
	chart_family  string      `json:"chart_family"`
	chart_id      string      `json:"chart_id"`
	chart_name    string      `json:"chart_name"`
	chart_type    string      `json:"chart_type"`
	hostname      string      `json:"hostname"`
	id            string      `json:"id"`
	labels        interface{} `json:"labels"`
	name          string      `json:"name"`
	prefix        string      `json:"prefix"`
	timestamp     int64       `json:"timestamp"`
	units         string      `json:"units"`
	value         float64     `json:"value"`
}

func server() {
	// listen on a port
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println(err)
		return
	}
	for {
		// accept a connection
		c, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		// handle the connection
		go handleServerConnection(c)
	}
}

func handleServerConnection(c net.Conn) {
	reader := bufio.NewReader(c)
	tp := textproto.NewReader(reader)
	var rcv map[string]interface{}

	defer c.Close()
	metrics := 0
	for {
		// read one line (ended with \n or \r\n)
		line, _ := tp.ReadLineBytes()
		tp.ReadLineBytes()
		json.Unmarshal(line, &rcv)

		//tp.ReadLine()
		metrics++
		fmt.Println(rcv)
		// do something with data here, concat, handle and etc...
	}
	fmt.Println(metrics)
	//fmt.Println(metrics)

	// we create a decoder that reads directly from the socket
	//d := json.NewDecoder(c)
	//var rcv interface{}
	//
	//err := d.Decode(&rcv)
	//if err == nil {
	//	fmt.Println(rcv)
	//
	//	//file, _ := json.MarshalIndent(rcv, "", " ")
	//	//_ = ioutil.WriteFile("test.json", file, 0644)
	//}
	//c.Close()

}

func main() {
	go server()
	//let the server goroutine run forever
	var input string
	fmt.Scanln(&input)
}
