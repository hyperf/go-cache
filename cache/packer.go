package cache

import "encoding/json"

type JsonPacker struct {
}

func (j *JsonPacker) Pack(data any) (string, error) {
	res, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	return string(res), nil
}

func (j *JsonPacker) UnPack(raw string, data any) error {
	return json.Unmarshal([]byte(raw), data)
}
