package IrisAPIs

import (
	"encoding/json"
	"fmt"
	"github.com/xormplus/xorm"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type PbsDataEntry struct {
	Region              *string
	Source              *string
	Area                *string
	UID                 *string `xorm:"pk, unique, 'uid'"`
	Direction           *string
	Longitude           *float64
	Latitude            *float64
	EntryTimestamp      *time.Time
	LastUpdateTimestamp *time.Time
	Information         *string
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

func FetchPbsFromServer() ([]PbsDataEntry, error) {
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
	ret := make([]PbsDataEntry, 0, 0)
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
			Region:              PString(v.Region),
			Source:              PString(v.Srcdetail),
			Area:                PString(v.AreaNm),
			UID:                 PString(v.UID),
			Direction:           PString(v.Direction),
			Longitude:           &lon,
			Latitude:            &lat,
			EntryTimestamp:      &happened,
			LastUpdateTimestamp: &modified,
			Information:         PString(v.Comment),
		}
		ret = append(ret, pbsDataEntry)
	}

	return ret, nil
}

func UpdateDatabase(e *xorm.Engine, data []PbsDataEntry) error {
	length := len(data)
	updated := 0
	skipped := 0
	inserted := 0
	for i, v := range data {
		target := &PbsDataEntry{
			UID: v.UID,
		}
		exist, err := e.Get(target)
		if err != nil {
			return err
		}

		if exist {
			//Do update and scan for modification
			if *target.Information != *v.Information {
				//Put old stuff into history
				fmt.Printf("Detected uid %s changed, writing to history and replace newer data", *target.UID)
				fmt.Printf("Information changed from : \n%s \n -----to : \n%s\n", *target.Information, *v.Information)
				fmt.Printf("Update timestamp changed from %s to %s", *target.LastUpdateTimestamp, *v.LastUpdateTimestamp)
				_, err := e.Insert(&PbsHistoryEntry{
					UID:                 v.UID,
					LastUpdateTimestamp: v.LastUpdateTimestamp,
					Information:         v.Information,
				})
				if err != nil {
					fmt.Printf("Error writing history : %s\n", err.Error())
				}
				//Update table to new data
				_, err = e.In("uid", v.UID).Update(v)
				if err != nil {
					return err
				}
				updated += 1
			} else {
				skipped += 1
			}

		} else {
			inserted += 1
			//Do insert
			_, err = e.Insert(v)
			if err != nil {
				return err
			}
		}
		fmt.Printf("Total: %d\tNow : %d\tUpdated : %d\tSkipped: %d\tInserted: %d\n", length, i, updated, skipped, inserted)
	}
	return nil
}

func PrintHistory(e *xorm.Engine) {

}
