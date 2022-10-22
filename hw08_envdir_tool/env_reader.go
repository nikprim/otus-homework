package main

import (
	"bufio"
	"os"
	"path"
	"strings"
)

type Environment map[string]EnvValue

type EnvValue struct {
	Value      string
	NeedRemove bool
}

func ReadDir(dir string) (Environment, error) {
	dirEntries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	envs := make(Environment)

	for _, dirEntry := range dirEntries {
		if dirEntry.IsDir() {
			continue
		}

		fileInfo, _ := dirEntry.Info()
		fileName := dirEntry.Name()

		if strings.Contains(dirEntry.Name(), "=") {
			continue
		}

		if fileInfo.Size() == 0 {
			envs[fileName] = EnvValue{NeedRemove: true}
			continue
		}

		file, err := os.OpenFile(path.Join(dir, "/", fileName), os.O_RDONLY, 0644)
		if err != nil {
			return nil, err
		}

		scanner := bufio.NewScanner(file)
		scanner.Scan()

		if err := scanner.Err(); err != nil {
			_ = file.Close()
			return nil, err
		}

		value := getValueFromLine(scanner.Text())
		if value == "" {
			envs[fileName] = EnvValue{NeedRemove: true}
			_ = file.Close()
			continue
		}

		envs[fileName] = EnvValue{Value: value}
		_ = file.Close()
	}

	return envs, nil
}

func getValueFromLine(line string) string {
	value := strings.ReplaceAll(line, string([]byte{0x00}), "\n")

	for strings.HasSuffix(value, " ") || strings.HasSuffix(value, "\t") {
		value = strings.TrimRight(value, "\t")
		value = strings.TrimRight(value, " ")
	}

	return value
}

func (e Environment) toStrings() []string {
	envs := make([]string, 0, len(e))

	for name, envValue := range e {
		if name == "" {
			continue
		}

		value := ""
		if !envValue.NeedRemove {
			value = envValue.Value
		}

		envs = append(envs, strings.Join([]string{name, "=", value}, ""))
	}

	return envs
}
