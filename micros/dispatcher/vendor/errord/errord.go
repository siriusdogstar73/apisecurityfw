package errord

import (
	"log"

	"constants"

	"context"

	"github.com/rsocket/rsocket-go"
	"github.com/rsocket/rsocket-go/payload"
)

func GetErrorPayload(sUuid string) string {
	// Connect to rsocket server
	cli, err := rsocket.Connect().
		SetupPayload(payload.NewString("Hello", "World")).
		Transport(
			rsocket.TCPClient().SetHostAndPort(
				constants.DOCKER_GW_ERROR_IP, 7883).Build()).
		Start(context.Background())
	if err != nil {
		log.Println(err)
	}

	// Simple Request Response.
	// Send request
	result, err := cli.RequestResponse(
		payload.NewString(
			sUuid, "Metadata")).Block(context.Background())
	if err != nil {
		log.Println(err)
	}
	if err != nil {
		log.Println(err)
	}

	return result.DataUTF8()
}

func SetErrorPayloadAudit(sUuid string) {
	// Connect to rsocket server
	cli, err := rsocket.Connect().
		SetupPayload(payload.NewString("Hello", "World")).
		Transport(rsocket.TCPClient().SetHostAndPort(constants.DOCKER_GW_ERROR_IP, 7883).Build()).
		Start(context.Background())
	if err != nil {
		log.Println(err)
	}

	// Simple FireAndForget.
	cli.FireAndForget(payload.NewString(sUuid, "Metadata"))
	if err != nil {
		log.Println(err)
	}
}
