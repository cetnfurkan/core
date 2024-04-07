package config

type Database struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	Extra    map[string]any
}
