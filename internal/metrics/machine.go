package metrics

import (
	"os/user"
	"runtime"
	"strconv"
	"sync"
	"syscall"

	"github.com/rahulkumarparida/roxkv/internal/utils"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
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
	TotalRam      uint64
	FreeRam       uint64
	UsedPercentge float64
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
	Total     uint64
	Free      uint64
	Avaliable uint64
	Err       error
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
	GoVersion        string
	OS               string
	Arch             string
	CPUs             string
	Goroutines       string
	AllocatedMemMB   uint64
	TotalAllocatedMB uint64
	SystemMemMB      uint64
	HeapAllocMB      uint64
	GCCycles         uint32
}

func GetRuntimeStats() RuntimeStats {
	// 1. Gather general runtime and CPU metrics
	goVersion := runtime.Version()
	osName := runtime.GOOS
	arch := runtime.GOARCH
	cpus := strconv.Itoa(runtime.NumCPU())
	goroutines := strconv.Itoa(runtime.NumGoroutine())

	// 2. Gather memory statistics
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)

	// Conversion constant for Megabytes
	const MB = 1024 * 1024

	// 3. Populate and return the struct with string values
	return RuntimeStats{
		GoVersion:        goVersion,
		OS:               osName,
		Arch:             arch,
		CPUs:             cpus,
		Goroutines:       goroutines,
		AllocatedMemMB:   ms.Alloc / MB,
		TotalAllocatedMB: ms.TotalAlloc / MB,
		SystemMemMB:      ms.Sys / MB,
		HeapAllocMB:      ms.HeapAlloc / MB,
		GCCycles:         ms.NumGC,
	}
}
