package agentstructs

const (
	SUPPORTED_OS_MACOS    = "macOS"
	SUPPORTED_OS_WINDOWS  = "Windows"
	SUPPORTED_OS_LINUX    = "Linux"
	SUPPORTED_OS_CHROME   = "Chrome"
	SUPPORTED_OS_WEBSHELL = "WebShell"
)

const (
	SUPPORTED_UI_FEATURE_TASK_PROCESS_INTERACTIVE_TASKS = "task:process_interactive_tasks"
	SUPPORTED_UI_FEATURE_TASK_RESPONSE_INTERACTIVE      = "task_response:interactive"
	SUPPORTED_UI_FEATURE_CALLBACK_TABLE_EXIT            = "callback_table:exit"
	SUPPORTED_UI_FEATURE_FILE_BROWSER_LIST              = "file_browser:list"
	SUPPORTED_UI_FEATURE_FILE_BROWSER_REMOVE            = "file_browser:remove"
	SUPPORTED_UI_FEATURE_FILE_BROWSER_UPLOAD            = "file_browser:upload"
	SUPPORTED_UI_FEATURE_FILE_BROWSER_DOWNLOAD          = "file_browser:download"
	SUPPORTED_UI_FEATURE_PROCESS_BROWSER_LIST           = "process_browser:list"
	SUPPORTED_UI_FEATURE_PROCESS_BROWSER_KILL           = "process_browser:kill"
	SUPPORTED_UI_FEATURE_PROCESS_BROWSER_INJECT         = "process_browser:inject"
	SUPPORTED_UI_FEATURE_PROCESS_BROWSER_STEAL_TOKEN    = "process_browser:steal_token"
	SUPPORTED_UI_FEATURE_PROCESS_BROWSER_LIST_TOKENS    = "process_browser:list_tokens"
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
