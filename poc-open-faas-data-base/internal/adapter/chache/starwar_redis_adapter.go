package chache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	cacheModel "handler/function/internal/adapter/chache/model"
	model "handler/function/internal/application/model"
	"handler/function/internal/application/port/out"
	"handler/function/pkg"
	"log"
	"net/http"
	"strconv"
)

var _ out.StarwarCache = (*StarwarRedisAdapter)(nil)

type StarwarRedisAdapter struct {
	client       *redis.Client
	cacheOptions *cacheModel.CacheOptions
}

func NewStarwarRedisAdapter(c *redis.Client, o *cacheModel.CacheOptions) (*StarwarRedisAdapter, error) {
	return &StarwarRedisAdapter{client: c, cacheOptions: o}, nil
}

func (s StarwarRedisAdapter) SaveCharacter(character *model.CharacterDetail, ctx context.Context) error {

	key := strconv.Itoa(character.Id.Id)
	log.Printf("Storing in redis cache %+v\n", character)

	jsonCharacterModel := &cacheModel.CharacterCache{
		Id:        character.Id.Id,
		Name:      character.Character.Name,
		Height:    character.Character.Height,
		Mass:      character.Character.Mass,
		HairColor: character.Character.HairColor,
		SkinColor: character.Character.SkinColor,
		EyeColor:  character.Character.EyeColor,
		BirthYear: character.Character.BirthYear,
		Gender:    character.Character.Gender,
		Homeworld: character.Character.Homeworld,
		Created:   character.Character.Created,
		Edited:    character.Character.Edited,
		Url:       character.Character.Url,
	}

	characterJson, errJson := json.Marshal(jsonCharacterModel)

	if errJson != nil {
		log.Printf("Error parsing struc to json character %s\n", errJson.Error())
		return pkg.GenericException{
			Msj:        fmt.Sprintf("error accessing to cache: %s\n", errJson.Error()),
			StatusCode: http.StatusInternalServerError,
		}
	}

	val, err := s.client.Set(ctx, key, characterJson, s.cacheOptions.Ttl).Result()

	if err != nil {
		log.Printf("error storing in redis cache: %s\n", err.Error())
		return &pkg.GenericException{
			Msj:        fmt.Sprintf("error accessing to cache: %s", err.Error()),
			StatusCode: http.StatusInternalServerError,
		}
	}

	log.Printf("Value from redis: %s\n", val)

	return nil
}

func (s StarwarRedisAdapter) FindCharacterById(character *model.CharacterIdentifier, ctx context.Context) (*model.CharacterDetail, error) {

	log.Printf("Searching value in redis by id: %d\n", character.Id)

	key := strconv.Itoa(character.Id)

	val, err := s.client.Get(ctx, key).Result()

	switch {
	case err == redis.Nil:
		fmt.Printf("key %s does not exist\n", key)
		return nil, nil

	case err != nil:
		fmt.Println("Get failed", err)
		return nil, pkg.GenericException{
			Msj:        fmt.Sprintf("error accessing to cache: %s\n", err.Error()),
			StatusCode: http.StatusInternalServerError,
		}
	}
	jsonCharacterModel := &cacheModel.CharacterCache{}

	parsingJsonError := json.Unmarshal([]byte(val), jsonCharacterModel)

	if parsingJsonError != nil {
		log.Printf("Error parsing json response from redis: %s\n", parsingJsonError.Error())
		return nil, pkg.GenericException{
			Msj:        fmt.Sprintf("error parsing json response from redis: %s\n", parsingJsonError.Error()),
			StatusCode: http.StatusInternalServerError,
		}
	}

	return &model.CharacterDetail{
		Id: &model.CharacterIdentifier{Id: jsonCharacterModel.Id},
		Character: &model.Character{
			Name:      jsonCharacterModel.Name,
			Height:    jsonCharacterModel.Height,
			Mass:      jsonCharacterModel.Mass,
			HairColor: jsonCharacterModel.HairColor,
			SkinColor: jsonCharacterModel.SkinColor,
			EyeColor:  jsonCharacterModel.EyeColor,
			BirthYear: jsonCharacterModel.BirthYear,
			Gender:    jsonCharacterModel.Gender,
			Homeworld: jsonCharacterModel.Homeworld,
			Created:   jsonCharacterModel.Created,
			Edited:    jsonCharacterModel.Edited,
			Url:       jsonCharacterModel.Url,
		},
	}, nil
}
