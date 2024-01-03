package rabbitmq

import "time"

const (
	MYTHIC_EXCHANGE                        = "mythic_exchange"
	MYTHIC_TOPIC_EXCHANGE                  = "mythic_topic_exchange"
	RETRY_CONNECT_DELAY                    = 5 * time.Second
	TIME_FORMAT_STRING_YYYY_MM_DD          = "2006-01-02"
	TIME_FORMAT_STRING_YYYY_MM_DD_HH_MM_SS = "2006-01-02 15:04:05 Z07"
	RPC_TIMEOUT                            = 30 * time.Second
	TASK_STATUS_CONTAINER_DOWN             = "Error: Container Down"
)

type CallbackPortType = string

const containerVersion = "v1.1.2"

const (
	CALLBACK_PORT_TYPE_SOCKS       CallbackPortType = "socks"
	CALLBACK_PORT_TYPE_RPORTFWD                     = "rpfwd"
	CALLBACK_PORT_TYPE_INTERACTIVE                  = "interactive"
)

// Direct fanout rabbitmq routes where Mythic is consuming messages, but others can also listen in and consume
const (
	// payload routes
	//	Syncing information about the payload type to Mythic
	//		send PayloadTypeSyncMessage to this route
	PT_SYNC_ROUTING_KEY       = "pt_sync"
	PT_RPC_RESYNC_ROUTING_KEY = "pt_rpc_resync"
	//	Result of asking a container to build a new payload
	//		send PayloadBuildResponse to this route
	PT_BUILD_RESPONSE_ROUTING_KEY = "pt_build_response"
	//	Result of asking a container to build a new c2profile-only payload for hot-swapping c2s
	//		send PayloadBuildC2Response to this route
	PT_BUILD_C2_RESPONSE_ROUTING_KEY = "pt_c2_build_response"
	//	Result of asking a container to do a pre-flight check for a task
	// 		send PTTTaskOPSECPreTaskMessageResponse to this route
	PT_TASK_OPSEC_PRE_CHECK_RESPONSE = "pt_task_opsec_pre_check_response"
	//	Result of asking a container to process a tasking request
	// 		send PTTaskCreateTaskingMessageResponse to this route
	PT_TASK_CREATE_TASKING_RESPONSE = "pt_task_create_tasking_response"
	//	Result of asking a container to do a post-flight check for a task before an agent picks it up
	//		send PTTaskOPSECPostTaskMessageResponse to this route
	PT_TASK_OPSEC_POST_CHECK_RESPONSE = "pt_task_opsec_post_check_response"
	//	Result of handling a task's completion function
	//		send PTTaskCompletionHandlerMessageResponse to this route
	PT_TASK_COMPLETION_FUNCTION_RESPONSE = "pt_task_completion_function_response"
	// c2 routes
	//		send C2SyncMessages to this route
	C2_SYNC_ROUTING_KEY       = "c2_sync"
	C2_RPC_RESYNC_ROUTING_KEY = "c2_rpc_resync"
	TR_SYNC_ROUTING_KEY       = "tr_sync"
	TR_RPC_RESYNC_ROUTING_KEY = "tr_rpc_resync"
)

// Direct fanout rabbitmq routes exclusively for other containers to listen and process
const (
	// emit routes
	//	Information specifically for SIEMs to ingest
	EMIT_SIEM_LOG_ROUTING_KEY = "emit_siem_log"
	//	Information specifically for sending to webhook endpoints in other platforms
	EMIT_WEBHOOK_ROUTING_KEY = "emit_webhook"
)

// Direct fanout rabbitmq routes where the container is consuming messages and responding back to Mythic, but others can also listen in and consume.
// These aren't RPC routes because these could take a long time and we don't want to block it. These will have the format of "containerName_"
// prepended to the constants below
//
//	ex: for the apfell agent, it would be "apfell_payload_build"
const (
	//	send PayloadBuildMessage to this route for agent containers to pick up
	//		send PayloadBuildResponse to PAYLOAD_BUILD_RESPONSE_ROUTING_KEY for Mythic to process response
	//
	PAYLOAD_BUILD_ROUTING_KEY = "payload_build"

	//
	PAYLOAD_BUILD_C2_ROUTING_KEY = "payload_c2_build"
	//	Sending a pre-flight tasking check to this queue
	PT_TASK_OPSEC_PRE_CHECK = "pt_task_opsec_pre_check"
	//
	PT_TASK_CREATE_TASKING = "pt_task_create_tasking"
	//
	PT_TASK_OPSEC_POST_CHECK = "pt_task_opsec_post_check"
	//
	PT_RPC_COMMAND_DYNAMIC_QUERY_FUNCTION = "pt_command_dynamic_query_function"
	//
	PT_RPC_COMMAND_TYPEDARRAY_PARSE = "pt_command_typedarray_parse"
	//
	PT_TASK_COMPLETION_FUNCTION = "pt_task_completion_function"
	//
	PT_TASK_PROCESS_RESPONSE          = "pt_task_process_response"
	PT_TASK_PROCESS_RESPONSE_RESPONSE = "pt_task_process_response_response"
)

// Routes where container is consuming messages and responding back to Mythic
//
//	These are exclusive to the container and not able for other containers to listen in on
//	These all have "containerName_" prepended to the constants below
//		ex: for the http profile, it would be "http_c2_rpc_opsec_check"
const (
	//

	C2_RPC_OPSEC_CHECKS_ROUTING_KEY = "c2_rpc_opsec_check"
	//
	C2_RPC_CONFIG_CHECK_ROUTING_KEY = "c2_rpc_config_check"
	//
	C2_RPC_GET_IOC_ROUTING_KEY = "c2_rpc_get_ioc"
	//
	C2_RPC_SAMPLE_MESSAGE_ROUTING_KEY = "c2_rpc_sample_message"
	//
	C2_RPC_REDIRECTOR_RULES_ROUTING_KEY = "c2_rpc_redirector_rules"
	//
	C2_RPC_START_SERVER_ROUTING_KEY = "c2_rpc_start_server"
	//
	C2_RPC_STOP_SERVER_ROUTING_KEY = "c2_rpc_stop_server"
	//
	C2_RPC_GET_SERVER_DEBUG_OUTPUT = "c2_rpc_get_server_debug_output"
	//
	C2_RPC_HOST_FILE = "c2_rpc_host_file"
	//
	C2_RPC_GET_FILE = "c2_rpc_get_file"
	//
	C2_RPC_REMOVE_FILE = "c2_rpc_remove_file"
	//
	C2_RPC_LIST_FILE = "c2_rpc_list_file"
	//
	C2_RPC_WRITE_FILE = "c2_rpc_write_file"
	//
	TR_RPC_GENERATE_KEYS = "tr_rpc_generate_keys"
	//
	TR_RPC_CONVERT_FROM_MYTHIC_C2_FORMAT = "tr_rpc_from_mythic_c2"
	//
	TR_RPC_CONVERT_TO_MYTHIC_C2_FORMAT = "tr_rpc_to_mythic_c2"
	//
	TR_RPC_ENCRYPT_BYTES = "tr_rpc_encrypt_bytes"
	//
	TR_RPC_DECRYPT_BYTES = "tr_rpc_decrypt_bytes"
)

// RPC Routes where Mythic is consuming messages and responding back to the container
//
//	These are exclusive to mythic and not able for other containers to listen in on
const (
	// MYTHIC_RPC_FILE_CREATE file operations
	MYTHIC_RPC_FILE_CREATE      = "mythic_rpc_file_create"
	MYTHIC_RPC_FILE_SEARCH      = "mythic_rpc_file_search"
	MYTHIC_RPC_FILE_UPDATE      = "mythic_rpc_file_update"
	MYTHIC_RPC_FILE_GET_CONTENT = "mythic_rpc_file_get_content"
	MYTHIC_RPC_FILE_REGISTER    = "mythic_rpc_file_register"
	// MYTHIC_RPC_PAYLOAD_CREATE_FROM_UUID payload operations
	MYTHIC_RPC_PAYLOAD_CREATE_FROM_UUID    = "mythic_rpc_payload_create_from_uuid"
	MYTHIC_RPC_PAYLOAD_CREATE_FROM_SCRATCH = "mythic_rpc_payload_create_from_scratch"
	MYTHIC_RPC_PAYLOAD_SEARCH              = "mythic_rpc_payload_search"
	MYTHIC_RPC_PAYLOAD_GET_PAYLOAD_CONTENT = "mythic_rpc_payload_get_content"
	MYTHIC_RPC_PAYLOAD_UPDATE_BUILD_STEP   = "mythic_rpc_payload_update_build_step"
	MYTHIC_RPC_PAYLOAD_ADD_COMMAND         = "mythic_rpc_payload_add_command"
	MYTHIC_RPC_PAYLOAD_REMOVE_COMMAND      = "mythic_rpc_payload_remove_command"
	// MYTHIC_RPC_TASK_SEARCH task operations
	MYTHIC_RPC_TASK_SEARCH                    = "mythic_rpc_task_search"
	MYTHIC_RPC_TASK_DISPLAY_TO_REAL_ID_SEARCH = "mythic_rpc_task_display_to_real_id_search"
	MYTHIC_RPC_TASK_UPDATE                    = "mythic_rpc_task_update"
	MYTHIC_RPC_TASK_CREATE_SUBTASK            = "mythic_rpc_task_create_subtask"
	MYTHIC_RPC_TASK_CREATE_SUBTASK_GROUP      = "mythic_rpc_task_create_group"
	// MYTHIC_RPC_RESPONSE_SEARCH response operations
	MYTHIC_RPC_RESPONSE_SEARCH = "mythic_rpc_response_search"
	MYTHIC_RPC_RESPONSE_CREATE = "mythic_rpc_response_create"
	// MYTHIC_RPC_COMMAND_SEARCH command operations
	MYTHIC_RPC_COMMAND_SEARCH = "mythic_rpc_command_search"
	// MYTHIC_RPC_CALLBACK_CREATE callback operations
	MYTHIC_RPC_CALLBACK_CREATE                    = "mythic_rpc_callback_create"
	MYTHIC_RPC_CALLBACK_SEARCH                    = "mythic_rpc_callback_search"
	MYTHIC_RPC_CALLBACK_EDGE_SEARCH               = "mythic_rpc_callback_edge_search"
	MYTHIC_RPC_CALLBACK_DISPLAY_TO_REAL_ID_SEARCH = "mythic_rpc_callback_display_to_real_id_search"
	MYTHIC_RPC_CALLBACK_ADD_COMMAND               = "mythic_rpc_callback_add_command"
	MYTHIC_RPC_CALLBACK_REMOVE_COMMAND            = "mythic_rpc_callback_remove_command"
	MYTHIC_RPC_CALLBACK_SEARCH_COMMAND            = "mythic_rpc_callback_search_command"
	MYTHIC_RPC_CALLBACK_UPDATE                    = "mythic_rpc_callback_update"
	MYTHIC_RPC_CALLBACK_ENCRYPT_BYTES             = "mythic_rpc_callback_encrypt_bytes"
	MYTHIC_RPC_CALLBACK_DECRYPT_BYTES             = "mythic_rpc_callback_decrypt_bytes"
	// MYTHIC_RPC_AGENTSTORAGE_CREATE agent storage operations
	MYTHIC_RPC_AGENTSTORAGE_CREATE = "mythic_rpc_agentstorage_create"
	MYTHIC_RPC_AGENTSTORAGE_SEARCH = "mythic_rpc_agentstorage_search"
	MYTHIC_RPC_AGENTSTORAGE_REMOVE = "mythic_rpc_agentstorage_remove"
	// MYTHIC_RPC_PROCESS_CREATE process operations
	MYTHIC_RPC_PROCESS_CREATE = "mythic_rpc_process_create"
	MYTHIC_RPC_PROCESS_SEARCH = "mythic_rpc_process_search"
	// MYTHIC_RPC_ARTIFACT_CREATE artifact operations
	MYTHIC_RPC_ARTIFACT_CREATE = "mythic_rpc_artifact_create"
	MYTHIC_RPC_ARTIFACT_SEARCH = "mythic_rpc_artifact_search"
	// MYTHIC_RPC_KEYLOG_CREATE keylog operations
	MYTHIC_RPC_KEYLOG_CREATE = "mythic_rpc_keylog_create"
	MYTHIC_RPC_KEYLOG_SEARCH = "mythic_rpc_keylog_search"
	// MYTHIC_RPC_CREDENTIAL_CREATE credential operations
	MYTHIC_RPC_CREDENTIAL_CREATE = "mythic_rpc_credential_create"
	MYTHIC_RPC_CREDENTIAL_SEARCH = "mythic_rpc_credential_search"
	// MYTHIC_RPC_EVENTLOG_CREATE event log operations
	MYTHIC_RPC_EVENTLOG_CREATE = "mythic_rpc_eventlog_create"
	// MYTHIC_RPC_FILEBROWSER_CREATE filebrowser operations
	MYTHIC_RPC_FILEBROWSER_CREATE = "mythic_rpc_filebrowser_create"
	MYTHIC_RPC_FILEBROWSER_REMOVE = "mythic_rpc_filebrowser_remove"
	// MYTHIC_RPC_PAYLOADONHOST_CREATE payload on host operations
	MYTHIC_RPC_PAYLOADONHOST_CREATE = "mythic_rpc_payloadonhost_create"
	// MYTHIC_RPC_CALLBACKTOKEN_CREATE callback token operations
	MYTHIC_RPC_CALLBACKTOKEN_CREATE = "mythic_rpc_callbacktoken_create"
	MYTHIC_RPC_CALLBACKTOKEN_REMOVE = "mythic_rpc_callbacktoken_remove"
	// MYTHIC_RPC_TOKEN_CREATE token operations
	MYTHIC_RPC_TOKEN_CREATE = "mythic_rpc_token_create"
	MYTHIC_RPC_TOKEN_REMOVE = "mythic_rpc_token_remove"
	// MYTHIC_RPC_PROXY_START proxy operations
	MYTHIC_RPC_PROXY_START = "mythic_rpc_proxy_start"
	MYTHIC_RPC_PROXY_STOP  = "mythic_rpc_proxy_stop"
	// MYTHIC_RPC_OTHER_SERVICES_RPC
	MYTHIC_RPC_OTHER_SERVICES_RPC = "mythic_rpc_other_service_rpc"
	// blank
	MYTHIC_RPC_BLANK = "mythic_rpc_blank"
)
