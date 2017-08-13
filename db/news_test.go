package db

import (
	"testing"

	"github.com/alde/fusion/config"

	"github.com/stretchr/testify/assert"
)

func TestReadNews(t *testing.T) {
	db := New(config.DatabaseConfig{
		Name:     "fusion",
		Host:     "127.0.0.1",
		Port:     3306,
		User:     "root",
		Password: "password",
	})
	news, err := db.News(0, 1)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, 1, len(news))
}
