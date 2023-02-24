package model

import "handler/function/internal/application/model"

type CreaterCharacterRequest struct {
	Name      string `json:"name"`
	Height    string `json:"height"`
	Mass      string `json:"mass"`
	HairColor string `json:"hair_color"`
	SkinColor string `json:"skin_color"`
	EyeColor  string `json:"eye_color"`
	BirthYear string `json:"birth_year"`
	Gender    string `json:"gender"`
	Homeworld string `json:"homeworld"`
	Created   string `json:"created"`
	Edited    string `json:"edited"`
	Url       string `json:"url"`
}

func (r CreaterCharacterRequest) ToDomain() *model.Character {
	return &model.Character{
		Name:      r.Name,
		Height:    r.Height,
		Mass:      r.Mass,
		HairColor: r.HairColor,
		SkinColor: r.SkinColor,
		EyeColor:  r.EyeColor,
		BirthYear: r.BirthYear,
		Gender:    r.Gender,
		Homeworld: r.Homeworld,
		Created:   r.Created,
		Edited:    r.Edited,
		Url:       r.Url,
	}
}

type CreateCharacterResponse struct {
	Id int `json:"id"`
}

type FindCharacterRequest struct {
	Id        int    `json:"Id"`
	Name      string `json:"name"`
	Height    string `json:"height"`
	Mass      string `json:"mass"`
	HairColor string `json:"hair_color"`
	SkinColor string `json:"skin_color"`
	EyeColor  string `json:"eye_color"`
	BirthYear string `json:"birth_year"`
	Gender    string `json:"gender"`
	Homeworld string `json:"homeworld"`
	Created   string `json:"created"`
	Edited    string `json:"edited"`
	Url       string `json:"url"`
}

func CreateResponseFromDomain(c *model.CharacterIdentifier) *CreateCharacterResponse {
	return &CreateCharacterResponse{
		Id: c.Id,
	}
}

func FindResponseFromDomain(c *model.CharacterDetail) *FindCharacterRequest {
	return &FindCharacterRequest{
		Id:        c.Id.Id,
		Name:      c.Character.Name,
		Height:    c.Character.Height,
		Mass:      c.Character.Mass,
		HairColor: c.Character.HairColor,
		SkinColor: c.Character.SkinColor,
		EyeColor:  c.Character.EyeColor,
		BirthYear: c.Character.BirthYear,
		Gender:    c.Character.Gender,
		Homeworld: c.Character.Homeworld,
		Created:   c.Character.Created,
		Edited:    c.Character.Edited,
		Url:       c.Character.Url,
	}
}
