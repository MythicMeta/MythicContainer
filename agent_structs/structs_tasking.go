package agentstructs

type FileBrowserTask struct {
	Path     string `json:"path" mapstructure:"path"`
	FullPath string `json:"full_path" mapstructure:"full_path"`
	Filename string `json:"file" mapstructure:"file"`
	Host     string `json:"host" mapstructure:"host"`
}
