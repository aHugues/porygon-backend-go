package config

import "fmt"

// DatabaseConfig represents the configuration for a SQL Database
type DatabaseConfig struct {
	Driver   string `json:"driver"`
	User     string `json:"user"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     int    `json:"int"`
	Database string `json:"database"`
}

// ToEndpoint returns the SQL endpoint from the configuration
// e.g. "username:password@tcp(127.0.0.1:3306)/test"
func (c *DatabaseConfig) ToEndpoint() string {
	return fmt.Sprintf("postgres://%s:%s@%s/%s", c.User, c.Password, c.Host, c.Database)
}
