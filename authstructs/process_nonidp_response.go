package authstructs

type ProcessNonIDPResponseMessage struct {
	ContainerName string            `json:"container_name"`
	ServerName    string            `json:"server_name"`
	NonIDPName    string            `json:"idp_name"`
	RequestValues map[string]string `json:"request_values"`
}
type ProcessNonIDPResponseMessageResponse struct {
	SuccessfulAuthentication bool   `json:"successful_authentication"`
	Error                    string `json:"error"`
	Email                    string `json:"email"`
}
