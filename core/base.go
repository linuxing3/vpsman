package core

// XrayConfig 结构体
type XrayConfig struct {
	Log       LogLevel         `json:"log"`
	Inbounds  []InBoundConfig  `json:"inbounds"`
	Outbounds []OutBoundConfig `json:"outbounds"`
}

// LogLevel 结构体
type LogLevel struct {
	LogLevel string `json: "loglevel"`
}

// InBoundConfig 结构体
type InBoundConfig struct {
	OutBoundConfig
	Port           int                        `json:"port"`
	Settings       InBoundSettingConfig       `json:"settings"`
	StreamSettings InBoundStreamSettingConfig `json:"streamSettings"`
}

// OutBoundConfig 结构体
type OutBoundConfig struct {
	Protocol string `json:"protocol"`
}

// InBoundSettingConfig 结构体
type InBoundSettingConfig struct {
	Clients    []InBoundSettingClientConfig   `json:"clients"`
	Decryption string                         `json:"decryption"`
	Fallbacks  []InBoundSettingFallbackConfig `json:"fallbacks"`
}

// InBoundSettingClientConfig 结构体
type InBoundSettingClientConfig struct {
	Id       string `json:"id"`
	Flow     string `json:"flow"`
	Password string `json:"password"`
	Level    int    `json:"level"`
	Email    string `json:"email"`
}

// InBoundSettingFallbackConfig 结构体
type InBoundSettingFallbackConfig struct {
	Path string `json:"path"`
	Xver int    `json:"xver"`
	Dest int    `json:"dest"`
}

// InBoundStreamSettingConfig 结构体
type InBoundStreamSettingConfig struct {
	Network      string            `json:"network"`
	Security     string            `json:"security"`
	XtlsSettings XtlsSettingConfig `json:"xtlsSettings"`
	TcpSettings  TcpSettingConfig  `json:"tcpSettings"`
	WsSettings   WsSettingConfig   `json:"wsSettings"`
}

// XtlsSettingConfig 结构体
type XtlsSettingConfig struct {
	Alpn         []string            `json:"alpn"`
	Certificates []CertificateConfig `json:"certificates"`
}

// TcpSettingConfig 结构体
type TcpSettingConfig struct {
	AcceptProxyProtocol bool `json:"acceptProxyProtocol"`
	Header              struct {
		Type    string `json:"type"`
		Request struct {
			Path []string `json:"path"`
		} `json:"request"`
	} `json:"header"`
}

type WsSettingConfig struct {
	Path string `json:"path"`
}

// CertificateConfig 结构体
type CertificateConfig struct {
	CertificateFile string `json:"certificateFile"`
	KeyFile         string `json:"keyFile"`
}

// TrojanConfig struct
type TrojanConfig struct {
	RunType    string   `json:"run_type"`
	LocalAddr  string   `json:"local_addr"`
	LocalPort  int      `json:"local_port"`
	RemoteAddr string   `json:"remote_addr"`
	RemotePort int      `json:"remote_port"`
	Password   []string `json:"password"`
	LogLevel   int      `json:"log_level"`
}

// SSL 结构体
type SSL struct {
	Cert          string   `json:"cert"`
	Cipher        string   `json:"cipher"`
	CipherTls13   string   `json:"cipher_tls13"`
	Alpn          []string `json:"alpn"`
	ReuseSession  bool     `json:"reuse_session"`
	SessionTicket bool     `json:"session_ticket"`
	Curves        string   `json:"curves"`
	Sni           string   `json:"sni"`
}

// TCP 结构体
type TCP struct {
	NoDelay      bool `json:"no_delay"`
	KeepAlive    bool `json:"keep_alive"`
	ReusePort    bool `json:"reuse_port"`
	FastOpen     bool `json:"fast_open"`
	FastOpenQlen int  `json:"fast_open_qlen"`
}
