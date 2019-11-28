package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	conf, err := NewConfiguration("config", "..")
	assert.NoError(t, err)

	t.Log(conf)
}
