package main

import (
	"gotest.tools/assert"
	"testing"
)

func TestGetDsnParts(t *testing.T) {
	dsn := "redis://w:wwww@localhost/2"
	parts := GetDsnParts(dsn)

	assert.Equal(t, parts.Host, "localhost")
	assert.Equal(t, parts.Port, 6379)
	assert.Equal(t, parts.Username, "w")
	assert.Equal(t, parts.Password, "wwww")
	assert.Equal(t, parts.Db, "2")
}

func TestGetDsnParts_CaseTwo(t *testing.T) {
	dsn := "redis://localhost/2"
	parts := GetDsnParts(dsn)

	assert.Equal(t, parts.Host, "localhost")
	assert.Equal(t, parts.Port, 6379)
	assert.Equal(t, parts.Username, "")
	assert.Equal(t, parts.Password, "")
	assert.Equal(t, parts.Db, "2")
}

func TestGetDsnParts_CaseThree(t *testing.T) {
	dsn := "redis://localhost:4569/2"
	parts := GetDsnParts(dsn)

	assert.Equal(t, parts.Host, "localhost")
	assert.Equal(t, parts.Port, 4569)
	assert.Equal(t, parts.Username, "")
	assert.Equal(t, parts.Password, "")
	assert.Equal(t, parts.Db, "2")
}

func TestGetDsnParts_CaseFour(t *testing.T) {
	dsn := "redis://first:second@localhost:4569/4"
	parts := GetDsnParts(dsn)

	assert.Equal(t, parts.Host, "localhost")
	assert.Equal(t, parts.Port, 4569)
	assert.Equal(t, parts.Username, "first")
	assert.Equal(t, parts.Password, "second")
	assert.Equal(t, parts.Db, "4")
}
