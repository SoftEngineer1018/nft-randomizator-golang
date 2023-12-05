package main

import "encoding/json"

func UnmarshalOpenseaMetadata(data []byte) (OpenseaMetadata, error) {
	var r OpenseaMetadata
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *OpenseaMetadata) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type OpenseaMetadata struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Image       string      `json:"image"`
	Attributes  []Attribute `json:"attributes"`
}
