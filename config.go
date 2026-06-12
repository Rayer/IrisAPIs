package main

import (
	"net/http"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Host     string           `yaml:"host"`
	Version  string           `yaml:"version"`
	Services []ServiceChecker `yaml:"-"`
}

// rawConfig mirrors the YAML structure
type rawConfig struct {
	Host     string        `yaml:"host"`
	Version  string        `yaml:"version"`
	Services []rawService  `yaml:"services"`
}

type rawService struct {
	Name    string        `yaml:"name"`
	URL     string        `yaml:"url"`
	Type    string        `yaml:"type"`    // "http" or "tcp"; auto-detected if empty
	Timeout time.Duration `yaml:"timeout"` // e.g. "10s"
}

func (c *Config) UnmarshalYAML(value *yaml.Node) error {
	var raw rawConfig
	if err := value.Decode(&raw); err != nil {
		return err
	}
	c.Host = raw.Host
	c.Version = raw.Version
	c.Services = make([]ServiceChecker, len(raw.Services))
	for i, s := range raw.Services {
		c.Services[i] = buildChecker(s)
	}
	return nil
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		Host:    "0.0.0.0:8800",
		Version: "dev",
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func buildChecker(s rawService) ServiceChecker {
	timeout := s.Timeout
	if timeout == 0 {
		timeout = 10 * time.Second
	}

	// Explicit type or auto-detect
	checkType := s.Type
	if checkType == "" {
		if len(s.URL) >= 6 && s.URL[:6] == "tcp://" {
			checkType = "tcp"
		} else {
			checkType = "http"
		}
	}

	switch checkType {
	case "tcp":
		addr := s.URL
		if len(addr) >= 6 && addr[:6] == "tcp://" {
			addr = addr[6:]
		}
		return &TCPChecker{Name: s.Name, Addr: addr, Timeout: timeout}
	default:
		return &HTTPChecker{
			Name:    s.Name,
			URL:     s.URL,
			Timeout: timeout,
			client:  &http.Client{Timeout: timeout},
		}
	}
}
