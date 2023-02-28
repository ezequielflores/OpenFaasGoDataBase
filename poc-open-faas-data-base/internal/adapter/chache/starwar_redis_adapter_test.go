package chache

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/assert"
	cacheModel "handler/function/internal/adapter/chache/model"
	"handler/function/internal/application/model"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
)

var character = model.Character{
	Name:      "Darth Ezequiel",
	Height:    "202",
	Mass:      "136",
	HairColor: "none",
	SkinColor: "white",
	EyeColor:  "yellow",
	BirthYear: "41.9BBY",
	Gender:    "male",
	Homeworld: "https://swapi.dev/api/planets/1/",
	Created:   "2014-12-10T15:18:20.704000Z",
	Edited:    "2014-12-20T21:17:50.313000Z",
	Url:       "https://swapi.dev/api/people/4/",
}

var jsonCharacter = "{\n\"Id\": 1,\n\"name\": \"Darth Ezequiel\",\n\"height\": \"202\",\n\"mass\": \"136\",\n\"hair_color\": \"none\",\n\"skin_color\": \"white\",\n\"eye_color\": \"yellow\",\n\"birth_year\": \"41.9BBY\",\n\"gender\": \"male\",\n\"homeworld\": \"https://swapi.dev/api/planets/1/\",\n\"created\": \"2014-12-10T15:18:20.704000Z\",\n\"edited\": \"2014-12-20T21:17:50.313000Z\",\n\"url\": \"https://swapi.dev/api/people/4/\"\n}"

var CharacterIdentifier = model.CharacterIdentifier{
	Id: 1,
}
var CharacterDetail = model.CharacterDetail{
	Character: &character,
	Id:        &CharacterIdentifier,
}

var CharacterCacheModel = &cacheModel.CharacterCache{
	Id:        CharacterDetail.Id.Id,
	Name:      CharacterDetail.Character.Name,
	Height:    CharacterDetail.Character.Height,
	Mass:      CharacterDetail.Character.Mass,
	HairColor: CharacterDetail.Character.HairColor,
	SkinColor: CharacterDetail.Character.SkinColor,
	EyeColor:  CharacterDetail.Character.EyeColor,
	BirthYear: CharacterDetail.Character.BirthYear,
	Gender:    CharacterDetail.Character.Gender,
	Homeworld: CharacterDetail.Character.Homeworld,
	Created:   CharacterDetail.Character.Created,
	Edited:    CharacterDetail.Character.Edited,
	Url:       CharacterDetail.Character.Url,
}

var cacheOptions = cacheModel.CacheOptions{
	Ttl: time.Duration(10),
}

func TestNewStarwarRedisAdapter(t *testing.T) {
	redisCliMock, _ := redismock.NewClientMock()

	adapter, err := NewStarwarRedisAdapter(redisCliMock, &cacheOptions)

	assert.NotNil(t, adapter)
	assert.Nil(t, err)
}

func TestStarwarRedisAdapter_SaveCharacter(t *testing.T) {
	req := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/starwar/characters",
		nil,
	)

	ctx := req.Context()

	redisCliMock, mock := redismock.NewClientMock()
	key := strconv.Itoa(CharacterDetail.Id.Id)
	characterJson, _ := json.Marshal(CharacterCacheModel)
	mock.ExpectSet(key, characterJson, cacheOptions.Ttl).SetVal("1")

	adapter, _ := NewStarwarRedisAdapter(redisCliMock, &cacheOptions)

	err := adapter.SaveCharacter(&CharacterDetail, ctx)

	assert.Nil(t, err)
}

func TestStarwarRedisAdapter_SaveCharacterError(t *testing.T) {
	req := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/starwar/characters",
		nil,
	)

	ctx := req.Context()

	redisCliMock, mock := redismock.NewClientMock()
	key := strconv.Itoa(CharacterDetail.Id.Id)
	characterJson, _ := json.Marshal(CharacterCacheModel)
	mock.ExpectSet(key, characterJson, cacheOptions.Ttl).SetErr(fmt.Errorf("Generic Error"))

	adapter, _ := NewStarwarRedisAdapter(redisCliMock, &cacheOptions)

	err := adapter.SaveCharacter(&CharacterDetail, ctx)

	assert.NotNil(t, err)
}

func TestStarwarRedisAdapter_FindCharacterById(t *testing.T) {

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/starwar/characters",
		nil,
	)

	ctx := req.Context()

	redisCliMock, mock := redismock.NewClientMock()
	key := strconv.Itoa(CharacterDetail.Id.Id)

	mock.ExpectGet(key).SetVal(jsonCharacter)

	adapter, _ := NewStarwarRedisAdapter(redisCliMock, &cacheOptions)

	characterDetail, err := adapter.FindCharacterById(&CharacterIdentifier, ctx)

	assert.Nil(t, err)
	assert.Equal(t, &CharacterDetail, characterDetail)
}

func TestStarwarRedisAdapter_FindCharacterByIdNotFound(t *testing.T) {

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/starwar/characters",
		nil,
	)

	ctx := req.Context()

	redisCliMock, mock := redismock.NewClientMock()
	key := strconv.Itoa(CharacterDetail.Id.Id)

	mock.ExpectGet(key).RedisNil()

	adapter, _ := NewStarwarRedisAdapter(redisCliMock, &cacheOptions)

	characterDetail, err := adapter.FindCharacterById(&CharacterIdentifier, ctx)

	assert.Nil(t, err)
	assert.Nil(t, characterDetail)
}

func TestStarwarRedisAdapter_FindCharacterByIdError(t *testing.T) {

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/starwar/characters",
		nil,
	)

	ctx := req.Context()

	redisCliMock, mock := redismock.NewClientMock()
	key := strconv.Itoa(CharacterDetail.Id.Id)

	mock.ExpectGet(key).SetErr(fmt.Errorf("generic Error"))

	adapter, _ := NewStarwarRedisAdapter(redisCliMock, &cacheOptions)

	characterDetail, err := adapter.FindCharacterById(&CharacterIdentifier, ctx)

	assert.NotNil(t, err)
	assert.Nil(t, characterDetail)
}
