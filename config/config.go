package config

// ZapLoggerConfig 日志配置
type (
	// ZapLoggerConfig 日志配置
	ZapLoggerConfig struct {
		Level       string `yaml:"level" json:"level"`       //日志等级 debug->info->warn->error
		FileName    string `yaml:"fileName" json:"fileName"` //文件名
		FilePath    string `yaml:"filePath" json:"filePath"`
		MaxSize     int64  `yaml:"maxSize" json:"maxSize"`       //单位为MB
		MaxBackups  int64  `yaml:"maxBackups" json:"maxBackups"` //最多保留日志备份
		MaxAge      int64  `yaml:"maxAge" json:"maxAge"`         //备份最大生命周期 0为长期保存, 单位:天
		Compress    bool   `yaml:"compress" json:"compress"`     //是否压缩
		ShowConsole bool   `yaml:"showConsole" json:"showConsole"`
	}
	// nacos 配置中心文件
	NacosConfig struct {
		NamespaceID         string      `yaml:"namespaceId"`
		TimeoutMs           uint64      `yaml:"timeoutMs"`
		NotLoadCacheAtStart bool        `yaml:"notLoadCacheAtStart"`
		LogDir              string      `yaml:"logDir"`
		CacheDir            string      `yaml:"cacheDir"`
		LogLevel            string      `yaml:"logLevel"`
		RotateTime          string      `yaml:"rotateTime"`
		MaxAge              int64       `yaml:"maxAge"`
		NacosServer         nacosServer `yaml:"nacosServer"`
	}
	nacosServer struct {
		IP   string `yaml:"ip"`
		Port uint64 `yaml:"port"`
	}
)
