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
	Message     string
}

type ServiceManagement interface {
	RegisterService(service ServiceDescriptor) error
	CheckAllServerStatus() []ServiceStatusRet
	CheckServerStatus(name ServiceStatus) ServiceStatusRet
	RegisterServices(services []ServiceDescriptor) error
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
			ServiceType: getTypeName(v),
			Message: func() string {
				if err != nil {
					return err.Error()
				}
				return ""
			}(),
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

func (s *ServiceManagementContext) RegisterServices(services []ServiceDescriptor) error {
	for _, v := range services {
		err := s.RegisterService(v)
		if err != nil {
			return err
		}
	}
	return nil
}

func NewServiceManagementContext() ServiceManagement {
	ret := &ServiceManagementContext{services: make(map[string]ServiceDescriptor)}
	return ret
}

func getTypeName(inVar interface{}) string {
	if t := reflect.TypeOf(inVar); t.Kind() == reflect.Ptr {
		return t.Elem().Name()
	} else {
		return t.Name()
	}
}
