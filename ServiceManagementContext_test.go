package IrisAPIs

import (
	"fmt"
	"testing"
)

func TestServiceManagementContext_CheckAllServerStatus(t *testing.T) {
	s := NewServiceManagementContext()
	_ = s.RegisterService(&WebServiceDescriptor{
		Name:    "Wordpress",
		PingUrl: "https://www.rayer.idv.tw/blog/wp-admin/install.php",
	})

	_ = s.RegisterService(&WebServiceDescriptor{
		Name:    "API",
		PingUrl: "https://api.rayer.idv.tw/ping",
	})

	_ = s.RegisterService(&WebServiceDescriptor{
		Name:    "Jenkins(Web)",
		PingUrl: "https://jenkins.rayer.idv.tw/login",
	})

	_ = s.RegisterService(&WebServiceDescriptor{
		Name:    "SupposedFail",
		PingUrl: "https://aa.cc.dde",
	})

	_ = s.RegisterService(&WebServiceDescriptor{
		Name:    "SupposedFail2",
		PingUrl: "https://api.rayer.idv.tw/notexist",
	})
	fmt.Printf("%+v", s.CheckAllServerStatus())
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
