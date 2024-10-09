package c2structs

// C2_SAMPLE_MESSAGE STRUCTS

// C2SampleMessageMessage - Generate sample C2 Traffic based on this configuration so that the
// operator and developer can more easily troubleshoot
type C2SampleMessageMessage struct {
	C2Parameters
}

// C2SampleMessageResponse - Provide a string representation of the C2 Traffic that the corresponding
// C2SampleMessageMessage configuration would generate
type C2SampleMessageResponse struct {
	Success               bool   `json:"success"`
	Error                 string `json:"error"`
	Message               string `json:"message"`
	RestartInternalServer bool   `json:"restart_internal_server"`
}
