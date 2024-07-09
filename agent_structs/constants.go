package agentstructs

const (
	SUPPORTED_OS_MACOS    = "macOS"
	SUPPORTED_OS_WINDOWS  = "Windows"
	SUPPORTED_OS_LINUX    = "Linux"
	SUPPORTED_OS_CHROME   = "Chrome"
	SUPPORTED_OS_WEBSHELL = "WebShell"
)

type AgentType string

const (
	AgentTypeAgent          AgentType = "agent"
	AgentTypeWrapper                  = "wrapper"
	AgentTypeService                  = "service"
	AgentTypeCommandAugment           = "command_augment"
)

type MessageFormat string

const (
	MessageFormatJSON MessageFormat = "json"
	MessageFormatXML                = "xml"
)
