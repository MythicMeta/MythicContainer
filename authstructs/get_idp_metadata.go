package authstructs

type GetIDPMetadataMessage struct {
	ContainerName string `json:"container_name"`
	ServerName    string `json:"server_name"`
	IDPName       string `json:"idp_name"`
}
type GetIDPMetadataMessageResponse struct {
	Success  bool   `json:"success"`
	Error    string `json:"error"`
	Metadata string `json:"metadata"`
}
