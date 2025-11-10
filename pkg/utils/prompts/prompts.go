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

func IsOk(input string) bool {
	return slices.Contains(OkChoice, strings.ToLower(input))
}

func NotEmptyString(s string) error {
	if len(s) > 0 {
		return nil
	}
	return fmt.Errorf("empty string is not allowed")
}

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
