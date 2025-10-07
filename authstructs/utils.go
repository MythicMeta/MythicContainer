package authstructs

import (
	"fmt"
	"github.com/MythicMeta/MythicContainer/utils/helpers"
	"github.com/MythicMeta/MythicContainer/utils/sharedStructs"
	"sync"
)

type AuthDefinition struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	// SemVer is a specific semantic version tracker you can use for your payload type
	SemVer                string   `json:"semver"`
	IDPServices           []string `json:"idp_services"`
	NonIDPServices        []string `json:"non_idp_services"`
	GetIDPMetadata        func(GetIDPMetadataMessage) GetIDPMetadataMessageResponse
	GetIDPRedirect        func(GetIDPRedirectMessage) GetIDPRedirectMessageResponse
	ProcessIDPResponse    func(ProcessIDPResponseMessage) ProcessIDPResponseMessageResponse
	GetNonIDPMetadata     func(GetNonIDPMetadataMessage) GetNonIDPMetadataMessageResponse
	GetNonIDPRedirect     func(GetNonIDPRedirectMessage) GetNonIDPRedirectMessageResponse
	ProcessNonIDPResponse func(ProcessNonIDPResponseMessage) ProcessNonIDPResponseMessageResponse
	// Subscriptions - don't bother here, this will be auto filled out on syncing
	Subscriptions            []string                                                                                  `json:"subscriptions"`
	OnContainerStartFunction func(sharedStructs.ContainerOnStartMessage) sharedStructs.ContainerOnStartMessageResponse `json:"-"`
}

// REQUIRED, Don't Modify
type allAuthData struct {
	mutex          sync.RWMutex
	rpcMethods     []sharedStructs.RabbitmqRPCMethod
	directMethods  []sharedStructs.RabbitmqDirectMethod
	authDefinition AuthDefinition
}

var (
	AllAuthData containerAuthData
)

type containerAuthData struct {
	AuthMap map[string]*allAuthData
}

func (r *containerAuthData) GetAllNames() []string {
	names := []string{}
	for key, _ := range r.AuthMap {
		if key != "" && !helpers.StringSliceContains(names, key) {
			names = append(names, key)
		}
	}
	return names
}
func (r *containerAuthData) Get(name string) *allAuthData {
	if r.AuthMap == nil {
		r.AuthMap = make(map[string]*allAuthData)
	}
	if existingC2Data, ok := r.AuthMap[name]; !ok {
		newC2Data := allAuthData{}
		r.AuthMap[name] = &newC2Data
		return &newC2Data
	} else {
		return existingC2Data
	}
}
func (r *allAuthData) AddAuthDefinition(def AuthDefinition) {
	r.authDefinition = def
}
func (r *allAuthData) GetAuthDefinition() AuthDefinition {
	return r.authDefinition
}
func (r *allAuthData) SetSubscriptions(subs []string) {
	r.authDefinition.Subscriptions = subs
}
func (r *allAuthData) SetName(name string) {
	r.authDefinition.Name = name
}
func (r *allAuthData) GetRoutingKey(routingKey string) string {
	return fmt.Sprintf("%s_%s", r.authDefinition.Name, routingKey)
}
func (r *allAuthData) AddRPCMethod(m sharedStructs.RabbitmqRPCMethod) {
	r.mutex.Lock()
	r.rpcMethods = append(r.rpcMethods, m)
	r.mutex.Unlock()
}
func (r *allAuthData) GetRPCMethods() []sharedStructs.RabbitmqRPCMethod {
	return r.rpcMethods
}
func (r *allAuthData) AddDirectMethod(m sharedStructs.RabbitmqDirectMethod) {
	r.mutex.Lock()
	r.directMethods = append(r.directMethods, m)
	r.mutex.Unlock()
}
func (r *allAuthData) GetDirectMethods() []sharedStructs.RabbitmqDirectMethod {
	return r.directMethods
}
