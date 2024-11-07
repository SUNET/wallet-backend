package model

import (
	"vc/pkg/logger"
)

// APIServer holds the api server configuration
type APIServer struct {
	Addr       string            `yaml:"addr" validate:"required"`
	PublicKeys map[string]string `yaml:"public_keys"`
	TLS        TLS               `yaml:"tls" validate:"omitempty"`
	BasicAuth  BasicAuth         `yaml:"basic_auth"`
}

// TLS holds the tls configuration
type TLS struct {
	Enabled      bool   `yaml:"enabled"`
	CertFilePath string `yaml:"cert_file_path" validate:"required"`
	KeyFilePath  string `yaml:"key_file_path" validate:"required"`
}

// Mongo holds the database configuration
type Mongo struct {
	URI string `yaml:"uri" validate:"required"`
}

// Kafka holds the kafka configuration that is common for the entire system
type Kafka struct {
	Enabled bool     `yaml:"enabled"`
	Brokers []string `yaml:"brokers" validate:"required"`
}

// KeyValue holds the key/value configuration
type KeyValue struct {
	Addr     string `yaml:"addr" validate:"required"`
	DB       int    `yaml:"db" validate:"required"`
	Password string `yaml:"password" validate:"required"`
	PDF      PDF    `yaml:"pdf" validate:"required"`
}

// Log holds the log configuration
type Log struct {
	Level      string `yaml:"level"`
	FolderPath string `yaml:"folder_path"`
}

// Common holds the common configuration
type Common struct {
	HTTPProxy  string   `yaml:"http_proxy"`
	Production bool     `yaml:"production"`
	Log        Log      `yaml:"log"`
	Mongo      Mongo    `yaml:"mongo" validate:"required"`
	Tracing    OTEL     `yaml:"tracing" validate:"required"`
	Queues     Queues   `yaml:"queues" validate:"omitempty"`
	KeyValue   KeyValue `yaml:"key_value" validate:"omitempty"`
	QR         QRCfg    `yaml:"qr" validate:"omitempty"`
	Kafka      Kafka    `yaml:"kafka" validate:"omitempty"`
}

// SMT Spares Merkel Tree configuration
type SMT struct {
	UpdatePeriodicity int    `yaml:"update_periodicity" validate:"required"`
	InitLeaf          string `yaml:"init_leaf" validate:"required"`
}

// GRPCServer holds the rpc configuration
type GRPCServer struct {
	Addr     string `yaml:"addr" validate:"required"`
	Insecure bool   `yaml:"insecure"`
}

// PDF holds the pdf configuration (special Ladok case)
type PDF struct {
	KeepSignedDuration   int `yaml:"keep_signed_duration"`
	KeepUnsignedDuration int `yaml:"keep_unsigned_duration"`
}

// QRCfg holds the qr configuration
type QRCfg struct {
	BaseURL       string `yaml:"base_url" validate:"required"`
	RecoveryLevel int    `yaml:"recovery_level" validate:"required,min=0,max=3"`
	Size          int    `yaml:"size" validate:"required"`
}

// Queues have the queue configuration
type Queues struct {
	SimpleQueue struct {
		VCPersistentSave struct {
			Name string `yaml:"name" validate:"required"`
		} `yaml:"vc_persistent_save" validate:"required"`
		VCPersistentGet struct {
			Name string `yaml:"name" validate:"required"`
		} `yaml:"vc_persistent_get" validate:"required"`
		VCPersistentDelete struct {
			Name string `yaml:"name" validate:"required"`
		} `yaml:"vc_persistent_delete" validate:"required"`
		VCPersistentReplace struct {
			Name string `yaml:"name" validate:"required"`
		} `yaml:"vc_persistent_replace" validate:"required"`
	} `yaml:"simple_queue" validate:"required"`
}

// JWTAttribute holds the jwt attribute configuration.
// In a later state this should be placed under authentic source in order to issue credentials based on that configuration.
type JWTAttribute struct {
	// Issuer of the token example: https://issuer.sunet.se
	Issuer string `yaml:"issuer" validate:"required"`

	// EnableNotBefore states the time not before which the token is valid
	EnableNotBefore bool `yaml:"enable_not_before"`

	// Valid duration of the token in seconds
	ValidDuration int64 `yaml:"valid_duration" validate:"required_with=EnableNotBefore"`

	// VerifiableCredentialType URL example: https://credential.sunet.se/identity_credential
	VerifiableCredentialType string `yaml:"verifiable_credential_type" validate:"required"`

	// Status status of the Verifiable Credential
	Status string `yaml:"status"`

	// Kid key id of the signing key
	Kid string `yaml:"kid"`
}

// Issuer holds the issuer configuration
type Issuer struct {
	APIServer      APIServer    `yaml:"api_server" validate:"required"`
	Identifier     string       `yaml:"identifier" validate:"required"`
	GRPCServer     GRPCServer   `yaml:"grpc_server" validate:"required"`
	SigningKeyPath string       `yaml:"signing_key_path" validate:"required"`
	JWTAttribute   JWTAttribute `yaml:"jwt_attribute" validate:"required"`
}

// Registry holds the registry configuration
type Registry struct {
	APIServer  APIServer  `yaml:"api_server" validate:"required"`
	SMT        SMT        `yaml:"smt" validate:"required"`
	GRPCServer GRPCServer `yaml:"grpc_server" validate:"required"`
}

// Persistent holds the persistent storage configuration
type Persistent struct {
	APIServer APIServer `yaml:"api_server" validate:"required"`
}

// MockAS holds the mock as configuration
type MockAS struct {
	APIServer    APIServer `yaml:"api_server" validate:"required"`
	DatastoreURL string    `yaml:"datastore_url" validate:"required"`
}

// Verifier holds the verifier configuration
type Verifier struct {
	APIServer  APIServer  `yaml:"api_server" validate:"required"`
	GRPCServer GRPCServer `yaml:"grpc_server" validate:"required"`
}

// Datastore holds the datastore configuration
type Datastore struct {
	APIServer  APIServer  `yaml:"api_server" validate:"required"`
	GRPCServer GRPCServer `yaml:"grpc_server" validate:"required"`
}

// BasicAuth holds the basic auth configuration
type BasicAuth struct {
	Users   map[string]string `yaml:"users"`
	Enabled bool              `yaml:"enabled"`
}

// APIGW holds the datastore configuration
type APIGW struct {
	APIServer APIServer `yaml:"api_server" validate:"required"`
}

// OTEL holds the opentelemetry configuration
type OTEL struct {
	Addr    string `yaml:"addr" validate:"required"`
	Type    string `yaml:"type" validate:"required"`
	Timeout int64  `yaml:"timeout" default:"10"`
}

// UI holds the user-interface configuration
type UI struct {
	APIServer                         APIServer `yaml:"api_server" validate:"required"`
	Username                          string    `yaml:"username" validate:"required"`
	Password                          string    `yaml:"password" validate:"required"`
	SessionCookieAuthenticationKey    string    `yaml:"session_cookie_authentication_key" validate:"required"`
	SessionStoreEncryptionKey         string    `yaml:"session_store_encryption_key" validate:"required"`
	SessionInactivityTimeoutInSeconds int       `yaml:"session_inactivity_timeout_in_seconds" validate:"required"`
	Services                          struct {
		APIGW struct {
			BaseURL string `yaml:"base_url"`
		} `yaml:"apigw"`
		MockAS struct {
			BaseURL string `yaml:"base_url"`
		} `yaml:"mockas"`
		Verifier struct {
			BaseURL string `yaml:"base_url"`
		} `yaml:"verifier"`
	} `yaml:"services"`
}

// CredentialType holds the configuration for the credential type
type CredentialType struct {
	Profile string `yaml:"profile" validate:"required"`
}

// NotificationEndpoint holds the configuration for the notification endpoint
type NotificationEndpoint struct {
	URL string `yaml:"url" validate:"required"`
}

// AuthenticSourceEndpoint holds the configuration for the authentic source
type AuthenticSourceEndpoint struct {
	URL string `yaml:"url" validate:"required"`
}

// SignatureServiceEndpoint holds the configuration for the signature service
type SignatureServiceEndpoint struct {
	URL string `yaml:"url" validate:"required"`
}

// RevocationServiceEndpoint holds the configuration for the revocation service
type RevocationServiceEndpoint struct {
	URL string `yaml:"url" validate:"required"`
}

// AuthenticSource holds the configuration for the authentic source
type AuthenticSource struct {
	CountryCode               string                    `yaml:"country_code" validate:"required,iso3166_1_alpha2"`
	NotificationEndpoint      NotificationEndpoint      `yaml:"notification_endpoint" validate:"required"`
	AuthenticSourceEndpoint   AuthenticSourceEndpoint   `yaml:"authentic_source_endpoint" validate:"required"`
	SignatureServiceEndpoint  SignatureServiceEndpoint  `yaml:"signature_service_endpoint" validate:"required"`
	RevocationServiceEndpoint RevocationServiceEndpoint `yaml:"revocation_service_endpoint" validate:"required"`
	CredentialTypes           map[string]CredentialType `yaml:"credential_types" validate:"required"`
}

// Cfg is the main configuration structure for this application
type Cfg struct {
	Common           Common                     `yaml:"common"`
	AuthenticSources map[string]AuthenticSource `yaml:"authentic_sources" validate:"omitempty"`
	APIGW            APIGW                      `yaml:"apigw" validate:"omitempty"`
	Issuer           Issuer                     `yaml:"issuer" validate:"omitempty"`
	Verifier         Verifier                   `yaml:"verifier" validate:"omitempty"`
	Datastore        Datastore                  `yaml:"datastore" validate:"omitempty"`
	Registry         Registry                   `yaml:"registry" validate:"omitempty"`
	Persistent       Persistent                 `yaml:"persistent" validate:"omitempty"`
	MockAS           MockAS                     `yaml:"mock_as" validate:"omitempty"`
	UI               UI                         `yaml:"ui" validate:"omitempty"`
}

func (cfg *Cfg) IsAsyncEnabled(log *logger.Log) bool {
	enabled := cfg.Common.Kafka.Enabled
	if !enabled {
		log.Info("EventPublisher disabled in config")
	}
	return enabled
}
