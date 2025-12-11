package custombrowserstructs

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/MythicMeta/MythicContainer/utils/helpers"
	"github.com/MythicMeta/MythicContainer/utils/sharedStructs"
)

type CUSTOMBROWSER_TYPE string

const (
	CUSTOMBROWSER_TYPE_FILE = "file"
)

type CUSTOMBROWSER_TABLE_COLUMN_TYPE string

const (
	CUSTOMBROWSER_TABLE_COLUMN_TYPE_STRING = "string"
	CUSTOMBROWSER_TABLE_COLUMN_TYPE_NUMBER = "number"
	CUSTOMBROWSER_TABLE_COLUMN_TYPE_DATE   = "date"
	CUSTOMBROWSER_TABLE_COLUMN_TYPE_SIZE   = "size"
)

type CustomBrowserTableColumn struct {
	Key                string                          `json:"key"`
	Name               string                          `json:"name"`
	FillWidth          bool                            `json:"fillWidth"`
	Width              int64                           `json:"width"`
	DisableSort        bool                            `json:"disableSort"`
	DisableDoubleClick bool                            `json:"disableDoubleClick"`
	DisableFilterMenu  bool                            `json:"disableFilterMenu"`
	Type               CUSTOMBROWSER_TABLE_COLUMN_TYPE `json:"type"`
}
type CustomBrowserRowAction struct {
	Name            string `json:"name"`
	UIFeature       string `json:"ui_feature"`
	Icon            string `json:"icon"`
	Color           string `json:"color"`
	SupportsFile    bool   `json:"supports_file"`
	SupportsFolder  bool   `json:"supports_folder"`
	OpenDialog      bool   `json:"openDialog"`
	GetConfirmation bool   `json:"getConfirmation"`
}
type CustomBrowserExtraTableTaskingInput struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	DisplayName string `json:"display_name"`
	Required    bool   `json:"required"`
}

type CustomBrowserDefinition struct {
	Name                       string                                `json:"name"`
	Description                string                                `json:"description"`
	Author                     string                                `json:"author"`
	SemVer                     string                                `json:"semver"`
	Type                       CUSTOMBROWSER_TYPE                    `json:"type"`
	PathSeparator              string                                `json:"separator"`
	Columns                    []CustomBrowserTableColumn            `json:"columns"`
	DefaultVisibleColumns      []string                              `json:"default_visible_columns"`
	IndicatePartialListingInUI bool                                  `json:"indicate_partial_listing"`
	ShowCurrentPathAboveTable  bool                                  `json:"show_current_path"`
	RowActions                 []CustomBrowserRowAction              `json:"row_actions"`
	ExtraTableTaskingInputs    []CustomBrowserExtraTableTaskingInput `json:"extra_table_inputs"`

	ExportFunction           CustomBrowserExportFunction                                                               `json:"export_function"`
	OnContainerStartFunction func(sharedStructs.ContainerOnStartMessage) sharedStructs.ContainerOnStartMessageResponse `json:"-"`
}
type CustomBrowserExportFunction func(message ExportFunctionMessage) ExportFunctionMessageResponse

func (f CustomBrowserExportFunction) MarshalJSON() ([]byte, error) {
	if f != nil {
		return json.Marshal("function defined")
	} else {
		return json.Marshal("")
	}
}

type ExportFunctionMessage struct {
	TreeType         string `json:"tree_type"`
	ContainerName    string `json:"container_name"`
	Host             string `json:"host"`
	Path             string `json:"path"`
	OperationID      int    `json:"operation_id"`
	OperatorID       int    `json:"operator_id"`
	OperatorUsername string `json:"operator_username"`
	CallbackGroup    string `json:"callback_group"`
}
type ExportFunctionMessageResponse struct {
	Success           bool   `json:"success"`
	Error             string `json:"error"`
	CompletionMessage string `json:"completion_message"`
	OperationID       int    `json:"operation_id"`
	TreeType          string `json:"tree_type"`
}

// REQUIRED, Don't Modify
type allCustomBrowserData struct {
	mutex                   sync.RWMutex
	rpcMethods              []sharedStructs.RabbitmqRPCMethod
	directMethods           []sharedStructs.RabbitmqDirectMethod
	custombrowserDefinition CustomBrowserDefinition
}

var (
	AllCustomBrowserData containerCustomBrowserData
)

type containerCustomBrowserData struct {
	CustomBrowserMap map[string]*allCustomBrowserData
}

func (r *containerCustomBrowserData) GetAllNames() []string {
	names := []string{}
	for key, _ := range r.CustomBrowserMap {
		if key != "" && !helpers.StringSliceContains(names, key) {
			names = append(names, key)
		}
	}
	return names
}
func (r *containerCustomBrowserData) Get(name string) *allCustomBrowserData {
	if r.CustomBrowserMap == nil {
		r.CustomBrowserMap = make(map[string]*allCustomBrowserData)
	}
	if existingC2Data, ok := r.CustomBrowserMap[name]; !ok {
		newC2Data := allCustomBrowserData{}
		r.CustomBrowserMap[name] = &newC2Data
		return &newC2Data
	} else {
		return existingC2Data
	}
}
func (r *allCustomBrowserData) AddCustomBrowserDefinition(def CustomBrowserDefinition) {
	r.custombrowserDefinition = def
}
func (r *allCustomBrowserData) GetCustomBrowserDefinition() CustomBrowserDefinition {
	return r.custombrowserDefinition
}
func (r *allCustomBrowserData) SetName(name string) {
	r.custombrowserDefinition.Name = name
}
func (r *allCustomBrowserData) GetRoutingKey(routingKey string) string {
	return fmt.Sprintf("%s_%s", r.custombrowserDefinition.Name, routingKey)
}
func (r *allCustomBrowserData) AddRPCMethod(m sharedStructs.RabbitmqRPCMethod) {
	r.mutex.Lock()
	r.rpcMethods = append(r.rpcMethods, m)
	r.mutex.Unlock()
}
func (r *allCustomBrowserData) GetRPCMethods() []sharedStructs.RabbitmqRPCMethod {
	return r.rpcMethods
}
func (r *allCustomBrowserData) AddDirectMethod(m sharedStructs.RabbitmqDirectMethod) {
	r.mutex.Lock()
	r.directMethods = append(r.directMethods, m)
	r.mutex.Unlock()
}
func (r *allCustomBrowserData) GetDirectMethods() []sharedStructs.RabbitmqDirectMethod {
	return r.directMethods
}
