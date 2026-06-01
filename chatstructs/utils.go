package chatstructs

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/MythicMeta/MythicContainer/utils/helpers"
	"github.com/MythicMeta/MythicContainer/utils/sharedStructs"
)

type ChatModelConfigurationOptionType string

const (
	CHAT_MODEL_CONFIGURATION_OPTION_TYPE_STRING ChatModelConfigurationOptionType = "string"
	CHAT_MODEL_CONFIGURATION_OPTION_TYPE_NUMBER                                  = "number"
	CHAT_MODEL_CONFIGURATION_OPTION_TYPE_CHOICE                                  = "choice"
)

type ChatModelConfigurationOptionChoice struct {
	// Label - the human-readable option shown to operators in the Mythic UI.
	Label string `json:"label"`
	// Value - the value written into ChatRequestMessage.Config when the option is selected.
	Value string `json:"value"`
	// Description - optional helper text describing this choice.
	Description string `json:"description,omitempty"`
}

type ChatModelConfigurationOption struct {
	// Name - the config key sent to the chat container in ChatRequestMessage.Config.
	Name string `json:"name"`
	// DisplayName - the human-readable label shown for this option in the Mythic UI.
	DisplayName string `json:"display_name"`
	// Type - how the UI should render the option, such as string, number, or choice.
	Type ChatModelConfigurationOptionType `json:"type"`
	// Description - helper text that explains what the operator should supply.
	Description string `json:"description"`
	// Required - indicates whether operators must provide a value before using the model.
	Required bool `json:"required"`
	// DefaultValue - the value shown before operators configure this option. The type should match Type.
	DefaultValue interface{} `json:"default_value"`
	// Choices - selectable values when Type is choice.
	Choices []ChatModelConfigurationOptionChoice `json:"choices,omitempty"`
}

type ChatModelMetadata struct {
	// Provider - the backing service or model provider, such as litellm, openai, anthropic, or a local engine name.
	Provider string `json:"provider,omitempty"`
	// ConfigurationOptions - per-chat configuration fields Mythic can render for operators and send in ChatRequestMessage.Config.
	ConfigurationOptions []ChatModelConfigurationOption `json:"configuration_options,omitempty"`
	// RequiredUserSecrets - Mythic user secret names required before this model can be used.
	RequiredUserSecrets []string `json:"required_user_secrets,omitempty"`
	// OptionalUserSecrets - Mythic user secret names this model can use when present.
	OptionalUserSecrets []string `json:"optional_user_secrets,omitempty"`
	// RequiredChannelAPITokenScopes - API token scopes required on the AI chat channel for this model's tool access.
	RequiredChannelAPITokenScopes []string `json:"required_channel_api_token_scopes,omitempty"`
	// AdditionalItems - model-specific metadata not covered by the typed fields above. Prefer typed fields when one applies.
	AdditionalItems map[string]interface{} `json:"-"`
}

func (c ChatModelMetadata) MarshalJSON() ([]byte, error) {
	type Alias ChatModelMetadata
	baseBytes, err := json.Marshal(Alias(c))
	if err != nil {
		return nil, err
	}
	baseMap := map[string]interface{}{}
	if err = json.Unmarshal(baseBytes, &baseMap); err != nil {
		return nil, err
	}
	for key, value := range c.AdditionalItems {
		if _, exists := baseMap[key]; !exists {
			baseMap[key] = value
		}
	}
	return json.Marshal(baseMap)
}

type ChatModelDefinition struct {
	// Name - the model name operators select when sending a prompt to this chat container.
	Name string `json:"name"`
	// Description - human-readable summary of what this model does and when operators should use it.
	Description string `json:"description"`
	// Metadata - typed details about provider configuration, user secrets, UI config fields, defaults, and required scopes.
	Metadata ChatModelMetadata `json:"metadata"`
}

type ChatDefinition struct {
	Name          string `json:"name"`
	Description   string `json:"description"`
	SemVer        string `json:"semver"`
	Models        []ChatModelDefinition
	Subscriptions []string `json:"subscriptions"`
	// ChatFunction receives a context that is cancelled when Mythic sends a chat_cancel message for the request.
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
