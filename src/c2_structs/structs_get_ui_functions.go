package c2structs

type C2_GET_UI_FUNCTIONS_STATUS = string

type C2GetUiFunctionsMessage struct {
	Name string `json:"c2_profile_name"`
}

type C2GetUiFunctionsMessageResponse struct {
	Success   bool     `json:"success"`
	Error     string   `json:"error"`
	Functions []string `json:"message"`
}
