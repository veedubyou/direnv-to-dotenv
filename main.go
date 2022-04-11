package main

import (
	"flag"
	"fmt"
	"github.com/atotto/clipboard"
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
	printOption := flag.Bool("print", false,
		"Prints the env statements instead of copying to the clipboard. By default the statements are copied to the clipboard, set this to true to show them in the terminal instead")
	flag.Parse()

	direnvOutput, err := runDirenvAllow()
	must(err)

	envVars, err := filterExportedEnvVars(direnvOutput)
	must(err)

	envStatements := resolveEnvKeysToStatements(envVars)

	if !*printOption {
		copyEnvOutput(envStatements)
	} else {
		printEnvOutput(envStatements)
	}
}

func copyEnvOutput(envOutput []string) {
	joinedOutput := strings.Join(envOutput, "\n")
	must(clipboard.WriteAll(joinedOutput))
	fmt.Println("Copied env vars to clipboard")
}

func printEnvOutput(envOutput []string) {
	joinedOutput := strings.Join(envOutput, "\n")
	fmt.Println("Paste in the following env into Jetbrains:")
	printHeader()
	fmt.Println(joinedOutput)
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

func resolveEnvKeysToStatements(envVarKeys []string) []string {
	envStatements := []string{}
	for _, envVarKey := range envVarKeys {
		value := mustGetEnv(envVarKey)
		envLine := fmt.Sprintf("%s=%s", envVarKey, value)
		envStatements = append(envStatements, envLine)
	}

	return envStatements
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
