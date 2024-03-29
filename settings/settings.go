package settings

import (
	"fmt"

	"github.com/fsnotify/fsnotify"

	"github.com/spf13/viper"
)

var Conf = new(AppConfig)

type MySQLConfig struct {
	Host           string `mapstructure:"host"`
	Port           int    `mapstructure:"port"`
	User           string `mapstructure:"user"`
	Password       string `mapstructure:"password"`
	DbName         string `mapstructure:"dbname"`
	Max_Open_Conns int    `mapstructure:"max_open_conns"`
	Max_Idle_Conns int    `mapstructure:"max_idle_conns"`
}

type RedisConfig struct {
	Host      string `mapstructure:"host"`
	Port      int    `mapstructure:"port"`
	Password  string `mapstructure:"password"`
	Pool_Size int    `mapstructure:"pool_size"`
	DB        int    `mapstructure:"db"`
}

type AppConfig struct {
	Host        string       `mapstructure:"host"`
	Name        string       `mapstructure:"name"`
	Mode        string       `mapstructure:"mode"`
	Port        int          `mapstructure:"port"`
	StartTime   string       `mapstructure:"start_time"`
	MachineID   int64        `mapstructure:"machineID"`
	MySQLConfig *MySQLConfig `mapstructure:"mysql"`
	LogConfig   *LogConfig   `mapstructure:"log"`
	RedisConfig *RedisConfig `mapstructure:"redis"`
}

type LogConfig struct {
	Level       string `mapstructure:"level"`
	FileName    string `mapstructure:"filename"`
	Max_Size    int    `mapstructure:"max_size"`
	Max_Age     int    `mapstructure:"max_age"`
	Max_Backups int    `mapstructure:"max_backups"`
}

func Init() (err error) {
	viper.SetConfigName("./conf/config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err = viper.ReadInConfig()
	if err != nil {
		fmt.Printf("viper.ReadInConfig failed, err: %v\n", err)
		return
	}
	if err = viper.Unmarshal(&Conf); err != nil {
		fmt.Printf("viper.ReadInConfig failed, err: %v\n", err)
		return
	}
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("配置文件修改了...正在热加载...")
		if err = viper.Unmarshal(&Conf); err != nil {
			fmt.Printf("viper.ReadInConfig failed, err: %v\n", err)
			return
		}
	})
	return
}
