// File: internal/config/validate.go
// Purpose: Validate common torrc fields before saving

package config

import (
	"errors"
	"net"
	"strconv"
	"strings"
)

// ValidateBandwidth ensures the bandwidth string is numeric and has a valid suffix
func ValidateBandwidth(value string) error {
	value = strings.ToUpper(strings.TrimSpace(value))
	validSuffixes := []string{"KB", "MB", "GB", "TB"}
	for _, suffix := range validSuffixes {
		if strings.HasSuffix(value, suffix) {
			n := strings.TrimSuffix(value, suffix)
			n = strings.TrimSpace(n)
			if _, err := strconv.Atoi(n); err != nil {
				return errors.New("invalid numeric bandwidth value")
			}
			return nil
		}
	}
	return errors.New("invalid bandwidth format (e.g. 5MB, 100KB)")
}

// ValidatePortMapping checks that the value is of form: "80 127.0.0.1:8080"
func ValidatePortMapping(value string) error {
	parts := strings.Fields(value)
	if len(parts) != 2 {
		return errors.New("port mapping must be in format: <port> <host:port>")
	}
	if _, err := strconv.Atoi(parts[0]); err != nil {
		return errors.New("invalid virtual port")
	}
	addrParts := strings.Split(parts[1], ":")
	if len(addrParts) != 2 {
		return errors.New("invalid address format")
	}
	if net.ParseIP(addrParts[0]) == nil && addrParts[0] != "localhost" {
		return errors.New("invalid IP or hostname")
	}
	if _, err := strconv.Atoi(addrParts[1]); err != nil {
		return errors.New("invalid target port")
	}
	return nil
}

// ValidateOnionDir ensures the path looks reasonable and doesn't allow path tricks
func ValidateOnionDir(path string) error {
	if strings.Contains(path, "..") || strings.Contains(path, "~") {
		return errors.New("invalid HiddenServiceDir path")
	}
	if len(path) < 3 {
		return errors.New("path too short")
	}
	return nil
}

// ValidateIPOrLocalhost ensures the address is valid
func ValidateIPOrLocalhost(addr string) error {
	if addr == "localhost" {
		return nil
	}
	ip := net.ParseIP(addr)
	if ip == nil {
		return errors.New("invalid IP address")
	}
	return nil
}
