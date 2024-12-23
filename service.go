package main

import (
	"context"

	"github.com/arcward/edx12"
)

func EdiToJsonService(input []byte) any {

	rawMessage, _ := edx12.Read(input)
	message, _ := rawMessage.Message(context.Background())

	return message
}
