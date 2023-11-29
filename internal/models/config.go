package models

type ConfigDB struct {
	Path     string
	DbType   string
	Host     string
	Port     string
	DBName   string
	User     string
	Password string
}

type ConfigApp struct {
	Port      string
	JwtSecret string
}
