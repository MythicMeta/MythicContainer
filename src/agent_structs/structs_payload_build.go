package agentstructs

// PAYLOAD_BUILD STRUCTS

type PayloadBuildMessage struct {
	PayloadType string   `json:"payload_type"`
	CommandList []string `json:"commands"`
	// build param name : build value
	BuildParameters map[string]interface{}  `json:"build_parameters"`
	C2Profiles      []PayloadBuildC2Profile `json:"c2profiles"`
	WrappedPayload  *[]byte                 `json:"wrapped_payload,omitempty"`
	SelectedOS      string                  `json:"selected_os"`
	PayloadUUID     string                  `json:"uuid"`
	PayloadFileUUID string                  `json:"payload_file_uuid"`
}

type PayloadBuildC2Profile struct {
	Name  string `json:"name"`
	IsP2P bool   `json:"is_p2p"`
	// parameter name: parameter value
	Parameters map[string]interface{} `json:"parameters"`
}

type PAYLOAD_BUILD_STATUS = string

const (
	PAYLOAD_BUILD_STATUS_SUCCESS PAYLOAD_BUILD_STATUS = "success"
	PAYLOAD_BUILD_STATUS_ERROR                        = "error"
)

type PayloadBuildResponse struct {
	PayloadUUID        string    `json:"uuid"`
	Success            bool      `json:"success"`
	Payload            *[]byte   `json:"payload,omitempty"`
	UpdatedCommandList *[]string `json:"updated_command_list,omitempty"`
	BuildStdErr        string    `json:"build_stderr"`
	BuildStdOut        string    `json:"build_stdout"`
	BuildMessage       string    `json:"build_message"`
}
