package c2structs

// C2_GET_IOC STRUCTS

// C2GetIOCMessage given the following C2 configuration, determine the IOCs that a defender should look for
type C2GetIOCMessage struct {
	C2Parameters
}

// IOC identify the type of ioc with Type and the actual IOC value
// An example could be a Type of URL with the actual IOC value being the configured callback URL with URI parameters
type IOC struct {
	Type string `json:"type" mapstructure:"type"`
	IOC  string `json:"ioc" mapstructure:"ioc"`
}

// C2GetIOCMessageResponse the resulting set of IOCs that a defender should look out for based on the
// C2GetIOCMessage configuration
type C2GetIOCMessageResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	IOCs    []IOC  `json:"iocs"`
}
