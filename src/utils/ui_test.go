package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBellSkipper(t *testing.T) {
	bs := &BellSkipper{}

	n, err := bs.Write([]byte{7}) // ASCII BEL
	assert.NoError(t, err)
	assert.Equal(t, 0, n, "Bell byte should be swallowed")

	n, err = bs.Write([]byte("hello"))
	assert.NoError(t, err)
	assert.Equal(t, 5, n)
}

func TestNewPromptUISearcher(t *testing.T) {
	items := []string{"dev.conf", "prod.conf", "staging.conf"}
	searcher := NewPromptUISearcher(items)

	assert.True(t, searcher("dev", 0), "Prefix match")
	assert.True(t, searcher("PROD", 1), "Case-insensitive match")
	assert.True(t, searcher("", 2), "Empty query matches everything")
	assert.False(t, searcher("nope", 0), "Non-matching query returns false")
}
