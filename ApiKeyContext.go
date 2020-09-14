package IrisAPIs

import (
	"fmt"
	"github.com/xormplus/xorm"
	"math/rand"
	"time"
)

type ApiKeyOperations interface {
	IssueApiKey(application string, useInHeader bool, useInQuery bool) (string, error)
	ValidateApiKey(key string, embeddedIn KEY_EMBEDDED_IN) bool
}

type ApiKeyDataModel struct {
	Id          int `xorm:"autoincr"`
	Key         *string
	UseInHeader *bool
	UseInQuery  *bool
	Application *string //Should be a type
	Issuer      *string
	IssueDate   time.Time `xorm:"created"`
	Valid       *bool
}

type KEY_EMBEDDED_IN int

const (
	HEADER       = 1
	QUERY_STRING = 2
	BOTH         = 3
)

func (d *ApiKeyDataModel) TableName() string {
	return "iris_api_key"
}

type ApiKeyContext struct {
	DB *xorm.Engine
}

func NewApiKeyOperations(DB DatabaseContext) ApiKeyOperations {
	return &ApiKeyContext{DB: DB.DbObject}
}

func (a *ApiKeyContext) IssueApiKey(application string, useInHeader bool, useInQuery bool) (string, error) {
	db := a.DB
	key := a.generateRandomString(16)
	issuer := "auto"
	_, err := db.Insert(&ApiKeyDataModel{
		Id:          0,
		Key:         &key,
		UseInHeader: &useInHeader,
		UseInQuery:  &useInQuery,
		Application: &application,
		Issuer:      &issuer,
	})

	if err != nil {
		return "", err
	}

	return key, nil
}

func (a *ApiKeyContext) ValidateApiKey(key string, embeddedIn KEY_EMBEDDED_IN) bool {
	db := a.DB
	got, err := db.Get(&ApiKeyDataModel{
		Key: &key,
	})

	if err != nil {
		fmt.Println("Error : " + err.Error())
		return false
	}

	return got
}

func (a *ApiKeyContext) generateRandomString(length int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	s := make([]rune, length)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}
