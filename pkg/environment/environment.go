package environment

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Environment string

const (
	EnvDevelopment Environment = "development"
	EnvProduction  Environment = "production"
)

type EnvVars struct {
	GoEnv              Environment
	TursoDbUrl         string
	TursoDbAuthToken   string
	TwilioWhatsAppFrom string
	WorldTidesApiKey   string
	ServerPort         int
}

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

func (e Environment) LoadEnv() (EnvVars, error) {
	err := godotenv.Load()
	if err != nil {
		return EnvVars{}, fmt.Errorf("Error loading .env file: %v\n", err)
	}

	envFileName := fmt.Sprintf(".env.%s", e)

	err = godotenv.Load(envFileName)
	if err != nil {
		return EnvVars{}, fmt.Errorf("Error loading %s file: %v\n", envFileName, err)
	}

	missingEnvs := []string{}

	TURSO_DB_URL := os.Getenv("TURSO_DB_URL")
	if len(TURSO_DB_URL) == 0 {
		missingEnvs = append(missingEnvs, "TURSO_DB_URL")
	}

	TURSO_DB_AUTH_TOKEN := os.Getenv("TURSO_DB_AUTH_TOKEN")
	if len(TURSO_DB_AUTH_TOKEN) == 0 {
		missingEnvs = append(missingEnvs, "TURSO_DB_AUTH_TOKEN")
	}

	TWILIO_WHATSAPP_FROM := os.Getenv("TWILIO_WHATSAPP_FROM")
	if len(TWILIO_WHATSAPP_FROM) == 0 {
		missingEnvs = append(missingEnvs, "TWILIO_WHATSAPP_FROM")
	}

	WORLDTIDES_API_KEY := os.Getenv("WORLDTIDES_API_KEY")
	if len(WORLDTIDES_API_KEY) == 0 {
		missingEnvs = append(missingEnvs, "WORLDTIDES_API_KEY")
	}

	SERVER_PORT := os.Getenv("SERVER_PORT")
	serverPort, err := strconv.Atoi(SERVER_PORT)
	if err != nil {
		serverPort = 8080
	}

	if len(missingEnvs) > 0 {
		return EnvVars{}, fmt.Errorf("Failed to load env. Missing variables: %v", missingEnvs)
	}

	return EnvVars{
		GoEnv:              e,
		TursoDbUrl:         TURSO_DB_URL,
		TursoDbAuthToken:   TURSO_DB_AUTH_TOKEN,
		TwilioWhatsAppFrom: TWILIO_WHATSAPP_FROM,
		WorldTidesApiKey:   WORLDTIDES_API_KEY,
		ServerPort:         serverPort,
	}, nil
}
