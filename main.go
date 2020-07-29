package main

import (
	"github.com/alexflint/go-arg"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"http-go-server/swagger/models"
	"http-go-server/swagger/restapi"
	"http-go-server/swagger/restapi/operations"
	"log"
	"os"
)

type cliArgs struct {
	Port int `arg:"-p,help:port to listen to"`
}

var (
	args = &cliArgs{
		Port: 8081,
	}
)

func getHostnameHandler(params operations.GetHostnameParams) middleware.Responder {
	payload, err := os.Hostname()

	if err != nil {
		errPayload := &models.Error{
			Code: 500,
			Message: swag.String("failed to retrieve hostname"),
		}

		return operations.NewGetHostnameDefault(500).WithPayload(errPayload)
	}
	return operations.NewGetHostnameOK().WithPayload(payload)
}

func main() {
	arg.MustParse(args)

	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		log.Fatalln(err)
	}

	api := operations.NewHelloAPI(swaggerSpec)
	server := restapi.NewServer(api)
	defer server.Shutdown()

	server.Port = args.Port

	api.GetHostnameHandler = operations.GetHostnameHandlerFunc(getHostnameHandler)

	if err := server.Serve(); err != nil {
		log.Fatalln(err)
	}
}
