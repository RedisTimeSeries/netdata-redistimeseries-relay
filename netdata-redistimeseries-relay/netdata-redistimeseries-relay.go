package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math"

	//"github.com/timescale/tsbs_generate_queriess/cmd/tsbs_generate_queries/databases/redistimeseries"
	redistimeseries "github.com/RedisTimeSeries/redistimeseries-go"
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
	timestamp     float64     `json:"timestamp"`
	units         string      `json:"units"`
	value         float64     `json:"value"`
}

func server() {
	// listen on a port
	var client = redistimeseries.NewClient("localhost:6379", "nohelp", nil)

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
		go handleServerConnection(c, client)
	}
}

func handleServerConnection(c net.Conn, client *redistimeseries.Client) {
	defer c.Close()

	reader := bufio.NewReader(c)
	tp := textproto.NewReader(reader)
	//var rcv map[string]interface{}
	var rcv map[string]interface{}

	defer c.Close()
	for {
		// read one line (ended with \n or \r\n)
		line, _ := tp.ReadLineBytes()
		tp.ReadLineBytes()
		json.Unmarshal(line, &rcv)
		prefix := rcv["prefix"]
		chart_type := rcv["chart_type"]
		chart_family := rcv["chart_family"]
		chart_name := rcv["chart_name"]
		hostname := rcv["hostname"]
		name := rcv["name"]
		id := rcv["id"]
		timestamp := int64(math.Round(rcv["timestamp"].(float64)))
		value := rcv["value"].(float64)
		//fmt.Println(reflect.TypeOf(timestamp))
		keyName := fmt.Sprintf("%s:%s:%s:%s:%s:%s:%s", prefix, chart_type, chart_family, chart_name, id, name, hostname)
		//fmt.Println(keyName, timestamp, value, rcv)
		client.Add(keyName, timestamp, value)
	}

}

func main() {
	go server()
	//let the server goroutine run forever
	var input string
	fmt.Scanln(&input)
}
