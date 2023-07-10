package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDetermineFileExtensionShouldCorrectlyPickTheFileExtension(t *testing.T) {
	cases := map[string]struct {
		extension string
		ok        bool
	}{
		"hello/world.jpg":             {"jpg", true},
		"hello/world.pNg":             {"png", true},
		"hello/world.gif":             {"", false},
		"'hello/world.jpg'":           {"jpg", true},
		"\"hello/world.png\"":         {"png", true},
		"\"'\"'hellp/world.jpg'\"'\"": {"jpg", true},
	}

	for path, expected := range cases {
		actualExtension, actualOk := determineFileExtension(path, []string{"jpg", "png"})

		assert.Equal(t, expected.extension, actualExtension)
		assert.Equal(t, expected.ok, actualOk)
	}
}
