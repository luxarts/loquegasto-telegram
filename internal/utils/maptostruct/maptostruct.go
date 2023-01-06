package maptostruct

import (
	"encoding/json"
)

func Convert(src interface{}, dst interface{}) error {
	jsonBody, err := json.Marshal(src)
	if err != nil {
		return err
	}
	err = json.Unmarshal(jsonBody, dst)
	if err != nil {
		return err
	}

	return nil
}
