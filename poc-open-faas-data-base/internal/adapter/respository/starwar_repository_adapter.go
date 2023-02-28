package respository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	repositoryModel "handler/function/internal/adapter/respository/model"
	"handler/function/internal/application/model"
	"handler/function/internal/application/port/out"
	"handler/function/pkg"
	"log"
	"net/http"
)

var _ out.StarwarRepository = (*StarwarRepositoryAdapter)(nil)

var insertCharacter = `INSERT INTO starwar.character (
                               name,
                               height,
                               mass,
                               hair_color,
                               skin_color,
                               eye_color,
                               birth_year,
                               gender,
                               homewor_ld,
                               created,
                               edited,
                               url) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12) RETURNING Id;`

var selectCharacter = `SELECT  name,
                               height,
                               mass,
                               hair_color,
                               skin_color,
                               eye_color,
                               birth_year,
                               gender,
                               homewor_ld,
                               created,
                               edited,
                               url
							FROM starwar.character
							WHERE id = $1;`

type StarwarRepositoryAdapter struct {
	pool *pgxpool.Pool
}

func NewStarwarRepositoryAdapter(p *pgxpool.Pool) (*StarwarRepositoryAdapter, error) {
	return &StarwarRepositoryAdapter{pool: p}, nil
}

func (a *StarwarRepositoryAdapter) CreateCharacter(character *model.Character, ctx context.Context) (*model.CharacterIdentifier, error) {
	insertResponse := repositoryModel.CharacterIdentifier{}
	err := a.pool.QueryRow(ctx,
		insertCharacter,
		character.Name,
		character.Height,
		character.Mass,
		character.HairColor,
		character.SkinColor,
		character.EyeColor,
		character.BirthYear,
		character.Gender,
		character.Homeworld,
		character.Created,
		character.Edited,
		character.Url).
		Scan(&insertResponse.Id)

	if err != nil {
		log.Printf("Error creating a new character %s\n", err.Error())
		return nil, &pkg.GenericException{
			StatusCode: http.StatusInternalServerError,
			Msj:        fmt.Sprintf("Error creating a new character %s\n", err.Error()),
		}
	}

	log.Printf("New character Id: %d\n", insertResponse.Id)

	return &model.CharacterIdentifier{
		Id: insertResponse.Id,
	}, nil
}

func (a *StarwarRepositoryAdapter) FindCharacterById(character *model.CharacterIdentifier, ctx context.Context) (*model.CharacterDetail, error) {

	log.Printf("FindCharacterById: %d\n", character.Id)

	findResponse := repositoryModel.CharacterRepository{}
	err := a.pool.QueryRow(ctx, selectCharacter, character.Id).
		Scan(
			&findResponse.Name,
			&findResponse.Height,
			&findResponse.Mass,
			&findResponse.HairColor,
			&findResponse.SkinColor,
			&findResponse.EyeColor,
			&findResponse.BirthYear,
			&findResponse.Gender,
			&findResponse.Homeworld,
			&findResponse.Created,
			&findResponse.Edited,
			&findResponse.Url,
		)

	if err != nil {
		log.Printf("Error finding a new character: %s\n", err.Error())
		return nil, pkg.GenericException{
			StatusCode: http.StatusNotFound,
			Msj:        fmt.Sprintf("Character Not Found: %s\n", err.Error()),
		}
	}

	log.Printf("Found character: %+v\n", findResponse)

	return &model.CharacterDetail{
		Id: character,
		Character: &model.Character{
			Name:      findResponse.Name,
			Height:    findResponse.Height,
			Mass:      findResponse.Mass,
			HairColor: findResponse.HairColor,
			SkinColor: findResponse.SkinColor,
			EyeColor:  findResponse.EyeColor,
			BirthYear: findResponse.BirthYear,
			Gender:    findResponse.Gender,
			Homeworld: findResponse.Homeworld,
			Created:   findResponse.Created,
			Edited:    findResponse.Edited,
			Url:       findResponse.Url,
		},
	}, nil
}
