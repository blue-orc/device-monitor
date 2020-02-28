package pythonJobRunner

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

type TrainingStatus struct {
	Step              string
	TrainingScript    string
	BatchSize         string
	NumberOfFiles     string
	NumberOfWorkers   string
	CurrentFileIndex  string
	CurrentImageIndex string
	ImagesPerFile     string
	Loss              string
	Status            string
	CurrentEpoch      string
	Epochs            string
	Layers            string
	Depth             string
	LearningRate      string
}

var Status TrainingStatus

func GetStatusJSON() ([]byte, error) {
	sBytes, err := json.Marshal(Status)
	if err != nil {
		return sBytes, err
	}
	return sBytes, nil
}

func Run(distributed bool, distNodes int, distGpus int) {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Println("main: %s", err.Error())
	}
	if Status.Status == "Running" {
		return
	}

	Status.Status = "Running"
	//cmd := exec.Command("python3", "/home/ubuntu/go/src/gpu-demonstration-api/python-job-runner/scripts/pytorch-training-gpu.py")
	var cmd *exec.Cmd

	if distributed {
		cmd = exec.Command("python3", wd+"/distributed-gpu.py", "-n", fmt.Sprintf("%d", distNodes), "-g", fmt.Sprintf("%d", distGpus), "-nr", fmt.Sprintf("%d", 0), "--epochs", fmt.Sprintf("%", 10))
	} else {
		cmd = exec.Command("python3", wd+"/cifar-gpu.py")
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println("Run Python Script Error: " + err.Error())
		Status.Status = "Finished"
		return
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		fmt.Println("Run Python Script Error: " + err.Error())
		Status.Status = "Finished"
		return
	}

	err = cmd.Start()
	if err != nil {
		fmt.Println("Run Python Script Error: " + err.Error())
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
		res := strings.Split(txt, ":")
		if res[0] == "Step" {
			Status.Step = res[1]
		} else if res[0] == "TrainingScript" {
			Status.TrainingScript = res[1]
		} else if res[0] == "BatchSize" {
			Status.BatchSize = res[1]
		} else if res[0] == "NumberOfFiles" {
			Status.NumberOfFiles = res[1]
		} else if res[0] == "CurrentFileIndex" {
			Status.CurrentFileIndex = res[1]
		} else if res[0] == "CurrentImageIndex" {
			Status.CurrentImageIndex = res[1]
		} else if res[0] == "ImagesPerFile" {
			Status.ImagesPerFile = res[1]
		} else if res[0] == "Loss" {
			Status.Loss = res[1]
		} else if res[0] == "Status" {
			Status.Status = res[1]
		} else if res[0] == "CurrentEpoch" {
			Status.CurrentEpoch = res[1]
		} else if res[0] == "Epochs" {
			Status.Epochs = res[1]
		} else if res[0] == "Layers" {
			Status.Layers = res[1]
		} else if res[0] == "Depth" {
			Status.Depth = res[1]
		} else if res[0] == "LearningRate" {
			Status.LearningRate = res[1]
		} else if res[0] == "NumberOfWorkers" {
			Status.NumberOfWorkers = res[1]
		} else {
			fmt.Println(scanner.Text())
		}
	}
}
