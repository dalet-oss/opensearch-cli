package print

import (
	"bytes"
	"encoding/json"
	"github.com/opensearch-project/opensearch-go/v4"
	"io"
)
import "bitbucket.org/ooyalaflex/opensearch-cli/pkg/utils/logging"

var log = logging.Logger()

func MarshalJSONOrDie[T any](marshallable T) []byte {
	jsonBytes, err := json.MarshalIndent(marshallable, "", "    ")
	if err != nil {
		log.Fatal().Err(err)
	}
	return jsonBytes
}

func RawResponse(r *opensearch.Response) {
	bodyBytes, readErr := io.ReadAll(r.Body)
	if readErr != nil {
		log.Fatal().Msgf("fail to read response body:%v", readErr)
	} else {
		log.Info().Msgf("resp code: %d\n", r.StatusCode)
		var prettyJSON bytes.Buffer
		if err := json.Indent(&prettyJSON, bodyBytes, "", "    "); err != nil {
			log.Info().Msgf("response:\n%s", bodyBytes)
		} else {
			log.Info().Msgf("response:\n%s\n", prettyJSON.String())
		}
	}

}
