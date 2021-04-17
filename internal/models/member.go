package models

import (
	"encoding/json"
	"io"

	"github.com/shayan-7/goshortener/internal/db"

	"golang.org/x/crypto/bcrypt"
)

type Member struct {
	ID       int    `json:"id" gorm:"primaryKey"`
	Username string `json:"username" gorm:"unique"`
	Password string `json:"password"`
}

func (m *Member) ToJSON(w io.Writer) error {
	enc := json.NewEncoder(w)
	return enc.Encode(m)
}

// FromJSON takes an io.Reader and converts the content of Reader to Item
// value if it's convertable
func (m *Member) FromJSON(r io.Reader) error {
	dec := json.NewDecoder(r)
	return dec.Decode(m)
}

func (Member) TableName() string {
	return "member"
}

// CreateUserRecord creates a user record in the database
func (m *Member) CreateRecord() error {
	result := db.GlobalDB.Create(&m)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// HashPassword encrypts user password
func (m *Member) HashPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}
	m.Password = string(bytes)
	return nil
}

// CheckPassword checks user password
func (m *Member) CheckPassword(providedPassword string) error {
	err := bcrypt.CompareHashAndPassword(
		[]byte(m.Password),
		[]byte(providedPassword),
	)
	if err != nil {
		return err
	}

	return nil
}
