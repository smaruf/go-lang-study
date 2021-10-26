package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

var (
	HostConfig = NewHostAddressConfig().Init()
)

type HostAddressConfig struct {
	Port                     string `yaml:"port"`
	RemoteServiceHost        string `yaml:"remote_service_host"`
	RemoteAccountHost        string `yaml:"remote_account_host"`
	RemoteTransferHost       string `yaml:"remote_transfer_host"`
	RemoteSearchHost         string `yaml:"remote_search_host"`
	RemoteMessagingHost      string `yaml:"remote_messaging_host"`
	RedisHosts               string `yaml:"redis_hosts"`
	RedisPassword            string `yaml:"redis_password"`
	RedisUser                string `yaml:"redis_user"`
	KafkaHosts               string `yaml:"kafka_hosts"`
	KafkaTopicAuditTail      string `yaml:"kafka_topic_audit_tail"`
	KafkaGroupAuditTail      string `yaml:"kafka_group_audit_tail"`
	KafkaTopicAccountStatus  string `yaml:"kafka_topic_account_status"`
	KafkaGroupAccountStatus  string `yaml:"kafka_group_account_status"`
	KafkaTopicTransferStatus string `yaml:"kafka_topic_transfer_status"`
	KafkaGroupTransferStatus string `yaml:"kafka_group_transfer_status"`
	KafkaTopicNotifyStatus   string `yaml:"kafka_topic_notify_status"`
	KafkaGroupNotifyStatus   string `yaml:"kafka_group_notify_status"`
	kafkaTopicMsgStatus      string `yaml:"kafka_topic_msg_status"`
	kafkaGroupMsgStatus      string `yaml:"kafka_group_msg_status"`
}

func (c *HostAddressConfig) Init() *HostAddressConfig {
	yamlFile, err := ioutil.ReadFile("app.yaml")
	if err != nil {
		fmt.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		fmt.Printf("Unmarshal: %v", err)
	}
	return c
}
func NewHostAddressConfig() *HostAddressConfig {
	return &HostAddressConfig{}
}
