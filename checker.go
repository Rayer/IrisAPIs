package main

import (
	"fmt"
	"net"
	"net/http"
	"time"
)

// ServiceChecker knows how to check a service's health
type ServiceChecker interface {
	GetName() string
	Check() ServiceStatus
}

// ServiceStatus is the result of a health check
type ServiceStatus struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
	Latency string `json:"latency"`
}

// --- HTTP Checker ---

type HTTPChecker struct {
	Name    string
	URL     string
	Timeout time.Duration
	client  *http.Client
}

func (h *HTTPChecker) GetName() string { return h.Name }

func (h *HTTPChecker) Check() ServiceStatus {
	start := time.Now()
	req, err := http.NewRequest("GET", h.URL, nil)
	if err != nil {
		return ServiceStatus{
			Name: h.Name, Type: "http", Status: "ERROR",
			Message: err.Error(), Latency: time.Since(start).String(),
		}
	}

	resp, err := h.client.Do(req)
	latency := time.Since(start)
	if err != nil {
		return ServiceStatus{
			Name: h.Name, Type: "http", Status: "DOWN",
			Message: err.Error(), Latency: latency.String(),
		}
	}
	resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		return ServiceStatus{
			Name: h.Name, Type: "http", Status: "OK",
			Latency: latency.String(),
		}
	}
	return ServiceStatus{
		Name: h.Name, Type: "http", Status: "DOWN",
		Message: fmt.Sprintf("HTTP %d", resp.StatusCode),
		Latency: latency.String(),
	}
}

// --- TCP Checker ---

type TCPChecker struct {
	Name    string
	Addr    string
	Timeout time.Duration
}

func (t *TCPChecker) GetName() string { return t.Name }

func (t *TCPChecker) Check() ServiceStatus {
	start := time.Now()
	conn, err := net.DialTimeout("tcp", t.Addr, t.Timeout)
	latency := time.Since(start)
	if err != nil {
		return ServiceStatus{
			Name: t.Name, Type: "tcp", Status: "DOWN",
			Message: err.Error(), Latency: latency.String(),
		}
	}
	conn.Close()
	return ServiceStatus{
		Name: t.Name, Type: "tcp", Status: "OK",
		Latency: latency.String(),
	}
}
