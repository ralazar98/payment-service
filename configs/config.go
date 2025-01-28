package configs

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	App      AppConfig      `mapstructure:"app"`
	Database DatabaseConfig `mapstructure:"database"`
	RabbitMQ RabbitMQConfig `mapstructure:"rabbitmq"`
}

type AppConfig struct {
	Port string `mapstructure:"port"`
}

type DatabaseConfig struct {
	Host       string `mapstructure:"host"`
	Port       string `mapstructure:"port"` //TODO: мы можем указывать это сразу в инт, как ты проводишь в подключении к бд, это касается и других типов данных, мы можем сразу принимать time.Time или любые другие типы данных
	DBName     string `mapstructure:"dbname"`
	Username   string `ymapstructure:"username"`
	DBPassword string `mapstructure:"password"`
}

type RabbitMQConfig struct {
	NameOfQueue string `mapstructure:"nameOfQueue"`
	RabbitUrl   string `mapstructure:"rabbitUrl"`
}

func LoadConfig() (*Config, error) {
	//TODO:Перенести путь
	//TODO:Верно
	viper.AddConfigPath("/app/configs")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}
	return &config, nil

}
