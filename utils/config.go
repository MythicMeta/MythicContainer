package utils

import (
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	// container configuration
	DebugLevel string
	// rabbitmq configuration
	RabbitmqHost     string
	RabbitmqPort     uint
	RabbitmqUser     string
	RabbitmqPassword string
	RabbitmqVHost    string
	// mythic information
	MythicServerHost     string
	MythicServerPort     uint
	MythicServerGRPCPort uint
	// Webhook configuration
	WebhookDefaultURL      string
	WebhookDefaultChannel  string
	WebhookFeedbackChannel string
	WebhookCallbackChannel string
	WebhookStartupChannel  string
	WebhookAlertChannel    string
	WebhookCustomChannel   string
}

var (
	MythicConfig = Config{}
)

func init() {
	mythicEnv := viper.New()
	// mythic config
	mythicEnv.SetDefault("debug_level", "staging")
	mythicEnv.SetDefault("mythic_server_grpc_port", 17444)
	mythicEnv.SetDefault("mythic_server_port", 17443)
	// rabbitmq configuration
	mythicEnv.SetDefault("rabbitmq_host", "mythic_rabbitmq")
	mythicEnv.SetDefault("rabbitmq_port", 5672)
	mythicEnv.SetDefault("rabbitmq_user", "mythic_user")
	mythicEnv.SetDefault("rabbitmq_password", "")
	mythicEnv.SetDefault("rabbitmq_vhost", "mythic_vhost")
	// webhook configuration
	mythicEnv.SetDefault("webhook_default_url", "")
	mythicEnv.SetDefault("webhook_default_channel", "")
	mythicEnv.SetDefault("webhook_default_feedback_channel", "")
	mythicEnv.SetDefault("webhook_default_callback_channel", "")
	mythicEnv.SetDefault("webhook_default_startup_channel", "")
	mythicEnv.SetDefault("webhook_default_alert_channel", "")
	mythicEnv.SetDefault("webhook_default_custom_channel", "")
	// pull in environment variables and configuration from .env if needed
	mythicEnv.SetConfigName(".env")
	mythicEnv.SetConfigType("env")
	mythicEnv.AddConfigPath(getCwdFromExe())
	mythicEnv.AutomaticEnv()
	if !fileExists(filepath.Join(getCwdFromExe(), ".env")) {
		_, err := os.Create(filepath.Join(getCwdFromExe(), ".env"))
		if err != nil {
			log.Fatalf("[-] .env doesn't exist and couldn't be created: %v", err)
		}
	}
	if err := mythicEnv.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Fatalf("[-] Error while reading in .env file: %v", err)
		} else {
			log.Fatalf("[-] Error while parsing .env file: %v", err)
		}
	}
	setConfigFromEnv(mythicEnv)
}

func setConfigFromEnv(mythicEnv *viper.Viper) {

	MythicConfig.DebugLevel = mythicEnv.GetString("debug_level")
	// rabbitmq configuration
	MythicConfig.RabbitmqHost = mythicEnv.GetString("rabbitmq_host")
	MythicConfig.RabbitmqPort = mythicEnv.GetUint("rabbitmq_port")
	MythicConfig.RabbitmqUser = mythicEnv.GetString("rabbitmq_user")
	MythicConfig.RabbitmqPassword = mythicEnv.GetString("rabbitmq_password")
	MythicConfig.RabbitmqVHost = mythicEnv.GetString("rabbitmq_vhost")
	// mythic information
	MythicConfig.MythicServerPort = mythicEnv.GetUint("mythic_server_port")
	MythicConfig.MythicServerGRPCPort = mythicEnv.GetUint("mythic_server_grpc_port")
	MythicConfig.MythicServerHost = mythicEnv.GetString("mythic_server_host")
	if MythicConfig.MythicServerHost == "" {
		log.Fatalf("[-] Missing MYTHIC_SERVER_HOST environment variable point to mythic server IP")
	}
	// webhook configuration
	MythicConfig.WebhookDefaultURL = mythicEnv.GetString("webhook_default_url")
	MythicConfig.WebhookDefaultChannel = mythicEnv.GetString("webhook_default_channel")
	MythicConfig.WebhookFeedbackChannel = mythicEnv.GetString("webhook_default_feedback_channel")
	MythicConfig.WebhookCallbackChannel = mythicEnv.GetString("webhook_default_callback_channel")
	MythicConfig.WebhookStartupChannel = mythicEnv.GetString("webhook_default_startup_channel")
	MythicConfig.WebhookAlertChannel = mythicEnv.GetString("webhook_default_alert_channel")
	MythicConfig.WebhookCustomChannel = mythicEnv.GetString("webhook_default_custom_channel")
}

func getCwdFromExe() string {
	exe, err := os.Executable()
	if err != nil {
		log.Fatalf("[-] Failed to get path to current executable: %v", err)
	}
	return filepath.Dir(exe)
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return !info.IsDir()
}
