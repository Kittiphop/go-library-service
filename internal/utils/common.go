package utils

import (
	"log"
	"os"
)

func RequiredEnv(key string) string {
	env, ok := os.LookupEnv(key)
	if !ok {
		log.Fatalf("required env %s not set", key)
	}
	return env
}