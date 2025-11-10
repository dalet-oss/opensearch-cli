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

func MarshalJSONOrDie[T any](marshallable T) []byte {
	jsonBytes, err := json.MarshalIndent(marshallable, "", "    ")
	if err != nil {
		log.Fatal().Err(err)
	}
	return jsonBytes
}

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
