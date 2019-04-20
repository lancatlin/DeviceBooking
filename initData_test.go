package main

import (
	"testing"
)

func TestMakeIDSlice(t *testing.T) {
	t.Log("test make id slice")
	t.Log(makeIDSlice("Test", 10, 4))
	for i, v := range makeIDSlice("ST", 70, 3) {
		t.Logf("%d: %s\n", i, v)
	}
	for i, v := range makeIDSlice("T", 13, 2) {
		t.Logf("%d: %s\n", i, v)
	}
	for i, v := range makeIDSlice("CB", 32, 3) {
		t.Logf("%d: %s\n", i, v)
	}
	for i, v := range makeIDSlice("C", 90, 3) {
		t.Logf("%d: %s\n", i, v)
	}
}
