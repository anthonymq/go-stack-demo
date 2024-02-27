package common

import (
	"os"

	"github.com/anthonymq/go-stack-demo/logger"
)

type secretKey string

const SPOTIFY_CLIENT_KEY = secretKey("SPOTIFY_CLIENT_KEY")
const SPOTIFY_CLIENT_SECRET = secretKey("SPOTIFY_CLIENT_SECRET")
const DEEZER_CLIENT_KEY = secretKey("DEEZER_CLIENT_KEY")
const DEEZER_CLIENT_SECRET = secretKey("DEEZER_CLIENT_SECRET")

func GetSecret(envKey secretKey) string {
	envVar, ok := os.LookupEnv(string(envKey))
	if !ok {
		logger.Get().Fatal("Error loading secrets")
	}
	return envVar
}
