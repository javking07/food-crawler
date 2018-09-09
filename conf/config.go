package conf

// Config is application config
type AppConfig struct {
	Server struct {
		Port string `json:"port"`
		Cert string `json:"cert"`
	} `json:"server"`
	Logging struct {
		Level string `json:"level"`
	} `jon:"logging"`
}
