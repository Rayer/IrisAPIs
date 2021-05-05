package IrisAPIs

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/xormplus/builder"
	"github.com/xormplus/xorm"
	"math/rand"
	"sync"
	"time"
)

type ApiKeyService interface {
	IssueApiKey(application string, useInHeader bool, useInQuery bool, issuer string, privileged bool) (string, error)
	ValidateApiKey(key string, embeddedIn ApiKeyLocation) (apiKeyRef int, privilegeLevel ApiKeyPrivilegeLevel)
	RecordActivity(path string, method string, key string, location ApiKeyLocation, ip string)
	GetAllKeys() ([]*ApiKeyDataModel, error)
	GetKeyModelById(id int) (*ApiKeyDataModel, error)
	SetExpire(keyId int, expire bool) error
	GetKeyUsageById(id int, from *time.Time, to *time.Time) ([]*ApiKeyAccess, error)
	GetKeyUsageByPath(path string, exactMatch bool, from *time.Time, to *time.Time) ([]*ApiKeyAccess, error)
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
	Expiration  *time.Time
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
	ApiKeyNotValid     ApiKeyPrivilegeLevel = -2
	ApiKeyExpired                           = -1
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

type ApiKeyCache interface {
	GetData(key string) *ApiKeyDataModel
	WriteData(key string, data *ApiKeyDataModel)
	Invalidate(key string)
}

type ApiKeyMapCache map[string]*ApiKeyDataModel

func (a *ApiKeyMapCache) GetData(key string) *ApiKeyDataModel {
	return (*a)[key]
}

func (a *ApiKeyMapCache) WriteData(key string, data *ApiKeyDataModel) {
	if a.GetData(key) == data {
		return
	}
	m := sync.Mutex{}
	m.Lock()
	defer m.Unlock()
	(*a)[key] = data
}

func (a *ApiKeyMapCache) Invalidate(key string) {
	m := sync.Mutex{}
	m.Lock()
	defer m.Unlock()
	delete(*a, key)
}

type ApiKeyContext struct {
	DB               *xorm.Engine
	Ip2NationService *IpNationContext
	cache            ApiKeyCache
}

func NewApiKeyService(DB *DatabaseContext) ApiKeyService {
	cache := make(ApiKeyMapCache)
	return &ApiKeyContext{DB: DB.DbObject, Ip2NationService: NewIpNationContext(DB), cache: &cache}
}

func (a *ApiKeyContext) IssueApiKey(application string, useInHeader bool, useInQuery bool, issuer string, privileged bool) (string, error) {
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

	_, err := db.Insert(&ApiKeyDataModel{
		Id:          nil,
		Key:         &key,
		UseInHeader: &useInHeader,
		UseInQuery:  &useInQuery,
		Application: &application,
		Issuer:      &issuer,
		Privileged:  PBool(privileged),
	})

	if err != nil {
		return "", err
	}

	return key, nil
}

func (a *ApiKeyContext) ValidateApiKey(key string, embeddedIn ApiKeyLocation) (int, ApiKeyPrivilegeLevel) {
	if key == "" {
		return -1, ApiKeyNotPresented
	}

	dataModel, err := a.GetKeyModelByKey(key)

	if err != nil {
		fmt.Println("Error : " + err.Error())
		return -1, ApiKeyNotValid
	}

	if dataModel == nil {
		return -1, ApiKeyNotValid
	}

	if dataModel.Expiration != nil {
		return -1, ApiKeyExpired
	}

	if *dataModel.Privileged {
		return *dataModel.Id, ApiKeyPrivileged
	} else {
		return *dataModel.Id, ApiKeyNormal
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

	apiKeyEntity, err := a.GetKeyModelByKey(key)
	if err != nil {
		log.Warnf("Database issue : %s", err.Error())
		return
	}

	if apiKeyEntity == nil {
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

func (a *ApiKeyContext) GetAllKeys() ([]*ApiKeyDataModel, error) {
	var ret []*ApiKeyDataModel
	err := a.DB.Find(&ret)

	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (a *ApiKeyContext) GetKeyModelById(id int) (*ApiKeyDataModel, error) {

	ret := &ApiKeyDataModel{Id: &id}
	find, err := a.DB.Get(ret)
	if err != nil {
		return nil, err
	}
	if !find {
		return nil, nil
	}

	return ret, nil
}

func (a *ApiKeyContext) GetKeyModelByKey(key string) (*ApiKeyDataModel, error) {
	if v := a.cache.GetData(key); v != nil {
		return v, nil
	}
	apiKeyEntity := &ApiKeyDataModel{
		Key: &key,
	}
	found, err := a.DB.Get(apiKeyEntity)
	if err != nil {
		log.Warnf("Database issue : %s", err.Error())
		return nil, err
	}

	if !found {
		log.Warnf("Api key not found : %s", key)
		return nil, nil
	}

	a.cache.WriteData(key, apiKeyEntity)
	return apiKeyEntity, nil

}

func (a *ApiKeyContext) GetKeyUsageById(id int, from *time.Time, to *time.Time) ([]*ApiKeyAccess, error) {
	//ret := make([]*ApiKeyAccess, 0)
	var ret []*ApiKeyAccess
	chain := a.DB.Where("api_key_ref = ?", id)
	if from != nil {
		chain.And("timestamp > ?", from)
	}
	if to != nil {
		chain.And("timestamp < ?", to)
	}
	err := chain.Find(&ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (a *ApiKeyContext) GetKeyUsageByPath(path string, exactMatch bool, from *time.Time, to *time.Time) ([]*ApiKeyAccess, error) {
	//ret := make([]*ApiKeyAccess, 0)
	var ret []*ApiKeyAccess
	var chain *xorm.Session
	if exactMatch {
		//chain = a.DB.Where("fullpath = ?", path)
		chain = a.DB.Where(builder.Eq{"fullpath": path})
	} else {
		chain = a.DB.Where(builder.Like{"fullpath", path})
	}
	if from != nil {
		chain.And("timestamp > ?", from)
	}
	if to != nil {
		chain.And("timestamp < ?", to)
	}
	err := chain.Find(&ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (a *ApiKeyContext) SetExpire(id int, setExpire bool) error {
	entity := &ApiKeyDataModel{Id: &id}
	exist, err := a.DB.Get(entity)
	if err != nil {
		return err
	}
	if !exist {
		return errors.New("bean id not found")
	}

	if setExpire {
		if entity.Expiration != nil {
			return errors.New("already expired")
		}
		now := time.Now()
		entity.Expiration = &now
	} else {
		entity.Expiration = nil
	}

	_, err = a.DB.Cols("expiration").ID(id).Update(entity)
	a.cache.Invalidate(*entity.Key)

	return err
}
