package postgres

import (
	"fmt"

	"github.com/rogeecn/atom/container"
	"github.com/rogeecn/atom/utils/opt"
)

const DefaultPrefix = "Postgres"

func DefaultProvider() container.ProviderContainer {
	return container.ProviderContainer{
		Provider: Provide,
		Options: []opt.Option{
			opt.Prefix(DefaultPrefix),
		},
	}
}

type Config struct {
	Username     string
	Password     string
	Database     string
	Host         string
	Port         uint
	SslMode      string
	TimeZone     string
	Prefix       string // 表前缀
	Singular     bool   // 是否开启全局禁用复数，true表示开启
	MaxIdleConns int    // 空闲中的最大连接数
	MaxOpenConns int    // 打开到数据库的最大连接数
}

func (m *Config) EmptyDsn() string {
	dsnTpl := "host=%s user=%s password=%s port=%d dbname=%s sslmode=%s TimeZone=%s"
	if m.SslMode == "" {
		m.SslMode = "disable"
	}

	if m.TimeZone == "" {
		m.TimeZone = "Asia/Shanghai"
	}

	return fmt.Sprintf(dsnTpl, m.Host, m.Username, m.Password, m.Port, m.Database, m.SslMode, m.TimeZone)
}

// DSN connection dsn
func (m *Config) DSN() string {
	dsnTpl := "host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s"

	if m.SslMode == "" {
		m.SslMode = "disable"
	}

	if m.TimeZone == "" {
		m.TimeZone = "Asia/Shanghai"
	}
	return fmt.Sprintf(dsnTpl, m.Host, m.Username, m.Password, m.Database, m.Port, m.SslMode, m.TimeZone)
}
