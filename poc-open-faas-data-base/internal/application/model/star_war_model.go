package model

type Character struct {
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

type CharacterIdentifier struct {
	Id int
}

type CharacterDetail struct {
	Id        *CharacterIdentifier
	Character *Character
}
