package metrics

import (
	"fmt"
	"os"
	"os/user"
	"runtime"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/rahulkumarparida/roxkv/internal/utils"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/shirou/gopsutil/v4/net"
)

func GetComputerUsage() utils.MonitorComputeStat {
	// runtime package

	//user information
	u, err := user.Current()
	if utils.HandleError("Error while fetching user data", err) {
		return utils.MonitorComputeStat{}
	}
	username := u.Username

	// os informations
	os := runtime.GOOS
	architechture := runtime.GOARCH
	cpus := runtime.NumCPU()

	//Memory informations
	var m runtime.MemStats

	runtime.ReadMemStats(&m)

	virtualmem, err := mem.VirtualMemory()
	totalram := virtualmem.Total
	freeram := virtualmem.Free
	usedpercentage := virtualmem.UsedPercent

	totalgoroutines := runtime.NumGoroutine()

	computeusage := utils.MonitorComputeStat{
		Username:        username,
		Os:              os,
		Architecture:    architechture,
		Cpus:            cpus,
		TotalRam:        totalram,
		FreeRam:         freeram,
		UsedRamPercent:  usedpercentage,
		TotalGoRoutines: totalgoroutines,
	}
	return computeusage

}

func GetCPUUsage() string {
	mu := sync.Mutex{}
	mu.Lock()
	percentage, err := cpu.Percent(0, false)
	mu.Unlock()
	if utils.HandleError("Error while fetching user data", err) {
		return ""
	}

	return strconv.FormatFloat(percentage[0], 'f', -1, 64)
}

type RAM struct {
	TotalRam      uint64  `json:"totalRam"`
	FreeRam       uint64  `json:"freeRam"`
	UsedPercentge float64 `json:"usedPercentge"`
}

func GetRAMUsage() RAM {
	var m runtime.MemStats

	runtime.ReadMemStats(&m)

	virtualmem, err := mem.VirtualMemory()
	if utils.HandleError("Error while fetching user data", err) {
		return RAM{}
	}

	totalram := virtualmem.Total
	freeram := virtualmem.Free
	usedpercentage := virtualmem.UsedPercent

	return RAM{
		TotalRam:      totalram,
		FreeRam:       freeram,
		UsedPercentge: usedpercentage,
	}

}

type DISK struct {
	Total     uint64 `json:"total"`
	Free      uint64 `json:"free"`
	Avaliable uint64 `json:"avaliable"`
	Err       error  `json:"err"`
}

func GetDiskUsage(path string) DISK {
	var stat syscall.Statfs_t
	err := syscall.Statfs(path, &stat)
	if err != nil {
		return DISK{0, 0, 0, err}
	}

	const GB = 1024 * 1024 * 1024
	// Calculate sizes in Gigabytes
	total := (stat.Blocks * uint64(stat.Bsize)) / GB
	free := (stat.Bfree * uint64(stat.Bsize)) / GB
	available := (stat.Bavail * uint64(stat.Bsize)) / GB

	return DISK{total, free, available, nil}
}

type RuntimeStats struct {
	GoVersion        string `json:"goVersion"`
	OS               string `json:"os"`
	User 			 string `json:"user"`
	Arch             string `json:"arch"`
	CPUs             string `json:"cpus"`
	Goroutines       string `json:"goroutines"`
	AllocatedMemMB   uint64 `json:"allocatedMemMB"`
	TotalAllocatedMB uint64 `json:"totalAllocatedMB"`
	SystemMemMB      uint64 `json:"systemMemMB"`
	ConnectedUsers   int 	`json:"connectedclient"`
	ServerUptime     float64 `json:"serverUptime"`
}

func GetRuntimeStats() RuntimeStats {
	// 1. Gather general runtime and CPU metrics
	goVersion := runtime.Version()
	osName := runtime.GOOS
	arch := runtime.GOARCH
	cpus := strconv.Itoa(runtime.NumCPU())
	goroutines := strconv.Itoa(runtime.NumGoroutine())
	user , err := os.Hostname()
	totalClients := GetConnectedClients()
	Uptime := GetServerUptime()

	if err != nil {
		user = "N/A"
	}


	// 2. Gather memory statistics
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)

	// Conversion constant for Megabytes
	const MB = 1024 * 1024

	// 3. Populate and return the struct with string values
	return RuntimeStats{
		GoVersion:        goVersion,
		OS:               osName,
		User:			  user,		
		Arch:             arch,
		CPUs:             cpus,
		Goroutines:       goroutines,
		AllocatedMemMB:   ms.Alloc / MB,
		TotalAllocatedMB: ms.TotalAlloc / MB,
		SystemMemMB:      ms.Sys / MB,
		ConnectedUsers: len(totalClients),
		ServerUptime:     Uptime,
	}
}


type NetworkStat struct{
	DownloadSpeed  float64 `json:"downloadspeed"` 
	UploadSpeed	float64	`json:"uploadspeed"`
	Name string
}

func NetworkStatistics(netStats chan<- NetworkStat) {
	prevStats , err := net.IOCounters(true) 

	if err != nil {
		fmt.Printf("Error getting initial stats: %v\n", err)
		return 
	}
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for range ticker.C{
		currentStats ,err := net.IOCounters(true)
		if err != nil {
			fmt.Printf("Error getting stats: %v\n", err)
			continue
		}

		for i := range currentStats {
			prev := getInterfaceStats(prevStats, currentStats[i].Name)
			current := currentStats[i]

			// Calculate differences and convert bytes to bits (* 8)
			diffRx := current.BytesRecv - prev.BytesRecv
			diffTx := current.BytesSent - prev.BytesSent

			rxSpeedKbps := float64(diffRx*8) / 1024 
			txSpeedKbps := float64(diffTx*8) / 1024

			netStats <- NetworkStat{
			DownloadSpeed : rxSpeedKbps,
			UploadSpeed : txSpeedKbps,
			Name : current.Name,
			}
		
			fmt.Printf("[%s] Down: %.2f Kbps | Up: %.2f Kbps\n", current.Name, rxSpeedKbps, txSpeedKbps)
		}

		prevStats = currentStats

	}

}

func getInterfaceStats(stats []net.IOCountersStat, name string) net.IOCountersStat {
	for _, stat := range stats {
		if stat.Name == name {
			return stat
		}
	}
	return net.IOCountersStat{}
}