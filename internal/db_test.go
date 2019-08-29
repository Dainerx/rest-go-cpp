package internal

import "testing"

func TestInitDB(t *testing.T) {
	err := InitDB()
	defer Db.Close()
	if err != nil {
		t.Fatalf("InitDB() failed: %v", err)
	}
}
