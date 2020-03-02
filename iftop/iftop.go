package iftop

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"strconv"
	"strings"
)

type IftopStatus struct {
	Status      string
	RecieveRate float64
}

var Status IftopStatus

func Run() {
	if Status.Status == "Running" {
		return
	}

	Status.Status = "Running"
	//cmd := exec.Command("python3", "/home/ubuntu/go/src/gpu-demonstration-api/python-job-runner/scripts/pytorch-training-gpu.py")
	var cmd *exec.Cmd

	cmd = exec.Command("iftop", "-t")

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println("Iftop Run: " + err.Error())
		Status.Status = "Finished"
		return
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		fmt.Println("Iftop Run: " + err.Error())
		Status.Status = "Finished"
		return
	}

	err = cmd.Start()
	if err != nil {
		fmt.Println("Iftop Run: " + err.Error())
		Status.Status = "Finished"
		return
	}

	go updateStatus(stdout)
	go updateStatus(stderr)
	cmd.Wait()
	Status.Status = "Finished"
}

func updateStatus(r io.Reader) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		txt := scanner.Text()
		res := strings.Index(txt, "Total receive rate:")
		if res == -1 {
			continue
		}
		sub1 := txt[res:len(txt)]
		res = strings.Index(sub1, "b")
		sub2 := sub1[res+1 : len(sub1)]
		res = strings.Index(sub2, "b")
		sub3 := sub2[0:res]
		sub3 = strings.TrimSpace(sub3)
		fmt.Println(sub3)
		fmt.Println(string(sub3[len(sub3)-1]))
		size := string(sub3[len(sub3)-1])
		if size == "K" {
			v, err := strconv.ParseFloat(sub3[0:len(sub3)-1], 64)
			if err != nil {
				fmt.Println("Update status: " + err.Error())
			}
			v = v * 1000
			Status.RecieveRate = v
		} else if size == "M" {
			v, err := strconv.ParseFloat(sub3[0:len(sub3)-1], 64)
			if err != nil {
				fmt.Println("Update status: " + err.Error())
			}
			v = v * 1000000
			Status.RecieveRate = v
		} else if size == "G" {
			v, err := strconv.ParseFloat(sub3[0:len(sub3)-1], 64)
			if err != nil {
				fmt.Println("Update status: " + err.Error())
			}
			v = v * 1000000000
			Status.RecieveRate = v
		}
		fmt.Println(Status.RecieveRate)
	}
}
