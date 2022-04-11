package main

import (
	"fmt"
	"github.com/veedubyou/xerr"
	"os"
	"os/exec"
	"strings"
)

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	direnvOutput, err := runDirenvAllow()
	must(err)

	envVars, err := filterExportedEnvVars(direnvOutput)
	must(err)

	jetbrainsOutput := ""
	for _, envVarKey := range envVars {
		value := mustGetEnv(envVarKey)
		jetbrainsOutput += fmt.Sprintf("%s=%s\n", envVarKey, value)
	}

	fmt.Println("Paste in the following env into Jetbrains:")
	printHeader()
	fmt.Println(jetbrainsOutput)
	printFooter()
}

func mustGetEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic(fmt.Sprintf("Key %s doesn't exist", key))
	}

	return value
}

func filterExportedEnvVars(direnvOutput string) ([]string, error) {
	direnvOutputLines := strings.Split(direnvOutput, "\n")

	for _, line := range direnvOutputLines {
		exportPrefix := "direnv: export "
		if strings.HasPrefix(line, exportPrefix) {
			singleLineEnvVars := strings.TrimPrefix(line, exportPrefix)
			singleLineEnvVars = strings.Replace(singleLineEnvVars, "+", "", -1)

			return strings.Split(singleLineEnvVars, " "), nil
		}
	}

	return nil, xe.Field("direnv-output", direnvOutput).
		Error("Could not find direnv export line")
}

func runDirenvAllow() (string, error) {
	cmd := exec.Command("bash", "-i")
	cmd.Stdin = strings.NewReader("direnv allow\n")

	shellOutputBytes, err := cmd.CombinedOutput()
	shellOutput := string(shellOutputBytes)

	if err != nil {
		return "", xe.Field("shell-output", shellOutput).
			Error("Direnv allow command failed")
	}

	if len(shellOutput) == 0 {
		return "", xe.Error("Shell output from direnv allow is empty")
	}

	return shellOutput, nil
}

const BANNER_LINE = "**************************************"

func printHeader() {
	header := BANNER_LINE + "\n" + "Begin env output" + "\n" + BANNER_LINE
	fmt.Println(header)
}

func printFooter() {
	footer := BANNER_LINE + "\n" + "End env output" + "\n" + BANNER_LINE
	fmt.Println(footer)
}
