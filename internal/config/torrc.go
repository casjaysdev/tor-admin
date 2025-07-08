// File: internal/config/torrc.go
// Purpose: Read, write, and modify torrc config safely with comment support

package config

import (
	"bufio"
	"errors"
	"os"
	"strings"
)

type TorConfigEntry struct {
	RawLine   string
	IsComment bool
	Key       string
	Value     string
}

type TorConfig struct {
	Entries []TorConfigEntry
}

// LoadTorrc loads the torrc file and parses it into structured entries
func LoadTorrc(path string) (*TorConfig, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	cfg := &TorConfig{}
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		entry := parseLine(line)
		cfg.Entries = append(cfg.Entries, entry)
	}
	return cfg, scanner.Err()
}

// parseLine parses a line from torrc into a TorConfigEntry
func parseLine(line string) TorConfigEntry {
	trimmed := strings.TrimSpace(line)
	if trimmed == "" || strings.HasPrefix(trimmed, "#") {
		return TorConfigEntry{RawLine: line, IsComment: true}
	}

	fields := strings.Fields(trimmed)
	if len(fields) < 2 {
		return TorConfigEntry{RawLine: line, IsComment: false}
	}
	return TorConfigEntry{
		RawLine:   line,
		IsComment: false,
		Key:       fields[0],
		Value:     strings.Join(fields[1:], " "),
	}
}

// Get returns the value of a given config key, if present
func (tc *TorConfig) Get(key string) (string, bool) {
	for _, e := range tc.Entries {
		if !e.IsComment && e.Key == key {
			return e.Value, true
		}
	}
	return "", false
}

// Set updates or appends a key/value pair
func (tc *TorConfig) Set(key string, value string) {
	for i, e := range tc.Entries {
		if !e.IsComment && e.Key == key {
			tc.Entries[i].Value = value
			tc.Entries[i].RawLine = key + " " + value
			return
		}
	}
	tc.Entries = append(tc.Entries, TorConfigEntry{
		Key:     key,
		Value:   value,
		RawLine: key + " " + value,
	})
}

// Save writes the config back to file, preserving formatting
func (tc *TorConfig) Save(path string) error {
	backup := path + ".bak"
	if _, err := os.Stat(path); err == nil {
		_ = os.Rename(path, backup)
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	for _, e := range tc.Entries {
		_, err := f.WriteString(e.RawLine + "\n")
		if err != nil {
			return err
		}
	}
	return nil
}

// GetHiddenServiceDirs returns all configured HiddenServiceDir entries
func (tc *TorConfig) GetHiddenServiceDirs() ([]string, error) {
	var dirs []string
	for _, e := range tc.Entries {
		if !e.IsComment && e.Key == "HiddenServiceDir" {
			dirs = append(dirs, e.Value)
		}
	}
	if len(dirs) == 0 {
		return nil, errors.New("no HiddenServiceDir entries found")
	}
	return dirs, nil
}
