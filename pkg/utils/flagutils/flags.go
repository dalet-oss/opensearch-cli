package flagutils

import (
	"github.com/dalet-oss/opensearch-cli/pkg/utils/logging"
	"github.com/spf13/pflag"
	"slices"
)

var log = logging.Logger()

// GetStringFlag returns the value of the flag with the given name.
// If the flag is not set, it will return an error.
func GetStringFlag(flagSet *pflag.FlagSet, flagName string) string {
	value, err := flagSet.GetString(flagName)
	if err != nil {
		log.Fatal().Msgf("failed to get flag %s from flagset: %v", flagName, err)
	}
	return value
}

// GetBoolFlag retrieves the boolean value of a flag from the provided FlagSet using the specified flag name.
// It logs a fatal error and exits if the flag retrieval fails.
func GetBoolFlag(flagSet *pflag.FlagSet, flagName string) bool {
	value, err := flagSet.GetBool(flagName)
	if err != nil {
		log.Fatal().Msgf("failed to get flag %s from flagset: %v", flagName, err)

	}
	return value
}

// GetNotEmptyStringFlag retrieves the value of a string flag by name, ensuring it is non-empty, or terminates the program.
func GetNotEmptyStringFlag(flagSet *pflag.FlagSet, flagName string) string {
	flagValue := GetStringFlag(flagSet, flagName)
	if flagValue == "" {
		log.Fatal().Msgf("flag '--%s' is required to be non-empty value", flagName)
	}
	return flagValue
}

// GetStringFlagInSet retrieves the value of a string flag by name, ensuring it matches one of the allowed values or terminates the program.
func GetStringFlagInSet(flagSet *pflag.FlagSet, flagName string, allowedValues []string) string {
	value, err := flagSet.GetString(flagName)
	if err != nil {
		log.Fatal().Msgf("failed to get flag %s from flagset: %v", flagName, err)
	}
	idx := slices.IndexFunc(allowedValues, func(item string) bool {
		return item == value
	})
	if idx == -1 {
		log.Fatal().Msgf("flag '--%s' is required to be one of %v", flagName, allowedValues)
	}
	return value
}
