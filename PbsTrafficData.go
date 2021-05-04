package IrisAPIs

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/xormplus/xorm"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type PbsDataEntry struct {
	Region         *string
	Source         *string
	Area           *string
	UID            *string `xorm:"pk, unique, 'uid'"`
	Direction      *string
	Longitude      *float64
	Latitude       *float64
	EntryTimestamp *time.Time
}

func (p *PbsDataEntry) TableName() string {
	return "pbs_traffic_events"
}

type PbsHistoryEntry struct {
	UID                 *string    `xorm:"'uid'"`
	LastUpdateTimestamp *time.Time `xorm:"'update_timestamp'"`
	Information         *string
}

func (p *PbsHistoryEntry) TableName() string {
	return "pbs_traffic_history"
}

type PbsParseJsonResult struct {
	PbsDataEntry
	PbsHistoryEntry
}

type PbsTrafficDataService interface {
	FetchPbsFromServer(ctx context.Context) ([]PbsParseJsonResult, error)
	UpdateDatabase(ctx context.Context, data []PbsParseJsonResult, callback func(total int, now int, updated int, inserted int, skipped int)) error
	GetHistory(ctx context.Context, pastDuration time.Duration) (map[string][]PbsHistoryEntry, error)
	ScheduledWorker(ctx context.Context, updateRate time.Duration)
}

type PbsTrafficDataServiceImpl struct {
	engine *xorm.Engine
}

func NewPbsTrafficDataService(databaseContext *DatabaseContext) PbsTrafficDataService {
	return &PbsTrafficDataServiceImpl{engine: databaseContext.DbObject}
}

func (p *PbsTrafficDataServiceImpl) FetchPbsFromServer(ctx context.Context) ([]PbsParseJsonResult, error) {
	const (
		DataSource = "https://od.moi.gov.tw/MOI/v1/pbs"
	)
	var entries []struct {
		Region     string `json:"region"`
		Srcdetail  string `json:"srcdetail"`
		AreaNm     string `json:"areaNm"`
		UID        string `json:"UID"`
		Direction  string `json:"direction"`
		Y1         string `json:"y1"`
		HappenTime string `json:"happentime"`
		Roadtype   string `json:"roadtype"`
		Road       string `json:"road"`
		ModDttm    string `json:"modDttm"`
		Comment    string `json:"comment"`
		HappenDate string `json:"happendate"`
		X1         string `json:"x1"`
	}

	client := http.Client{
		Timeout: 3 * time.Second,
	}
	resp, err := client.Get(DataSource)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = resp.Body.Close()
	}()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &entries)
	if err != nil {
		return nil, err
	}

	//converv to data entry
	ret := make([]PbsParseJsonResult, 0, 0)
	for _, v := range entries {
		lon, err := strconv.ParseFloat(v.Y1, 64)
		if err != nil {
			lon = 0
		}
		lat, err := strconv.ParseFloat(v.X1, 64)
		if err != nil {
			lat = 0
		}
		modified, err := time.Parse("2006-01-02 15:04:05.99999999 Z07", v.ModDttm+" +08")
		if err != nil {
			modified = time.Time{}
		}
		//組出happendate / happentime
		happenedStr := fmt.Sprintf("%s %s", v.HappenDate, v.HappenTime)
		happened, err := time.Parse("2006-01-02 15:04:05.99999999 Z07", happenedStr+" +08")
		if err != nil {
			happened = time.Time{}
		}
		pbsDataEntry := PbsDataEntry{
			Region:         PString(v.Region),
			Source:         PString(v.Srcdetail),
			Area:           PString(v.AreaNm),
			UID:            PString(v.UID),
			Direction:      PString(v.Direction),
			Longitude:      &lon,
			Latitude:       &lat,
			EntryTimestamp: &happened,
		}
		pbsHIstoryEntry := PbsHistoryEntry{
			UID:                 pbsDataEntry.UID,
			LastUpdateTimestamp: &modified,
			Information:         PString(v.Comment),
		}

		ret = append(ret, PbsParseJsonResult{
			PbsDataEntry:    pbsDataEntry,
			PbsHistoryEntry: pbsHIstoryEntry,
		})
	}

	return ret, nil
}

func (p *PbsTrafficDataServiceImpl) UpdateDatabase(ctx context.Context, data []PbsParseJsonResult, progressCb func(total int, now int, updated int, inserted int, skipped int)) error {
	e := p.engine
	length := len(data)
	updated := 0
	skipped := 0
	inserted := 0
	for i, v := range data {
		target := &PbsDataEntry{
			UID: v.PbsDataEntry.UID,
		}
		exist, err := e.Get(target)
		if err != nil {
			return err
		}

		if exist {
			//Compare last timestamp (or message?) from history
			lastHistory := &PbsHistoryEntry{UID: v.PbsHistoryEntry.UID}
			exist, err := e.Desc("update_timestamp").Limit(1).Get(lastHistory)
			if !exist {
				return errors.Errorf("Database integraty error for UID : %s", *v.PbsDataEntry.UID)
			}
			if err != nil {
				return err
			}
			if *lastHistory.Information != *v.PbsHistoryEntry.Information {
				//fmt.Printf("Comparing uid %s (%s):\n", *lastHistory.UID, *lastHistory.Information)
				//fmt.Printf("Comparing uid %s (%s):\n", *v.PbsHistoryEntry.UID, *v.PbsHistoryEntry.Information)

				updated += 1
				_, err := e.Insert(v.PbsHistoryEntry)
				//fmt.Printf("\nupdated : %s\t%s (%s) \n", *v.PbsHistoryEntry.UID, *v.PbsHistoryEntry.Information, v.PbsHistoryEntry.LastUpdateTimestamp.Format(time.RFC3339))
				if err != nil {
					return err
				}
			} else {
				skipped += 1
			}

		} else {
			inserted += 1
			//Do insert
			_, err = e.Insert(v.PbsDataEntry)
			if err != nil {
				return err
			}
			_, err = e.Insert(v.PbsHistoryEntry)
			if err != nil {
				return err
			}
		}
		if progressCb != nil {
			progressCb(length, i+1, updated, inserted, skipped)
		}
	}
	return nil
}

func (p *PbsTrafficDataServiceImpl) GetHistory(ctx context.Context, pastDuration time.Duration) (map[string][]PbsHistoryEntry, error) {
	e := p.engine
	var result []PbsHistoryEntry
	err := e.Where(fmt.Sprintf("pbs_traffic_history.update_timestamp > NOW() - INTERVAL %v SECOND", pastDuration.Seconds())).Find(&result)
	if err != nil {
		return nil, err
	}

	ret := make(map[string][]PbsHistoryEntry)
	for _, v := range result {
		if _, exist := ret[*v.UID]; exist {
			ret[*v.UID] = append(ret[*v.UID], v)
		} else {
			ret[*v.UID] = []PbsHistoryEntry{v}
		}
	}

	return ret, nil
}

func (p *PbsTrafficDataServiceImpl) ScheduledWorker(ctx context.Context, updateRate time.Duration) {
	log.Infof("Starting PBS update service")
	go func() {
		timer := time.NewTicker(updateRate)
		defer timer.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-timer.C:
				log.Debug("Updating PBS from server")
				res, _ := p.FetchPbsFromServer(ctx)
				_ = p.UpdateDatabase(ctx, res, func(total int, now int, updated int, inserted int, skipped int) {
					if total == now {
						log.Debugf("Processed %d records, %d updated, %d inserted and %d skipped", now, updated, inserted, skipped)
					}
				})
			}
		}
	}()
}
