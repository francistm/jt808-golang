package generator

import (
	"os"
	"strings"
)

func UpdateWorkingDir()error {
	cwd, err := os.Getwd()

	if err != nil {
		return err
	}

	if strings.HasSuffix(cwd, "jt808") {
		return os.Chdir("..")
	}

	return nil
}
