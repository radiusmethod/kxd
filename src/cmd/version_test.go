package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersionCommandMetadata(t *testing.T) {
	assert.Equal(t, "version", versionCmd.Use)
	assert.Equal(t, []string{"v"}, versionCmd.Aliases)
	assert.NotEmpty(t, version)
}
