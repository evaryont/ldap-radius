package main

import (
	"log"

	"github.com/bronze1man/radius"
)

type radiusService struct{}

var radiusServer *radius.Server

func (p radiusService) RadiusHandle(request *radius.Packet) *radius.Packet {
	log.Printf("[auth] New connection, %s for user %s\n", request.Code.String(), request.GetUsername())
	npac := request.Reply()
	switch request.Code {
	case radius.AccessRequest:
		if checkCredentials(request.GetUsername(), request.GetPassword()) {
			npac.Code = radius.AccessAccept
			log.Printf("[auth] Is OK\n")
			return npac
		}
		log.Printf("[auth] Is incorrect, Go away!\n")
		npac.Code = radius.AccessReject
		npac.AVPs = append(npac.AVPs, radius.AVP{Type: radius.ReplyMessage, Value: []byte("Go away!")})
		return npac
	case radius.AccountingRequest:
		log.Printf("[acct] Accounting request, sending response\n")
		npac.Code = radius.AccountingResponse
		return npac
	default:
		log.Printf("[radius] Access rejected.\n")
		npac.Code = radius.AccessReject
		return npac
	}
}

func initRadius() {
	radiusServer = radius.NewServer(config.Radius.Listen, config.Radius.Secret, radiusService{})
	/* or you can convert it to a server that accept request from some hosts with different secrets
	cls := radius.NewClientList([]radius.Client{
		radius.NewClient("127.0.0.1", "secret1"),
		radius.NewClient("10.10.10.10", "secret2"),
	})
	s.WithClientList(cls)
	*/
}
