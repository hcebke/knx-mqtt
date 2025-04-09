package models

import "time"

// Config represents the top-level structure of the YAML configuration.
type Config struct {
	LogLevel                    string              `yaml:"loglevel"`
	OutgoingMqttMessage         OutgoingMqttMessage `yaml:"outgoingMqttMessage"`
	IgnoreUnknownGroupAddresses bool                `yaml:"ignoreUnknownGroupAddresses"`
	KNX                         KNXConfig           `yaml:"knx"`
	MQTT                        MQTTConfig          `yaml:"mqtt"`
}

const ValueType = "value"
const ValueWithUnitType = "value-with-unit"
const BytesType = "bytes"
const JsonType = "json"

type OutgoingMqttMessage struct {
	Type               string             `yaml:"type"`
	EmitUsingAddress   bool               `yaml:"emitUsingAddress"`
	EmitUsingName      bool               `yaml:"emitUsingName"`
	EmitValueAsString  bool               `yaml:"emitValueAsString"`
	IncludedJsonFields IncludedJsonFields `yaml:"includedJsonFields"`
}

type IncludedJsonFields struct {
	IncludeBytes bool `yaml:"bytes"`
	IncludeName  bool `yaml:"name"`
	IncludeValue bool `yaml:"value"`
	IncludeUnit  bool `yaml:"unit"`
}

// KNXLogConfig represents the KNX message logging configuration.
type KNXLogConfig struct {
	Enabled  bool          `yaml:"enabled"`
	Format   string        `yaml:"format"`
	File     string        `yaml:"file"`
	MaxSize  int           `yaml:"maxSize"`  // Maximum size in KiB before rotation
	MaxAge   time.Duration `yaml:"maxAge"`   // Maximum age before deletion
	MaxFiles int           `yaml:"maxFiles"` // Maximum number of files to keep
	Compress bool          `yaml:"compress"` // Whether to compress rotated files
}

// KNXConfig represents the KNX configuration section.
type KNXConfig struct {
	ETSExport                   string                 `yaml:"etsExport"`
	Endpoint                    string                 `yaml:"endpoint"`
	TunnelMode                  bool                   `yaml:"tunnelMode"`
	IgnoreUnknownGroupAddresses bool                   `yaml:"ignoreUnknownGroupAddresses"`
	EnableLogs                  bool                   `yaml:"enableLogs"`
	GaTranslation               FlatAddressTranslation `yaml:"translateFlatGroupAddresses"`
	KNXLog                      KNXLogConfig           `yaml:"knxLog"`
}

// MQTTConfig represents the MQTT configuration section.
type MQTTConfig struct {
	URL         string  `yaml:"url"`
	ClientID    *string `yaml:"clientId"`
	Username    *string `yaml:"username,omitempty"`
	Password    *string `yaml:"password,omitempty"`
	TLSKey      *string `yaml:"tlsKey,omitempty"`
	TLSCert     *string `yaml:"tlsCert,omitempty"`
	TLSCA       *string `yaml:"tlsCa,omitempty"`
	TopicPrefix string  `yaml:"topicPrefix"`
	Qos         byte    `yaml:"qos"`
	Retain      bool    `yaml:"retain"`
}
