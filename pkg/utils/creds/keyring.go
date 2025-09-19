package creds

import (
	"bitbucket.org/ooyalaflex/opensearch-cli/pkg/consts"
	"bitbucket.org/ooyalaflex/opensearch-cli/pkg/types"
	"github.com/google/uuid"
	"github.com/zalando/go-keyring"
	"log"
)

// PushToKeyring stores the encoded user credentials in the system's keyring under a unique identifier and service name.
func PushToKeyring(user, passw string) string {
	id, creds := BuildKeyringPair(user, passw)
	err := keyring.Set(consts.ServiceName, id, creds)
	if err != nil {
		log.Fatalf("unable to store data in the keyring:%v", err)
	}
	return id
}

// DeleteFromKeyring deletes the encoded user credentials from the system's keyring under a unique identifier and service name.'
func DeleteFromKeyring(id string) {
	err := keyring.Delete(consts.ServiceName, id)
	if err != nil {
		log.Fatalf("unable to delete data from the keyring:%v", err)
	}
}

// PullFromKeyring retrieves the encoded user credentials from the system's keyring under a unique identifier and service name.'
func PullFromKeyring(id string) types.Creds {
	secret, err := keyring.Get(consts.ServiceName, id)
	if err != nil {
		log.Fatal(err)
	}
	return decodeCreds(secret)
}

// BuildKeyringPair builds a unique identifier and encoded user credentials pair.
func BuildKeyringPair(username, password string) (id string, keyringEntry string) {
	v7, err := uuid.NewV7()
	if err != nil {
		log.Fatal(err)
	}
	return v7.String(), encodeCreds(username, password)
}
