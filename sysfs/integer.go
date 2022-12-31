package sysfs

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"syscall"
)

type IntegerHandler struct {
	Min     int
	Max     int
	Current int
}

func (h *IntegerHandler) LoadConfig(cfg map[string]interface{}) error {
	iVal, ok := cfg["min"]
	if ok {
		h.Min = iVal.(int)
	} else {
		h.Min = math.MinInt
	}
	iVal, ok = cfg["max"]
	if ok {
		h.Max = iVal.(int)
	} else {
		h.Max = math.MaxInt
	}
	iVal, ok = cfg["current"]
	if ok {
		h.Current = iVal.(int)
	} else {
		h.Current = 0
	}
	return nil
}

func (h *IntegerHandler) GetData() ([]byte, syscall.Errno) {
	return []byte(fmt.Sprintf("%d\n", h.Current)), 0
}

func (h *IntegerHandler) SetData(data []byte) syscall.Errno {
	str := strings.Trim(string(data), " \t\r\n")

	num, err := strconv.Atoi(str)
	if err != nil || num < h.Min || num > h.Max {
		return syscall.EINVAL
	}

	h.Current = num
	return 0
}
