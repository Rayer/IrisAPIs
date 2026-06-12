package main

import (
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/load"
	"github.com/shirou/gopsutil/v4/mem"
)

type SystemInfo struct {
	CPU    SystemCPU    `json:"cpu"`
	Memory SystemMemory `json:"memory"`
	Disk   SystemDisk   `json:"disk"`
	Load   SystemLoad   `json:"load"`
}

type SystemCPU struct {
	UsagePercent float64 `json:"usage_percent"`
	Cores        int     `json:"cores"`
}

type SystemMemory struct {
	TotalMB      int     `json:"total_mb"`
	UsedMB       int     `json:"used_mb"`
	AvailableMB  int     `json:"available_mb"`
	UsagePercent float64 `json:"usage_percent"`
}

type SystemDisk struct {
	TotalGB      float64 `json:"total_gb"`
	UsedGB       float64 `json:"used_gb"`
	FreeGB       float64 `json:"free_gb"`
	UsagePercent float64 `json:"usage_percent"`
}

type SystemLoad struct {
	Load1  float64 `json:"load_1m"`
	Load5  float64 `json:"load_5m"`
	Load15 float64 `json:"load_15m"`
}

func roundTo1(v float64) float64 {
	return float64(int(v*10+0.5)) / 10
}

func roundTo2(v float64) float64 {
	return float64(int(v*100+0.5)) / 100
}

func GetSystemInfo() SystemInfo {
	info := SystemInfo{}

	// CPU: percent over 1 second
	if percent, err := cpu.Percent(0, false); err == nil && len(percent) > 0 {
		info.CPU.UsagePercent = roundTo1(percent[0])
	}
	if cores, err := cpu.Counts(true); err == nil {
		info.CPU.Cores = cores
	}

	// Memory
	if v, err := mem.VirtualMemory(); err == nil {
		info.Memory.TotalMB = int(v.Total / 1024 / 1024)
		info.Memory.UsedMB = int(v.Used / 1024 / 1024)
		info.Memory.AvailableMB = int(v.Available / 1024 / 1024)
		info.Memory.UsagePercent = roundTo1(v.UsedPercent)
	}

	// Disk
	if v, err := disk.Usage("/"); err == nil {
		const GB = 1024 * 1024 * 1024
		info.Disk.TotalGB = roundTo2(float64(v.Total) / GB)
		info.Disk.UsedGB = roundTo2(float64(v.Used) / GB)
		info.Disk.FreeGB = roundTo2(float64(v.Free) / GB)
		info.Disk.UsagePercent = roundTo1(v.UsedPercent)
	}

	// Load average
	if v, err := load.Avg(); err == nil {
		info.Load.Load1 = roundTo2(v.Load1)
		info.Load.Load5 = roundTo2(v.Load5)
		info.Load.Load15 = roundTo2(v.Load15)
	}

	return info
}
