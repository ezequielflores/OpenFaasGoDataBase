package starwar

import (
	"context"
	"handler/function/internal/application/model"
	"handler/function/internal/application/port/in"
	"handler/function/internal/application/port/out"
)

var _ in.CreateCharacter = (*CreateCharacter)(nil)

type CreateCharacter struct {
	starwarRepository out.StarwarRepository
}

func NewCreateCharacterUseCase(starwarRepository out.StarwarRepository) *CreateCharacter {
	return &CreateCharacter{
		starwarRepository: starwarRepository,
	}
}

func (c *CreateCharacter) CreateCharacter(character *model.Character, ctx context.Context) (*model.CharacterIdentifier, error) {
	return c.starwarRepository.CreateCharacter(character, ctx)
}
