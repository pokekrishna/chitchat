package config

const (
	dbHost = "192.168.56.117"
	dbPort ="5432"
	dbName ="chitchat"
	dbUser = "postgres"
	dbPassword = "password"
	logLevel = 3
)


func GetDbHost() string {
	return dbHost
}

func GetDbPort() string {
	return dbPort
}

func GetDbName() string {
	return dbName
}

func GetDbUser() string {
	return dbUser
}

func GetDbPassword() string {
	return dbPassword
}

func GetLogLevel() int {
	return logLevel
}