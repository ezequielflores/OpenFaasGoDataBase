package in

import (
	"context"
	"handler/function/internal/application/model"
)

type CreateCharacter interface {
	CreateCharacter(character *model.Character, ctx context.Context) (*model.CharacterIdentifier, error)
}
