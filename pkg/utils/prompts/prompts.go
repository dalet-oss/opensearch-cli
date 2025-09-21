package prompts

import (
	"fmt"
	"github.com/manifoldco/promptui"
	"log"
	"slices"
	"strings"
)

var (
	OkChoice = []string{"ok", "y", "yes", "t", "true"}
	NoChoice = []string{"n", "no", "f", "false"}
)

func IsOk(input string) bool {
	return slices.Contains(OkChoice, strings.ToLower(input))
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
		Label:    label,
		Validate: QuestionValidateStd,
	}
	run, err := prompt.Run()
	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}
	return run
}

func SimplePrompt(label string) string {
	prompt := promptui.Prompt{
		Label: label,
	}
	run, err := prompt.Run()
	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
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
		log.Fatalf("Prompt failed %v\n", err)
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
		log.Fatalf("Prompt failed %v\n", err)
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
		log.Fatalf("Prompt failed %v\n", err)
	}
	return result
}
