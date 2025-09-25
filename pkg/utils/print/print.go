package print

import (
	"encoding/json"
	"github.com/opensearch-project/opensearch-go/v4"
	"io"
	"log"
)

func MarshalJSONOrDie[T any](marshallable T) []byte {
	bytes, err := json.MarshalIndent(marshallable, "", "    ")
	if err != nil {
		log.Fatal(err)
	}
	return bytes
}

func RawResponse(r *opensearch.Response) {
	bodyBytes, readErr := io.ReadAll(r.Body)
	if readErr != nil {
		log.Fatalf("fail to read response body:%v", readErr)
	} else {
		log.Printf("[code:%d] %s", r.StatusCode, bodyBytes)
	}

}
