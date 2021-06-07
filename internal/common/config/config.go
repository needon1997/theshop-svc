package config

var ServerConfig = SeverConfiguration{}

type SeverConfiguration struct {
	Host          string                   `mapstructure:"host"`
	Port          int                      `mapstructure:"port"`
	Name          string                   `mapstructure:"name"`
	LogConfig     LogConfiguration         `mapstructure:"log"`
	RedisConfig   RedisConfiguration       `mapstructure:"redis"`
	EmailConfig   EmailConfiguration       `mapstructure:"email"`
	ConsulConfig  ConsulConfiguration      `mapstructure:"consul"`
	MySqlConfig   MySqlConfiguration       `mapstructure:"mysql"`
	ServiceConfig ServiceNameConfiguration `mapstructure:"svc_name"`
}
type LogConfiguration struct {
	LogPath string `mapstructure:"log_path"`
}
type RedisConfiguration struct {
	Host       string `mapstructure:"host"`
	Port       string `mapstructure:"port"`
	Expiration int64  `mapstructure:"expiration"`
}
type EmailConfiguration struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Mime     string `mapstructure:"mime"`
}
type ConsulConfiguration struct {
	Name  string                   `mapstructure:"name"`
	Host  string                   `mapstructure:"host"`
	Port  int                      `mapstructure:"port"`
	Id    string                   `mapstructure:"id"`
	Tags  []string                 `mapstructure:"tags"`
	Url   string                   `mapstructure:"url"`
	Check ConsulCheckConfiguration `mapstructure:"check"`
}
type ConsulCheckConfiguration struct {
	CheckMethod string `mapstructure:"check_method"`
	Method      string `mapstructure:"method"`
	Interval    string `mapstructure:"interval"`
	Uri         string `mapstructure:"uri"`
}

type MySqlConfiguration struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Password string `mapstructure:"password"`
	Username string `mapstructure:"username"`
	RestURI  string `mapstructure:"rest_uri"`
}
type ServiceNameConfiguration struct {
	UserServiceName      string `mapstructure:"user_svc_name"`
	EmailServiceName     string `mapstructure:"email_svc_name"`
	GoodsServiceName     string `mapstructure:"goods_svc_name"`
	InventoryServiceName string `mapstructure:"inventory_svc_name"`
	OrderServiceName     string `mapstructure:"order_svc_name"`
}
