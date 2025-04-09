package conf

type Config struct {
	ZapConf *ZapConf `yaml:"zapConf"`
}

type ZapConf struct {
	Model         string `yaml:"model"`         //开发环境
	Level         string `yaml:"level"`         //日志级别
	Format        string `yaml:"format"`        //日志格式 json，console
	Director      string `yaml:"director"`      //日志输出目录
	EncodeLevel   string `yaml:"encodeLevel"`   //日志输出格式
	StacktraceKey string `yaml:"stacktraceKey"` //堆栈信息key
	MaxAge        int32  `yaml:"maxAge"`        //日志最大保存时间
	AddCaller     bool   `yaml:"showLine"`      //打印调用者信息
	AddCallerSkip int32  `yaml:"addCallerSkip"` //打印调用者信息的跳过层级
	LogInConsole  bool   `yaml:"logInConsole"`  //是否输出到控制台
	MaxSize       int32  `yaml:"maxSize"`       //单个日志文件最大大小,以MB为单位
	Compress      bool   `yaml:"compress"`      //是否压缩
	MaxBackups    int32  `yaml:"maxBackups"`    //最大备份数
	TimeRotation  int32  `yaml:"timeRotation"`  //时间轮转类型: "minute", "hour" 或 "day"
}
