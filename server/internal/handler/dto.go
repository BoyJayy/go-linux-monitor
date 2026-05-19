package handler

import "encoding/json"

func (e ErrorDTO) toString() string {
	b, err := json.Marshal(e)
	if err != nil {
		panic(err)
	}
	return string(b)
}
