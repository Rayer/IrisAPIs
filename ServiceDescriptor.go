package IrisAPIs

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strings"
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
	var err error
	gDockerClient, err = client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		log.Warnf("Failed to initialize Docker client: %v", err)
	}
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
	ret, err := d.client.ContainerList(ctx, container.ListOptions{
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

func (d *DockerComponentDescriptor) isImageNameTagMatch() bool {
	imageNameInDocker := d.containerParam.Image
	expectedImageName := d.ImageName
	expectedImageTag := d.ImageTag
	slice := strings.Split(imageNameInDocker, ":")
	imageName := slice[0]
	tagName := ""
	if len(slice) > 1 {
		tagName = slice[1]
	}
	if expectedImageTag == "" {
		//It doesn't care about tag name, just compare image name
		return imageName == expectedImageName
	} else {
		return imageName == expectedImageName && tagName == expectedImageTag
	}
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

	if d.isImageNameTagMatch() == false {
		return false, errors.Errorf("Container %s found, but image mismatch: %s, expected %s(tag: %s)", container.Names, container.Image, d.ImageName, d.ImageTag)
	}

	return true, nil
}

func (d *DockerComponentDescriptor) Install() error {
	return errors.New("operation not supported for DockerComponentDescriptor")
}

func (d *DockerComponentDescriptor) Shutdown() error {
	return errors.New("operation not supported for DockerComponentDescriptor")
}

func (d *DockerComponentDescriptor) Startup() error {
	return errors.New("operation not supported for DockerComponentDescriptor")
}

func (d *DockerComponentDescriptor) Restart() error {
	return errors.New("operation not supported for DockerComponentDescriptor")
}

func (d *DockerComponentDescriptor) Logs(ctx context.Context) (string, error) {
	options := container.LogsOptions{ShowStdout: true, ShowStderr: true}
	op, err := d.client.ContainerLogs(ctx, d.containerParam.ID, options)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = op.Close()
	}()
	logs, err := io.ReadAll(op)
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
	return errors.New("operation not supported for WebServiceDescriptor")
}

func (r *WebServiceDescriptor) Shutdown() error {
	return errors.New("operation not supported for WebServiceDescriptor")
}

func (r *WebServiceDescriptor) Startup() error {
	return errors.New("operation not supported for WebServiceDescriptor")
}

func (r *WebServiceDescriptor) Restart() error {
	return errors.New("operation not supported for WebServiceDescriptor")
}

func (r *WebServiceDescriptor) Logs(context.Context) (string, error) {
	return "", errors.New("operation not supported for WebServiceDescriptor")
}

type DatabaseComponentDescriptor struct {
	Name             string
	ConnectionString string
}

func (d *DatabaseComponentDescriptor) GetServiceName() string {
	return d.Name
}

func (d *DatabaseComponentDescriptor) IsAlive(ctx context.Context) (bool, error) {
	// This is a placeholder implementation
	// In a real implementation, you would check the database connection
	// For example, by pinging the database or executing a simple query
	return false, errors.New("database connectivity check not implemented")
}

func (d *DatabaseComponentDescriptor) Install() error {
	return errors.New("operation not supported for DatabaseComponentDescriptor")
}

func (d *DatabaseComponentDescriptor) Shutdown() error {
	return errors.New("operation not supported for DatabaseComponentDescriptor")
}

func (d *DatabaseComponentDescriptor) Startup() error {
	return errors.New("operation not supported for DatabaseComponentDescriptor")
}

func (d *DatabaseComponentDescriptor) Restart() error {
	return errors.New("operation not supported for DatabaseComponentDescriptor")
}

func (d *DatabaseComponentDescriptor) Logs(context.Context) (string, error) {
	return "", errors.New("operation not supported for DatabaseComponentDescriptor")
}
