package loggingstructs

import (
	"fmt"
	"github.com/MythicMeta/MythicContainer/utils/helpers"
	"github.com/MythicMeta/MythicContainer/utils/sharedStructs"
	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"runtime"
	"sync"
)

const EMIT_LOG_ROUTING_KEY_PREFIX = "emit_log"

type LOG_TYPE = string

const (
	LOG_TYPE_CALLBACK   LOG_TYPE = "new_callback"
	LOG_TYPE_CREDENTIAL          = "new_credential"
	LOG_TYPE_ARTIFACT            = "new_artifact"
	LOG_TYPE_TASK                = "new_task"
	LOG_TYPE_FILE                = "new_file"
	LOG_TYPE_PAYLOAD             = "new_payload"
	LOG_TYPE_KEYLOG              = "new_keylog"
	LOG_TYPE_RESPONSE            = "new_response"
)

type loggingMessageBase struct {
	OperationID      int      `json:"operation_id"`
	OperationName    string   `json:"operation_name"`
	OperatorUsername string   `json:"username"`
	Timestamp        string   `json:"timestamp"`
	ServerName       string   `json:"server_name"`
	Action           LOG_TYPE `json:"action"`
}

type LoggingDefinition struct {
	Name        string
	Description string
	// SemVer is a specific semantic version tracker you can use for your payload type
	SemVer                   string `json:"semver"`
	LogToFilePath            string
	LogLevel                 string
	LogMaxSizeInMB           int
	LogMaxBackups            int
	NewCallbackFunction      func(input NewCallbackLog)
	NewCredentialFunction    func(input NewCredentialLog)
	NewKeylogFunction        func(input NewKeylogLog)
	NewFileFunction          func(input NewFileLog)
	NewPayloadFunction       func(input NewPayloadLog)
	NewArtifactFunction      func(input NewArtifactLog)
	NewTaskFunction          func(input NewTaskLog)
	NewResponseFunction      func(input NewResponseLog)
	Subscriptions            []string
	OnContainerStartFunction func(sharedStructs.ContainerOnStartMessage) sharedStructs.ContainerOnStartMessageResponse
}

// REQUIRED, Don't Modify
type allLoggingData struct {
	mutex             sync.RWMutex
	rpcMethods        []sharedStructs.RabbitmqRPCMethod
	directMethods     []sharedStructs.RabbitmqDirectMethod
	loggingDefinition LoggingDefinition
	logger            *zerolog.Logger
}

var (
	AllLoggingData containerLoggingData
)

type containerLoggingData struct {
	LoggingMap map[string]*allLoggingData
}

func (r *containerLoggingData) GetAllNames() []string {
	names := []string{}
	for key, _ := range r.LoggingMap {
		if key != "" && !helpers.StringSliceContains(names, key) {
			names = append(names, key)
		}
	}
	return names
}
func (r *containerLoggingData) Get(name string) *allLoggingData {
	if r.LoggingMap == nil {
		r.LoggingMap = make(map[string]*allLoggingData)
	}
	if existingC2Data, ok := r.LoggingMap[name]; !ok {
		newC2Data := allLoggingData{}
		r.LoggingMap[name] = &newC2Data
		return &newC2Data
	} else {
		return existingC2Data
	}
}
func (r *allLoggingData) AddLoggingDefinition(def LoggingDefinition) {
	r.loggingDefinition = def
	var zl zerolog.Logger
	if def.LogToFilePath != "" {
		fileLogger := &lumberjack.Logger{
			Filename:   def.LogToFilePath,
			MaxSize:    def.LogMaxSizeInMB,
			MaxAge:     0, // don't remove files after x days
			MaxBackups: def.LogMaxBackups,
			LocalTime:  false, // use UTC times
			Compress:   true,
		}
		writers := io.MultiWriter(os.Stdout, fileLogger)
		zl = zerolog.New(writers)

	} else {
		zl = zerolog.New(os.Stdout)
	}
	switch def.LogLevel {
	case "warning":
		zl = zl.Level(zerolog.WarnLevel)
	case "info":
		zl = zl.Level(zerolog.InfoLevel)
	case "debug":
		zl = zl.Level(zerolog.DebugLevel)
	case "trace":
		zl = zl.Level(zerolog.TraceLevel)
	default:
		zl = zl.Level(zerolog.InfoLevel)
	}
	zl = zl.With().Timestamp().Logger()
	r.logger = &zl
}
func (r *allLoggingData) GetLoggingDefinition() LoggingDefinition {
	return r.loggingDefinition
}
func (r *allLoggingData) AddRPCMethod(m sharedStructs.RabbitmqRPCMethod) {
	r.mutex.Lock()
	r.rpcMethods = append(r.rpcMethods, m)
	r.mutex.Unlock()
}
func (r *allLoggingData) GetRPCMethods() []sharedStructs.RabbitmqRPCMethod {
	return r.rpcMethods
}
func (r *allLoggingData) AddDirectMethod(m sharedStructs.RabbitmqDirectMethod) {
	r.mutex.Lock()
	r.directMethods = append(r.directMethods, m)
	r.mutex.Unlock()
}
func (r *allLoggingData) GetDirectMethods() []sharedStructs.RabbitmqDirectMethod {
	return r.directMethods
}
func (r *allLoggingData) SetSubscriptions(subs []string) {
	r.loggingDefinition.Subscriptions = subs
}
func (r *allLoggingData) SetName(name string) {
	r.loggingDefinition.Name = name
}
func (r *allLoggingData) GetRoutingKey(routingKey string) string {
	return fmt.Sprintf("%s_%s", r.loggingDefinition.Name, routingKey)
}
func (r *allLoggingData) LogFatalError(err error, message string, messages ...interface{}) {
	if pc, _, line, ok := runtime.Caller(1); ok {
		if err == nil {
			r.logger.Error().Fields(messages).Msg(message)
			//logger.Error(errors.New(message), "", messages...)
		} else {
			r.logger.Error().Err(err).Fields(append([]interface{}{"func", runtime.FuncForPC(pc).Name(), "line", line}, messages...)).Msg(message)
			//logger.Error(err, message, append([]interface{}{"func", runtime.FuncForPC(pc).Name(), "line", line}, messages...)...)
		}
	} else {
		if err == nil {
			r.logger.Error().Fields(messages).Msg(message)
			//logger.Error(errors.New(message), "", messages...)
		} else {
			r.logger.Error().Err(err).Fields(messages).Msg(message)
			//logger.Error(err, message, messages...)
		}
	}
	os.Exit(1)
}
func (r *allLoggingData) LogTrace(message string, messages ...interface{}) {
	if pc, _, line, ok := runtime.Caller(1); ok {
		r.logger.Trace().Fields(append([]interface{}{"func", runtime.FuncForPC(pc).Name(), "line", line}, messages...)).Msg(message)
		//logger.V(2).Info(message, append([]interface{}{"func", runtime.FuncForPC(pc).Name(), "line", line}, messages...)...)
	} else {
		r.logger.Trace().Fields(messages).Msg(message)
		//logger.V(2).Info(message, messages...)
	}
}
func (r *allLoggingData) LogDebug(message string, messages ...interface{}) {
	if pc, _, line, ok := runtime.Caller(1); ok {
		r.logger.Debug().Fields(append([]interface{}{"func", runtime.FuncForPC(pc).Name(), "line", line}, messages...)).Msg(message)
		//logger.V(1).Info(message, append([]interface{}{"func", runtime.FuncForPC(pc).Name(), "line", line}, messages...)...)
	} else {
		r.logger.Debug().Fields(messages).Msg(message)
		//logger.V(1).Info(message, messages...)
	}
}
func (r *allLoggingData) LogInfo(message string, messages ...interface{}) {
	if pc, _, line, ok := runtime.Caller(1); ok {
		r.logger.Info().Fields(append([]interface{}{"func", runtime.FuncForPC(pc).Name(), "line", line}, messages...)).Msg(message)
		//logger.V(0).Info(message, append([]interface{}{"func", runtime.FuncForPC(pc).Name(), "line", line}, messages...)...)
	} else {
		r.logger.Info().Fields(messages).Msg(message)
		//logger.V(0).Info(message, messages...)
	}
}
func (r *allLoggingData) LogWarning(message string, messages ...interface{}) {
	if pc, _, line, ok := runtime.Caller(1); ok {
		r.logger.Warn().Fields(append([]interface{}{"func", runtime.FuncForPC(pc).Name(), "line", line}, messages...)).Msg(message)
		//logger.V(1).Info(message, append([]interface{}{"func", runtime.FuncForPC(pc).Name(), "line", line}, messages...)...)
	} else {
		r.logger.Warn().Fields(messages).Msg(message)
		//logger.V(1).Info(message, messages...)
	}
}
func (r *allLoggingData) LogError(err error, message string, messages ...interface{}) {
	if pc, _, line, ok := runtime.Caller(1); ok {
		if err == nil {
			r.logger.Error().Fields(messages).Msg(message)
			//logger.Error(errors.New(message), "", messages...)
		} else {
			r.logger.Error().Err(err).Fields(append([]interface{}{"func", runtime.FuncForPC(pc).Name(), "line", line}, messages...)).Msg(message)
			//logger.Error(err, message, append([]interface{}{"func", runtime.FuncForPC(pc).Name(), "line", line}, messages...)...)
		}
	} else {
		if err == nil {
			r.logger.Error().Fields(messages).Msg(message)
			//logger.Error(errors.New(message), "", messages...)
		} else {
			r.logger.Error().Err(err).Fields(messages).Msg(message)
			//logger.Error(err, message, messages...)
		}
	}

}
func GetRoutingKeyFor(logType string) string {
	return fmt.Sprintf("%s.%s", EMIT_LOG_ROUTING_KEY_PREFIX, logType)
}
