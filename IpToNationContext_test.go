package IrisAPIs

import (
	"context"
	"errors"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"reflect"
	"testing"
)

type IpToNationContextTestSuite struct {
	suite.Suite
	db              *DatabaseContext
	ipNationContext *IpNationContext
}

func TestIpToNationContextSuite(t *testing.T) {
	suite.Run(t, new(IpToNationContextTestSuite))
}

func (c *IpToNationContextTestSuite) SetupSuite() {
	c.db, _ = NewTestDatabaseContext(context.TODO())
	c.ipNationContext = NewIpNationContext(c.db)
	logrus.SetLevel(logrus.DebugLevel)
}

func (c *IpToNationContextTestSuite) SetupTest() {
	if c.db == nil {
		c.T().Errorf("Can't initialize database!")
		c.T().Fail()
	}
}

func (c *IpToNationContextTestSuite) Test_isCorrectIPAddress() {
	type args struct {
		ip string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Valid IP address 1",
			args: args{"1.3.11.23"},
			want: true,
		},
		{
			name: "Valid IP address 2",
			args: args{"125.22.14.88"},
			want: true,
		},
		{
			name: "Invalid IP address -- greater then 255",
			args: args{"2.114.257.221"},
			want: false,
		},
		{
			name: "Not even an IP address",
			args: args{
				ip: "hello.world!",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		c.Run(tt.name, func() {
			if got := isCorrectIPAddress(tt.args.ip); got != tt.want {
				assert.Equal(c.T(), got, tt.want)
				c.Errorf(errors.New("Fail in test"), "isCorrectIPAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}

func (c *IpToNationContextTestSuite) TestGetIPNation() {
	type args struct {
		ip string
	}
	tests := []struct {
		name    string
		args    args
		want    *IpNationCountries
		wantErr bool
	}{
		{
			name: "Normal test cases : Taiwan",
			args: args{
				ip: "1.169.18.15",
			},
			want: &IpNationCountries{
				Code:      "tw",
				IsoCode_2: "TW",
				IsoCode_3: "TWN",
				Country:   "Taiwan",
				Lat:       23.3,
				Lon:       121,
			},
			wantErr: false,
		},
		{
			name: "Private IP",
			args: args{
				ip: "10.0.1.4",
			},
			want: &IpNationCountries{
				Code:      "01",
				IsoCode_2: "",
				IsoCode_3: "",
				Country:   "Private",
				Lat:       0,
				Lon:       0,
			},
			wantErr: false,
		},
		{
			name: "Invalid IP address",
			args: args{
				ip: "11.165.283.217",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Localhost IP Test",
			args: args{
				ip: "127.0.0.1",
			},
			want: &IpNationCountries{
				Code:      "03",
				IsoCode_2: "",
				IsoCode_3: "",
				Country:   "(Reserved)Local Host",
				Lat:       0,
				Lon:       0,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		c.Run(tt.name, func() {
			got, err := c.ipNationContext.GetIPNation(tt.args.ip)
			if (err != nil) != tt.wantErr {
				assert.Equal(c.T(), tt.wantErr, err)
				c.Errorf(err, "GetIPNation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				assert.Equal(c.T(), got, tt.want)
				c.Errorf(errors.New(""), "GetIPNation() got = %v, want %v", got, tt.want)
			}
		})
	}
}
