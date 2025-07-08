// File: internal/bandwidth/math.go
// Purpose: Convert MB/GB per day/week/month into Tor's expected bytes/sec

package bandwidth

import (
	"errors"
	"fmt"
	"strings"
)

type Unit string
type Interval string

const (
	KB Unit = "KB"
	MB Unit = "MB"
	GB Unit = "GB"
	TB Unit = "TB"

	Daily   Interval = "daily"
	Weekly  Interval = "weekly"
	Monthly Interval = "monthly"
)

var unitMap = map[Unit]int64{
	KB: 1024,
	MB: 1024 * 1024,
	GB: 1024 * 1024 * 1024,
	TB: 1024 * 1024 * 1024 * 1024,
}

var intervalSeconds = map[Interval]int64{
	Daily:   86400,
	Weekly:  604800,
	Monthly: 2592000, // 30 days
}

// ToBytesPerSecond converts user input into Tor-compatible bytes/sec
func ToBytesPerSecond(amount int64, unit Unit, interval Interval) (int64, error) {
	mult, ok := unitMap[unit]
	if !ok {
		return 0, errors.New("invalid unit")
	}
	seconds, ok := intervalSeconds[interval]
	if !ok {
		return 0, errors.New("invalid interval")
	}
	totalBytes := amount * mult
	return totalBytes / seconds, nil
}

// PrettyPrintBytes formats a value in bytes/sec into a readable string
func PrettyPrintBytes(bytesPerSec int64) string {
	switch {
	case bytesPerSec >= unitMap[GB]:
		return fmt.Sprintf("%.2f GB/s", float64(bytesPerSec)/float64(unitMap[GB]))
	case bytesPerSec >= unitMap[MB]:
		return fmt.Sprintf("%.2f MB/s", float64(bytesPerSec)/float64(unitMap[MB]))
	case bytesPerSec >= unitMap[KB]:
		return fmt.Sprintf("%.2f KB/s", float64(bytesPerSec)/float64(unitMap[KB]))
	default:
		return fmt.Sprintf("%d B/s", bytesPerSec)
	}
}

// ParseUnit converts user input into a valid Unit
func ParseUnit(input string) (Unit, error) {
	input = strings.ToUpper(strings.TrimSpace(input))
	switch input {
	case "KB":
		return KB, nil
	case "MB":
		return MB, nil
	case "GB":
		return GB, nil
	case "TB":
		return TB, nil
	default:
		return "", errors.New("invalid unit")
	}
}

// ParseInterval converts user input into a valid Interval
func ParseInterval(input string) (Interval, error) {
	input = strings.ToLower(strings.TrimSpace(input))
	switch input {
	case "daily":
		return Daily, nil
	case "weekly":
		return Weekly, nil
	case "monthly":
		return Monthly, nil
	default:
		return "", errors.New("invalid interval")
	}
}
