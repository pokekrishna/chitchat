package config

const (
	dbHost     = "192.168.56.117"
	dbPort     = "5432"
	dbName     = "chitchat"
	dbUser     = "postgres"
	dbPassword = "password"
	logLevel   = 3
)

func DbHost() string {
	return dbHost
}

func DbPort() string {
	return dbPort
}

func DbName() string {
	return dbName
}

func DbUser() string {
	return dbUser
}

func DbPassword() string {
	return dbPassword
}

func LogLevel() int {
	return logLevel
}
