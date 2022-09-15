package conf

import (
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	Kafka *KafkaConfig `mapstructure:"kafka"`
	Tail  *TailConfig  `mapstructure:"tail"`
	Etcd  *EtcdConfig  `mapstructure:"etcd"`
}
type KafkaConfig struct {
	Address  []string `mapstructure:"address"`
	Topic    string   `mapstructure:"topic"`
	ChanSize int      `mapstructure:"chan_size"`
}
type TailConfig struct {
	CollectPath string `mapstructure:"collect_path"`
}
type EtcdConfig struct {
	Address    []string `mapstructure:"address"`
	CollectKey string   `mapstructure:"collect_key"`
}

var conf = new(Config)

func Init() (err error) {
	viper.SetConfigFile("./conf/config.yml")
	err = viper.ReadInConfig()
	if err != nil {
		log.Println("read config error", err)
		return err
	}
	err = viper.Unmarshal(conf)
	if err != nil {
		log.Println("unmarshal conf error", err)
		return err
	}
	return nil
}
func GetConf() *Config {
	return conf
}
