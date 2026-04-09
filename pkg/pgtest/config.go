package pgtest

import "github.com/spf13/viper"

type pgTestConfig struct {
	Username           string
	Password           string
	Database           string
	MigrationDirectory string // Относительно рабочей директории
	WorkingDirectory   string

	External bool
	Host     string // только для external
	Port     uint16 // только для external

	Image string
}

const (
	DefaultUsername           = "postgres"
	DefaultPassword           = "postgres"
	DefaultDatabase           = "default"
	DefaultImage              = "postgres:17.4"
	DefaultMigrationDirectory = "./migrations"
	DefaultWorkingDirectory   = "."
	DefaultExternal           = false
	DefaultHost               = "localhost"
	DefaultPort               = 5432
)

const (
	EnvUsername           = "PGTEST_USERNAME"
	EnvPassword           = "PGTEST_PASSWORD"
	EnvDatabase           = "PGTEST_DATABASE"
	EnvImage              = "PGTEST_IMAGE"
	EnvMigrationDirectory = "PGTEST_MIGRATION_DIRECTORY"
	EnvWorkingDirectory   = "PGTEST_WORKING_DIRECTORY"
	EnvExternal           = "PGTEST_EXTERNAL"
	EnvHost               = "PGTEST_HOST"
	EnvPort               = "PGTEST_PORT"
)

func configFromEnv() pgTestConfig {
	viper.SetDefault(EnvUsername, DefaultUsername)
	viper.SetDefault(EnvPassword, DefaultPassword)
	viper.SetDefault(EnvDatabase, DefaultDatabase)
	viper.SetDefault(EnvImage, DefaultImage)
	viper.SetDefault(EnvMigrationDirectory, DefaultMigrationDirectory)
	viper.SetDefault(EnvWorkingDirectory, DefaultWorkingDirectory)
	viper.SetDefault(EnvExternal, DefaultExternal)
	viper.SetDefault(EnvHost, DefaultHost)
	viper.SetDefault(EnvPort, DefaultPort)
	viper.AutomaticEnv()

	return pgTestConfig{
		Username:           viper.GetString(EnvUsername),
		Password:           viper.GetString(EnvPassword),
		Database:           viper.GetString(EnvDatabase),
		Image:              viper.GetString(EnvImage),
		MigrationDirectory: viper.GetString(EnvMigrationDirectory),
		WorkingDirectory:   viper.GetString(EnvWorkingDirectory),
		External:           viper.GetBool(EnvExternal),
		Host:               viper.GetString(EnvHost),
		Port:               uint16(viper.GetInt(EnvPort)),
	}
}
