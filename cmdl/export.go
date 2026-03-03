package cmdl

import (
	"errors"
	"os"
	"strconv"
)

type EnvironmentInput struct {
	Key     string
	Default string
	Must    bool
	Int     bool
}

var environmentMap map[string]string = make(map[string]string)

func ParseEnv(i []EnvironmentInput) error {
	if len(environmentMap) != 0 {
		return errors.New("Environment already parsed")
	}

	for _, input := range i {
		value := os.Getenv(input.Key)

		if value == "" {
			if input.Must {
				return errors.New("Environment variable " + input.Key + " not provided")
			}

			if input.Default != "" {
				if input.Int {
					_, err := strconv.Atoi(input.Default)
					if err != nil {
						return errors.New("Environment variable default " + input.Key + " failed to parse as int")
					}
				}

				environmentMap[input.Key] = input.Default
			}

			continue
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
	return environmentMap[s]
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
	Key string
	Opt bool
	Int bool
}

var flagMap map[string]string = make(map[string]string)

func ParseFlags(i []FlagInput) error {
	if len(flagMap) != 0 {
		return errors.New("Flags already parsed")
	}

	args := os.Args[1:]
	for _, flag := range i {
		for i := 0; i < len(args)-1; i += 2 {
			if args[i] == flag.Key {
				if flag.Int {
					_, err := strconv.Atoi(args[i+1])
					if err != nil {
						return errors.New("Flag " + flag.Key + " failed to parse as int")
					}
				}

				flagMap[flag.Key] = args[i+1]
			}
		}

		_, ok := flagMap[flag.Key]
		if !flag.Opt && !ok {
			return errors.New("Flag " + flag.Key + " not provided")
		}
	}

	return nil
}

func GetFlag(s string) string {
	return flagMap[s]
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
