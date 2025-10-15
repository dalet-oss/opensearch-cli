package print

import (
	"bytes"
	"encoding/json"
	"github.com/opensearch-project/opensearch-go/v4"
	"io"
	"log"
)

func MarshalJSONOrDie[T any](marshallable T) []byte {
	jsonBytes, err := json.MarshalIndent(marshallable, "", "    ")
	if err != nil {
		log.Fatal(err)
	}
	return jsonBytes
}

func RawResponse(r *opensearch.Response) {
	bodyBytes, readErr := io.ReadAll(r.Body)
	if readErr != nil {
		log.Fatalf("fail to read response body:%v", readErr)
	} else {
		log.Printf("resp code: %d\n", r.StatusCode)
		var prettyJSON bytes.Buffer
		if err := json.Indent(&prettyJSON, bodyBytes, "", "    "); err != nil {
			log.Printf("response:\n%s", bodyBytes)
		} else {
			log.Printf("response:\n%s\n", prettyJSON.String())
		}
	}

}
