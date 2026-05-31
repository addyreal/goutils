package cmdl

import (
	"errors"
	"os"
	"strconv"
)

func handleMissingInput(def string, req bool) (string, bool) {
	if def != "" {
		return def, true
	}

	if req {
		return "", false
	}

	return "", true
}

type EnvironmentInput struct {
	Key      string
	Default  string
	Required bool
	Int      bool
}

var environmentMap map[string]string = make(map[string]string)

func ParseEnv(i []EnvironmentInput) error {
	if len(environmentMap) != 0 {
		return errors.New("Environment already parsed")
	}

	for _, input := range i {
		value := os.Getenv(input.Key)

		if value == "" {
			def, ok := handleMissingInput(input.Default, input.Required)
			if !ok {
				return errors.New("Environment variable " + input.Key + " not provided")
			}

			value = def
		}

		if input.Int {
			_, err := strconv.Atoi(value)
			if err != nil {
				return errors.New("Environment variable " + input.Key + " failed to parse as int")
			}
		}

		environmentMap[input.Key] = value
	}

	return nil
}

func GetEnv(s string) string {
	val, ok := environmentMap[s]
	if !ok {
		return ""
	}

	return val
}

func GetEnvInt(s string) int {
	val, ok := environmentMap[s]
	if !ok {
		return 0
	}

	res, err := strconv.Atoi(val)
	if err != nil {
		panic(err)
	}

	return res
}

type FlagInput struct {
	Key      string
	Default  string
	Optional bool
	Int      bool
}

var flagMap map[string]string = make(map[string]string)

func ParseFlags(i []FlagInput) error {
	if len(flagMap) != 0 {
		return errors.New("Flags already parsed")
	}

	args := os.Args[1:]
	for _, input := range i {
		var value string
		for i := 0; i+1 < len(args)-1; i += 2 {
			if args[i] == input.Key {
				value = args[i+1]
			}
		}

		if value == "" {
			def, ok := handleMissingInput(input.Default, !input.Optional)
			if !ok {
				return errors.New("Flag " + input.Key + " not provided")
			}

			value = def
		}

		if input.Int {
			_, err := strconv.Atoi(value)
			if err != nil {
				return errors.New("Flag " + input.Key + " failed to parse as int")
			}
		}

		flagMap[input.Key] = value
	}

	return nil
}

func GetFlag(s string) string {
	val, ok := flagMap[s]
	if !ok {
		return ""
	}

	return val
}

func GetFlagInt(s string) int {
	val, ok := flagMap[s]
	if !ok {
		return 0
	}

	res, err := strconv.Atoi(val)
	if err != nil {
		panic(err)
	}

	return res
}
