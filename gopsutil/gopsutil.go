package gopsutil

import (
	"encoding/json"
	"fmt"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/net"
	"gpu-demonstration-api/models"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

type CPUStatus struct {
	Memory     *mem.VirtualMemoryStat
	CPUPercent []float64
}

func GetCPUMemoryUtilizationJSON() ([]byte, error) {
	var c models.CPUStatus
	pct, err := GetCPUPercent()
	if err != nil {
		return nil, err
	}
	mem, err := GetVirtualMemory()
	if err != nil {
		return nil, err
	}
	c.CPUPercent = pct
	c.Memory = mem

	cpuBytes, err := json.Marshal(c)
	if err != nil {
		fmt.Println("devicemonitor.GetCPUMemoryUtilizationJSON: " + err.Error())
		return cpuBytes, err
	}
	return cpuBytes, nil
}

func GetCPUInfo() ([]cpu.InfoStat, error) {
	cpu, err := cpu.Info()
	if err != nil {
		fmt.Println("devicemonitor.GetCPUInfoJSON: " + err.Error())
		return nil, err
	}
	return cpu, nil
}

func GetCPUPercent() ([]float64, error) {
	var dur time.Duration
	percent, err := cpu.Percent(dur, false)
	if err != nil {
		fmt.Println("devicemonitor.GetCPUPercent: " + err.Error())
		return nil, err
	}
	return percent, nil
}

func GetVirtualMemory() (*mem.VirtualMemoryStat, error) {
	mem, err := mem.VirtualMemory()
	if err != nil {
		fmt.Println("devicemonitor.GetVirtualMemory: " + err.Error())
		return mem, err
	}

	return mem, nil
}

func GetNetIO() (i []net.IOCountersStat, err error) {
	i, err = net.IOCounters(false)
	if err != nil {
		fmt.Println("gopsutil.GetNetIO: " + err.Error())
		return
	}
	return
}

func GetNetIOJSON() (ioBytes []byte, err error) {
	ioStat, err := GetNetIO()
	if err != nil {
		return
	}

	ioBytes, err = json.Marshal(ioStat)
	if err != nil {
		fmt.Println("gopsutil.GetNetIOJSON: " + err.Error())
		return
	}
	return
}

func GetDiskIO() (i map[string]disk.IOCountersStat, err error) {
	i, err = disk.IOCounters()
	if err != nil {
		fmt.Println("gopsutil.GetDiskIO: " + err.Error())
		return
	}
	return
}

func GetDiskIOJSON() (ioBytes []byte, err error) {
	ioStat, err := GetDiskIO()
	if err != nil {
		return
	}

	ioBytes, err = json.Marshal(ioStat)
	if err != nil {
		fmt.Println("gopsutil.GetDiskIOJSON: " + err.Error())
		return
	}
	return
}
