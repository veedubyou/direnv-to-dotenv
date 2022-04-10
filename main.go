package main

import (
	"bufio"
	"fmt"
	"github.com/veedubyou/xerr"
	"os"
	"strings"
)

func main() {
	fmt.Println("Please paste in your direnv output:")
	scanner := bufio.NewScanner(os.Stdin)

	direnvOutput := []string{}
	for scanner.Scan() {
		direnvOutput = append(direnvOutput, scanner.Text())
	}

	envVars, err := filterExportedEnvVars(direnvOutput)
	if err != nil {
		panic(err)
	}

	jetbrainsOutput := "Paste in the following env into Jetbrains:\n"
	for _, envVarKey := range envVars {
		value := mustGetEnv(envVarKey)
		jetbrainsOutput += fmt.Sprintf("%s=%s\n", envVarKey, value)
	}

	fmt.Println(jetbrainsOutput)
}

func mustGetEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic(fmt.Sprintf("Key %s doesn't exist", key))
	}

	return value
}

func filterExportedEnvVars(direnvOutput []string) ([]string, error) {
	for _, line := range direnvOutput {
		exportPrefix := "direnv: export "
		if strings.HasPrefix(line, exportPrefix) {
			singleLineEnvVars := strings.TrimPrefix(line, exportPrefix)
			singleLineEnvVars = strings.Replace(singleLineEnvVars, "+", "", -1)

			return strings.Split(singleLineEnvVars, " "), nil
		}
	}

	return nil, xe.Error("Could not find direnv export line")
}
