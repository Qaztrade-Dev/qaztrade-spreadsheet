package jsondomain

import "encoding/json"

type UserAttrs struct {
	OrgName string `json:"org_name"`
}

func DecodeUserAttrs(jsonBytes []byte) (*UserAttrs, error) {
	var attrs UserAttrs

	if err := json.Unmarshal(jsonBytes, &attrs); err != nil {
		return nil, err
	}

	return &attrs, nil
}

func EncodeUserAttrs(attrs *UserAttrs) ([]byte, error) {
	return json.Marshal(attrs)
}
