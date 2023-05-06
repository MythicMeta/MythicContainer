package agentstructs

// PAYLOAD_BUILD STRUCTS

// PayloadBuildMessage - A structure of the build information the user provided to generate an instance of the payload type.
// This information gets passed to your payload type's build function.
type PayloadBuildMessage struct {
	// PayloadType - the name of the payload type for the build
	PayloadType string `json:"payload_type"`
	// Filename - the name of the file the user originally supplied for this build
	Filename string `json:"filename"`
	// CommandList - the list of commands the user selected to include in the build
	CommandList []string `json:"commands"`
	// build param name : build value
	// BuildParameters - map of param name -> build value from the user for the build parameters defined
	// File type build parameters are supplied as a string UUID to use with MythicRPC for fetching file contents
	// Array type build parameters are supplied as []string{}
	BuildParameters map[string]interface{} `json:"build_parameters"`
	// C2Profiles - list of C2 profiles selected to include in the payload and their associated parameters
	C2Profiles []PayloadBuildC2Profile `json:"c2profiles"`
	// WrappedPayload - bytes of the wrapped payload if one exists
	WrappedPayload *[]byte `json:"wrapped_payload,omitempty"`
	// SelectedOS - the operating system the user selected when building the agent
	SelectedOS string `json:"selected_os"`
	// PayloadUUID - the Mythic generated UUID for this payload instance
	PayloadUUID string `json:"uuid"`
	// PayloadFileUUID - The Mythic generated File UUID associated with this payload
	PayloadFileUUID string `json:"payload_file_uuid"`
}

// PayloadBuildC2Profile - A structure of the selected C2 Profile information the user selected to build into a payload.
type PayloadBuildC2Profile struct {
	Name  string `json:"name"`
	IsP2P bool   `json:"is_p2p"`
	// parameter name: parameter value
	// Parameters - this is an interface of parameter name -> parameter value from the associated C2 profile.
	// The types for the various parameter names can be found by looking at the build parameters in the Mythic UI.
	Parameters map[string]interface{} `json:"parameters"`
}

type PAYLOAD_BUILD_STATUS = string

const (
	PAYLOAD_BUILD_STATUS_SUCCESS PAYLOAD_BUILD_STATUS = "success"
	PAYLOAD_BUILD_STATUS_ERROR                        = "error"
)

// PayloadBuildResponse - The result of calling a payload type's build function. This returns not only the actual
// payload bytes, but surrounding metadata such as updated filenames, command lists, and stdout/stderr messages.
type PayloadBuildResponse struct {
	// PayloadUUID - The UUID associated with this payload
	PayloadUUID string `json:"uuid"`
	// Success - was this build process successful or not
	Success bool `json:"success"`
	// UpdatedFilename - Optionally updated filename based on build parameters to more closely match the return file type
	UpdatedFilename *string `json:"updated_filename,omitempty"`
	// Payload - the raw bytes of the payload that was compiled/created
	Payload *[]byte `json:"payload,omitempty"`
	// UpdatedCommandList - if you want to adjust the list of commands in this payload from what the user provided,
	// provide the updated list of command names here
	UpdatedCommandList *[]string `json:"updated_command_list,omitempty"`
	// BuildStdErr - build stderr message to associate with the build
	BuildStdErr string `json:"build_stderr"`
	// BuildStdOut - build stdout message to associate with the build
	BuildStdOut string `json:"build_stdout"`
	// BuildMessage - general message to associate with the build. Usually not as verbose as the stdout/stderr.
	BuildMessage string `json:"build_message"`
}
