package dotenv

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// Env defines Key and Value
type Env struct {
	Key   string
	Value string
}

var readFile = ioutil.ReadFile

// File sets the env variables by a given file
func File(filename string) error {
	envs, err := parseFile(filename)
	if err != nil {
		return err
	}
	err = setEnvs(envs)
	if err != nil {
		return err
	}
	return nil
}

// setEnvs takes []Env and sets all entries
func setEnvs(envs []Env) error {
	for _, e := range envs {
		err := os.Setenv(e.Key, e.Value)
		if err != nil {
			return err
		}
	}
	return nil
}

// parseFile takes a filename as input and returns the values
// as a slice.
// The file should be in the form:
// KEY1=VALUE1
// KEY2=VALUE2
// Spaces or empty lines are allowed.
// Comments begin with #
func parseFile(filename string) ([]Env, error) {
	var e []Env
	lines, err := loadLines(filename)
	if err != nil {
		return e, err
	}
	for _, l := range lines {
		env, err := parseLine(l)
		if err != nil {
			fmt.Println(err)
			continue
		}
		e = append(e, env)
	}
	return e, nil
}

func loadLines(filename string) ([]string, error) {
	var lines []string
	b, err := readFile(filename)
	if err != nil {
		return lines, err
	}
	splitLines := strings.Split(string(b), "\n")
	for _, l := range splitLines {
		if l == "" || strings.HasPrefix(l, "#") {
			continue
		}
		lines = append(lines, l)
	}
	return lines, err
}

func parseLine(line string) (Env, error) {
	var e Env
	s := strings.Split(line, "=")
	if len(s) != 2 {
		return Env{}, errors.New("No valid line!")
	}
	e.Key = strings.TrimSpace(s[0])
	e.Value = strings.TrimSpace(s[1])
	return e, nil
}
