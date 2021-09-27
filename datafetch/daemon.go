package main

import (
	"github.com/sevlyar/go-daemon"
)

func GetDaemon() *daemon.Context {
	return &daemon.Context{
		PidFileName: "datafetch.pid",
		PidFilePerm: 0644,
		LogFileName: "datafetch.log",
		LogFilePerm: 0644,
		WorkDir:     "./",
		Args:        nil,
		Umask:       027,
	}
}
