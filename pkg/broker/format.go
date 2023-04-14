package broker

import (
	"fmt"
	"math"
)

func FormatFileSize(byteSize int) string {
	if float64(byteSize) >= math.Pow(10, 9) {
		return FormatDecimals(float64(byteSize)/math.Pow(10, 9), 0) + "GB"
	}
	if float64(byteSize) >= math.Pow(10, 6) {
		return FormatDecimals(float64(byteSize)/math.Pow(10, 6), 0) + "MB"
	}
	if float64(byteSize) >= math.Pow(10, 3) {
		return FormatDecimals(float64(byteSize)/math.Pow(10, 3), 0) + "KB"
	}
	return fmt.Sprintf("%vB", byteSize)
}

func FormatDecimals(duration float64, decimals int) string {
	return fmt.Sprintf("%.2f", duration)
}
