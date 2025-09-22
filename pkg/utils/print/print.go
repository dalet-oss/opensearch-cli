package print

import (
	"encoding/json"
	"log"
)

func MarshalJSONOrDie[T any](marshallable T) []byte {
	bytes, err := json.MarshalIndent(marshallable, "", "    ")
	if err != nil {
		log.Fatal(err)
	}
	return bytes
}
