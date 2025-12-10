package mythicrpc

import (
	"encoding/json"

	"github.com/MythicMeta/MythicContainer/logging"
	"github.com/MythicMeta/MythicContainer/rabbitmq"
)

type MythicRPCCallbackUpdateMessage struct {
	AgentCallbackID                   *string   `json:"agent_callback_id"` // required
	CallbackID                        *int      `json:"callback_id"`
	TaskID                            *int      `json:"task_id"`
	EncryptionKey                     *[]byte   `json:"encryption_key,omitempty"`
	DecryptionKey                     *[]byte   `json:"decryption_key,omitempty"`
	CryptoType                        *string   `json:"crypto_type,omitempty"`
	User                              *string   `json:"user,omitempty"`
	Host                              *string   `json:"host,omitempty"`
	PID                               *int      `json:"pid,omitempty"`
	ExtraInfo                         *string   `json:"extra_info,omitempty"`
	SleepInfo                         *string   `json:"sleep_info,omitempty"`
	Ip                                *string   `json:"ip,omitempty"`
	IPs                               *[]string `json:"ips,omitempty"`
	ExternalIP                        *string   `json:"external_ip,omitempty"`
	IntegrityLevel                    *int      `json:"integrity_level,omitempty"`
	Os                                *string   `json:"os,omitempty"`
	Domain                            *string   `json:"domain,omitempty"`
	Architecture                      *string   `json:"architecture,omitempty"`
	Description                       *string   `json:"description,omitempty"`
	ProcessName                       *string   `json:"process_name,omitempty"`
	UpdateLastCheckinTime             *bool     `json:"update_last_checkin_time,omitempty"`
	UpdateLastCheckinTimeViaC2Profile *string   `json:"update_last_checkin_time_via_c2_profile,omitempty"`
	Cwd                               *string   `json:"cwd,omitempty"`
	ImpersonationContext              *string   `json:"impersonation_context,omitempty"`
	Dead                              *bool     `json:"dead,omitempty"`
}
type MythicRPCCallbackUpdateMessageResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

func SendMythicRPCCallbackUpdate(input MythicRPCCallbackUpdateMessage) (*MythicRPCCallbackUpdateMessageResponse, error) {
	response := MythicRPCCallbackUpdateMessageResponse{}
	if responseBytes, err := rabbitmq.RabbitMQConnection.SendRPCStructMessage(
		rabbitmq.MYTHIC_EXCHANGE,
		rabbitmq.MYTHIC_RPC_CALLBACK_UPDATE,
		input,
	); err != nil {
		logging.LogError(err, "Failed to send RPC message")
		return nil, err
	} else if err := json.Unmarshal(responseBytes, &response); err != nil {
		logging.LogError(err, "Failed to parse response back to struct", "response", response)
		return nil, err
	} else {
		return &response, nil
	}
}
