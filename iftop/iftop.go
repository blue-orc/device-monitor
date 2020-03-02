package iftop

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"strings"
)

type IftopStatus struct {
	Status      string
	RecieveRate string
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
		res := strings.Index(txt, "Total send rate:")
		if res == -1 {
			continue
		}
		sub1 := txt[res:len(txt)]
		res = strings.Index(sub1, "b")
		sub2 := txt[res:len(sub1)]
		res = strings.Index(sub2, "b")
		sub3 := txt[0:res]
		fmt.Println(sub3)
	}
}
