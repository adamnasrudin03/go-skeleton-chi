package configs

import "time"

type Configs struct {
	App   AppConfig
	DB    DbConfig
	Redis RedisConfig
}

type AppConfig struct {
	Name          string `json:"name"`
	Env           string `json:"env"`
	Port          string `json:"port"`
	BasicUsername string `json:"basic_username"`
	BasicPassword string `json:"basic_password"`
}

type DbConfig struct {
	Host        string `json:"host"`
	Port        string `json:"port"`
	DbName      string `json:"db_name"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	DbIsMigrate bool   `json:"db_is_migrate"`
	DebugMode   bool   `json:"debug_mode"`
}

type RedisConfig struct {
	Host                string        `json:"host"`
	Port                int           `json:"port"`
	Password            string        `json:"password"`
	Database            int           `json:"database"`
	Master              string        `json:"master"`
	PoolSize            int           `json:"pool_size"`
	PoolTimeout         int           `json:"pool_timeout"`
	MinIdleConn         int           `json:"min_idle_conn"`
	DefaultCacheTimeOut time.Duration `json:"default_cache_time_out"`
}
