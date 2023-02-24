package in

import (
	"context"
	"handler/function/internal/application/model"
)

type FindCharacter interface {
	FindCharacter(character *model.CharacterIdentifier, ctx context.Context) (*model.CharacterDetail, error)
}
