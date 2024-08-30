package authstructs

type GetIDPRedirectMessage struct {
	ContainerName  string            `json:"container_name"`
	ServerName     string            `json:"server_name"`
	IDPName        string            `json:"idp_name"`
	RequestURL     string            `json:"request_url"`
	RequestHeaders map[string]string `json:"request_headers"`
	RequestCookies map[string]string `json:"request_cookies"`
	RequestQuery   map[string]string `json:"request_query"`
}
type GetIDPRedirectMessageResponse struct {
	Success         bool              `json:"success"`
	Error           string            `json:"error"`
	RedirectURL     string            `json:"redirect_url"`
	RedirectHeaders map[string]string `json:"redirect_headers"`
	RedirectCookies map[string]string `json:"redirect_cookies"`
}
