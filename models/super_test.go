package models

import "fmt"

//Struct de parametros
type params struct {
	OriginalID      int
	Name            string
	FullName        string
	PlaceOfBirth    string
	FirstAppearance string
	AlterEgos       string
	Publisher       string
	Alignment       string
	Gender          string
	Race            string
	HeightFeet      string
	HeightCm        int
	WeightLb        string
	WeightKg        int
	EyeColor        string
	HairColor       string
	Occupation      string
	Base            string
	Image           string
	Intelligence    int
	Strength        int
	Speed           int
	Durability      int
	Power           int
	Combat          int
}

//Parametros para teste
var param1 params = params{
	1,
	"Little Super",
	"Little Super, The First",
	"The Void",
	"right now",
	"She who began, He who began, Nyarlatotep",
	"Non-Existing Comics",
	"Neutral",
	"All",
	"Extra Existential Being",
	"9'99",
	999,
	"900 lb",
	99999,
	"Rainbow",
	"Black",
	"Testing Value",
	"This very file",
	"https://pbs.twimg.com/media/EEWqCjLUwAAGZi_?format=png&name=240x240",
	90,
	20,
	100,
	999,
	1,
	0,
}

func (ms *ModelSuite) Test_Super() {
	//Struct de casos de teste
	tcases := []struct {
		params params
	}{
		{param1},
	}
	//Para cada caso
	for i, tcase := range tcases {
		super := &Super{}
		super.OriginalID = tcase.params.OriginalID
		super.Name = tcase.params.Name
		super.FullName = tcase.params.FullName
		super.PlaceOfBirth = tcase.params.PlaceOfBirth
		super.FirstAppearance = tcase.params.FirstAppearance
		super.AlterEgos = tcase.params.AlterEgos
		super.Publisher = tcase.params.Publisher
		super.Alignment = tcase.params.Alignment
		super.Gender = tcase.params.Gender
		super.Race = tcase.params.Race
		super.HeightFeet = tcase.params.HeightFeet
		super.HeightCm = tcase.params.HeightCm
		super.WeightLb = tcase.params.WeightLb
		super.WeightKg = tcase.params.WeightKg
		super.EyeColor = tcase.params.EyeColor
		super.HairColor = tcase.params.HairColor
		super.Occupation = tcase.params.Occupation
		super.Base = tcase.params.Base
		super.Image = tcase.params.Image
		super.Intelligence = tcase.params.Intelligence
		super.Strength = tcase.params.Strength
		super.Speed = tcase.params.Speed
		super.Durability = tcase.params.Durability
		super.Power = tcase.params.Power
		super.Combat = tcase.params.Combat
		_, err := ms.DB.ValidateAndCreate(super)
		if err != nil {
			ms.Fail(fmt.Sprintf("Erro ao tentar validar e criar super no caso de teste %d Err: %s", i, err.Error()))
		}
	}
}
