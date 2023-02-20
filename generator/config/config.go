package config

type Data struct {
	Remote           bool   `env:"REMOTE"`
	InputFolder      string `env:"INPUT_FOLDER"`
	OutputFolder     string `env:"OUTPUT_FOLDER"`
	GoogleAnalytics  string `env:"GA"`
	BaseURL          string `env:"BASE_URL"`
	OutputPostFolder string `env:"POST_FOLDER"`

	RemoteAddr string `env:"ADDR"`
	RemotePort string `env:"PORT"`
	User       string `env:"USER"`
	Password   string `env:"PASSWORD"`
	KeyFile    string `env:"KEY_FILE"`
	KeyStr     string `env:"KEY_STR"`

	BuildAllPosts bool `env:"BUILD_ALL" envDefault:"true"`
	Clean         bool // todo: impl clean dir
	Backup        bool
	BackupDir     string
}
