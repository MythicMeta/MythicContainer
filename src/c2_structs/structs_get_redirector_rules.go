package c2structs

// C2_REDIRECTOR_RULES STRUCTS

type C2_GET_REDIRECTOR_RULE_STATUS = string

type C2GetRedirectorRuleMessage struct {
	Name       string                 `json:"c2_profile_name"`
	Parameters map[string]interface{} `json:"parameters"`
}

type C2GetRedirectorRuleMessageResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Message string `json:"message"`
}
