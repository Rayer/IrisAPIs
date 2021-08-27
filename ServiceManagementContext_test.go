package IrisAPIs

import (
	"github.com/docker/distribution/uuid"
	log "github.com/sirupsen/logrus"
	"testing"
)

func TestServiceManagementContext_CheckAllServerStatus(t *testing.T) {
	s := NewServiceManagement()
	err := s.RegisterPresetServices(nil)
	if err != nil {
		log.Warn(err.Error())
	}
	result := s.CheckAllServerStatus(nil)
	for _, v := range result {
		t.Logf("%s - %s - %s - %s - %s\n", v.ID, v.Name, v.Status, v.ServiceType, v.Message)
	}
}

func TestServiceManagementContext_RegisterService(t *testing.T) {
	type fields struct {
		services map[uuid.UUID]ServiceDescriptor
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
				services: make(map[uuid.UUID]ServiceDescriptor),
			},
			args: args{
				service: &WebServiceDescriptor{
					Name:    "Wordpress",
					PingUrl: "https://rayer.idv.tw/blog",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ServiceManagementContext{
				services: tt.fields.services,
			}
			if err := s.RegisterService(nil, tt.args.service); (err != nil) != tt.wantErr {
				t.Errorf("RegisterService() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
