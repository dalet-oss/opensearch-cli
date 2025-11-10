package prompts

import (
	"fmt"
	"github.com/manifoldco/promptui"
	"slices"
	"strings"
)
import "github.com/dalet-oss/opensearch-cli/pkg/utils/logging"

var log = logging.Logger()

var (
	OkChoice = []string{"y", "yes"}
	NoChoice = []string{"n", "no"}
)

// IsOk checks if the given input string matches any value in the predefined OkChoice slice (case-insensitive).
func IsOk(input string) bool {
	return slices.Contains(OkChoice, strings.ToLower(input))
}

// NotEmptyString validates that the provided string is not empty and returns an error if the string is empty.
func NotEmptyString(s string) error {
	if len(s) > 0 {
		return nil
	}
	return fmt.Errorf("empty string is not allowed")
}

// QuestionValidateStd validates the input string based on predefined "yes" and "no" choices and returns an error if invalid.
func QuestionValidateStd(input string) error {
	choices := slices.Concat(OkChoice, NoChoice)
	if len(input) == 0 {
		return nil
	} else {
		if slices.Contains(choices, strings.ToLower(input)) {
			return nil
		} else {
			return fmt.Errorf("invalid choice '%s',available options are:%s", input, strings.Join(choices, ""))
		}
	}
}

// QuestionPrompt displays a prompt for a yes/no question using a predefined set of acceptable choices and returns the user input.
func QuestionPrompt(label string) string {
	prompt := promptui.Prompt{
		Label:    fmt.Sprintf("%s[%s]", label, strings.Join(slices.Concat(OkChoice, NoChoice), "/")),
		Validate: QuestionValidateStd,
	}
	run, err := prompt.Run()
	if err != nil {
		log.Fatal().Msgf("Prompt failed %v\n", err)
	}
	return run
}

// SimplePrompt displays a prompt with the given label and returns the user input as a string. Logs an error and exits on failure.
func SimplePrompt(label string) string {
	prompt := promptui.Prompt{
		Label: label,
	}
	run, err := prompt.Run()
	if err != nil {
		log.Fatal().Msgf("Prompt failed %v\n", err)
	}
	return run
}

// ValidatedPrompt prompts the user with a label and validates the input using a provided validation function.
func ValidatedPrompt(label string, validate func(input string) error) string {
	prompt := promptui.Prompt{
		Label:    label,
		Validate: validate,
	}
	run, err := prompt.Run()
	if err != nil {
		log.Fatal().Msgf("Prompt failed %v\n", err)
	}
	return run
}

// SecretPrompt displays a masked input prompt to securely receive user input and returns the entered string.
func SecretPrompt(label string) string {
	prompt := promptui.Prompt{
		Label: label,
		Mask:  '*',
	}
	run, err := prompt.Run()
	if err != nil {
		log.Fatal().Msgf("Prompt failed %v\n", err)
	}
	return run
}

// SelectivePrompt displays a prompt with a list of options and returns the selected option.
func SelectivePrompt(label string, options []string) string {
	prompt := promptui.Select{
		Label: label,
		Items: options,
	}
	_, result, err := prompt.Run()
	if err != nil {
		log.Fatal().Msgf("Prompt failed %v\n", err)
	}
	return result
}
