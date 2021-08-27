package IrisAPIs

import (
	"github.com/docker/distribution/uuid"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
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
	RegisterService(ctx context.Context, service ServiceDescriptor) error
	CheckAllServerStatus(ctx context.Context) []ServiceStatusRet
	CheckServerStatus(ctx context.Context, id uuid.UUID) (ServiceStatusRet, error)
	GetLogs(ctx context.Context, id uuid.UUID) (string, error)
	RegisterServices(ctx context.Context, services []ServiceDescriptor) error
	RegisterPresetServices(ctx context.Context) error
}

type ServiceManagementContext struct {
	services map[uuid.UUID]ServiceDescriptor
}

func (s *ServiceManagementContext) CheckAllServerStatus(ctx context.Context) []ServiceStatusRet {
	ret := make([]ServiceStatusRet, 0)
	for k := range s.services {
		s, _ := s.CheckServerStatus(ctx, k)
		ret = append(ret, s)
	}
	return ret
}

func (s *ServiceManagementContext) CheckServerStatus(ctx context.Context, id uuid.UUID) (ServiceStatusRet, error) {

	service := s.getService(id)
	if service == nil {
		return ServiceStatusRet{}, errors.Errorf("no such service bound on id : %s", id.String())
	}

	alive, err := service.IsAlive(ctx)
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

func (s *ServiceManagementContext) RegisterService(_ context.Context, service ServiceDescriptor) error {
	if service == nil {
		return nil
	}
	s.services[uuid.Generate()] = service
	return nil
}

func (s *ServiceManagementContext) RegisterServices(ctx context.Context, services []ServiceDescriptor) error {
	for _, v := range services {
		err := s.RegisterService(ctx, v)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *ServiceManagementContext) GetLogs(ctx context.Context, id uuid.UUID) (string, error) {
	service := s.getService(id)
	if service == nil {
		return "", errors.Errorf("No such service with id have been found : %s", id)
	}
	return service.Logs(ctx)

}

func (s *ServiceManagementContext) getService(id uuid.UUID) ServiceDescriptor {
	if service, exists := s.services[id]; !exists {
		return nil
	} else {
		return service
	}
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

func (s *ServiceManagementContext) RegisterPresetServices(ctx context.Context) error {
	return s.RegisterServices(ctx, []ServiceDescriptor{
		NewDockerComponentDescriptor(ctx, "Iris Mainframe API", "APIService", "rayer/iris-apis", "release"),
		NewDockerComponentDescriptor(ctx, "Iris Mainframe API(Test)", "APIService-Test", "rayer/iris-apis", "latest"),
		NewDockerComponentDescriptor(ctx, "OneIndex Service", "oneindex-service", "setzero/oneindex", ""),
		NewDockerComponentDescriptor(ctx, "Jenkins Docker Service", "jenkins-service", "jenkins/jenkins", "alpine"),
		NewDockerComponentDescriptor(ctx, "AppleNCCMonitor", "AppleProductMonitor", "rayer/apple-product-monitor", ""),
		//NewDockerComponentDescriptor(ctx, "MTProxy", "mtproxy", "telegrammessenger/proxy", ""),
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
