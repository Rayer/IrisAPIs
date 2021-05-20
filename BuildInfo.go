package IrisAPIs

import (
	"context"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

type BuildInfo struct {
	ImageTag        string `yaml:"image_tag"`
	CreateTimestamp int64  `yaml:"create_timestamp"`
	JenkinsLink     string `yaml:"jenkins_link"`
	TimeZone        string `yaml:"time_zone"`
	TimeZoneOffset  int    `yaml:"time_zone_offset"`
}

type BuildInfoService interface {
	GetBuildInfo(ctx context.Context) *BuildInfo
}

type BuildInfoServiceImpl struct {
	cached *BuildInfo
}

func NewBuildInfoService() BuildInfoService {
	return &BuildInfoServiceImpl{}
}

func (b *BuildInfoServiceImpl) GetBuildInfo(ctx context.Context) *BuildInfo {
	log := GetLogger(ctx)
	if b.cached != nil {
		return b.cached
	}

	file := "./release-info.yaml"
	log.Debugf("Read configuration from %s", file)
	out, err := ioutil.ReadFile(file)
	if err != nil {
		log.Warningf("Fail to read %s!", file)
		return b.createDefaultBuildInfo()
	}
	var ret BuildInfo
	err = yaml.Unmarshal(out, &ret)
	if err != nil {
		log.Warningf("%s corrupted!", file)
		return b.createDefaultBuildInfo()
	}
	b.cached = &ret
	return b.cached
}

func (b *BuildInfoServiceImpl) createDefaultBuildInfo() *BuildInfo {
	b.cached = &BuildInfo{
		ImageTag:        "<No info>",
		CreateTimestamp: 0,
		JenkinsLink:     "<No info>",
		TimeZone:        "",
		TimeZoneOffset:  0,
	}
	return b.cached
}
