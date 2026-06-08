package c2structs

type C2_HOST_FILE_STATUS = string

type C2HostFileMessage struct {
	AgentFileID   string `json:"agent_file_id"`
	HostURL       string `json:"host_url"`
	Remove        bool   `json:"remove"`
	DownloadToken string `json:"download_token,omitempty"`
	Filename      string `json:"filename,omitempty"`
}

type C2HostFileMessageResponse struct {
	Success     bool   `json:"success"`
	Error       string `json:"error"`
	AgentFileID string `json:"agent_file_id"`
	HostURL     string `json:"host_url"`
}

type C2HostFilesMessage struct {
	Name  string              `json:"c2_profile_name"`
	Files []C2HostFileMessage `json:"files"`
}

type C2HostFilesMessageResponse struct {
	Success               bool                        `json:"success"`
	Error                 string                      `json:"error"`
	Results               []C2HostFileMessageResponse `json:"results,omitempty"`
	RestartInternalServer bool                        `json:"restart_internal_server,omitempty"`
}
