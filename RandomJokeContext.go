package IrisAPIs

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"net/http"
)

/*
{"id":102,"type":"general","setup":"Did you hear the one about the guy with the broken hearing aid?","punchline":"Neither did he."}
*/
type RandomJoke struct {
	Id        int
	Type      string
	Setup     string
	Punchline string
}

func fetchRandomJoke(ctx context.Context) (*RandomJoke, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	logger := GetLogger(ctx)
	client := &http.Client{Transport: tr}
	resp, err := client.Get("https://official-joke-api.appspot.com/random_joke")
	if err != nil {
		logger.Error("error:", err)
		return nil, err
	}
	defer resp.Body.Close()
	body := resp.Body

	ret := RandomJoke{}
	err = json.NewDecoder(body).Decode(&ret)
	if err != nil {
		return nil, err
	}

	return &ret, nil
}
