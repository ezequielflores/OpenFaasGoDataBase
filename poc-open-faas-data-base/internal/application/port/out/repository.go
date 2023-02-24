package out

import (
	"context"
	"handler/function/internal/application/model"
)

type StarwarRepository interface {
	CreateCharacter(character *model.Character, ctx context.Context) (*model.CharacterIdentifier, error)
	FindCharacterById(character *model.CharacterIdentifier, ctx context.Context) (*model.CharacterDetail, error)
}
