package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	controllerModel "handler/function/internal/adapter/controller/model"
	"handler/function/internal/application/model"
	"net/http"
	"net/http/httptest"
	"testing"
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
var CreateCharacterIdentifier = model.CharacterIdentifier{
	Id: 1,
}

var CharacterDetail = model.CharacterDetail{
	Character: &character,
	Id:        &CreateCharacterIdentifier,
}

type CreateCharacterControllerMock struct {
	mock.Mock
}
type FindCharacterControllerMock struct {
	mock.Mock
}

func (c *CreateCharacterControllerMock) CreateCharacter(
	character *model.Character,
	ctx context.Context) (*model.CharacterIdentifier, error) {

	args := c.Called(character, ctx)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	firstParameter := args.Get(0)
	characterDetail := firstParameter.(*model.CharacterIdentifier)

	return characterDetail, nil

}

func (f *FindCharacterControllerMock) FindCharacter(
	character *model.CharacterIdentifier,
	ctx context.Context) (*model.CharacterDetail, error) {

	args := f.Called(character, ctx)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	firstParameter := args.Get(0)
	characterDetail := firstParameter.(*model.CharacterDetail)

	return characterDetail, nil

}

func TestNewStarWarController(t *testing.T) {
	controllerMock := CreateCharacterControllerMock{}
	characterControllerMock := FindCharacterControllerMock{}

	controller := NewStarWarController(&controllerMock, &characterControllerMock)
	assert.NotNil(t, controller)

}

func TestGetRequestBody(t *testing.T) {

	request := controllerModel.CreaterCharacterRequest{
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

	marshal, _ := json.Marshal(request)

	reader := bytes.NewReader(marshal)

	newRequest := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/starwar/characters",
		reader,
	)

	requestModel, err := GetRequestBody(newRequest)

	assert.Nil(t, err)
	assert.NotNil(t, requestModel)

}

func TestGetRequestBodyError(t *testing.T) {

	data := []byte("Test data error")

	reader := bytes.NewReader(data)
	newRequest := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/starwar/characters",
		reader,
	)

	requestModel, err := GetRequestBody(newRequest)

	assert.Nil(t, requestModel)
	assert.NotNil(t, err)

}

func TestStarWarController_CreateStarWarCharacter(t *testing.T) {
	request := controllerModel.CreaterCharacterRequest{
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

	marshal, _ := json.Marshal(request)

	reader := bytes.NewReader(marshal)

	newRequest := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/starwar/characters",
		reader,
	)

	response := httptest.NewRecorder()

	createControllerMock := CreateCharacterControllerMock{}
	findCharacterControllerMock := FindCharacterControllerMock{}

	createControllerMock.On("CreateCharacter", mock.IsType(&model.Character{}), newRequest.Context()).
		Return(&CreateCharacterIdentifier, nil)

	controller := NewStarWarController(&createControllerMock, &findCharacterControllerMock)

	controller.CreateStarWarCharacter(response, newRequest)
	assert.EqualValues(t, http.StatusOK, response.Result().StatusCode)
}

func TestStarWarController_CreateStarWarCharacterErrorBody(t *testing.T) {

	data := []byte("Test data error")

	reader := bytes.NewReader(data)
	newRequest := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/starwar/characters",
		reader,
	)

	response := httptest.NewRecorder()

	createControllerMock := CreateCharacterControllerMock{}
	findCharacterControllerMock := FindCharacterControllerMock{}

	createControllerMock.On("CreateCharacter", mock.IsType(&model.Character{}), newRequest.Context()).
		Return(&CreateCharacterIdentifier, nil)

	controller := NewStarWarController(&createControllerMock, &findCharacterControllerMock)

	controller.CreateStarWarCharacter(response, newRequest)
	assert.EqualValues(t, http.StatusBadRequest, response.Result().StatusCode)
}

func TestStarWarController_CreateStarWarCharacterError(t *testing.T) {
	request := controllerModel.CreaterCharacterRequest{
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

	marshal, _ := json.Marshal(request)

	reader := bytes.NewReader(marshal)

	newRequest := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/starwar/characters",
		reader,
	)

	response := httptest.NewRecorder()

	createControllerMock := CreateCharacterControllerMock{}
	findCharacterControllerMock := FindCharacterControllerMock{}

	createControllerMock.On("CreateCharacter", mock.IsType(&model.Character{}), newRequest.Context()).
		Return(nil, fmt.Errorf("generic error"))

	controller := NewStarWarController(&createControllerMock, &findCharacterControllerMock)

	controller.CreateStarWarCharacter(response, newRequest)
	assert.EqualValues(t, http.StatusInternalServerError, response.Result().StatusCode)

}

func TestStarWarController_FindStarWarCharacterOk(t *testing.T) {

	newRequest := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/starwar/characters/1",
		nil,
	)

	response := httptest.NewRecorder()

	createControllerMock := CreateCharacterControllerMock{}
	findCharacterControllerMock := FindCharacterControllerMock{}

	findCharacterControllerMock.On("FindCharacter", mock.IsType(&model.CharacterIdentifier{}), newRequest.Context()).
		Return(&CharacterDetail, nil)

	controller := NewStarWarController(&createControllerMock, &findCharacterControllerMock)

	controller.FindStarWarCharacter(response, newRequest)
	assert.EqualValues(t, http.StatusOK, response.Result().StatusCode)

}

func TestStarWarController_FindStarWarCharacterInvalidUrl(t *testing.T) {

	newRequest := httptest.NewRequest(
		http.MethodPost,
		"/",
		nil,
	)

	response := httptest.NewRecorder()

	createControllerMock := CreateCharacterControllerMock{}
	findCharacterControllerMock := FindCharacterControllerMock{}

	findCharacterControllerMock.On("FindCharacter", mock.IsType(&model.CharacterIdentifier{}), newRequest.Context()).
		Return(&CharacterDetail, nil)

	controller := NewStarWarController(&createControllerMock, &findCharacterControllerMock)

	controller.FindStarWarCharacter(response, newRequest)
	assert.EqualValues(t, http.StatusBadRequest, response.Result().StatusCode)

}

func TestStarWarController_FindStarWarCharacterError(t *testing.T) {

	newRequest := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/starwar/characters/1",
		nil,
	)

	response := httptest.NewRecorder()

	createControllerMock := CreateCharacterControllerMock{}
	findCharacterControllerMock := FindCharacterControllerMock{}

	findCharacterControllerMock.On("FindCharacter", mock.IsType(&model.CharacterIdentifier{}), newRequest.Context()).
		Return(nil, fmt.Errorf("generic error"))

	controller := NewStarWarController(&createControllerMock, &findCharacterControllerMock)

	controller.FindStarWarCharacter(response, newRequest)
	assert.EqualValues(t, http.StatusInternalServerError, response.Result().StatusCode)

}
