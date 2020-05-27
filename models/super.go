package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
)

// Super is used by pop to map your supers database table to your go code.
type Super struct {
	ID              uuid.UUID `json:"id" db:"id"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
	OriginalID      int       `json:"original_id" db:"original_id"`
	Name            string    `json:"name" db:"name"`
	FullName        string    `json:"full_name" db:"full_name"`
	PlaceOfBirth    string    `json:"place_of_birth" db:"place_of_birth"`
	FirstAppearance string    `json:"first_appearance" db:"first_appearance"`
	AlterEgos       string    `json:"alter_egos" db:"alter_egos"`
	Publisher       string    `json:"publisher" db:"publisher"`
	Alignment       string    `json:"alignment" db:"alignment"`
	Gender          string    `json:"gender" db:"gender"`
	Race            string    `json:"race" db:"race"`
	HeightFeet      string    `json:"height_feet" db:"height_feet"`
	HeightCm        int       `json:"height_cm" db:"height_cm"`
	WeightLb        string    `json:"weight_lb" db:"weight_lb"`
	WeightKg        int       `json:"weight_kg" db:"weight_kg"`
	EyeColor        string    `json:"eye_color" db:"eye_color"`
	HairColor       string    `json:"hair_color" db:"hair_color"`
	Occupation      string    `json:"occupation" db:"occupation"`
	Base            string    `json:"base" db:"base"`
	Image           string    `json:"image" db:"image"`
	Intelligence    int       `json:"intelligence" db:"intelligence"`
	Strength        int       `json:"strength" db:"strength"`
	Speed           int       `json:"speed" db:"speed"`
	Durability      int       `json:"durability" db:"durability"`
	Power           int       `json:"power" db:"power"`
	Combat          int       `json:"combat" db:"combat"`
}

// String is not required by pop and may be deleted
func (s Super) String() string {
	js, _ := json.Marshal(s)
	return string(js)
}

// Supers is not required by pop and may be deleted
type Supers []Super

// String is not required by pop and may be deleted
func (s Supers) String() string {
	js, _ := json.Marshal(s)
	return string(js)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (s *Super) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (s *Super) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (s *Super) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
