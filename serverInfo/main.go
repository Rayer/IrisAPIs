package main

import (
	"IrisAPIs"
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"time"
)

func main() {
	// output current time zone
	fmt.Print("Local time zone ")
	zone, offset := time.Now().Zone()
	fmt.Println(time.Now().Zone())
	fmt.Println(time.Now().Format("2006-01-02T15:04:05.000 MST"))
	fmt.Println(os.Hostname())

	filename := "release-info.yaml"
	buildInfo := IrisAPIs.BuildInfoFile{
		ImageTag:        os.Getenv("IMAGE_TAG"),
		CreateTimestamp: time.Now().Unix(),
		JenkinsLink:     os.Getenv("JENKINS_LINK"),
		TimeZone:        zone,
		TimeZoneOffset:  offset / 3600,
	}

	fmt.Println("Writing file : " + filename)
	fmt.Printf("Build info will be written into file : %+v", buildInfo)
	out, err := yaml.Marshal(&buildInfo)
	if err != nil {
		fmt.Println("Error while generating build info yaml file:", err.Error())
		os.Exit(1)
	}
	err = ioutil.WriteFile(filename, out, 0644)
	if err != nil {
		fmt.Println("Error while writing build info yaml file:", err.Error())
		os.Exit(1)
	}
}
