package IrisAPIs

import (
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
}

type ServiceManagement interface {
	RegisterService(service ServiceDescriptor) error
	CheckAllServerStatus() []ServiceStatusRet
	CheckServerStatus(name ServiceStatus) ServiceStatusRet
	//RegisterServices(services []ServiceDescriptor) error
}

type ServiceManagementContext struct {
	services map[string]ServiceDescriptor
}

func (s *ServiceManagementContext) CheckAllServerStatus() []ServiceStatusRet {
	ret := make([]ServiceStatusRet, 0)
	for k, v := range s.services {
		alive, err := v.IsAlive()
		ret = append(ret, ServiceStatusRet{
			Status: func() ServiceStatus {
				if err != nil {
					return StatusInternalError
				}
				if alive {
					return StatusOK
				}

				return StatusDown
			}(),
			Name:        k,
			ServiceType: getType(v),
		})
	}
	return ret
}

func (s *ServiceManagementContext) CheckServerStatus(name ServiceStatus) ServiceStatusRet {
	panic("implement me")
}

func (s *ServiceManagementContext) RegisterService(service ServiceDescriptor) error {
	if _, exist := s.services[service.GetServiceName()]; exist {
		return errors.New("duplicated service name is found!")
	}
	s.services[service.GetServiceName()] = service
	return nil
}

func NewServiceManagementContext() ServiceManagement {
	ret := &ServiceManagementContext{services: make(map[string]ServiceDescriptor)}
	_ = ret.RegisterService(&WebServiceDescriptor{
		Name:    "Wordpress",
		PingUrl: "https://www.rayer.idv.tw/blog/wp-admin/install.php",
	})

	_ = ret.RegisterService(&WebServiceDescriptor{
		Name:    "API",
		PingUrl: "https://api.rayer.idv.tw/ping",
	})

	_ = ret.RegisterService(&WebServiceDescriptor{
		Name:    "Jenkins(Web)",
		PingUrl: "https://jenkins.rayer.idv.tw/login",
	})

	_ = ret.RegisterService(&WebServiceDescriptor{
		Name:    "SupposedFail",
		PingUrl: "https://aa.cc.dde",
	})

	_ = ret.RegisterService(&WebServiceDescriptor{
		Name:    "SupposedFail2",
		PingUrl: "https://api.rayer.idv.tw/notexist",
	})

	return ret
}

func getType(myvar interface{}) string {
	if t := reflect.TypeOf(myvar); t.Kind() == reflect.Ptr {
		return "*" + t.Elem().Name()
	} else {
		return t.Name()
	}
}
