package authstructs

type ProcessIDPResponseMessage struct {
	ContainerName  string            `json:"container_name"`
	ServerName     string            `json:"server_name"`
	IDPName        string            `json:"idp_name"`
	RequestURL     string            `json:"request_url"`
	RequestHeaders map[string]string `json:"request_headers"`
	RequestCookies map[string]string `json:"request_cookies"`
	RequestQuery   map[string]string `json:"request_query"`
	RequestBody    string            `json:"request_body"`
}
type ProcessIDPResponseMessageResponse struct {
	SuccessfulAuthentication bool   `json:"successful_authentication"`
	Error                    string `json:"error"`
	Email                    string `json:"email"`
}
