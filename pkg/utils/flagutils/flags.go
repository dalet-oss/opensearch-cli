package flagutils

import (
	"github.com/spf13/pflag"
	"log"
	"slices"
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

func GetNotEmptyStringFlag(flagSet *pflag.FlagSet, flagName string) string {
	flagValue := GetStringFlag(flagSet, flagName)
	if flagValue == "" {
		log.Fatalf("flag '--%s' is required to be non-empty value", flagName)
	}
	return flagValue
}

func GetStringFlagInSet(flagSet *pflag.FlagSet, flagName string, allowedValues []string) string {
	value, err := flagSet.GetString(flagName)
	if err != nil {
		log.Fatalf("failed to get flag %s from flagset: %v", flagName, err)
	}
	idx := slices.IndexFunc(allowedValues, func(item string) bool {
		return item == value
	})
	if idx == -1 {
		log.Fatalf("flag '--%s' is required to be one of %v", flagName, allowedValues)
	}
	return value
}

func GetStringFlagValidated(flagSet *pflag.FlagSet, flagName string, f func(s string) bool) string {
	value, err := flagSet.GetString(flagName)
	if err != nil {
		log.Fatalf("failed to get flag %s from flagset: %v", flagName, err)
	}
	if !f(value) {
		log.Fatalf("flag '--%s' is not a valid value", flagName)
	}
	return value
}
