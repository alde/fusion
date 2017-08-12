package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadConfigFile(t *testing.T) {
	wd, _ := os.Getwd()
	conf := DefaultConfig()
	ReadConfigFile(conf, fmt.Sprintf("%s/config_test.yml", wd))

	assert.Equal(t, "mysql", conf.Database.Host)
	assert.Equal(t, int16(3306), conf.Database.Port)
}

func TestReadEnvironment(t *testing.T) {
	conf := DefaultConfig()
	os.Setenv("DATABASE_HOST", "postgres")
	os.Setenv("DATABASE_PORT", "5432")

	ReadEnvironment(conf)

	assert.Equal(t, "postgres", conf.Database.Host)
	assert.Equal(t, int16(5432), conf.Database.Port)
}
