package creds

import (
	"bitbucket.org/ooyalaflex/opensearch-cli/pkg/consts"
	"bitbucket.org/ooyalaflex/opensearch-cli/pkg/types"
	"log"
	"strings"
)
import "encoding/base64"

// encodeCreds encodes the user credentials into a base64 encoded string.
func encodeCreds(username, password string) string {
	return strings.Join([]string{encodeb64(username), encodeb64(password)}, consts.CredSeparator)
}

// encodeb64 encodes the string into a base64 encoded string.
func encodeb64(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}

// decodeCreds decodes the base64 encoded user credentials into a user credentials object.
func decodeCreds(creds string) types.Creds {
	parts := strings.Split(creds, consts.CredSeparator)
	return types.Creds{
		Username: decodeb64(parts[0]),
		Password: decodeb64(parts[1]),
	}
}

// decodeb64 decodes the base64 encoded string into a string.
func decodeb64(s string) string {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		log.Fatal(err)
	}
	return string(data)
}
