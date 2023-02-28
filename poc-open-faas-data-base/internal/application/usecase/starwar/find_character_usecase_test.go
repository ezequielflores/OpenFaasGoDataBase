package starwar

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

var CharacterIdentifier = model.CharacterIdentifier{
	Id: 1,
}
var CharacterDetail = model.CharacterDetail{
	Character: &character,
	Id:        &CharacterIdentifier,
}

type StarwarRepositoryMock struct {
	mock.Mock
}

type StarwarCacheMock struct {
	mock.Mock
}

func (s *StarwarRepositoryMock) CreateCharacter(
	character *model.Character,
	ctx context.Context) (*model.CharacterIdentifier, error) {

	args := s.Called(character, ctx)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	firstParameter := args.Get(0)
	characterIdentifier := firstParameter.(*model.CharacterIdentifier)

	return characterIdentifier, nil

}

func (s *StarwarRepositoryMock) FindCharacterById(
	character *model.CharacterIdentifier,
	ctx context.Context) (*model.CharacterDetail, error) {

	args := s.Called(character, ctx)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	firstParameter := args.Get(0)
	characterDetail := firstParameter.(*model.CharacterDetail)

	return characterDetail, nil
}

func (s *StarwarCacheMock) SaveCharacter(character *model.CharacterDetail, ctx context.Context) error {

	args := s.Called(character, ctx)

	if args.Error(0) != nil {
		return args.Error(0)
	}

	return nil
}

func (s *StarwarCacheMock) FindCharacterById(characterIdentifier *model.CharacterIdentifier, ctx context.Context) (*model.CharacterDetail, error) {

	args := s.Called(characterIdentifier, ctx)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	firstParameter := args.Get(0)
	if firstParameter == nil {
		return nil, nil
	}
	characterDetail := firstParameter.(*model.CharacterDetail)

	return characterDetail, nil
}

func TestFindCharacterCacheOk(t *testing.T) {

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/starwar/characters/5",
		nil,
	)

	ctx := req.Context()

	cacheMock := new(StarwarCacheMock)
	cacheMock.On("FindCharacterById", mock.IsType(&model.CharacterIdentifier{}), ctx).
		Return(&CharacterDetail, nil)

	repositoryMock := new(StarwarRepositoryMock)
	repositoryMock.On("FindCharacterById", mock.IsType(&model.CharacterIdentifier{}), ctx).
		Return(&CharacterDetail, nil)

	useCase := NewFindCharacterUseCase(repositoryMock, cacheMock)
	identifierParam := &model.CharacterIdentifier{
		Id: 1,
	}
	character, err := useCase.FindCharacter(identifierParam, ctx)

	assert.Nil(t, err)
	assert.NotNil(t, character)
	assert.Equal(t, &CharacterDetail, character)

}

func TestFindCharacterCacheNotFound(t *testing.T) {

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/starwar/characters/5",
		nil,
	)

	ctx := req.Context()

	cacheMock := new(StarwarCacheMock)
	cacheMock.On("FindCharacterById", mock.IsType(&model.CharacterIdentifier{}), ctx).
		Return(nil, nil)
	cacheMock.On("SaveCharacter", mock.IsType(&model.CharacterDetail{}), ctx).
		Return(nil)

	repositoryMock := new(StarwarRepositoryMock)
	repositoryMock.On("FindCharacterById", mock.IsType(&model.CharacterIdentifier{}), ctx).
		Return(&CharacterDetail, nil)

	useCase := NewFindCharacterUseCase(repositoryMock, cacheMock)
	identifierParam := &model.CharacterIdentifier{
		Id: 1,
	}

	character, err := useCase.FindCharacter(identifierParam, ctx)

	assert.Nil(t, err)
	assert.NotNil(t, character)
	assert.Equal(t, &CharacterDetail, character)

}

func TestFindCharacterCacheError(t *testing.T) {

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/starwar/characters/5",
		nil,
	)

	ctx := req.Context()

	cacheMock := new(StarwarCacheMock)
	cacheMock.On("FindCharacterById", mock.IsType(&model.CharacterIdentifier{}), ctx).
		Return(nil, fmt.Errorf("generic error"))
	cacheMock.On("SaveCharacter", mock.IsType(&model.CharacterDetail{}), ctx).
		Return(nil)

	repositoryMock := new(StarwarRepositoryMock)
	repositoryMock.On("FindCharacterById", mock.IsType(&model.CharacterIdentifier{}), ctx).
		Return(&CharacterDetail, nil)

	useCase := NewFindCharacterUseCase(repositoryMock, cacheMock)
	identifierParam := &model.CharacterIdentifier{
		Id: 1,
	}

	character, err := useCase.FindCharacter(identifierParam, ctx)

	assert.Nil(t, character)
	assert.Equal(t, "generic error", err.Error())

}

func TestFindCharacterFindRepositoryError(t *testing.T) {

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/starwar/characters/5",
		nil,
	)

	ctx := req.Context()

	cacheMock := new(StarwarCacheMock)
	cacheMock.On("FindCharacterById", mock.IsType(&model.CharacterIdentifier{}), ctx).
		Return(nil, nil)
	cacheMock.On("SaveCharacter", mock.IsType(&model.CharacterDetail{}), ctx).
		Return(nil)

	repositoryMock := new(StarwarRepositoryMock)
	repositoryMock.On("FindCharacterById", mock.IsType(&model.CharacterIdentifier{}), ctx).
		Return(nil, fmt.Errorf("generic error"))

	useCase := NewFindCharacterUseCase(repositoryMock, cacheMock)
	identifierParam := &model.CharacterIdentifier{
		Id: 1,
	}

	character, err := useCase.FindCharacter(identifierParam, ctx)

	assert.Nil(t, character)
	assert.Equal(t, "generic error", err.Error())

}

func TestFindCharacterSaveRepositoryError(t *testing.T) {

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/starwar/characters/5",
		nil,
	)

	ctx := req.Context()

	cacheMock := new(StarwarCacheMock)
	cacheMock.On("FindCharacterById", mock.IsType(&model.CharacterIdentifier{}), ctx).
		Return(nil, nil)
	cacheMock.On("SaveCharacter", mock.IsType(&model.CharacterDetail{}), ctx).
		Return(fmt.Errorf("generic error"))

	repositoryMock := new(StarwarRepositoryMock)
	repositoryMock.On("FindCharacterById", mock.IsType(&model.CharacterIdentifier{}), ctx).
		Return(&CharacterDetail, nil)

	useCase := NewFindCharacterUseCase(repositoryMock, cacheMock)
	identifierParam := &model.CharacterIdentifier{
		Id: 1,
	}

	character, err := useCase.FindCharacter(identifierParam, ctx)

	assert.Nil(t, character)
	assert.Equal(t, "generic error", err.Error())

}
