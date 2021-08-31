package IrisAPIs

import (
	"bytes"
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type ServiceDescriptorSuite struct {
	suite.Suite
	dockerClient      *client.Client
	dockerAvailable   bool
	dockerImgId       string
	dockerContainerId string
}

func TestServiceDescriptorSuite(t *testing.T) {
	suite.Run(t, new(ServiceDescriptorSuite))
}

func (s *ServiceDescriptorSuite) SetupSuite() {
	log.SetLevel(log.DebugLevel)
	var err error
	s.dockerAvailable = true
	s.dockerClient, err = client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		log.Warnf("Error creating docker client! (%s)", err.Error())
		s.dockerAvailable = false
	}

	_, err = s.dockerClient.Ping(context.TODO())
	if err != nil {
		log.Warnf("Error creating docker client! (%s)", err.Error())
		s.dockerAvailable = false
	}

	if !s.dockerAvailable {
		log.Warn("Docker client initialization failed, will skip all docker test cases.")
	} else {
		//Download and init docker tester
		p, err := s.dockerClient.ImagePull(context.TODO(), "docker.io/rayer/chatbot-server:latest", types.ImagePullOptions{})
		defer func() {
			err := p.Close()
			if err != nil {
				log.Warnf("Error while closing ImagePull stream : %s", err.Error())
			}
		}()
		buf := new(bytes.Buffer)
		_, _ = buf.ReadFrom(p)

		if err != nil {
			//log.Warnf("Fail to pull image : %s", err.Error())
			panic(err)
		} else {
			log.Info("Finished pulling test docker image : ", buf.String())
		}

		cr, err := s.dockerClient.ContainerCreate(context.TODO(), &container.Config{
			Image: "rayer/chatbot-server:latest",
			//Cmd: []string{"echo", "Running UT Docker container..."},
		}, nil, nil, "UTDocker")
		if err != nil {
			//log.Warnf("Fail to create container from image : %s", err.Error())
			panic(err)
		} else {
			log.Info("Created container with ID :", cr.ID)
			s.dockerContainerId = cr.ID
		}

		err = s.dockerClient.ContainerStart(context.TODO(), cr.ID, types.ContainerStartOptions{})

		if err != nil {
			panic(err)
		}

		log.Info("Container created!")

	}

}

func (s *ServiceDescriptorSuite) TearDownSuite() {
	if s.dockerClient != nil {
		//s.dockerClient.ImageRemove(context.TODO(), )
		err := s.dockerClient.ContainerRemove(context.TODO(), s.dockerContainerId, types.ContainerRemoveOptions{Force: true})
		if err != nil {
			log.Warnf("Error removing container %s! - %s", s.dockerContainerId, err.Error())
		}
		_ = s.dockerClient.Close()
	}
}

func (s *ServiceDescriptorSuite) TestDockerComponentDescriptor_IsAlive() {
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
				Name:          "TestDocker",
				ContainerName: "UTDocker",
				ImageName:     "rayer/chatbot-server",
				ImageTag:      "latest",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Wrong Container Name",
			fields: fields{
				Name:          "TestDockerNone",
				ContainerName: "UTDockerNotHere",
				ImageName:     "rayer/chatbot-server",
				ImageTag:      "latest",
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Wrong Image Name",
			fields: fields{
				Name:          "TestDockerWrongName",
				ContainerName: "UTDocker",
				ImageName:     "rayer/chatbot-server-Wrong",
				ImageTag:      "latest",
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Success With Tag Waive",
			fields: fields{
				Name:          "TestDockerTagWaive",
				ContainerName: "UTDocker",
				ImageName:     "rayer/chatbot-server",
				ImageTag:      "",
			},
			want:    true,
			wantErr: false,
		},
	}

	if !s.dockerAvailable {
		s.T().Skip("Can't find docker instance, skip this test")
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			d := NewDockerComponentDescriptor(context.Background(), tt.fields.Name, tt.fields.ContainerName, tt.fields.ImageName, tt.fields.ImageTag)

			got, err := d.IsAlive(context.TODO())
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

func (s *ServiceDescriptorSuite) TestWebServiceDescriptor_IsAlive() {
	type fields struct {
		Name    string
		PingUrl string
	}
	tests := []struct {
		name    string
		fields  fields
		want    bool
		wantErr bool
	}{
		{
			name: "NormalWebSuccess",
			fields: fields{
				Name:    "Hinet Service",
				PingUrl: "https://www.hinet.net",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "NormalWebFail",
			fields: fields{
				Name:    "RandomFQDN",
				PingUrl: "https://aa.nbnb.ears.c",
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Normal404",
			fields: fields{
				Name:    "Hinet Service 404",
				PingUrl: "https://www.hinet.net/acccc",
			},
			want:    false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			r := &WebServiceDescriptor{
				Name:    tt.fields.Name,
				PingUrl: tt.fields.PingUrl,
			}
			got, err := r.IsAlive(nil)
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

func (s *ServiceDescriptorSuite) TestDockerComponentDescriptor_Logs() {
	type fields struct {
		Name          string
		ContainerName string
		ImageName     string
		ImageTag      string
	}
	tests := []struct {
		name    string
		fields  fields
		haveLog bool
		wantErr bool
	}{
		{
			name: "LogSuccess",
			fields: fields{
				Name:          "TestDocker",
				ContainerName: "UTDocker",
				ImageName:     "rayer/chatbot-server",
				ImageTag:      "latest",
			},
			haveLog: true,
			wantErr: false,
		},
		{
			name: "NoContainer",
			fields: fields{
				Name:          "TestDockerNotHere",
				ContainerName: "UTDockerNotHere",
				ImageName:     "rayer/chatbot-server",
				ImageTag:      "latest",
			},
			haveLog: false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			d := NewDockerComponentDescriptor(context.Background(), tt.fields.Name, tt.fields.ContainerName, tt.fields.ImageName, tt.fields.ImageTag)
			time.Sleep(2 * time.Second)
			got, err := d.Logs(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("Logs() error = %v, wantErr %v", err, tt.wantErr)
			}
			if len(got) > 1 != tt.haveLog {
				t.Errorf("Expected haveLog == %v but it's %v", tt.haveLog, len(got) > 1)
			}
		})
	}
}
