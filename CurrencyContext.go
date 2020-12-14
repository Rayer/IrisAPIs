package IrisAPIs

import (
	"encoding/json"
	"flag"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/xormplus/xorm"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type CurrencyService interface {
	Convert(from string, to string, amount float64) (float64, error)
	GetMostRecentCurrencyDataRaw() (string, error)
	SyncToDb() error
	CurrencySyncRoutine()
	CurrencySyncWorker() (*CurrencySyncResult, error)
}

type CurrencyContext struct {
	ApiKey                 string
	Db                     *xorm.Engine
	UpdateAfterLastSuccess int
	UpdateAfterLastFail    int
	cachedBatch            *CurrencyBatch
	cachedConvert          map[string]currencyConvertCache
}

func NewCurrencyContext(apiKey string, db *DatabaseContext) CurrencyService {
	return &CurrencyContext{
		ApiKey:                 apiKey,
		Db:                     db.DbObject,
		UpdateAfterLastSuccess: 43200,
		UpdateAfterLastFail:    10800,
		cachedConvert:          make(map[string]currencyConvertCache),
	}
}

func NewTestCurrencyContext() *CurrencyContext {
	dc, err := NewTestDatabaseContext()
	if dc == nil || err != nil {
		return nil
	}
	//Load from env
	apiKey, exist := os.LookupEnv("FIXERIO_KEY")
	if !exist || apiKey == "" {
		apiKey = *flag.String("fixerio_key", "", "fixer io key")
	}

	//Fetch configuration file. It usually only exists in local test environment
	if apiKey == "" {
		apiKey = NewConfiguration().FixerIoApiKey
	}

	if apiKey == "" {
		return nil
	}

	return &CurrencyContext{
		ApiKey:                 apiKey,
		Db:                     dc.DbObject,
		UpdateAfterLastSuccess: 43200,
		UpdateAfterLastFail:    10800,
		cachedConvert:          make(map[string]currencyConvertCache),
	}
}

func NewCurrencyContextWithConfig(c *Configuration, db *DatabaseContext) CurrencyService {
	return &CurrencyContext{
		ApiKey:                 c.FixerIoApiKey,
		Db:                     db.DbObject,
		UpdateAfterLastSuccess: c.FixerIoLastFetchSuccessfulPeriod,
		UpdateAfterLastFail:    c.FixerIoLastFetchFailedPeriod,
		cachedConvert:          make(map[string]currencyConvertCache),
	}
}

type currencyConvertCache struct {
	rate  float64
	batch int64
}

type CurrencyBatch struct {
	Batch   int64     `xorm:"autoincr"`
	Exec    time.Time `xorm:"created"`
	Raw     string
	Success bool
	Host    string
}

type CurrencyEntry struct {
	Symbol string `xorm:"varchar(3)"`
	Base   string `xorm:"varchar(3)"`
	Batch  int64
	Rate   float64
}

func (c *CurrencyContext) Convert(from string, to string, amount float64) (float64, error) {
	if c.cachedBatch == nil {
		return 0, errors.New("no data found, please check network or DB")
	}

	key := from + to

	if val, ok := c.cachedConvert[key]; ok {
		if val.batch == c.cachedBatch.Batch {
			return amount * val.rate, nil
		}
	}

	//read data and refresh cache
	fromCurrency := &CurrencyEntry{
		Symbol: from,
		Batch:  c.cachedBatch.Batch,
	}
	toCurrency := &CurrencyEntry{
		Symbol: to,
		Batch:  c.cachedBatch.Batch,
	}
	_, err := c.Db.Get(fromCurrency)
	if err != nil {
		return 0, err
	}
	_, err = c.Db.Get(toCurrency)
	if err != nil {
		return 0, err
	}
	ratio := 1.0 / fromCurrency.Rate * toCurrency.Rate
	c.cachedConvert[key] = currencyConvertCache{
		rate:  ratio,
		batch: c.cachedBatch.Batch,
	}
	return amount * ratio, nil
}

func (c *CurrencyContext) GetMostRecentCurrencyDataRaw() (string, error) {
	if c.cachedBatch != nil {
		log.Debugf("Read from cache...")
		return c.cachedBatch.Raw, nil
	}

	lastSuccess := &CurrencyBatch{}
	_, err := c.Db.Where("success=?", 1).Desc("exec").Limit(1).Get(lastSuccess)
	if err != nil {
		return "", err
	}
	return lastSuccess.Raw, nil
}

func (c *CurrencyContext) SyncToDb() error {
	log.Debugf("Trying connecting to currency server....")
	resp, err := http.Get("http://data.fixer.io/api/latest?access_key=" + c.ApiKey)
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

	err = c.saveCurrencyEntries(string(raw))
	if err != nil {
		return err
	}

	return nil
}

func (c *CurrencyContext) saveCurrencyEntries(raw string) error {
	var data map[string]interface{}
	err := json.Unmarshal([]byte(raw), &data)

	hostName, err := os.Hostname()

	if err != nil {
		log.Warnf("Error while getting host name : %s", err.Error())
		err = nil
	}

	batch := &CurrencyBatch{
		Raw:     raw,
		Success: data["success"].(bool),
		Host:    hostName,
	}
	_, err = c.Db.Insert(batch)

	c.cachedBatch = batch

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

	aff, err := c.Db.Insert(entries)

	log.Printf("aff : %v", aff)
	return err
}

type CurrencySyncResult struct {
	lastSyncTime    time.Time
	lastSyncSuccess bool
}

func (c *CurrencyContext) CurrencySyncRoutine() {
	go func() {
		for {
			log.Infoln("Starting another round of CurrencySyncWorker...")
			_, err := c.CurrencySyncWorker()
			if err != nil {
				log.Warnf("CurrencySyncWorker ends with an error : %s", err.Error())
			}
		}
	}()
}

func (c *CurrencyContext) CurrencySyncWorker() (*CurrencySyncResult, error) {

	log.Printf("Database Object : %+v", c.Db)

	lastSuccess := &CurrencyBatch{}
	lastFail := &CurrencyBatch{}
	var err error

	_, err = c.Db.Where("success=?", 1).Desc("exec").Limit(1).Get(lastSuccess)
	_, err = c.Db.Where("success=?", 0).Desc("exec").Limit(1).Get(lastFail)

	if err != nil {
		return nil, err
	}

	log.Infoln("Fetched last exec time from database.")
	log.Infof("Last success : %s", lastSuccess.Exec)
	log.Infof("Last failed : %s", lastFail.Exec)

	var next time.Time

	if lastFail.Batch > lastSuccess.Batch {
		next = lastFail.Exec.Add(time.Second * time.Duration(c.UpdateAfterLastFail))
	} else {
		next = lastSuccess.Exec.Add(time.Second * time.Duration(c.UpdateAfterLastSuccess))
	}

	c.cachedBatch = lastSuccess

	invoke := next.Sub(time.Now())
	timer := time.NewTimer(invoke)
	defer timer.Stop()

	log.Infof("DB Sync will be executed after : %v", invoke)

	<-timer.C
	err = c.SyncToDb()

	log.Infof("DB Sync has been executed : %v", time.Now())
	if err != nil {
		log.Warnf("DB Sync ends with an error : %s", err)
	}

	return &CurrencySyncResult{
		lastSyncTime:    time.Now(),
		lastSyncSuccess: err != nil,
	}, nil
}
