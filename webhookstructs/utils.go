package webhookstructs

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/MythicMeta/MythicContainer/utils/sharedStructs"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/MythicMeta/MythicContainer/config"
	"github.com/MythicMeta/MythicContainer/logging"
	"github.com/MythicMeta/MythicContainer/utils/helpers"
)

const EMIT_WEBHOOK_ROUTING_KEY_PREFIX = "emit_webhook"

type WEBHOOK_TYPE = string

const (
	WEBHOOK_TYPE_NEW_CALLBACK WEBHOOK_TYPE = "new_callback"
	WEBHOOK_TYPE_NEW_FEEDBACK              = "new_feedback"
	WEBHOOK_TYPE_NEW_STARTUP               = "new_startup"
	WEBHOOK_TYPE_NEW_ALERT                 = "new_alert"
	WEBHOOK_TYPE_NEW_CUSTOM                = "new_custom"
)

type WebhookDefinition struct {
	Name        string
	Description string
	// SemVer is a specific semantic version tracker you can use for your payload type
	SemVer                   string `json:"semver"`
	WebhookURL               string
	WebhookChannel           string
	NewFeedbackFunction      func(input NewFeedbackWebookMessage)
	NewCallbackFunction      func(input NewCallbackWebookMessage)
	NewStartupFunction       func(input NewStartupWebhookMessage)
	NewAlertFunction         func(input NewAlertWebhookMessage)
	NewCustomFunction        func(input NewCustomWebhookMessage)
	Subscriptions            []string
	OnContainerStartFunction func(sharedStructs.ContainerOnStartMessage) sharedStructs.ContainerOnStartMessageResponse
}

var tr = &http.Transport{
	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	MaxIdleConns:    10,
	MaxConnsPerHost: 10,
	//IdleConnTimeout: 1 * time.Nanosecond,
}
var HTTPClient = &http.Client{
	Timeout:   5 * time.Second,
	Transport: tr,
}

// REQUIRED, Don't Modify
type allWebhookData struct {
	mutex             sync.RWMutex
	rpcMethods        []sharedStructs.RabbitmqRPCMethod
	directMethods     []sharedStructs.RabbitmqDirectMethod
	webhookDefinition WebhookDefinition
}

var (
	AllWebhookData containerWebhookData
)

type containerWebhookData struct {
	WebhookMap map[string]*allWebhookData
}

func (r *containerWebhookData) GetAllNames() []string {
	names := []string{}
	for key, _ := range r.WebhookMap {
		if key != "" && !helpers.StringSliceContains(names, key) {
			names = append(names, key)
		}
	}
	return names
}
func (r *containerWebhookData) Get(name string) *allWebhookData {
	if r.WebhookMap == nil {
		r.WebhookMap = make(map[string]*allWebhookData)
	}
	if existingC2Data, ok := r.WebhookMap[name]; !ok {
		newC2Data := allWebhookData{}
		r.WebhookMap[name] = &newC2Data
		return &newC2Data
	} else {
		return existingC2Data
	}
}
func (r *allWebhookData) AddWebhookDefinition(def WebhookDefinition) {
	r.webhookDefinition = def
}
func (r *allWebhookData) GetWebhookDefinition() WebhookDefinition {
	return r.webhookDefinition
}
func (r *allWebhookData) SetSubscriptions(subs []string) {
	r.webhookDefinition.Subscriptions = subs
}
func (r *allWebhookData) SetName(name string) {
	r.webhookDefinition.Name = name
}
func (r *allWebhookData) GetRoutingKey(routingKey string) string {
	return fmt.Sprintf("%s_%s", r.webhookDefinition.Name, routingKey)
}
func (r *allWebhookData) AddRPCMethod(m sharedStructs.RabbitmqRPCMethod) {
	r.mutex.Lock()
	r.rpcMethods = append(r.rpcMethods, m)
	r.mutex.Unlock()
}
func (r *allWebhookData) GetRPCMethods() []sharedStructs.RabbitmqRPCMethod {
	return r.rpcMethods
}
func (r *allWebhookData) AddDirectMethod(m sharedStructs.RabbitmqDirectMethod) {
	r.mutex.Lock()
	r.directMethods = append(r.directMethods, m)
	r.mutex.Unlock()
}
func (r *allWebhookData) GetDirectMethods() []sharedStructs.RabbitmqDirectMethod {
	return r.directMethods
}
func (r *allWebhookData) AddWebhookURL(url string) {
	r.mutex.Lock()
	r.webhookDefinition.WebhookURL = url
	r.mutex.Unlock()
}
func (r *allWebhookData) GetWebhookURL(input interface{}, channelType WEBHOOK_TYPE) string {
	switch channelType {
	case WEBHOOK_TYPE_NEW_FEEDBACK:
		msg := input.(NewFeedbackWebookMessage)
		if msg.OperationWebhook != "" {
			return msg.OperationWebhook
		}
	case WEBHOOK_TYPE_NEW_CALLBACK:
		msg := input.(NewCallbackWebookMessage)
		if msg.OperationWebhook != "" {
			return msg.OperationWebhook
		}
	case WEBHOOK_TYPE_NEW_STARTUP:
		msg := input.(NewStartupWebhookMessage)
		if msg.OperationWebhook != "" {
			return msg.OperationWebhook
		}
	case WEBHOOK_TYPE_NEW_ALERT:
		msg := input.(NewAlertWebhookMessage)
		if msg.OperationWebhook != "" {
			return msg.OperationWebhook
		}
	case WEBHOOK_TYPE_NEW_CUSTOM:
		msg := input.(NewCustomWebhookMessage)
		if msg.OperationWebhook != "" {
			return msg.OperationWebhook
		}
	default:
	}
	// allow the environment to override the program definition
	if config.MythicConfig.WebhookDefaultURL != "" {
		return config.MythicConfig.WebhookDefaultURL
	} else {
		return r.webhookDefinition.WebhookURL
	}
}
func (r *allWebhookData) AddWebhookChannel(channel string) {
	r.mutex.Lock()
	r.webhookDefinition.WebhookChannel = channel
	r.mutex.Unlock()
}
func (r *allWebhookData) GetWebhookChannel(input interface{}, channelType WEBHOOK_TYPE) string {
	switch channelType {
	case WEBHOOK_TYPE_NEW_FEEDBACK:
		msg := input.(NewFeedbackWebookMessage)
		if config.MythicConfig.WebhookFeedbackChannel != "" {
			return config.MythicConfig.WebhookFeedbackChannel
		} else if msg.OperationChannel != "" {
			return msg.OperationChannel
		}
	case WEBHOOK_TYPE_NEW_CALLBACK:
		msg := input.(NewCallbackWebookMessage)
		if config.MythicConfig.WebhookCallbackChannel != "" {
			return config.MythicConfig.WebhookCallbackChannel
		} else if msg.OperationChannel != "" {
			return msg.OperationChannel
		}
	case WEBHOOK_TYPE_NEW_STARTUP:
		msg := input.(NewStartupWebhookMessage)
		if config.MythicConfig.WebhookStartupChannel != "" {
			return config.MythicConfig.WebhookStartupChannel
		} else if msg.OperationChannel != "" {
			return msg.OperationChannel
		}
	case WEBHOOK_TYPE_NEW_ALERT:
		msg := input.(NewAlertWebhookMessage)
		if config.MythicConfig.WebhookAlertChannel != "" {
			return config.MythicConfig.WebhookAlertChannel
		} else if msg.OperationChannel != "" {
			return msg.OperationChannel
		}
	case WEBHOOK_TYPE_NEW_CUSTOM:
		msg := input.(NewCustomWebhookMessage)
		if config.MythicConfig.WebhookCustomChannel != "" {
			return config.MythicConfig.WebhookCustomChannel
		} else if msg.OperationChannel != "" {
			return msg.OperationChannel
		}
	default:
		logging.LogError(nil, "unknown webhook type when getting webhook channel", "type", channelType)
	}
	if config.MythicConfig.WebhookDefaultChannel != "" {
		return config.MythicConfig.WebhookDefaultChannel
	} else {
		return r.webhookDefinition.WebhookChannel
	}
}
func GetRoutingKeyFor(webhookType string) string {
	return fmt.Sprintf("%s.%s", EMIT_WEBHOOK_ROUTING_KEY_PREFIX, webhookType)
}
func SubmitWebRequest(method string, url string, body interface{}) ([]byte, int, error) {
	messageBytes, err := json.Marshal(body)
	if err != nil {
		logging.LogError(err, "Failed to marshal new webhook message")
		return nil, 0, err
	}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(messageBytes))
	if err != nil {
		logging.LogError(err, "Failed to make new http request")
		return nil, 0, err
	}
	req.ContentLength = int64(len(messageBytes))
	resp, err := HTTPClient.Do(req)
	if err != nil {
		logging.LogError(err, "Failed to make http request")
		return nil, 0, err
	}
	defer resp.Body.Close()
	resultBody, err := io.ReadAll(resp.Body)
	if err != nil {
		logging.LogError(err, "Failed to get response from webhook")
		return nil, resp.StatusCode, err
	}
	//logging.LogDebug("webhook output", "response", body)
	return resultBody, resp.StatusCode, nil
}
