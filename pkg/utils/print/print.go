package print

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/opensearch-project/opensearch-go/v4"
	"io"
)
import "github.com/dalet-oss/opensearch-cli/pkg/utils/logging"

var log = logging.Logger()

// MarshalJSONOrDie marshals a given object into indented JSON, logging a fatal error and exiting on failure.
func MarshalJSONOrDie[T any](marshallable T) []byte {
	jsonBytes, err := json.MarshalIndent(marshallable, "", "    ")
	if err != nil {
		log.Fatal().Err(err)
	}
	return jsonBytes
}

// RawResponse reads the raw response body from an OpenSearch response and formats it into a human-readable string.
// It includes the HTTP status code and pretty-printed JSON of the response body if parsing succeeds.
// If the raw body cannot be read or formatted, the function returns an error string with appropriate details.
func RawResponse(r *opensearch.Response) string {
	bodyBytes, readErr := io.ReadAll(r.Body)
	if readErr != nil {
		return fmt.Sprintf("fail to read response body:%v", readErr)
	} else {
		var prettyJSON bytes.Buffer
		if err := json.Indent(&prettyJSON, bodyBytes, "", "    "); err != nil {
			return fmt.Sprintf("fail to pretty print response body:%v\nbody:%s", err, string(bodyBytes))
		}
		return fmt.Sprintf("code: %d\nresponse: %s", r.StatusCode, prettyJSON.String())
	}

}
