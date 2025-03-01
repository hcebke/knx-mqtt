package msg

import (
	mqttgo "github.com/eclipse/paho.mqtt.golang"
)

type MQTTMessage struct {
	message mqttgo.Message
}

func NewMQTT(m mqttgo.Message) *MQTTMessage {
	return &MQTTMessage{message: m}
}

func (m MQTTMessage) Topic() string {
	return m.message.Topic()
}

func (m MQTTMessage) Bytes() []byte {
	return m.message.Payload()
}
