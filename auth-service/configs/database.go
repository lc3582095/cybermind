package configs

type DatabaseConfig struct {
    Host     string
    Port     int
    User     string
    Password string
    DBName   string
    SSLMode  string
}

func GetDatabaseConfig() *DatabaseConfig {
    return &DatabaseConfig{
        Host:     "dbconn.sealosbja.site",
        Port:     37550,
        User:     "postgres",
        Password: "wkzhx7jn",
        DBName:   "postgres",
        SSLMode:  "disable",
    }
} 