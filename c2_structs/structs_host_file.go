package c2structs

type C2_HOST_FILE_STATUS = string

type C2HostFileMessage struct {
	Name     string `json:"c2_profile_name"`
	FileUUID string `json:"file_uuid"`
	HostURL  string `json:"host_url"`
	Remove   bool   `json:"remove"`
}

type C2HostFileMessageResponse struct {
	Success               bool   `json:"success"`
	Error                 string `json:"error"`
	RestartInternalServer bool   `json:"restart_internal_server"`
}
