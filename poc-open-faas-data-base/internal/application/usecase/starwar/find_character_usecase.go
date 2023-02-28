package starwar

import (
	"context"
	"handler/function/internal/application/model"
	"handler/function/internal/application/port/in"
	"handler/function/internal/application/port/out"
)

var _ in.FindCharacter = (*FindCharacter)(nil)

type FindCharacter struct {
	starwarRepository out.StarwarRepository
	starwarCache      out.StarwarCache
}

func NewFindCharacterUseCase(
	starwarRepository out.StarwarRepository,
	starwarCache out.StarwarCache) *FindCharacter {

	return &FindCharacter{
		starwarRepository: starwarRepository,
		starwarCache:      starwarCache,
	}
}

func (f FindCharacter) FindCharacter(characterIdentifier *model.CharacterIdentifier, ctx context.Context) (*model.CharacterDetail, error) {

	findResult, err := f.starwarCache.FindCharacterById(characterIdentifier, ctx)

	if err != nil {
		return nil, err
	}

	if findResult != nil {
		return findResult, nil
	}

	characterDetail, err := f.starwarRepository.FindCharacterById(characterIdentifier, ctx)

	if err != nil {
		return nil, err
	}

	errorSaveCache := f.starwarCache.SaveCharacter(characterDetail, ctx)
	if errorSaveCache != nil {
		return nil, errorSaveCache
	}

	return characterDetail, nil
}
