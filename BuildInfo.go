package IrisAPIs

type BuildInfoFile struct {
	ImageTag        string `yaml:"image_tag"`
	CreateTimestamp int64  `yaml:"create_timestamp"`
	JenkinsLink     string `yaml:"jenkins_link"`
	TimeZone        string `yaml:"time_zone"`
	TimeZoneOffset  int    `yaml:"time_zone_offset"`
}
