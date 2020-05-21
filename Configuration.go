package IrisAPIs

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"reflect"
)

type Configuration struct {
	FixerIoApiKey    string `doc:"API Key of fixer.io, you can get one on its website"`
	ConnectionString string `doc:"Connection string to database."`
	DatabaseType     string `doc:"Database Type, for example, mysql"`
	LogLevel         string `doc:"Log Level, 0 for debug and 7 for info"`
}

func (c *Configuration) LoadConfiguration() error {
	viper.SetConfigName("iris-apis")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("/etc/iris/")

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
