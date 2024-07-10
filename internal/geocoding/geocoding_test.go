package geocoding

import (
	"testing"
)

func TestWrongAPIKey(t *testing.T) {
	_, err := GetCoordinates("Москва", "")

	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestWrongCityName(t *testing.T) {
	c, _ := GetCoordinates("Addsds", "925d1cb191ea87f8275e56f301cf1f9d")
	if c != nil {
		t.Errorf("Expected nil, got %v", c)
	}
}

func TestCorrectCityName(t *testing.T) {
	c, _ := GetCoordinates("Москва", "925d1cb191ea87f8275e56f301cf1f9d")
	if c == nil {
		t.Errorf("Expected not nil, got nil")
	}
}
