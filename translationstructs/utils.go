package translationstructs

import (
	"fmt"
	"github.com/MythicMeta/MythicContainer/utils/helpers"
	"sync"
)

type RabbitmqRPCMethod struct {
	RabbitmqRoutingKey         string
	RabbitmqProcessingFunction func([]byte) interface{}
}
type RabbitmqDirectMethod struct {
	RabbitmqRoutingKey         string
	RabbitmqProcessingFunction func([]byte)
}

// REQUIRED, Don't Modify
type allPayloadData struct {
	mutex             sync.RWMutex
	containerVersion  string
	rpcMethods        []RabbitmqRPCMethod
	directMethods     []RabbitmqDirectMethod
	payloadDefinition TranslationContainer
}

var (
	AllTranslationData containerPayloadData
)

type containerPayloadData struct {
	PayloadMap map[string]*allPayloadData
}

func (r *containerPayloadData) GetAllPayloadTypeNames() []string {
	names := []string{}
	for key, _ := range r.PayloadMap {
		if key != "" && !helpers.StringSliceContains(names, key) {
			names = append(names, key)
		}

	}
	return names
}
func (r *containerPayloadData) Get(name string) *allPayloadData {
	if r.PayloadMap == nil {
		r.PayloadMap = make(map[string]*allPayloadData)
	}
	if existingC2Data, ok := r.PayloadMap[name]; !ok {
		newC2Data := allPayloadData{}
		r.PayloadMap[name] = &newC2Data
		return &newC2Data
	} else {
		return existingC2Data
	}
}
func (r *allPayloadData) AddPayloadDefinition(payloadDef TranslationContainer) {
	r.payloadDefinition = payloadDef
}
func (r *allPayloadData) GetPayloadDefinition() TranslationContainer {

	return r.payloadDefinition
}
func (r *allPayloadData) AddContainerVersion(ver string) {
	r.containerVersion = ver
}
func (r *allPayloadData) GetPayloadName() string {
	return r.payloadDefinition.Name
}
func (r *allPayloadData) GetAuthor() string {
	return r.payloadDefinition.Author
}
func (r *allPayloadData) GetDescription() string {
	return r.payloadDefinition.Description
}
func (r *allPayloadData) GetContainerVersion() string {
	return r.containerVersion
}
func (r *allPayloadData) AddRPCMethod(m RabbitmqRPCMethod) {
	r.mutex.Lock()
	r.rpcMethods = append(r.rpcMethods, m)
	r.mutex.Unlock()
}
func (r *allPayloadData) GetRPCMethods() []RabbitmqRPCMethod {
	return r.rpcMethods
}
func (r *allPayloadData) AddDirectMethod(m RabbitmqDirectMethod) {
	r.mutex.Lock()
	r.directMethods = append(r.directMethods, m)
	r.mutex.Unlock()
}
func (r *allPayloadData) GetDirectMethods() []RabbitmqDirectMethod {
	return r.directMethods
}
func (r *allPayloadData) GetRoutingKey(baseRoutingKey string) string {
	return fmt.Sprintf("%s_%s", r.GetPayloadName(), baseRoutingKey)
}
