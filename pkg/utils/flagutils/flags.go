package flagutils

import (
	"github.com/spf13/pflag"
	"log"
)

func GetStringFlag(flagSet *pflag.FlagSet, flagName string) string {
	value, err := flagSet.GetString(flagName)
	if err != nil {
		log.Fatalf("failed to get flag %s from flagset: %v", flagName, err)
	}
	return value
}

func GetBoolFlag(flagSet *pflag.FlagSet, flagName string) bool {
	value, err := flagSet.GetBool(flagName)
	if err != nil {
		log.Fatalf("failed to get flag %s from flagset: %v", flagName, err)

	}
	return value
}
