package conf

// AppConfig is application config
type AppConfig struct {
	Server   *ServerConfig   `json:"server" yaml:"server"`
	Database *DatabaseConfig `json:"database" yaml:"database"`
	Cache    *CacheConfig    `json:"cache" yaml:"cache"`
	Logging  *LoggingConfig  `json:"logging" yaml:"logging"`
}

type ServerConfig struct {
	Port string `json:"port" yaml:"port"`
	Cert string `json:"cert" yaml:"cert"`
	Key  string `json:"cert" yaml:"key"`
	TLS  bool   `json:"tls" yaml:"tls"`
}

type DatabaseConfig struct {
	Type         string `json:"type yaml:"type"`
	Host         string `json:"host" yaml:"host"`
	Port         int    `json:"port" yaml:"port"`
	User         string `json:"user" yaml:"user"`
	Password     string `json:"password" yaml:"password"`
	DatabaseName string `json:"databaseName" yaml:"databaseName"`
	SslMode      string `json:"sslMode" yaml:"sslMode"`
	SslFactory   string `json:"sslFactory" yaml:"sslFactory"`
}

type CacheConfig struct {
	Size int // config size in bytes
}

type LoggingConfig struct {
	Level string `json:"level" yaml:"level"`
}

// SaneDefaultsLocal provides base config for testing
func SaneDefaultsLocal(databaseType string) *AppConfig {
	var config = &AppConfig{
		Server: &ServerConfig{
			Port: "8090",
			Cert: "certs/cert.crt",
			Key:  "certs/cert.key",
			TLS:  false,
		},
		Database: &DatabaseConfig{
			Type:         databaseType,
			Host:         "127.0.0.1",
			Port:         9042,
			User:         "Username",
			Password:     "Password",
			DatabaseName: "testspace",
			SslMode:      "disable",
			SslFactory:   "org.postgresql.ssl.NonValidatingFactory",
		},
		Cache: &CacheConfig{
			Size: 1000 * 1000,
		},
		Logging: &LoggingConfig{
			Level: "WARN",
		},
	}

	return config
}
