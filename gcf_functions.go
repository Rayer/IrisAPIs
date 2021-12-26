package IrisAPIs

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

// PubSubMessage is the payload of a Pub/Sub event. Please refer to the docs for
// additional information regarding Pub/Sub events.
type PubSubMessage struct {
	Data []byte `json:"data"`
}

type UpdateCurrencyDatabaseSecretPathPayload struct {
	SecretPath string `json:"secret_path"`
}

type UpdateCurrencyDatabaseConfig struct {
	ApiKey           string `yaml:"api_key"`
	ConnectingString string `yaml:"connecting_string"`
}

func UpdateCurrencyDatabase(ctx context.Context, m PubSubMessage) error {
	secretPathPayload := UpdateCurrencyDatabaseSecretPathPayload{}
	err := json.Unmarshal(m.Data, &secretPathPayload)
	if err != nil {
		return errors.Errorf("update PubSubMessage doesn't contain correct Secret Path data : %v", err)
	}

	b, err := ioutil.ReadFile(secretPathPayload.SecretPath)
	if err != nil {
		return errors.Errorf("secret path brought from PubSubMessage is incorrect, which is %v, err : %v", secretPathPayload.SecretPath, err)
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
