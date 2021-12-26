package IrisAPIs

import (
	"context"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

// PubSubMessage is the payload of a Pub/Sub event. Please refer to the docs for
// additional information regarding Pub/Sub events.
type PubSubMessage struct {
	SecretPath string `json:"secret_path"`
}

type UpdateCurrencyDatabaseConfig struct {
	ApiKey           string `yaml:"api_key"`
	ConnectingString string `yaml:"connecting_string"`
}

func UpdateCurrencyDatabase(ctx context.Context, m PubSubMessage) error {

	b, err := ioutil.ReadFile(m.SecretPath)
	if err != nil {
		return errors.Errorf("secret path brought from PubSubMessage is incorrect, which is %v, err : %v", m.SecretPath, err)
	}

	config := UpdateCurrencyDatabaseConfig{}
	err = yaml.Unmarshal(b, &config)
	if err != nil {
		return errors.Errorf("failed to unmarshal secret!")
	}

	db, err := NewDatabaseContext(config.ConnectingString, false, nil)
	if err != nil {
		return errors.Errorf("failed to initialize database, err : %v", err)
	}

	c := NewCurrencyContextWithConfig(config.ApiKey, 100, 100, db)
	err = c.SyncToDb()
	if err != nil {
		return errors.Errorf("sync to db failed, err : %v", err)
	}

	return nil
}
