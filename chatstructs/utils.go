package chatstructs

import (
	"context"
	"fmt"
	"sync"

	"github.com/MythicMeta/MythicContainer/utils/helpers"
	"github.com/MythicMeta/MythicContainer/utils/sharedStructs"
)

type ChatModelDefinition struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Metadata    map[string]interface{} `json:"metadata"`
}

type ChatDefinition struct {
	Name                     string `json:"name"`
	Description              string `json:"description"`
	SemVer                   string `json:"semver"`
	Models                   []ChatModelDefinition
	Subscriptions            []string `json:"subscriptions"`
	ChatFunction             func(context.Context, ChatRequestMessage)
	OnContainerStartFunction func(context.Context, sharedStructs.ContainerOnStartMessage) sharedStructs.ContainerOnStartMessageResponse `json:"-"`
}

type allChatData struct {
	mutex          sync.RWMutex
	rpcMethods     []sharedStructs.RabbitmqRPCMethod
	directMethods  []sharedStructs.RabbitmqDirectMethod
	chatDefinition ChatDefinition
}

var (
	AllChatData containerChatData
)

type containerChatData struct {
	ChatMap map[string]*allChatData
}

func (r *containerChatData) GetAllNames() []string {
	names := []string{}
	for key := range r.ChatMap {
		if key != "" && !helpers.StringSliceContains(names, key) {
			names = append(names, key)
		}
	}
	return names
}

func (r *containerChatData) Get(name string) *allChatData {
	if r.ChatMap == nil {
		r.ChatMap = make(map[string]*allChatData)
	}
	if existingChatData, ok := r.ChatMap[name]; !ok {
		newChatData := allChatData{}
		r.ChatMap[name] = &newChatData
		return &newChatData
	} else {
		return existingChatData
	}
}

func (r *allChatData) AddChatDefinition(def ChatDefinition) {
	r.chatDefinition = def
}

func (r *allChatData) GetChatDefinition() ChatDefinition {
	return r.chatDefinition
}

func (r *allChatData) SetSubscriptions(subs []string) {
	r.chatDefinition.Subscriptions = subs
}

func (r *allChatData) SetName(name string) {
	r.chatDefinition.Name = name
}

func (r *allChatData) GetRoutingKey(routingKey string) string {
	return fmt.Sprintf("%s_%s", r.chatDefinition.Name, routingKey)
}

func (r *allChatData) AddRPCMethod(m sharedStructs.RabbitmqRPCMethod) {
	r.mutex.Lock()
	r.rpcMethods = append(r.rpcMethods, m)
	r.mutex.Unlock()
}

func (r *allChatData) GetRPCMethods() []sharedStructs.RabbitmqRPCMethod {
	return r.rpcMethods
}

func (r *allChatData) AddDirectMethod(m sharedStructs.RabbitmqDirectMethod) {
	r.mutex.Lock()
	r.directMethods = append(r.directMethods, m)
	r.mutex.Unlock()
}

func (r *allChatData) GetDirectMethods() []sharedStructs.RabbitmqDirectMethod {
	return r.directMethods
}
