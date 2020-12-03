package IrisAPIs

import (
	"github.com/docker/docker/client"
	"github.com/pkg/errors"
	"net/http"
)

type ServiceDescriptor interface {
	GetServiceName() string
	IsAlive() (bool, error)
	Install() error
	Shutdown() error
	Startup() error
	Restart() error
}

type DockerComponentDescriptor struct {
	Name            string
	DockerImageName string
	dockerClient    *client.Client
}

func (d *DockerComponentDescriptor) GetServiceName() string {
	return d.Name
}

func (d *DockerComponentDescriptor) IsAlive() (bool, error) {
	return true, nil
}

func (d *DockerComponentDescriptor) Install() error {
	return errors.New("not support this operation")
}

func (d *DockerComponentDescriptor) Shutdown() error {
	return errors.New("not support this operation")
}

func (d *DockerComponentDescriptor) Startup() error {
	return errors.New("not support this operation")
}

func (d *DockerComponentDescriptor) Restart() error {
	return errors.New("not support this operation")
}

type WebServiceDescriptor struct {
	Name    string
	PingUrl string
}

func (r *WebServiceDescriptor) GetServiceName() string {
	return r.Name
}

func (r *WebServiceDescriptor) IsAlive() (bool, error) {
	resp, err := http.Get(r.PingUrl)
	if err != nil {
		return false, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		return false, nil
	}
	return true, nil
}

func (r *WebServiceDescriptor) Install() error {
	return errors.New("not support this operation")
}

func (r *WebServiceDescriptor) Shutdown() error {
	return errors.New("not support this operation")
}

func (r *WebServiceDescriptor) Startup() error {
	return errors.New("not support this operation")
}

func (r *WebServiceDescriptor) Restart() error {
	return errors.New("not support this operation")
}

type DatabaseComponentDescriptor struct {
	Name string
}

func (d *DatabaseComponentDescriptor) GetServiceName() string {
	panic("implement me")
}

func (d *DatabaseComponentDescriptor) IsAlive() (bool, error) {
	panic("implement me")
}

func (d *DatabaseComponentDescriptor) Install() error {
	return errors.New("not support this operation")
}

func (d *DatabaseComponentDescriptor) Shutdown() error {
	return errors.New("not support this operation")
}

func (d *DatabaseComponentDescriptor) Startup() error {
	return errors.New("not support this operation")
}

func (d *DatabaseComponentDescriptor) Restart() error {
	return errors.New("not support this operation")
}
