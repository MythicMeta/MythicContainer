package MythicContainer

import (
	"github.com/MythicMeta/MythicContainer/logging"
	"github.com/MythicMeta/MythicContainer/rabbitmq"
	"os"
)

type MythicServices = string

const (
	MythicServicePayload              MythicServices = "payload"
	MythicServiceLogger               MythicServices = "logger"
	MythicServiceWebhook              MythicServices = "webhook"
	MythicServiceC2                   MythicServices = "c2"
	MythicServiceTranslationContainer MythicServices = "translation"
	MythicServiceEventing             MythicServices = "eventing"
	MythicServiceAuth                 MythicServices = "auth"
)

func StartAndRunForever(services []MythicServices) {
	if len(services) == 0 {
		logging.LogError(nil, "Must supply at least one MythicService to start")
		os.Exit(0)
	}

	rabbitmq.Initialize()
	rabbitmq.StartServices(services)

	forever := make(chan bool)
	<-forever
}
