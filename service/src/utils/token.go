package utils

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	RCode = "BcEfghJkLMPquvxYzisntapASwordpWQ"
)

func Gen8RCode() string {
	var r = ""
	t := time.Now().UnixMicro()
	s := fmt.Sprintf("%v", t)
	v := 0
	for i := 16; i > 0; i = i - 2 {
		k1, _ := strconv.Atoi(s[i-1 : i])
		k2, _ := strconv.Atoi(s[i-2 : i-1])
		k := k1 ^ k2
		if v == 0 {
			v = k
		}
		if i <= 8 {
			k = k ^ v
		}
		r = r + RCode[k:k+1]
	}
	return r
}

func GetMdTitle(data string) string {
	ls := strings.Split(data, "\n")
	for _, item := range ls {
		if strings.HasPrefix(item, "#  ") {
			return item[3:]
		}
	}
	return "undefined"
}

func GetMdTags(data, theme string) string {
	return ""
}