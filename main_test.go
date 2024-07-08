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

func TestParseChunk(t *testing.T) {
	tx := make(chan map[string][]int, 1)
	chunk := []byte("Gent;13.0\nDeinze;-88.0\nGent;-13.0\n")

	parseBuffer(tx, chunk)

	result := <-tx

	if result["Gent"][0] != 130 || result["Gent"][1] != -130 {
		t.Error("Expected temperatures in Gent to be [130 -130], but got", result["Gent"])
	}
	if result["Deinze"][0] != -880 {
		t.Error("Expected temperatures in Deinze to be [-880], but got", result["Deinze"])
	}
}
