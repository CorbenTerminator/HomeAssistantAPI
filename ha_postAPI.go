package homeassistant

import (
	"fmt"
	"encoding/json"
)

const (
	Sensor = "sensor"
	BinarySensor = "binary_sensor"
	AlarmPanel = "alarm_control_panel"
	Light = "light"
	Switch = "switch"
	Script = "script"
)

var StatePath = "states/%s.%s"


type StatePost struct {
	State string	`json:"state"`
	Attributes struct {
		FriendlyName	string	`json:"friendly_name"`
		UnitMeasure	string	`json:"unit_of_measurement"`
	} `json:"attributes"`
}

func (ha *HomeAssistant) SendState(sensorType string, entityId string, st StatePost) error {
	path := fmt.Sprintf(StatePath, sensorType, entityId)
	jsonObj, err := json.Marshal(st)
	if err != nil {
		return err
	}
	//the resp body is not needed at this moment, there is not useful information there
	_, err1 := ha.httpPost(path, jsonObj)
	if err1 != nil {
		return err
	}
	return nil
}