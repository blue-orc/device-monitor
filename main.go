package main

import (
	"bufio"
	"bytes"
	"device-monitor/gopsutil"
	"device-monitor/pythonJobRunner"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var cpu bool
var disk bool
var gpu bool
var net bool
var training bool
var distributed bool

var distNodes int
var distGpus int

func main() {
	var ip string
	flag.StringVar(&ip, "ip", "unset", "String representing ip address for api")

	var port string
	flag.StringVar(&port, "port", "unset", "String representing port for api")

	flag.BoolVar(&cpu, "cpu", false, "Enable CPU monitor")
	flag.BoolVar(&gpu, "gpu", false, "Enable GPU monitor")
	flag.BoolVar(&net, "net", false, "Enable net monitor")
	flag.BoolVar(&disk, "disk", false, "Enable disk monitor")
	flag.BoolVar(&training, "training", false, "Enable training monitor")
	flag.BoolVar(&distributed, "distributed", false, "Distributed monitor script")

	flag.IntVar(&distNodes, "distributed nodes", 1, "Number of distributed nodes")
	flag.IntVar(&distGpus, "distributed gpus", 1, "Number of distributed gpus per node")

	flag.Parse()

	if ip == "" || ip == "unset" {
		log.Fatal("ip flag is not set")
	}

	if port == "" || port == "unset" {
		log.Fatal("port flag is not set")
	}

	if gpu {
		go gopsutil.GPUMonitorInit()
	}

	go send(ip + ":" + port)
	fmt.Println("Device monitor started")
	if training {
		checkRunScript()
	}
	select {} // block forever
}

func checkRunScript() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("Run Training Script -> ")
		text, _ := reader.ReadString('\n')
		// convert CRLF to LF
		text = strings.Replace(text, "\n", "", -1)
		if distributed {
			pythonJobRunner.Run(true, distNodes, distGpus)
		} else {
			pythonJobRunner.Run(false, distNodes, distGpus)
		}
	}
}

func send(dest string) {
	if disk {
		dJSON, err := gopsutil.GetDiskIOJSON()
		if err != nil {
			fmt.Println("send disk: " + err.Error())
			return
		}
		dReader := bytes.NewReader(dJSON)
		_, err = http.Post("http://"+dest+"/disk", "application/json", dReader)
		if err != nil {
			fmt.Println("send disk: " + err.Error())
		}
	}

	if net {
		nJSON, err := gopsutil.GetNetIOJSON()
		if err != nil {
			fmt.Println("send net: " + err.Error())
			return
		}
		nReader := bytes.NewReader(nJSON)
		_, err = http.Post("http://"+dest+"/net", "application/json", nReader)
		if err != nil {
			fmt.Println("send net: " + err.Error())
			return
		}
	}

	if cpu {
		cJSON, err := gopsutil.GetCPUMemoryUtilizationJSON()
		if err != nil {
			fmt.Println("send cpu: " + err.Error())
			return
		}
		cReader := bytes.NewReader(cJSON)
		_, err = http.Post("http://"+dest+"/cpu", "application/json", cReader)
		if err != nil {
			fmt.Println("send cpu: " + err.Error())
			return
		}
	}

	if gpu {
		gJSON, err := gopsutil.GetGPUStatusJSON()
		if err != nil {
			fmt.Println("send gpu: " + err.Error())
			return
		}
		gReader := bytes.NewReader(gJSON)
		_, err = http.Post("http://"+dest+"/gpu", "application/json", gReader)
		if err != nil {
			fmt.Println("send gpu: " + err.Error())
			return
		}
	}

	if training {
		tJSON, err := pythonJobRunner.GetStatusJSON()
		if err != nil {
			fmt.Println("send training: " + err.Error())
			return
		}
		tReader := bytes.NewReader(tJSON)
		_, err = http.Post("http://"+dest+"/training", "application/json", tReader)
		if err != nil {
			fmt.Println("send training: " + err.Error())
			return
		}
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
