package model

type ControllerCharacter struct {
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
