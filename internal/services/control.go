// File: internal/services/control.go
// Purpose: Control the Tor service via systemctl or launchctl

package services

import (
	"bytes"
	"errors"
	"os/exec"
	"runtime"
	"strings"
)

type Action string

const (
	Start   Action = "start"
	Stop    Action = "stop"
	Restart Action = "restart"
	Status  Action = "status"
)

func RunServiceAction(serviceName string, action Action) (string, error) {
	switch runtime.GOOS {
	case "linux":
		return runSystemctl(serviceName, action)
	case "darwin":
		return runLaunchctl(serviceName, action)
	default:
		return "", errors.New("unsupported OS")
	}
}

func runSystemctl(service string, action Action) (string, error) {
	cmd := exec.Command("systemctl", string(action), service)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func runLaunchctl(service string, action Action) (string, error) {
	var cmd *exec.Cmd
	switch action {
	case Start:
		cmd = exec.Command("launchctl", "start", service)
	case Stop:
		cmd = exec.Command("launchctl", "stop", service)
	case Restart:
		// Launchctl has no restart; emulate it
		stop := exec.Command("launchctl", "stop", service)
		start := exec.Command("launchctl", "start", service)
		var out1, out2 bytes.Buffer
		stop.Stdout, stop.Stderr = &out1, &out1
		start.Stdout, start.Stderr = &out2, &out2
		err1 := stop.Run()
		err2 := start.Run()
		return out1.String() + "\n" + out2.String(), combineErr(err1, err2)
	case Status:
		cmd = exec.Command("launchctl", "list")
	default:
		return "", errors.New("unsupported action")
	}
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()

	if action == Status {
		lines := strings.Split(out.String(), "\n")
		for _, line := range lines {
			if strings.Contains(line, service) {
				return line, nil
			}
		}
		return "Not found", nil
	}

	return out.String(), err
}

func combineErr(e1, e2 error) error {
	if e1 != nil && e2 != nil {
		return errors.New(e1.Error() + "; " + e2.Error())
	}
	if e1 != nil {
		return e1
	}
	return e2
}
