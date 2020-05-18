package IrisAPIs

import (
	"github.com/sirupsen/logrus"
	"reflect"
	"testing"
)

func Test_isCorrectIPAddress(t *testing.T) {
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
		t.Run(tt.name, func(t *testing.T) {
			if got := isCorrectIPAddress(tt.args.ip); got != tt.want {
				t.Errorf("isCorrectIPAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetIPNation(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetIPNation(tt.args.ip)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetIPNation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetIPNation() got = %v, want %v", got, tt.want)
			}
		})
	}
}