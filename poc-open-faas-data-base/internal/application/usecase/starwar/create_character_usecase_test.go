package starwar

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"handler/function/internal/application/model"
	"net/http"
	"net/http/httptest"
	"testing"
)

var characterCreateIdentifier = model.CharacterIdentifier{
	Id: 1,
}

type StarwarRepositoryCreateMock struct {
	mock.Mock
}

func (s *StarwarRepositoryCreateMock) CreateCharacter(character *model.Character, ctx context.Context) (*model.CharacterIdentifier, error) {
	args := s.Called(character, ctx)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	firstParameter := args.Get(0)
	characterIdentifier := firstParameter.(*model.CharacterIdentifier)

	return characterIdentifier, nil
}

func (s *StarwarRepositoryCreateMock) FindCharacterById(character *model.CharacterIdentifier, ctx context.Context) (*model.CharacterDetail, error) {
	args := s.Called(character, ctx)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	firstParameter := args.Get(0)
	characterDetail := firstParameter.(*model.CharacterDetail)

	return characterDetail, nil
}

func TestCreateCharacter_CreateCharacter(t *testing.T) {

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/starwar/characters",
		nil,
	)

	ctx := req.Context()

	createRepositoryMock := new(StarwarRepositoryCreateMock)
	createRepositoryMock.On("CreateCharacter", mock.IsType(&model.Character{}), ctx).
		Return(&characterCreateIdentifier, nil)

	useCase := NewCreateCharacterUseCase(createRepositoryMock)

	characterIdentifier, err := useCase.CreateCharacter(&character, ctx)

	assert.Nil(t, err)
	assert.Equal(t, &characterCreateIdentifier, characterIdentifier)

}
