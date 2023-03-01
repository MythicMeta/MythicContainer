# MythicContainer
GoLang package for creating Mythic Payload Types, C2 Profiles, Translation Services, WebHook listeners, and Loggers

## Install
```shell
go get http://github.com/MythicMeta/MythicContainer
```

## Usage

Import the package, `github.com/MythicMeta/MythicContainer` and run the service via `MythicContainer.StartAndRunForever`. 
You can import as many Payload Types, C2 Profiles, Translation Services, WebHook listeners, and Loggers as you want as long as you call their `Initialize()` function and supply their corresponding service type when starting as shown below.
```go
package main

import (
	httpfunctions "PoseidonContainer/http/c2functions"
	"PoseidonContainer/my_logger"
	"PoseidonContainer/my_webhooks"
	mytranslatorfunctions "PoseidonContainer/no_actual_translation/translationfunctions"
	poseidonfunctions "PoseidonContainer/poseidon/agentfunctions"
	poseidontcpfunctions "PoseidonContainer/poseidon_tcp/c2functions"
	servicewrapperfunctions "PoseidonContainer/service_wrapper/agentfunctions"
	"github.com/MythicMeta/MythicContainer"
)

func main() {
	// load up the agent functions directory so all the init() functions execute
	poseidonfunctions.Initialize()
	httpfunctions.Initialize()
	poseidontcpfunctions.Initialize()
	servicewrapperfunctions.Initialize()
	mytranslatorfunctions.Initialize()
	my_webhooks.Initialize()
	my_logger.Initialize()
	// sync over definitions and listen
	MythicContainer.StartAndRunForever([]MythicContainer.MythicServices{
		MythicContainer.MythicServicePayload,
		MythicContainer.MythicServiceC2,
		MythicContainer.MythicServiceTranslationContainer,
		MythicContainer.MythicServiceWebhook,
		MythicContainer.MythicServiceLogger,
	})
}
```