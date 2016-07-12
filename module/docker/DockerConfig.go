package docker



type TlsConfig struct {
	Enable   bool   `config:"enable"`
	CaPath   string `config:"ca_path"`
	CertPath string `config:"cert_path"`
	KeyPath  string `config:"key_path"`
}

type Conf struct {
	Socket string  `config:"socket"`
	Tls    TlsConfig
}

func GetDefaultConf() Conf {
	return Conf{
		Socket: "unix:///var/run/docker.sock",
		Tls : TlsConfig{
			Enable: false,
		},
	}
}