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



func GetComputerUsage() utils.MonitorComputeStat{
    // runtime package

	//user information
	u,err := user.Current()
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

	virtualmem,err := mem.VirtualMemory()
	totalram := virtualmem.Total
	freeram := virtualmem.Free
	usedpercentage := virtualmem.UsedPercent

	totalgoroutines := runtime.NumGoroutine()

	computeusage := utils.MonitorComputeStat{
		Username: username,
		Os: os,
		Architecture: architechture,
		Cpus: cpus,
		TotalRam: totalram,
		FreeRam: freeram,
		UsedRamPercent: usedpercentage,
		TotalGoRoutines: totalgoroutines,
	}
	return computeusage

}


func GetCPUUsage() string{
	mu := sync.Mutex{}
	mu.Lock()
	percentage , err:= cpu.Percent(0,false)
	mu.Unlock()
	 if utils.HandleError("Error while fetching user data", err) {
		return ""
	 }

	return strconv.FormatFloat(percentage[0],'f', -1, 641)
}

type RAM struct{
	TotalRam string
	FreeRam string
	UsedPercentge string 
}

func GetRAMUsage() RAM{
	var m runtime.MemStats

	runtime.ReadMemStats(&m)

	virtualmem,err := mem.VirtualMemory()
	if utils.HandleError("Error while fetching user data", err) {
		return RAM{}
	}

	totalram := strconv.FormatUint(virtualmem.Total,10)
	freeram := strconv.FormatUint(virtualmem.Free,10)
	usedpercentage := strconv.FormatFloat(virtualmem.UsedPercent,'f', -1, 641)

	
	return RAM{
		TotalRam: totalram,
		FreeRam: freeram,
		UsedPercentge: usedpercentage,
	}


}

type DISK struct{
	Total string
	Free string
	Avaliable string
	Err error 
}

func GetDiskUsage(path string) DISK{
	var stat syscall.Statfs_t
	err := syscall.Statfs(path, &stat)
	if err != nil {
		return DISK{"0", "0", "0", err}
	}

	const GB = 1024 * 1024 * 1024
	// Calculate sizes in Gigabytes
	total := strconv.FormatUint((stat.Blocks * uint64(stat.Bsize))/ GB,10)
	free := strconv.FormatUint((stat.Bfree * uint64(stat.Bsize))/GB,10)
	available := strconv.FormatUint((stat.Bavail * uint64(stat.Bsize))/GB,10) 

	return DISK{total, free, available, nil}
}

type RuntimeStats struct {
	GoVersion        string
	OS               string
	Arch             string
	CPUs             string
	Goroutines       string
	AllocatedMemMB   string
	TotalAllocatedMB string
	SystemMemMB      string
	HeapAllocMB      string
	GCCycles         string
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
		AllocatedMemMB:   strconv.FormatUint(ms.Alloc/MB, 10),
		TotalAllocatedMB: strconv.FormatUint(ms.TotalAlloc/MB, 10),
		SystemMemMB:      strconv.FormatUint(ms.Sys/MB, 10),
		HeapAllocMB:      strconv.FormatUint(ms.HeapAlloc/MB, 10),
		GCCycles:         strconv.FormatUint(uint64(ms.NumGC), 10),
	}
}
