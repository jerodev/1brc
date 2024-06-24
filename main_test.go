package main

import "testing"

func TestParseLine(t *testing.T) {
	city, temp := parseLine([]byte("Waregem;18.0"))
	if city != "Waregem" || temp != 180 {
		t.Errorf("Expected Waregem 180°, got %s %v°", city, temp)
	}

	city, temp = parseLine([]byte("Deerlijk;-9.7"))
	if city != "Deerlijk" || temp != -97 {
		t.Errorf("Expected Deerlijk -97°, got %s %v°", city, temp)
	}
}
