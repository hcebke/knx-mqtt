package protocols

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/pakerfeldt/knx-mqtt/models"
	"github.com/pakerfeldt/knx-mqtt/utils"
	"github.com/rs/zerolog/log"
	"github.com/vapourismo/knx-go/knx"
	"github.com/vapourismo/knx-go/knx/dpt"
)

type KnxEventHandler func(knx.GroupEvent)

func SubscribeKnx(knxClient KnxClient, handler KnxEventHandler) {
	go func() {
		for event := range knxClient.Inbound() {
			handler(event)
		}
	}()
}

func IncomingKnxEventHandler(mqttClient mqtt.Client, knxItems *models.KNX, mqttMessageCfg models.OutgoingMqttMessage, mqttCfg models.MQTTConfig, ignoreUnknownGroupAddresses bool) func(event knx.GroupEvent) {
	return func(event knx.GroupEvent) {
		incomingKnxEventHandler(event, mqttClient, knxItems, mqttMessageCfg, mqttCfg, ignoreUnknownGroupAddresses)
	}
}

func incomingKnxEventHandler(event knx.GroupEvent, mqttClient mqtt.Client, knxItems *models.KNX, mqttMessageCfg models.OutgoingMqttMessage, mqttCfg models.MQTTConfig, ignoreUnknownGroupAddresses bool) {
	index, exists := knxItems.GadToIndex[event.Destination.String()]
	if !exists && ignoreUnknownGroupAddresses {
		log.Debug().Str("address", event.Destination.String()).Msg("Ignoring unknown group address")
		return
	}
	if !exists && mqttMessageCfg.Type == "bytes" && mqttMessageCfg.EmitUsingAddress {
		log.Debug().Str("protocol", "knx").Str("address", event.Destination.String()).Msg("Incoming")
		mqttClient.Publish(mqttCfg.TopicPrefix+event.Destination.String(), mqttCfg.Qos, mqttCfg.Retain, event.Data)
		return
	} else if !exists {
		log.Info().Str("address", event.Destination.String()).Msg("Cannot read unknown address, update your KNX XML export")
		return
	}
	groupAddress := knxItems.GroupAddresses[index]
	datapoint, ok := dpt.Produce(groupAddress.Datapoint)
	if ok {
		datapoint.Unpack(event.Data)
		log.Debug().Str("protocol", "knx").Str("address", event.Destination.String()).Str("name", groupAddress.Name).Str("value", datapoint.String()).Msg("Incoming")

		payload, err := constructPayload(datapoint, groupAddress.Datapoint, mqttMessageCfg.EmitValueAsString, mqttMessageCfg.Type, &mqttMessageCfg.IncludedJsonFields, &groupAddress.Name)
		if err != nil {
			return
		}

		if mqttMessageCfg.EmitUsingAddress {
			mqttClient.Publish(mqttCfg.TopicPrefix+groupAddress.Address, mqttCfg.Qos, mqttCfg.Retain, payload)
		}
		if mqttMessageCfg.EmitUsingName {
			mqttClient.Publish(mqttCfg.TopicPrefix+groupAddress.FullName, mqttCfg.Qos, mqttCfg.Retain, payload)
		}
	} else {
		log.Warn().Str("address", event.Destination.String()).Msgf("Could not create datapoint %s for address %s", knxItems.GroupAddresses[index].Datapoint, event.Destination.String())
	}
}

func constructPayload(dpt dpt.Datapoint, dptType string, emitValueAsString bool, messageType string, jsonFields *models.IncludedJsonFields, addressName *string) (interface{}, error) {
	var payload interface{}
	if messageType == models.JsonType {
		outgoingJson := models.OutgoingMqttJson{}
		if jsonFields.IncludeBytes {
			base64 := base64.StdEncoding.EncodeToString(dpt.Pack())
			outgoingJson.Bytes = &base64
		}
		if jsonFields.IncludeName {
			outgoingJson.Name = addressName
		}
		if jsonFields.IncludeValue {
			if emitValueAsString {
				outgoingJson.Value = utils.StringWithoutSuffix(dpt)
			} else {
				outgoingJson.Value = utils.ExtractDatapointValue(dpt, dptType)
			}
		}
		if jsonFields.IncludeUnit {
			unit := dpt.Unit()
			outgoingJson.Unit = &unit
		}
		jsonBytes, err := json.Marshal(outgoingJson)
		if err != nil {
			log.Error().Str("error", fmt.Sprintf("%+v", err)).Msg("Failed to create outgoing JSON message")
			return nil, err
		}
		payload = string(jsonBytes)
	} else if messageType == models.ValueType {
		if emitValueAsString {
			payload = utils.StringWithoutSuffix(dpt)
		} else {
			payload = fmt.Sprintf("%v", utils.ExtractDatapointValue(dpt, dptType))
		}
	} else if messageType == models.ValueWithUnitType {
		payload = dpt.String()
	} else if messageType == models.BytesType {
		payload = dpt.Pack()
	}
	return payload, nil
}
