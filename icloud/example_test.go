package icloud_test

import (
	"context"
	"crypto/x509"
	"log"
	"os"

	"github.com/lukasmalkmus/icloud-go/icloud"
)

func Example() {
	var (
		keyID         = os.Getenv("ICLOUD_KEY_ID")
		container     = os.Getenv("ICLOUD_CONTAINER")
		rawPrivateKey = os.Getenv("ICLOUD_PRIVATE_KEY")
	)

	// 1. Parse the private key.
	privateKey, err := x509.ParseECPrivateKey([]byte(rawPrivateKey))
	if err != nil {
		log.Fatal(err)
	}

	// 2. Create the iCloud client.
	client, err := icloud.NewClient(container, keyID, privateKey, icloud.Development)
	if err != nil {
		log.Fatal(err)
	}

	// 3. Create a record.
	if _, err = client.Records.Modify(context.Background(), icloud.Public, icloud.RecordsRequest{
		Operations: []icloud.RecordOperation{
			{
				Type: icloud.Create,
				Record: icloud.Record{
					Type: "MyRecord",
					Fields: icloud.Fields{
						{
							Name:  "MyField",
							Value: "Hello, World!",
						},
						{
							Name:  "MyOtherField",
							Value: 1000,
						},
					},
				},
			},
		},
	}); err != nil {
		log.Fatal(err)
	}
}
