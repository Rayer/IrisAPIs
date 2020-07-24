package IrisAPIs

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type LotteryContext struct {
}

//以下的Struct name都取自於台彩source code的命名...
//威力彩
type SuperLotto638Result struct {
	AZone []int
	AZoneSorted []int
	BZone int
	Serial string
}

type SuperLottto638Reward struct {
	Reward int
	Description string
	Title string
}

func (s *SuperLotto638Result) RewardOf(aZone []int, bZone int) (*SuperLottto638Reward, error) {

	//Check aZone is unique and correct length
	aMap := make(map[int]int)
	for i, ia := range aZone {
		aMap[ia] = i
	}

	if len(aMap) != 6 {
		return nil, errors.New("In A Zone number count should be 6 and unique!")
	}

	aZoneCount := 0
	for _, ia := range aZone {
		for _, a := range s.AZone {
			if ia == a {
				aZoneCount++
			}
		}
	}

	var reward *SuperLottto638Reward
	if bZone == s.BZone {
		switch aZoneCount {
		case 1:
			reward = &SuperLottto638Reward{
				Reward: 100,
				Description: "第1區任1個+第2區",
				Title:  "普獎",
			}
		case 2:
			reward = &SuperLottto638Reward{
				Reward: 200,
				Description: "第1區任2個+第2區",
				Title:  "捌獎",
			}
		case 3:
			reward = &SuperLottto638Reward{
				Reward: 400,
				Description: "第1區任3個+第2區",
				Title:  "柒獎",
			}
		case 4:
			reward = &SuperLottto638Reward{
				Reward: 4000,
				Description: "第1區任4個+第2區",
				Title:  "伍獎",
			}
		case 5:
			reward = &SuperLottto638Reward{
				Reward: 150000,
				Description: "第1區任5個+第2區",
				Title:  "參獎",
			}
		case 6:
			reward = &SuperLottto638Reward{
				Reward: 999999999,
				Description: "第1區6個+第2區",
				Title:  "頭獎",
			}
		}
	} else {
		switch aZoneCount {
		case 3:
			reward = &SuperLottto638Reward{
				Reward: 100,
				Description: "第1區任3個",
				Title:  "玖獎",
			}
		case 4:
			reward = &SuperLottto638Reward{
				Reward: 800,
				Description: "第1區任4個",
				Title:  "陸獎",
			}
		case 5:
			reward = &SuperLottto638Reward{
				Reward: 20000,
				Description: "第1區任5個",
				Title:  "肆獎",
			}
		case 6:
			reward = &SuperLottto638Reward{
				Reward: 999999,
				Description: "第1區任6個",
				Title:  "貳獎",
			}
		}
	}

	if reward == nil {
		reward = &SuperLottto638Reward{
			Reward:      0,
			Description: "...",
			Title:       "沒中獎，再接再厲",
		}
	}

	return reward, nil
}


//All
type LotteryResult struct {
	SuperLotto638Result *SuperLotto638Result
}

func (l *LotteryContext) Fetch() (*LotteryResult, error) {
	url := "https://www.taiwanlottery.com.tw/index_new.aspx"
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	//Parse SuperLotto638
	superLotto638Result, err := l.parseSuperLotto638(doc)

	return &LotteryResult{
		SuperLotto638Result: superLotto638Result,
	}, nil
}

func (l *LotteryContext) parseSuperLotto638(doc *goquery.Document) (*SuperLotto638Result, error) {
	iconNode := doc.Find("div#contents_logo_02")
	if len(iconNode.Nodes) != 1 {
		return nil, errors.New(fmt.Sprintf("Wrong SuperLotto638 icon node : %d", len(iconNode.Nodes)))
	}
	parent := iconNode.Parent()
	aZoneBalls := make([]int, 6)

	var err error
	for i, a := range strings.Split(parent.Find("div.ball_tx").Text(), " ")[0:6] {
		aZoneBalls[i], err = strconv.Atoi(a)
		if err != nil {
			return nil, err
		}
	}
	aZoneBallsSorted := make([]int, 6)
	for i, a := range strings.Split(parent.Find("div.ball_tx").Text(), " ")[6:12] {
		aZoneBallsSorted[i], err = strconv.Atoi(a)
		if err != nil {
			return nil, err
		}
	}
	bZoneBall, err := strconv.Atoi(strings.Trim(parent.Find("div.ball_red").Text(), " "))

	//Parse date and serial
	//109/7/23&nbsp;第109000059期
	var serial string
	dateSerialString := strings.Trim(parent.Find("div.contents_mine_tx02").Find("span.font_black15").Text(), " ")
	r := regexp.MustCompile("第(\\d+)期")
	if find := r.FindStringSubmatch(dateSerialString); len(find) > 1 {
		serial = find[1]
	}

	sl638result := &SuperLotto638Result{
		AZone:       aZoneBalls,
		AZoneSorted: aZoneBallsSorted,
		BZone:       bZoneBall,
		Serial:      serial,
	}

	return sl638result, nil
}





