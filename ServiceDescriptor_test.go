package IrisAPIs

import (
	"github.com/docker/docker/client"
	"github.com/stretchr/testify/suite"
	"testing"
)

type ServiceDescriptorSuite struct {
	suite.Suite
}

func TestDockerComponentDescriptor_IsAlive(t *testing.T) {
	type fields struct {
		Name          string
		ContainerName string
		ImageName     string
		ImageTag      string
		client        *client.Client
	}
	tests := []struct {
		name    string
		fields  fields
		want    bool
		wantErr bool
	}{
		{
			name: "Generic Test",
			fields: fields{
				Name:          "IrisAPI",
				ContainerName: "APIService",
				ImageName:     "rayer/iris-apis",
				ImageTag:      "latest",
				client: func() *client.Client {
					ret, _ := client.NewEnvClient()
					return ret
				}(),
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DockerComponentDescriptor{
				Name:          tt.fields.Name,
				ContainerName: tt.fields.ContainerName,
				ImageName:     tt.fields.ImageName,
				ImageTag:      tt.fields.ImageTag,
				client:        tt.fields.client,
			}
			got, err := d.IsAlive()
			if (err != nil) != tt.wantErr {
				t.Errorf("IsAlive() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("IsAlive() got = %v, want %v", got, tt.want)
			}
		})
	}
}
