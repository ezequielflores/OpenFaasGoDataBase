package out

import (
	"context"
	"handler/function/internal/application/model"
)

type StarwarCache interface {
	SaveCharacter(character *model.CharacterDetail, ctx context.Context) error
	FindCharacterById(characterIdentifier *model.CharacterIdentifier, ctx context.Context) (*model.CharacterDetail, error)
}
