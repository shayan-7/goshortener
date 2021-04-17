package models

import (
	"os"
	"testing"

	"github.com/shayan-7/goshortener/internal/db"
	"github.com/stretchr/testify/assert"
)

func TestHashPassword(t *testing.T) {
	member := Member{
		Password: "secret",
	}

	err := member.HashPassword(member.Password)
	assert.NoError(t, err)

	os.Setenv("passwordHash", member.Password)
}

func TestCreateRecord(t *testing.T) {
	var memberResult Member

	err := db.InitDatabase()
	if err != nil {
		t.Error(err)
	}

	err = db.GlobalDB.AutoMigrate(&Member{})
	assert.NoError(t, err)

	member := Member{
		Username: "Test Member",
		Password: os.Getenv("passwordHash"),
	}

	err = member.CreateRecord()
	assert.NoError(t, err)

	db.GlobalDB.Where("username = ?", member.Username).Find(&memberResult)

	db.GlobalDB.Unscoped().Delete(&member)

	assert.Equal(t, "Test Member", memberResult.Username)
}

func TestCheckPassword(t *testing.T) {
	hash := os.Getenv("passwordHash")

	member := Member{
		Password: hash,
	}

	err := member.CheckPassword("secret")
	assert.NoError(t, err)
}
