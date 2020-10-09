package IrisAPIs

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/xormplus/xorm"
	"math/rand"
	"time"
)

type ApiKeyService interface {
	IssueApiKey(application string, useInHeader bool, useInQuery bool) (string, error)
	ValidateApiKey(key string, embeddedIn ApiKeyLocation) ApiKeyPrivilegeLevel
	RecordActivity(path string, method string, key string, location ApiKeyLocation, ip string)
}

type ApiKeyDataModel struct {
	Id          *int `xorm:"autoincr pk"`
	Key         *string
	UseInHeader *bool
	UseInQuery  *bool
	Application *string //Should be a type
	Issuer      *string
	IssueDate   *time.Time `xorm:"created"`
	Privileged  *bool
}

type ApiKeyAccess struct {
	Id        *int `xorm:"autoincr pk"`
	ApiKeyRef *int
	Fullpath  *string
	Method    *string
	Ip        *string
	Nation    *string
	Timestamp *time.Time
}

type ApiKeyLocation int

const (
	Header      ApiKeyLocation = 1
	QueryString                = 2
)

type ApiKeyPrivilegeLevel int

const (
	ApiKeyNotValid     ApiKeyPrivilegeLevel = -1
	ApiKeyNotPresented                      = 0
	ApiKeyNormal                            = 1
	ApiKeyPrivileged                        = 2
)

func (d *ApiKeyDataModel) TableName() string {
	return "iris_api_key"
}

func (a *ApiKeyAccess) TableName() string {
	return "iris_api_key_access"
}

type ApiKeyContext struct {
	DB               *xorm.Engine
	Ip2NationService *IpNationContext
}

func NewApiKeyService(DB *DatabaseContext) ApiKeyService {
	return &ApiKeyContext{DB: DB.DbObject, Ip2NationService: NewIpNationContext(DB)}
}

func (a *ApiKeyContext) IssueApiKey(application string, useInHeader bool, useInQuery bool) (string, error) {
	db := a.DB

	var key string
	for {
		key = a.generateRandomString(24)
		//Do collision test
		count, err := db.Count(&ApiKeyDataModel{
			Key: &key,
		})

		if err != nil {
			return "", err
		}

		if count == 0 {
			break
		}
	}
	issuer := "auto"

	_, err := db.Insert(&ApiKeyDataModel{
		Id:          nil,
		Key:         &key,
		UseInHeader: &useInHeader,
		UseInQuery:  &useInQuery,
		Application: &application,
		Issuer:      &issuer,
		Privileged:  PBool(false),
	})

	if err != nil {
		return "", err
	}

	return key, nil
}

func (a *ApiKeyContext) ValidateApiKey(key string, embeddedIn ApiKeyLocation) ApiKeyPrivilegeLevel {
	if key == "" {
		return ApiKeyNotPresented
	}

	db := a.DB
	dataModel := &ApiKeyDataModel{Key: &key}
	got, err := db.Get(dataModel)

	if err != nil {
		fmt.Println("Error : " + err.Error())
		return ApiKeyNotValid
	}

	if !got {
		return ApiKeyNotValid
	}

	if *dataModel.Privileged {
		return ApiKeyPrivileged
	} else {
		return ApiKeyNormal
	}
}

func (a *ApiKeyContext) generateRandomString(length int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	s := make([]rune, length)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}

func (a *ApiKeyContext) RecordActivity(path string, method string, key string, location ApiKeyLocation, ip string) {
	db := a.DB
	apiKeyEntity := &ApiKeyDataModel{
		Key: &key,
	}
	found, err := db.Get(apiKeyEntity)
	if err != nil {
		log.Warnf("Database issue : %s", err.Error())
		return
	}

	if !found {
		log.Warnf("Api key not found : %s", key)
		return
	}

	ipNationResult, err := a.Ip2NationService.GetIPNation(ip)
	nation := "NaN"
	if err != nil {
		log.Warnf("Error while trying translating ip address to nation : %s", err.Error())
	}

	if ipNationResult != nil {
		nation = ipNationResult.IsoCode_3
	}

	now := time.Now()
	keyAccess := &ApiKeyAccess{
		Id:        nil,
		ApiKeyRef: apiKeyEntity.Id,
		Fullpath:  &path,
		Method:    &method,
		Ip:        &ip,
		Nation:    &nation,
		Timestamp: &now,
	}

	_, err = db.Insert(keyAccess)

	if err != nil {
		log.Warnf("Database issue : %s", err.Error())
		return
	}

	log.Debugf("Saved %+v", keyAccess)
	return
}
