package IrisAPIs

import (
	"encoding/binary"
	"errors"
	log "github.com/sirupsen/logrus"
	"net"
	"regexp"
)

type IpNation struct {
	Ip int64
	Country string
}

type IpNationCountries struct {
	Code      string
	IsoCode_2 string
	IsoCode_3 string
	Country   string
	Lat       float32
	Lon       float32
}

//ip2int and int2ip comes from
//https://gist.github.com/ammario/649d4c0da650162efd404af23e25b86b
func ip2int(ip net.IP) uint32 {
	if len(ip) == 16 {
		return binary.BigEndian.Uint32(ip[12:16])
	}
	return binary.BigEndian.Uint32(ip)
}

func int2ip(nn uint32) net.IP {
	ip := make(net.IP, 4)
	binary.BigEndian.PutUint32(ip, nn)
	return ip
}

func isCorrectIPAddress(ip string) bool {
	ipReg := regexp.MustCompile(`^(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$`)
	return ipReg.MatchString(ip)
}

func GetIPNation(ip string) (*IpNationCountries, error){
	db := GetDatabaseContext().DbObject
	if isCorrectIPAddress(ip) == false {
		return nil, errors.New("Invalid IP address : " + ip)
	}
	ipNet := net.ParseIP(ip)
	ip2n := &IpNation{}
	_, err := db.Where("ip <= ?", ip2int(ipNet)).Desc("ip").Limit(1, 0).Get(ip2n)
	if err != nil {
		return nil, err
	}
	log.Debugf("From ip %s get %+v", ip, ip2n)

	nationInfo := &IpNationCountries{}
	_, err = db.Where("code = ?", ip2n.Country).Get(nationInfo)
	if err != nil {
		return nil, err
	}
	log.Debugf("From country %+v get %+v", ip2n, nationInfo)
	return nationInfo, nil
}