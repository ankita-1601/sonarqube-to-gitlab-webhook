package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringInSlice(t *testing.T) {
	testSlice := []string{"foo", "bar", "test"}
	testString := "test"
	testResult := StringInSlice(testString, testSlice)
	assert.True(t, testResult)

	testSlice1 := []string{"foo", "bar", "testa"}
	testString1 := "testb"
	testResult1 := StringInSlice(testString1, testSlice1)
	assert.False(t, testResult1)
}

func TestCleanNextLinksHeader(t *testing.T) {
	next := "<https://gitlab.example.local/api/v4/projects/8/issues/8/notes?page=1&per_page=3>"
	res := CleanNextLinksHeader(next)
	assert.Contains(t, res, "gitlab.example.local/api/v4/projects/8")
}
