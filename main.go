package main

import (
	"bytes"
	"device-monitor/gopsutil"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	var ip string
	flag.StringVar(&ip, "ip", "unset", "String representing ip address for api")

	var port string
	flag.StringVar(&port, "port", "unset", "String representing port for api")
	flag.Parse()

	if ip == "" || ip == "unset" {
		log.Fatal("ip flag is not set")
	}

	if port == "" || port == "unset" {
		log.Fatal("port flag is not set")
	}

	go send(ip + ":" + port)
	fmt.Println("Device monitor started")
	select {} // block forever
}

func send(dest string) {
	dJSON, err := gopsutil.GetDiskIOJSON()
	dReader := bytes.NewReader(dJSON)
	if err != nil {
		fmt.Println("send: " + err.Error())
		return
	}
	_, err = http.Post("http://"+dest+"/disk", "application/json", dReader)
	if err != nil {
		fmt.Println("send: " + err.Error())
	}
	nJSON, err := gopsutil.GetNetIOJSON()
	nReader := bytes.NewReader(nJSON)
	_, err = http.Post("http://"+dest+"/net", "application/json", nReader)
	if err != nil {
		fmt.Println("send: " + err.Error())
		return
	}

	time.Sleep(time.Second)
	send(dest)
}

func serverIP() string {
	ip, ok := os.LookupEnv("IP_ADDRESS")
	if !ok {
		panic("IP_ADDRESS missing from .env file")
	}
	return ip
}

func port() string {
	port, ok := os.LookupEnv("PORT")
	if !ok {
		panic("PORT missing from .env file")
	}
	return port
}
