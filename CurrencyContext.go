package IrisAPIs

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"time"
)

type CurrencyBatch struct {
	Batch   int64     `xorm:"autoincr"`
	Exec    time.Time `xorm:"created"`
	Raw     string
	Success bool
}

type CurrencyEntry struct {
	Symbol string `xorm:"varchar(3)"`
	Base   string `xorm:"varchar(3)"`
	Batch  int64
	Rate   float64
}

func GetMostRecentCurrencyDataRaw() (string, error) {
	db := GetDatabaseContext().DbObject
	lastSuccess := &CurrencyBatch{}
	_, err := db.Where("success=?", 1).Desc("exec").Limit(1).Get(lastSuccess)
	if err != nil {
		return "", err
	}
	return lastSuccess.Raw, nil
}

func SyncToDb() error {
	log.Debugf("Trying connecting to currency server.....")
	resp, err := http.Get("http://data.fixer.io/api/latest?access_key=676ac77e5ce5d4b9a57ee6464ff84433")
	if err != nil {
		return err
	}
	defer func() {
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
	var data map[string]interface{}
	err := json.Unmarshal([]byte(raw), &data)

	batch := &CurrencyBatch{
		Raw:     raw,
		Success: data["success"].(bool),
	}
	_, err = engine.Insert(batch)

	var entries = make([]*CurrencyEntry, 0, 0)

	if err != nil {
		return err
	}

	base := data["base"].(string)
	rates := data["rates"].(map[string]interface{})

	for sym, rate := range rates {
		entries = append(entries, &CurrencyEntry{
			Symbol: sym,
			Base:   base,
			Batch:  batch.Batch,
			Rate:   rate.(float64),
		})
	}

	aff, err := engine.Insert(entries)

	log.Printf("aff : %v", aff)
	return err
}

type CurrencySyncResult struct {
	lastSyncTime    time.Time
	lastSyncSuccess bool
}

func CurrencySyncRoutine() {
	go func() {
		for {
			log.Infoln("Starting another round of CurrencySyncWorker...")
			_, err := CurrencySyncWorker()
			if err != nil {
				log.Warnf("CurrencySyncWorker ends with an error : %s", err.Error())
			}
		}
	}()
}

func CurrencySyncWorker() (*CurrencySyncResult, error) {

	db := GetDatabaseContext().DbObject

	log.Printf("Database Object : %+v", db)

	lastSuccess := &CurrencyBatch{}
	lastFail := &CurrencyBatch{}
	var err error

	_, err = db.Where("success=?", 1).Desc("exec").Limit(1).Get(lastSuccess)
	_, err = db.Where("success=?", 0).Desc("exec").Limit(1).Get(lastFail)

	if err != nil {
		return nil, err
	}

	log.Infoln("Fetched last exec time from database.")
	log.Infof("Last success : %s", lastSuccess.Exec)
	log.Infof("Last failed : %s", lastFail.Exec)

	/*
		Rules :
		1. Exec every 12hr
		2. If last one is success, exec + 12hr
		3. If last one is failed, exec + 3hr
	*/

	var next time.Time

	if lastFail.Batch > lastSuccess.Batch {
		next = lastFail.Exec.Add(time.Hour * 3)
	} else {
		next = lastSuccess.Exec.Add(time.Hour * 12)
	}

	invoke := next.Sub(time.Now())
	timer := time.NewTimer(invoke)
	defer timer.Stop()

	log.Infof("DB Sync will be executed after : %v", invoke)

	<-timer.C
	err = SyncToDb()

	log.Infof("DB Sync has been executed : %v", time.Now())
	if err != nil {
		log.Warnf("DB Sync ends with an error : %s", err)
	}

	return &CurrencySyncResult{
		lastSyncTime:    time.Now(),
		lastSyncSuccess: err != nil,
	}, nil
}
