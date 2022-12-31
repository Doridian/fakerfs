package sysfs

import (
	"fmt"
	"strings"
	"syscall"
)

type ChoiceHandler struct {
	Choices  []string
	Selected string
}

func (h *ChoiceHandler) LoadConfig(cfg map[string]interface{}) error {
	choicesRaw := cfg["choices"].([]interface{})
	h.Choices = make([]string, 0, len(choicesRaw))
	for _, choice := range choicesRaw {
		h.Choices = append(h.Choices, choice.(string))
	}

	sel, ok := cfg["selected"]
	if ok {
		h.Selected = sel.(string)
	} else {
		h.Selected = h.Choices[0]
	}

	return nil
}

func (h *ChoiceHandler) GetData() ([]byte, syscall.Errno) {
	res := []string{}
	for _, choice := range h.Choices {
		if choice == h.Selected {
			choice = fmt.Sprintf("[%s]", choice)
		}
		res = append(res, choice)
	}
	return []byte(strings.Join(res, " ") + "\n"), 0
}

func (h *ChoiceHandler) SetData(data []byte) syscall.Errno {
	str := strings.Trim(string(data), " \t\r\n")

	for _, choice := range h.Choices {
		if choice == str {
			h.Selected = choice
			return 0
		}
	}

	return syscall.EINVAL
}
