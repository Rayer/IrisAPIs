package IrisAPIs

import (
	"github.com/docker/distribution/uuid"
	"github.com/docker/docker/client"
	"github.com/pkg/errors"
	"reflect"
)

type ServiceStatus string

const (
	StatusOK            ServiceStatus = "OK"
	StatusDown                        = "DOWN"
	StatusInternalError               = "INTERNAL_ERROR"
)

type ServiceStatusRet struct {
	Status      ServiceStatus
	Name        string
	ServiceType string
	Message     string
	ID          uuid.UUID
}

type ServiceManagement interface {
	RegisterService(service ServiceDescriptor) error
	CheckAllServerStatus() []ServiceStatusRet
	CheckServerStatus(id uuid.UUID) (ServiceStatusRet, error)
	RegisterServices(services []ServiceDescriptor) error
	RegisterPresetServices() error
}

type ServiceManagementContext struct {
	services map[uuid.UUID]ServiceDescriptor
}

func (s *ServiceManagementContext) CheckAllServerStatus() []ServiceStatusRet {
	ret := make([]ServiceStatusRet, 0)
	for k := range s.services {
		s, _ := s.CheckServerStatus(k)
		ret = append(ret, s)
	}
	return ret
}

func (s *ServiceManagementContext) CheckServerStatus(id uuid.UUID) (ServiceStatusRet, error) {

	if service, exists := s.services[id]; !exists {
		return ServiceStatusRet{}, errors.Errorf("no such service bound on id : %s", id.String())
	} else {
		alive, err := service.IsAlive()
		return ServiceStatusRet{
			ID: id,
			Status: func() ServiceStatus {
				if err != nil {
					return StatusInternalError
				}
				if alive {
					return StatusOK
				}

				return StatusDown
			}(),
			Name:        service.GetServiceName(),
			ServiceType: getTypeName(service),
			Message: func() string {
				if err != nil {
					return err.Error()
				}
				return ""
			}(),
		}, nil
	}

}

func (s *ServiceManagementContext) RegisterService(service ServiceDescriptor) error {
	//if _, exist := s.services[service.GetServiceName()]; exist {
	//	return errors.New("duplicated service name is found!")
	//}

	s.services[uuid.Generate()] = service
	return nil
}

func (s *ServiceManagementContext) RegisterServices(services []ServiceDescriptor) error {
	for _, v := range services {
		err := s.RegisterService(v)
		if err != nil {
			return err
		}
	}
	return nil
}

func NewServiceManagement() ServiceManagement {
	ret := &ServiceManagementContext{services: make(map[uuid.UUID]ServiceDescriptor)}
	return ret
}

func getTypeName(inVar interface{}) string {
	if t := reflect.TypeOf(inVar); t.Kind() == reflect.Ptr {
		return t.Elem().Name()
	} else {
		return t.Name()
	}
}

func (s *ServiceManagementContext) RegisterPresetServices() error {
	dc, err := client.NewEnvClient()
	if err != nil {
		return err
	}
	return s.RegisterServices([]ServiceDescriptor{
		&DockerComponentDescriptor{
			Name:          "IrisAPI",
			ContainerName: "APIService",
			ImageName:     "rayer/iris-apis",
			ImageTag:      "release",
			client:        dc,
		},
		&DockerComponentDescriptor{
			Name:          "IrisAPI-Test",
			ContainerName: "APIService-Test",
			ImageName:     "rayer/iris-apis",
			ImageTag:      "latest",
			client:        dc,
		},
		&DockerComponentDescriptor{
			Name:          "OneIndex",
			ContainerName: "oneindex-service",
			ImageName:     "setzero/oneindex",
			ImageTag:      "",
			client:        dc,
		},
		&DockerComponentDescriptor{
			Name:          "Jenkins-Docker",
			ContainerName: "jenkins-service",
			ImageName:     "jenkins/jenkins",
			ImageTag:      "alpine",
			client:        dc,
		},
		&DockerComponentDescriptor{
			Name:          "MTProxy",
			ContainerName: "mtproxy",
			ImageName:     "telegrammessenger/proxy",
			ImageTag:      "",
			client:        dc,
		},
		&WebServiceDescriptor{
			Name:    "WordPress",
			PingUrl: "https://www.rayer.idv.tw/blog/wp-admin/install.php",
		},
		&WebServiceDescriptor{
			Name:    "IrisAPI",
			PingUrl: "https://api.rayer.idv.tw/ping",
		},
		&WebServiceDescriptor{
			Name:    "IrisAPI-Test",
			PingUrl: "https://api.rayer.idv.tw/ping",
		},
		&WebServiceDescriptor{
			Name:    "Jenkins-Web",
			PingUrl: "https://jenkins.rayer.idv.tw/login",
		},
	})

}
