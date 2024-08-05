package utils

import "fmt"

func BytesSize2Str(n int64) string {
	var s string

	if n < 1024 {
		s = fmt.Sprintf("%dB", n)
	} else if n < 1024*1024 {
		s = fmt.Sprintf("%.2fKB", float64(n)/1024.0)
	} else if n < 1024*1024*1024 {
		s = fmt.Sprintf("%.2fMB", float64(n)/1024.0/1024.0)
	} else {
		s = fmt.Sprintf("%dB", n)
	}

	return s
}
