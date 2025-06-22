package environment

import (
	"fmt"

	"github.com/joho/godotenv"
)

type Environment string

const (
	EnvDevelopment Environment = "development"
	EnvProduction  Environment = "production"
)

func ParseEnvironment(envStr string) (Environment, error) {
	switch envStr {
	case string(EnvDevelopment):
		return EnvDevelopment, nil
	case string(EnvProduction):
		return EnvProduction, nil
	default:
		return "", fmt.Errorf("Failed to parse environment string: %s", envStr)
	}
}

func (e Environment) LoadEnv() error {
	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("Error loading .env file: %v\n", err)
	}

	envFileName := fmt.Sprintf(".env.%s", e)

	err = godotenv.Load(envFileName)
	if err != nil {
		return fmt.Errorf("Error loading %s file: %v\n", envFileName, err)
	}

	return nil
}
