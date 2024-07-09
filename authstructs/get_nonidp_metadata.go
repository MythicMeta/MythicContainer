package authstructs

type GetNonIDPMetadataMessage struct {
	ContainerName string `json:"container_name"`
	ServerName    string `json:"server_name"`
	NonIDPName    string `json:"nonidp_name"`
}
type GetNonIDPMetadataMessageResponse struct {
	Success  bool   `json:"success"`
	Error    string `json:"error"`
	Metadata string `json:"metadata"`
}
