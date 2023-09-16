package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func GlobDirectory(dirPath string) ([]string, error) {
	defer func() {

		if r := recover(); r != nil {
			fmt.Println("Serious error: cannot find any gosass for server type:", dirPath, r)
		}

	}()

	files, err := filepath.Glob(filepath.Join("./gosass/", dirPath, "*.sh"))
	if err != nil {
		return nil, err
	}

	return files, nil
}

func Getgosass(gosassList string) []string {
	var gosass []string

	if strings.Contains(gosassList, ",") {
		gosass = strings.Split(gosassList, ",")
	} else {
		gosass = append(gosass, gosassList)
	}

	scripts := []string{}

	for _, gosass := range gosass {
		scripts_found, err := GlobDirectory(gosass)
		if len(scripts_found) == 0 || err != nil {
			fmt.Println("No gosass found.", gosass, err)
			return []string{}
		}

		scripts = append(scripts, scripts_found...)
	}

	return scripts
}

func GetSharedgosass(name string) (string, error) {
	data, err := os.ReadFile(filepath.Join("./gosass/__shared/" + name + ".sh"))
	return string(data), err

}
