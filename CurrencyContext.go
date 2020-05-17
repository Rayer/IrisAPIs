package IrisAPIs

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"time"
)

type CurrencyBatch struct {
	Batch int64 `xorm:"autoincr"`
	Exec time.Time `xorm:"created"`
	Raw string
	Success bool
}

type CurrencyEntry struct {
	Symbol string `xorm:"varchar(3)"`
	Base string `xorm:"varchar(3)"`
	Batch int64
	Rate float64
}

func SyncToDb () error {
	log.Debugf("Trying connecting to currency server....")
	resp, err := http.Get("http://data.fixer.io/api/latest?access_key=676ac77e5ce5d4b9a57ee6464ff84433")
	if err != nil {
		return err
	}
	defer func () {
		_ = resp.Body.Close()
	}()
	raw, err := ioutil.ReadAll(resp.Body)
	var result map[string]interface{}
	err = json.Unmarshal(raw, &result)
	if err != nil {
		return err
	}

	err = saveCurrencyEntries(string(raw))
	if err != nil {
		return err
	}

	return nil
}

func saveCurrencyEntries(raw string) error {
	engine := GetDatabaseContext().DbObject

	batch := &CurrencyBatch{
		//ExecTime: time.Now(),
		Raw:      raw,
		Success:  true,
	}
	_, err := engine.Insert(batch)

	var entries = make([]*CurrencyEntry, 0, 0)
	var data map[string]interface{}
	json.Unmarshal([]byte(raw), &data)

	base := data["base"]
	rates := data["rates"].(map[string]interface{})

	for sym, rate := range rates {
		entries = append(entries, &CurrencyEntry{
			Symbol: sym,
			Base:   base.(string),
			Batch:  batch.Batch,
			Rate:   rate.(float64),
		})
	}

	aff, err := engine.Insert(entries)

	log.Printf("aff : %v", aff)
	return err
}