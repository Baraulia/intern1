package handlers

import (
	"fmt"
	"reflect"
	"strings"
)

const tagName = "validate"

type Validator interface {
	Validate(interface{}) error
}

type DefaultValidator struct {
}

func (v DefaultValidator) Validate(val interface{}) error {
	return nil
}

type Alpha2Validator struct {
}

func (v Alpha2Validator) Validate(val interface{}) error {
	if len(val.(string)) != 2 {
		return fmt.Errorf("Alpha2Validator: the number of characters in the field must be 2")
	}
	return nil
}

type Alpha3Validator struct {
}

func (v Alpha3Validator) Validate(val interface{}) error {
	if len(val.(string)) != 3 {
		return fmt.Errorf("Alpha2Validator: the number of characters in the field must be 3")
	}
	return nil
}

type RequiredValidator struct{}

func (v RequiredValidator) Validate(val interface{}) error {
	if val == 0 || val == "" {
		return fmt.Errorf("RequiredValidator: this field cannot be empty")
	}
	return nil
}

func getValidatorFromTag(tag string) Validator {
	args := strings.Split(tag, ",")
	switch args[0] {
	case "alpha_2":
		return Alpha2Validator{}
	case "alpha_3":
		return Alpha3Validator{}
	case "required":
		return RequiredValidator{}
	}
	return DefaultValidator{}
}

func validateStruct(s interface{}) map[string]string {
	var errs = make(map[string]string)
	v := reflect.ValueOf(s)
	for i := 0; i < v.NumField(); i++ {
		tag := v.Type().Field(i).Tag.Get(tagName)
		if tag == "" || tag == "-" {
			continue
		}
		validator := getValidatorFromTag(tag)
		err := validator.Validate(v.Field(i).Interface())
		if err != nil {
			errs[v.Type().Field(i).Name] = err.Error()
		}
	}
	return errs
}
