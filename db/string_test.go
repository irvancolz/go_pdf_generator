package db

import (
	"strings"
	"testing"
)

func GenreateError() error {
	return nil
}

func TestStrIncludes(t *testing.T) {
	origin := "lorem ipsum dol"
	GenreateError()
	t.Log(strings.Contains(origin, "ore"))
}
