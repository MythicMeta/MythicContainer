package mythicrpc

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

type PTTaskMessageTaskData struct {
	ID                                 int                       `json:"id"`
	DisplayID                          int                       `json:"display_id"`
	AgentTaskID                        string                    `json:"agent_task_id"`
	CommandName                        string                    `json:"command_name"`
	Params                             string                    `json:"params"`
	Timestamp                          string                    `json:"timestamp"`
	CallbackID                         int                       `json:"callback_id"`
	Status                             string                    `json:"status"`
	OriginalParams                     string                    `json:"original_params"`
	DisplayParams                      string                    `json:"display_params"`
	KeywordResolution                  []PTTaskKeywordResolution `json:"keyword_resolution"`
	Comment                            string                    `json:"comment"`
	Stdout                             string                    `json:"stdout"`
	Stderr                             string                    `json:"stderr"`
	Completed                          bool                      `json:"completed"`
	OpsecPreBlocked                    bool                      `json:"opsec_pre_blocked"`
	OpsecPreMessage                    string                    `json:"opsec_pre_message"`
	OpsecPreBypassed                   bool                      `json:"opsec_pre_bypassed"`
	OpsecPreBypassRole                 string                    `json:"opsec_pre_bypass_role"`
	OpsecPostBlocked                   bool                      `json:"opsec_post_blocked"`
	OpsecPostMessage                   string                    `json:"opsec_post_message"`
	OpsecPostBypassed                  bool                      `json:"opsec_post_bypassed"`
	OpsecPostBypassRole                string                    `json:"opsec_post_bypass_role"`
	ParentTaskID                       int                       `json:"parent_task_id"`
	OperatorUsername                   string                    `json:"operator_username"`
	SubtaskCallbackFunction            string                    `json:"subtask_callback_function"`
	SubtaskCallbackFunctionCompleted   bool                      `json:"subtask_callback_function_completed"`
	GroupCallbackFunction              string                    `json:"group_callback_function"`
	GroupCallbackFunctionCompleted     bool                      `json:"group_callback_function_completed"`
	CompletedCallbackFunction          string                    `json:"completed_callback_function"`
	CompletedCallbackFunctionCompleted bool                      `json:"completed_callback_function_completed"`
	SubtaskGroupName                   string                    `json:"subtask_group_name"`
	TaskingLocation                    string                    `json:"tasking_location"`
	ParameterGroupName                 string                    `json:"parameter_group_name"`
	TokenID                            int                       `json:"token_id"`
}

type PTTaskKeywordResolution struct {
	Raw            string   `json:"raw"`
	Keyword        string   `json:"keyword"`
	Selector       string   `json:"selector"`
	Field          string   `json:"field"`
	ValueType      string   `json:"value_type"`
	ExpandedValue  string   `json:"expanded_value"`
	ParameterNames []string `json:"parameter_names"`
}

func (t PTTaskMessageTaskData) RevertKeywords(parameter interface{}, parameterName ...string) string {
	name := ""
	if len(parameterName) > 0 {
		name = parameterName[0]
	}
	return revertKeywords(parameter, t.KeywordResolution, name)
}

func revertKeywords(parameter interface{}, keywordResolution []PTTaskKeywordResolution, parameterName string) string {
	reverted := keywordParameterToString(parameter)
	if len(keywordResolution) == 0 {
		return reverted
	}
	if parameterName != "" {
		for _, entry := range keywordResolution {
			if entry.ValueType == "structured" && keywordResolutionContainsParameter(entry.ParameterNames, parameterName) {
				return entry.Raw
			}
		}
		specificEntries := filterKeywordResolutionEntries(keywordResolution, parameterName, false)
		updated := applyKeywordResolutionStringEntries(reverted, specificEntries)
		if updated != reverted {
			return updated
		}
		return applyKeywordResolutionStringEntries(reverted, filterKeywordResolutionEntries(keywordResolution, "", true))
	}
	return applyKeywordResolutionStringEntries(reverted, keywordResolution)
}

func keywordParameterToString(parameter interface{}) string {
	switch typed := parameter.(type) {
	case string:
		return typed
	case []byte:
		return string(typed)
	default:
		marshaled, err := json.Marshal(typed)
		if err != nil {
			return fmt.Sprint(typed)
		}
		return string(marshaled)
	}
}

func filterKeywordResolutionEntries(keywordResolution []PTTaskKeywordResolution, parameterName string, globalsOnly bool) []PTTaskKeywordResolution {
	filtered := make([]PTTaskKeywordResolution, 0, len(keywordResolution))
	for _, entry := range keywordResolution {
		if entry.ValueType != "string" {
			continue
		}
		if globalsOnly {
			if len(entry.ParameterNames) == 0 {
				filtered = append(filtered, entry)
			}
			continue
		}
		if keywordResolutionContainsParameter(entry.ParameterNames, parameterName) {
			filtered = append(filtered, entry)
		}
	}
	return filtered
}

func applyKeywordResolutionStringEntries(parameter string, keywordResolution []PTTaskKeywordResolution) string {
	sort.SliceStable(keywordResolution, func(i, j int) bool {
		return len(keywordResolution[i].ExpandedValue) > len(keywordResolution[j].ExpandedValue)
	})
	reverted := parameter
	placeholders := make(map[string]string)
	for i, entry := range keywordResolution {
		if entry.ValueType != "string" || entry.ExpandedValue == "" {
			continue
		}
		placeholder := fmt.Sprintf("\x00MYTHIC_KEYWORD_%d\x00", i)
		placeholders[placeholder] = entry.Raw
		reverted = strings.ReplaceAll(reverted, entry.ExpandedValue, placeholder)
	}
	for placeholder, raw := range placeholders {
		reverted = strings.ReplaceAll(reverted, placeholder, raw)
	}
	return reverted
}

func keywordResolutionContainsParameter(parameterNames []string, parameterName string) bool {
	for _, existing := range parameterNames {
		if existing == parameterName {
			return true
		}
	}
	return false
}

// exporting a payload configuration
type PayloadConfiguration struct {
	Description        string                                `json:"description"`
	PayloadType        string                                `json:"payload_type" binding:"required"`
	C2Profiles         *[]PayloadConfigurationC2Profile      `json:"c2_profiles,omitempty"`
	BuildParameters    *[]PayloadConfigurationBuildParameter `json:"build_parameters,omitempty"`
	Commands           []string                              `json:"commands"`
	SelectedOS         string                                `json:"selected_os" binding:"required"`
	Filename           string                                `json:"filename" binding:"required"`
	WrappedPayloadUUID string                                `json:"wrapped_payload"`
	UUID               string                                `json:"uuid"`
	AgentFileID        string                                `json:"agent_file_id"`
	BuildPhase         string                                `json:"build_phase"`
}
type PayloadConfigurationC2Profile struct {
	Name       string                 `json:"c2_profile"`
	Parameters map[string]interface{} `json:"c2_profile_parameters"`
}
type PayloadConfigurationBuildParameter struct {
	Name  string      `json:"name" binding:"required"`
	Value interface{} `json:"value" binding:"required"`
}
