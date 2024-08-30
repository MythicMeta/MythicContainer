package agentstructs

import "time"

type PTCallbacksToCheck struct {
	ID               int       `json:"id"`
	DisplayID        int       `json:"display_id"`
	AgentCallbackID  string    `json:"agent_callback_id"`
	InitialCheckin   time.Time `json:"initial_checkin"`
	LastCheckin      time.Time `json:"last_checkin"`
	SleepInfo        string    `json:"sleep_info"`
	ActiveC2Profiles []string  `json:"active_c2_profiles"`
}
type PTCheckIfCallbacksAliveMessage struct {
	ContainerName string               `json:"container_name"`
	Callbacks     []PTCallbacksToCheck `json:"callbacks"`
}
type PTCallbacksToCheckResponse struct {
	ID    int  `json:"id"`
	Alive bool `json:"alive"`
}
type PTCheckIfCallbacksAliveMessageResponse struct {
	Success   bool   `json:"success"`
	Error     string `json:"error"`
	Callbacks []PTCallbacksToCheckResponse
}
