package IrisAPIs

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/url"
	"reflect"
)

type Configuration struct {
	Host                             string `doc:"Host name with port"`
	EnforceApiKey                    bool   `doc:"Enforce API Key check, default is true"`
	FixerIoApiKey                    string `doc:"API Key of fixer.io, you can get one on its website"`
	FixerIoLastFetchSuccessfulPeriod int    `doc:"Fetch interval for last successful fetch"`
	FixerIoLastFetchFailedPeriod     int    `doc:"Fetch interval for last fail fetch"`
	ConnectionString                 string `doc:"Connection string to database."`
	TestConnectionString             string `doc:"Connection string to Test Database"`
	TimeZone                         string `doc:"TimeZone, default is Asia/Taipei"`
	DatabaseType                     string `doc:"Database Type, for example, mysql"`
	LogLevel                         int    `doc:"Log Level, 6 for trace and 0 for fatal"`
	GRPCServerHost                   string `doc:"gRPC Server address for gRPC Server config, default is :8082"`
	GRPCServerTarget                 string `doc:"gRPC Server address for gRPC Client config, default is :8082"`
}

func NewConfiguration() *Configuration {
	ret := &Configuration{}
	_ = ret.LoadConfiguration()
	return ret
}

func (c *Configuration) LoadConfiguration() error {
	viper.SetConfigName("iris-apis")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("../config")
	viper.AddConfigPath("/etc/iris/")

	//Defaults
	viper.SetDefault("FixerIoLastFetchSuccessfulPeriod", 43200)
	viper.SetDefault("FixerIoLastFetchFailedPeriod", 10800)
	viper.SetDefault("Host", "localhost:8080")
	viper.SetDefault("EnforceApiKey", true)
	viper.SetDefault("TimeZone", "Asia/Taipei")
	viper.SetDefault("GRPCServerHost", ":8082")
	viper.SetDefault("GRPCServerTarget", ":8082")

	err := viper.ReadInConfig()

	if err != nil {
		logrus.Errorf("Not able to find configuration file.")
		fmt.Println(c.ExampleUsage())
		return err
	}

	err = viper.Unmarshal(c)

	if err != nil {
		return err
	}
	return nil
}

func (c *Configuration) ExampleUsage() string {
	v := reflect.ValueOf(c).Elem()
	ret := "Here is a template of content for iris-apis.yaml : \n\n"
	for i := 0; i < v.NumField(); i++ {
		name := v.Type().Field(i).Name
		tag := v.Type().Field(i).Tag.Get("doc")
		ret += fmt.Sprintf("# %s\n%s:\n", tag, name)
	}
	return ret
}

func (c *Configuration) SplitSchemeAndHost() (string, string, error) {
	u, err := url.ParseRequestURI(c.Host)
	if err != nil {
		return "", "", fmt.Errorf("Error parsing host : "+c.Host, err)
	}
	return u.Scheme, u.Host, nil
}
