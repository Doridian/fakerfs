package dev

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"syscall"
)

type FixedHandler struct {
	Value       []byte
	AllowReads  bool
	AllowWrites bool
}

func (h *FixedHandler) LoadConfig(cfg map[string]interface{}) error {
	var err error

	rawValue := cfg["value"].(string)
	valueType, ok := cfg["value_type"].(string)
	if !ok || valueType == "" {
		valueType = "string"
	}

	switch valueType {
	case "string":
		h.Value = []byte(rawValue)
	case "base64":
		h.Value, err = base64.StdEncoding.DecodeString(rawValue)
	case "hex":
		h.Value, err = hex.DecodeString(rawValue)
	default:
		return fmt.Errorf("unknown value_type: %s", valueType)
	}

	if err != nil {
		return err
	}

	bVal, ok := cfg["allow_reads"]
	if ok {
		h.AllowReads = bVal.(bool)
	} else {
		h.AllowReads = true
	}
	bVal, ok = cfg["allow_writes"]
	if ok {
		h.AllowWrites = bVal.(bool)
	} else {
		h.AllowWrites = false
	}
	return nil
}

func (h *FixedHandler) GetData() ([]byte, syscall.Errno) {
	if !h.AllowReads {
		return nil, syscall.EINVAL
	}
	return h.Value, 0
}

func (h *FixedHandler) SetData(data []byte) syscall.Errno {
	if !h.AllowWrites {
		return syscall.EINVAL
	}
	return 0
}
