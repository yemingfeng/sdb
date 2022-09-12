package config

type Config struct {
	Debug      bool       `default:"false" yaml:"debug"`
	Store      Store      `yaml:"store"`
	Collection Collection `yaml:"collection"`
}

type Store struct {
	Path      string `default:"./db/" yaml:"path"`
	DeleteOld bool   `default:"false" yaml:"delete_old"`
}

type Collection struct {
	LockerCount int `default:"128" yaml:"locker_count"`
}

func NewDefaultConfig() *Config {
	return &Config{
		Debug: true,
		Store: Store{
			Path: "./db/",
		},
		Collection: Collection{
			LockerCount: 128,
		},
	}
}

func NewTestConfig() *Config {
	return &Config{
		Debug: true,
		Store: Store{
			Path:      "./db/",
			DeleteOld: true,
		},
		Collection: Collection{
			LockerCount: 128,
		},
	}
}
