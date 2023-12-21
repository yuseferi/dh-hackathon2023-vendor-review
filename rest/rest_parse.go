package rest

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"io"
)

func ParseAndValidateBody[T any](body io.Reader, v *validator.Validate) (*T, error) {
	var data T
	err := json.NewDecoder(body).Decode(&data)
	if err != nil {
		return nil, err
	}
	err = v.Struct(&data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}
