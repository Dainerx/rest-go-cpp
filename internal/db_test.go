package internal

import "testing"

func TestInitDB(t *testing.T) {
	err := InitDB()
	if err != nil {
		t.Fatalf("InitDB() failed: %v", err)
	}
}
