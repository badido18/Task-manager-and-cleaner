package process

import (
	"syscall"
	"unsafe"
)

const (
	processQueryInfo = 1 << 10
	processVMRead    = 1 << 4
)

type ProcessMemoryCounters struct {
	cb                         uint32
	PageFaultCount             uint32
	PeakWorkingSetSize         uint64
	WorkingSetSize             uint64
	QuotaPeakPagedPoolUsage    uint64
	QuotaPagedPoolUsage        uint64
	QuotaPeakNonPagedPoolUsage uint64
	QuotaNonPagedPoolUsage     uint64
	PagefileUsage              uint64
	PeakPagefileUsage          uint64
}

// GetAllSnapshot kda mena melhih
func GetAllSnapshot() (pEntries []syscall.ProcessEntry32, r error) {

	handle, r := syscall.CreateToolhelp32Snapshot(syscall.TH32CS_SNAPPROCESS, 0)

	if r != nil {
		return
	}
	defer syscall.CloseHandle(handle)

	pe32 := syscall.ProcessEntry32{}

	cal := unsafe.Sizeof(pe32)
	pe32.Size = uint32(cal)

	r = syscall.Process32First(handle, &pe32)

	if r != nil {
		return
	}

	for {

		pEntries = append(pEntries, pe32)

		r = syscall.Process32Next(handle, &pe32)
		if r != nil {
			r = nil
			break
		}
	}

	return
}

func GetMemoryUsage(pid uint32) (p ProcessMemoryCounters, err error) {

	current, err := syscall.OpenProcess(processQueryInfo|processVMRead, false, 1964)

	if err != nil {
		return
	}
	defer syscall.CloseHandle(current)

	psapi := syscall.NewLazyDLL("psapi.dll")

	p.cb = uint32(unsafe.Sizeof(p))

	GetProcessMemoryInfo := psapi.NewProc("GetProcessMemoryInfo")
	_, _, err = GetProcessMemoryInfo.Call(uintptr(current), uintptr(unsafe.Pointer(&p)), uintptr(p.cb))

	return
}
