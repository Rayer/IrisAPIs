package IrisAPIs

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

type ServiceDescriptor interface {
	GetServiceName() string
	IsAlive(ctx context.Context) (bool, error)
	Install() error
	Shutdown() error
	Startup() error
	Restart() error
	Logs(ctx context.Context) (string, error)
}

var gDockerClient *client.Client

func init() {
	gDockerClient, _ = client.NewClientWithOpts(client.FromEnv)
}

type DockerComponentDescriptor struct {
	containerParam types.Container
	Name           string
	ContainerName  string
	ImageName      string
	ImageTag       string
	client         *client.Client
}

func NewDockerComponentDescriptor(ctx context.Context, name string, containerName string, imageName string, imageTag string) ServiceDescriptor {
	logger := GetLogger(ctx)
	if gDockerClient == nil {
		logger.Warnf("Docker client initialization failed! skipped service [%s]", name)
		return nil
	}
	ret := &DockerComponentDescriptor{Name: name, ContainerName: containerName, ImageName: imageName, ImageTag: imageTag, client: gDockerClient}
	err := ret.refreshDockerParameters(ctx)
	if err != nil {
		//refresh ID fail, it means currently it is not in docker container, but we still can create ServiceDescriptor
		logger.Warnf("Service %s is not running in docker daemon now, but still monitoring.", name)
	}

	return ret
}

func (d *DockerComponentDescriptor) refreshDockerParameters(ctx context.Context) error {

	f := filters.NewArgs()
	f.Add("name", fmt.Sprintf("^/%s$", d.ContainerName))
	ret, err := d.client.ContainerList(ctx, types.ContainerListOptions{
		All:     false,
		Filters: f,
	})

	if err != nil {
		return err
	}

	if len(ret) > 1 {
		//Not possible.... but still be there
		return errors.Errorf("Multiple container is found by this name : %s", d.ContainerName)
	}

	if len(ret) < 1 {
		return errors.Errorf("No docker with name %s with service %s is found, id is not refreshed!", d.ContainerName, d.Name)
	}

	d.containerParam = ret[0]
	return nil
}

func (d *DockerComponentDescriptor) GetServiceName() string {
	return d.Name
}

func (d *DockerComponentDescriptor) IsAlive(ctx context.Context) (bool, error) {
	err := d.refreshDockerParameters(ctx)
	if err != nil {
		return false, err
	}

	container := d.containerParam
	log.Infof("Get container info : %+v", container)

	imageName := func() string {
		if d.ImageTag == "" {
			return d.ImageName
		}
		return fmt.Sprintf("%s:%s", d.ImageName, d.ImageTag)
	}()

	if imageName != container.Image {
		return false, errors.Errorf("Continer %s found, but image mismatch : %s, expected %s", container.Names, container.Image, d.ImageName)
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

func (d *DockerComponentDescriptor) Logs(ctx context.Context) (string, error) {
	options := types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true}
	op, err := d.client.ContainerLogs(ctx, d.containerParam.ID, options)
	if err != nil {
		return "", err
	}
	logs, err := ioutil.ReadAll(op)
	defer func() {
		_ = op.Close()
	}()
	if err != nil {
		return "", err
	}
	return string(logs), nil
}

type WebServiceDescriptor struct {
	Name    string
	PingUrl string
}

func (r *WebServiceDescriptor) GetServiceName() string {
	return r.Name
}

func (r *WebServiceDescriptor) IsAlive(context.Context) (bool, error) {
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

func (r *WebServiceDescriptor) Logs(context.Context) (string, error) {
	return "", errors.New("not support this operation")
}

type DatabaseComponentDescriptor struct {
	Name             string
	ConnectionString string
}

func (d *DatabaseComponentDescriptor) GetServiceName() string {
	return d.Name
}

func (d *DatabaseComponentDescriptor) IsAlive(context.Context) (bool, error) {
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

func (d *DatabaseComponentDescriptor) Logs(context.Context) (string, error) {
	return "", errors.New("not support this operation")
}
