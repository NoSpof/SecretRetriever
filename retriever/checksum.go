package retriever

import (
	"hash/crc32"
	"log"

	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
)

func CheckChecksum(payload *secretmanagerpb.SecretPayload) {
	// Verify the data checksum.
	crc32c := crc32.MakeTable(crc32.Castagnoli)
	checksum := int64(crc32.Checksum(payload.Data, crc32c))
	if checksum != *payload.DataCrc32C {
		log.Panicln("Data corruption detected.")
	}
}
