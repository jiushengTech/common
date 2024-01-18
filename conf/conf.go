package conf

type Config struct {
	ZapConf *ZapConf `yaml:"zapConf"`
}

type ZapConf struct {
	Level         string `yaml:"level"`
	Format        string `yaml:"format"`
	Director      string `yaml:"director"`
	EncodeLevel   string `yaml:"encodeLevel"`
	StacktraceKey string `yaml:"stacktraceKey"`
	MaxAge        int32  `yaml:"maxAge"`
	ShowLine      bool   `yaml:"showLine"`
	LogInConsole  bool   `yaml:"logInConsole"`
	MaxSize       int32  `yaml:"maxSize"`
	Compress      bool   `yaml:"compress"`
	MaxBackups    int32  `yaml:"maxBackups"`
}
