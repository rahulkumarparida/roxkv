package metrics

import (
	"os/user"
	"runtime"

	"github.com/rahulkumarparida/roxkv/internal/utils"
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