package cpp

import (
	"bytes"
	"os/exec"
)

const CPP_DIRECTORY = "./cpp/"

func inputCorrect() bool {
	return false
}

func inputOverflow() bool {
	return false
}

func Run(script string) (string, error) {
	cmd := exec.Command(CPP_DIRECTORY + script)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return out.String(), nil
}
