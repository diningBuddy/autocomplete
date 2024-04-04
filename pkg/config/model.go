package config

type Properties struct {
	Addr                    string `yaml:"addr"`
	SearchAutocompleteRedis struct {
		Addr     string `yaml:"addr"`
		Password string `yaml:"password"`
	} `yaml:"searchAutocompleteRedis"`
}
