package generic

import (
	"os"
	"path/filepath"
	"strings"
)
import "github.com/dalet-oss/opensearch-cli/pkg/utils/logging"

var log = logging.Logger()

func SaveFile(filename string, data string) error {
	filePath := filename
	if strings.HasPrefix(filename, "~") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		filePath = filepath.Join(homeDir, filename[2:])
	}
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		mkdirErr := os.MkdirAll(filepath.Dir(filePath), 0755)
		if mkdirErr != nil {
			return mkdirErr
		}
		_, createErr := os.Create(filePath)
		if createErr != nil {
			return createErr
		}
	}
	return os.WriteFile(filePath, []byte(data), 0644)
}

func SaveFileOrDie(filename string, data string) {
	err := SaveFile(filename, data)
	if err != nil {
		log.Fatal().Err(err)
	}

}
