package eventingstructs

import (
	"fmt"
	"github.com/MythicMeta/MythicContainer/utils/helpers"
	"github.com/MythicMeta/MythicContainer/utils/sharedStructs"
	"sync"
)

type CustomFunctionDefinition struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	// SemVer is a specific semantic version tracker you can use for your payload type
	SemVer   string                                                                `json:"semver"`
	Function func(input NewCustomEventingMessage) NewCustomEventingMessageResponse `json:"-"`
}
type ConditionalCheckDefinition struct {
	Name        string                                                                              `json:"name"`
	Description string                                                                              `json:"description"`
	Function    func(input ConditionalCheckEventingMessage) ConditionalCheckEventingMessageResponse `json:"-"`
}

type EventingDefinition struct {
	Name                      string `json:"name"`
	Description               string `json:"description"`
	CustomFunctions           []CustomFunctionDefinition
	ConditionalChecks         []ConditionalCheckDefinition
	TaskInterceptFunction     func(input TaskInterceptMessage) TaskInterceptMessageResponse                             `json:"-"`
	ResponseInterceptFunction func(input ResponseInterceptMessage) ResponseInterceptMessageResponse                     `json:"-"`
	Subscriptions             []string                                                                                  `json:"subscriptions"`
	OnContainerStartFunction  func(sharedStructs.ContainerOnStartMessage) sharedStructs.ContainerOnStartMessageResponse `json:"-"`
}

// REQUIRED, Don't Modify
type allEventingData struct {
	mutex              sync.RWMutex
	rpcMethods         []sharedStructs.RabbitmqRPCMethod
	directMethods      []sharedStructs.RabbitmqDirectMethod
	eventingDefinition EventingDefinition
}

var (
	AllEventingData containerEventingData
)

type containerEventingData struct {
	EventingMap map[string]*allEventingData
}

func (r *containerEventingData) GetAllNames() []string {
	names := []string{}
	for key, _ := range r.EventingMap {
		if key != "" && !helpers.StringSliceContains(names, key) {
			names = append(names, key)
		}
	}
	return names
}
func (r *containerEventingData) Get(name string) *allEventingData {
	if r.EventingMap == nil {
		r.EventingMap = make(map[string]*allEventingData)
	}
	if existingC2Data, ok := r.EventingMap[name]; !ok {
		newC2Data := allEventingData{}
		r.EventingMap[name] = &newC2Data
		return &newC2Data
	} else {
		return existingC2Data
	}
}
func (r *allEventingData) AddEventingDefinition(def EventingDefinition) {
	r.eventingDefinition = def
}
func (r *allEventingData) GetEventingDefinition() EventingDefinition {
	return r.eventingDefinition
}
func (r *allEventingData) SetSubscriptions(subs []string) {
	r.eventingDefinition.Subscriptions = subs
}
func (r *allEventingData) SetName(name string) {
	r.eventingDefinition.Name = name
}
func (r *allEventingData) GetRoutingKey(routingKey string) string {
	return fmt.Sprintf("%s_%s", r.eventingDefinition.Name, routingKey)
}
func (r *allEventingData) AddRPCMethod(m sharedStructs.RabbitmqRPCMethod) {
	r.mutex.Lock()
	r.rpcMethods = append(r.rpcMethods, m)
	r.mutex.Unlock()
}
func (r *allEventingData) GetRPCMethods() []sharedStructs.RabbitmqRPCMethod {
	return r.rpcMethods
}
func (r *allEventingData) AddDirectMethod(m sharedStructs.RabbitmqDirectMethod) {
	r.mutex.Lock()
	r.directMethods = append(r.directMethods, m)
	r.mutex.Unlock()
}
func (r *allEventingData) GetDirectMethods() []sharedStructs.RabbitmqDirectMethod {
	return r.directMethods
}
