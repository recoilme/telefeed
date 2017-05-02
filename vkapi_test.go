package main

import (
	"testing"
)

func TestPosts(t *testing.T) {
	s := WallGet("cook_good")
	if len(s) != 20 {
		t.Error("Expected 20, got ", len(s))
	}
}

func TestGroup(t *testing.T) {
	test := "cook_good"
	groups := GroupsGetById(test)

	if len(groups) != 1 {
		t.Error("Expected 1, got ", len(groups))
	}
	if groups[0].ScreenName != test {
		t.Error("Expected "+test, groups[0].ScreenName)
	}
}
