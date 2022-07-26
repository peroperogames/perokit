package hardware_percent

import "testing"

func TestGetCpuPercent(t *testing.T) {
	percent := GetCpuPercent()
	if percent == 0 {
		t.Fail()
	}
	t.Log("CpuPercent: ", percent)
}

func TestGetMemPercent(t *testing.T) {
	percent := GetMemPercent()
	if percent == 0 {
		t.Fail()
	}
	t.Log("MemPercent: ", percent)
}

func TestGetDiskPercent(t *testing.T) {
	percent := GetDiskPercent()
	if percent == 0 {
		t.Fail()
	}
	t.Log("DiskPercent: ", percent)
}
