package main

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	apiserver "github.com/microsoft/commercial-marketplace-offer-deploy/cmd/apiserver/app"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/hosting"
)

var (
	configurationFilePath string = "."
)

func main() {
	app := apiserver.BuildApp(configurationFilePath)
	startOptions := &hosting.AppStartOptions{
		Port:      to.Ptr(8080),
		WebServer: true,
	}

	app.Start(startOptions)
}
