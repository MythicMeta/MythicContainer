package sharedStructs

type RabbitmqRPCMethod struct {
	RabbitmqRoutingKey         string
	RabbitmqProcessingFunction func([]byte) interface{}
}
type RabbitmqDirectMethod struct {
	RabbitmqRoutingKey         string
	RabbitmqProcessingFunction func([]byte)
}

type ContainerOnStartMessage struct {
	ContainerName string `json:"container_name"`
	OperationID   int    `json:"operation_id"`
	OperationName string `json:"operation_name"`
	ServerName    string `json:"server_name"`
	APIToken      string `json:"apitoken"`
}

type ContainerOnStartMessageResponse struct {
	ContainerName         string `json:"container_name"`
	EventLogInfoMessage   string `json:"stdout"`
	EventLogErrorMessage  string `json:"stderr"`
	RestartInternalServer bool   `json:"restart_internal_server"`
}

type ContainerRPCGetFileMessage struct {
	ContainerName string `json:"container_name"`
	Filename      string `json:"filename"`
}

type ContainerRPCGetFileMessageResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Message []byte `json:"message"`
}

type ContainerRPCListFileMessage struct {
	ContainerName string `json:"container_name"`
}

type ContainerRPCListFileMessageResponse struct {
	Success bool     `json:"success"`
	Error   string   `json:"error"`
	Files   []string `json:"files"`
}

type ContainerRPCRemoveFileMessage struct {
	ContainerName string `json:"container_name"`
	Filename      string `json:"filename"`
}

type ContainerRPCRemoveFileMessageResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

type ContainerRPCWriteFileMessage struct {
	ContainerName string `json:"container_name"`
	Filename      string `json:"filename"`
	Contents      []byte `json:"contents"`
}

type ContainerRPCWriteFileMessageResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Message string `json:"message"`
}
