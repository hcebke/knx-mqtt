package models

type OutgoingMqttJson struct {
	Bytes   *string `json:"bytes,omitempty"`
	Name    *string `json:"name,omitempty"`
	Value   any     `json:"value,omitempty"`
	Unit    *string `json:"unit,omitempty"`
	Command string  `json:"command"`
}
