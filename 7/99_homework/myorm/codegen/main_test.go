package main

import "testing"

func TestSearch(t *testing.T) {
	GenerateFile("../user/user.go")
	GenerateFile("../test/struct.go")
	GenerateFile("../test/struct2.go")
}
