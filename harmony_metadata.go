package main

import "encoding/json"

func UnmarshalHarmonyMetadata(data []byte) (HarmonyMetadata, error) {
	var r HarmonyMetadata
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *HarmonyMetadata) Marshal() ([]byte, error) {
	return json.Marshal(r)
}
type HarmonyMetadata struct {
	Background   string `json:"background,omitempty"`
	Skin         string `json:"skin,omitempty"`
	Eye          string `json:"eye,omitempty"`
	Mouth        string `json:"mouth,omitempty"`
	Clothes      string `json:"clothes,omitempty"`
	Glass        string `json:"glass,omitempty"`
	Hair         string `json:"hair,omitempty"`
	Hat          string `json:"hat,omitempty"`
	Neck         string `json:"neck,omitempty"`
	RandomObject string `json:"random_object,omitempty"`
}
