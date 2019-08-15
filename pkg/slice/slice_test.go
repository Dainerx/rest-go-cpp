package slice

import (
	"strconv"
	"strings"
	"testing"
)

func TestContainString(t *testing.T) {
	s := []string{"hello", "world", "!", "welcome", "betell", "rest"}
	got := ContainString(s, "betell")
	if got != true {
		t.Errorf("ContainString(betell) = %s; want true", strconv.FormatBool(got))
	}

	got = ContainString(s, "betel")
	if got != false {
		t.Errorf("ContainString(betell) = %s; want false", strconv.FormatBool(got))
	}

	got = ContainString(s, "WORLD")
	if got != false {
		t.Errorf("ContainString(WORLD) = %s; want false", strconv.FormatBool(got))
	}

	got = ContainString(s, strings.ToLower("WORLD"))
	if got != true {
		t.Errorf("ContainString(WORLD) = %s; want false", strconv.FormatBool(got))
	}
}

func TestContainInt(t *testing.T) {
	s := []int{46552, 2544, 25343654, -66, 2, 25}

	got := ContainInt(s, 25343654)
	if got != true {
		t.Errorf("ContainInt(25343654) = %s; want true", strconv.FormatBool(got))
	}

	got = ContainInt(s, 9)
	if got != false {
		t.Errorf("ContainInt(9) = %s; want false", strconv.FormatBool(got))
	}
}
