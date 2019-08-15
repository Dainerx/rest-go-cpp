package password

import (
	"strconv"
	"testing"
)

func TestCheckPasswordHash(t *testing.T) {
	hash, err := HashPass("Hello@World19")
	if err != nil {
		t.Error("HashPass(pass) failed.")
	}

	pass := "Hello@World19"
	got := CheckPassHash(pass, hash)
	if got != true {
		t.Errorf("CheckPassHash(WORLD) = %s; want true", strconv.FormatBool(got))
	}
}
