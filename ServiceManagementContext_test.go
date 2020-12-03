package IrisAPIs

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"testing"
)

func TestServiceManagementContext_CheckAllServerStatus(t *testing.T) {
	s := NewServiceManagementContext()
	err := s.RegisterPresetServices()
	if err != nil {
		log.Warn(err.Error())
	}
	result := s.CheckAllServerStatus()
	for _, v := range result {
		fmt.Printf("%s - %s - %s - %s\n", v.Name, v.Status, v.ServiceType, v.Message)
	}
}

func TestServiceManagementContext_RegisterService(t *testing.T) {
	type fields struct {
		services map[string]ServiceDescriptor
	}
	type args struct {
		service ServiceDescriptor
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "RegularAddWebService",
			fields: fields{
				services: make(map[string]ServiceDescriptor),
			},
			args: args{
				service: &WebServiceDescriptor{
					Name:    "Wordpress",
					PingUrl: "https://rayer.idv.tw/blog",
				},
			},
			wantErr: false,
		},
		{
			name: "DuplicatedName",
			fields: fields{
				services: map[string]ServiceDescriptor{
					"DuplicatedWordPress": &WebServiceDescriptor{
						Name:    "DuplicatedWordPress",
						PingUrl: "https://rayer.idv.tw/blog",
					},
				},
			},
			args: args{
				service: &WebServiceDescriptor{
					Name:    "DuplicatedWordPress",
					PingUrl: "https://rayer.idv.tw/blog",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ServiceManagementContext{
				services: tt.fields.services,
			}
			if err := s.RegisterService(tt.args.service); (err != nil) != tt.wantErr {
				t.Errorf("RegisterService() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
