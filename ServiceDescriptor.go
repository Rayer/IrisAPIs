package IrisAPIs

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
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
	Name          string
	ContainerName string
	ImageName     string
	ImageTag      string
	client        *client.Client
}

func (d *DockerComponentDescriptor) GetServiceName() string {
	return d.Name
}

func (d *DockerComponentDescriptor) IsAlive() (bool, error) {
	f := filters.NewArgs()
	f.Add("name", fmt.Sprintf("^/%s$", d.ContainerName))
	ret, err := d.client.ContainerList(context.TODO(), types.ContainerListOptions{
		All:     false,
		Filters: f,
	})

	if err != nil {
		return false, err
	}

	if len(ret) > 1 {
		//Not possible.... but still be there
		return false, errors.Errorf("Multiple container is found by this name : %s", d.ContainerName)
	}

	if len(ret) < 1 {
		return false, nil
	}

	container := ret[0]
	log.Infof("Get container info : %+v", container)

	imageName := func() string {
		if d.ImageTag == "" {
			return d.ImageName
		}
		return fmt.Sprintf("%s:%s", d.ImageName, d.ImageTag)
	}()

	if imageName != container.Image {
		return false, errors.Errorf("Continer %s found, but image mismatch : %s, expected %s", container.Names, container.Names, d.ImageName)
	}

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
	Name             string
	ConnectionString string
}

func (d *DatabaseComponentDescriptor) GetServiceName() string {
	return d.Name
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
