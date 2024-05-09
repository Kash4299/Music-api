package util

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func ParseString(value interface{}) string {
	str, ok := value.(string)
	if !ok {
		return str
	}
	str = strings.Replace(str, "\n", "", -1)
	str = strings.Trim(str, "\r\n")
	str = strings.TrimSpace(str)
	return str
}

func ParseStringToInt(str string) int {
	i, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}
	return i
}

func ParseFloat64(value interface{}) float64 {
	typeData := reflect.ValueOf(value)
	str := ""
	if typeData.Kind() == reflect.String {
		str = ParseString(value)
	} else {
		str = fmt.Sprintf("%v", value)
	}
	i, err := strconv.ParseFloat(str, 64)
	if err != nil {
		i = 0
	}
	return i
}

func ParseInt(str string) int {
	str = ParseString(str)
	i, err := strconv.Atoi(str)
	if err != nil {
		i = 0
	}
	return i
}

func ParseMapToString(value interface{}) (string, error) {
	bytes, err := json.Marshal(value)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func ParseAnyToAny(value any, dest any) error {
	bytes, err := json.Marshal(value)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(bytes, dest); err != nil {
		return err
	}
	return nil
}

func GetAudioDir() string {
	return "upload_file/audio"
}
