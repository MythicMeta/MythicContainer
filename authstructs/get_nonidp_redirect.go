package authstructs

type GetNonIDPRedirectMessage struct {
	ContainerName string `json:"container_name"`
	ServerName    string `json:"server_name"`
	NonIDPName    string `json:"nonidp_name"`
}
type GetNonIDPRedirectMessageResponse struct {
	Success       bool     `json:"success"`
	Error         string   `json:"error"`
	RequestFields []string `json:"request_fields"`
}
