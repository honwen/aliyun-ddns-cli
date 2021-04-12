package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSplitDomain001(t *testing.T) {
	rr, domain := splitDomain("a.example.com")

	assert.Equal(t, rr, "a")
	assert.Equal(t, domain, "example.com")
}

func TestSplitDomain002(t *testing.T) {
	rr, domain := splitDomain("example.com")

	assert.Equal(t, rr, "@")
	assert.Equal(t, domain, "example.com")
}

func TestSplitDomain003(t *testing.T) {
	rr, domain := splitDomain("*.example.com")

	assert.Equal(t, rr, "*")
	assert.Equal(t, domain, "example.com")
}

func TestSplitDomain004(t *testing.T) {
	rr, domain := splitDomain("*.a.example.com")

	assert.Equal(t, rr, "*.a")
	assert.Equal(t, domain, "example.com")
}

func TestSplitDomain005(t *testing.T) {
	rr, domain := splitDomain("*.b.a.example.com")

	assert.Equal(t, rr, "*.b.a")
	assert.Equal(t, domain, "example.com")
}
func TestSplitDomain006(t *testing.T) {
	rr, domain := splitDomain("a.example.co.kr")

	assert.Equal(t, rr, "a")
	assert.Equal(t, domain, "example.co.kr")
}

func TestSplitDomain007(t *testing.T) {
	rr, domain := splitDomain("*.a.example.co.kr")

	assert.Equal(t, rr, "*.a")
	assert.Equal(t, domain, "example.co.kr")
}

func TestSplitDomain008(t *testing.T) {
	rr, domain := splitDomain("example.co.kr")

	assert.Equal(t, rr, "@")
	assert.Equal(t, domain, "example.co.kr")
}
