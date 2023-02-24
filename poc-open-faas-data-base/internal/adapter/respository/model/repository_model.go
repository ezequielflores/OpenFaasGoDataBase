package model

type CharacterIdentifier struct {
	Id int
}

type CharacterRepository struct {
	Id        int
	Name      string
	Height    string
	Mass      string
	HairColor string
	SkinColor string
	EyeColor  string
	BirthYear string
	Gender    string
	Homeworld string
	Created   string
	Edited    string
	Url       string
}
