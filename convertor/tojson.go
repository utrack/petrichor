package convertor

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/utrack/petrichor/client/confc"
)

var goTypeToConfc = map[reflect.Kind]confc.SettingDescType{
	reflect.Bool:    confc.TypeBoolean,
	reflect.String:  confc.TypeString,
	reflect.Int:     confc.TypeInteger,
	reflect.Int64:   confc.TypeDuration,
	reflect.Float64: confc.TypeFloat,
	reflect.Slice:   confc.TypeEnum,
}

type toJson confc.TypedSettingDesc

var Json = toJson(confc.TypedSettingDesc{})

func parseTypedSettingsDesc(obj interface{}, m map[string]reflect.Kind) map[string]reflect.Kind {
	val := reflect.Indirect(reflect.ValueOf(obj))
	valType := val.Type()
	fieldsCount := val.NumField()
	for i := 0; i < fieldsCount; i++ {
		kind := valType.Field(i).Type.Kind()
		if kind != reflect.Struct {
			m[strings.ToLower(valType.Field(i).Name)] = kind
		} else {
			parseTypedSettingsDesc(val.Field(i).Interface(), m)
		}
	}
	return m
}

var tagsMap map[string]reflect.Kind

func init() {
	tagsMap = parseTypedSettingsDesc(confc.SettingDesc{}, make(map[string]reflect.Kind))
	tagsMap["json"] = reflect.String
}

func (j *toJson) addTags(fieldTags reflect.StructTag) {
	for tagMaybe := range tagsMap {
		if tag, ok := fieldTags.Lookup(tagMaybe); ok {
			switch tagMaybe {
			case "json":
				j.Name = tag
			case "tags":
				j.Tags = strings.Split(tag, ",")
			case "module":
				j.Module = tag
			case "description":
				j.Description = tag
			}
		}
	}
}

func (j toJson) Convert(obj interface{}) ([]byte, error) {
	val := reflect.Indirect(reflect.ValueOf(obj))
	if !val.IsValid() {
		return nil, errors.New("input object is nil")
	}
	valType := val.Type()
	fieldsCount := val.NumField()
	toJsonSlice := make([]toJson, 0)
	for i := 0; i < fieldsCount; i++ {
		confcType, ok := goTypeToConfc[valType.Field(i).Type.Kind()]
		if !ok {
			return nil, fmt.Errorf("unsupported type: %v", valType.Field(i).Type.String())
		}
		c := toJson{}
		c.Name = valType.Field(i).Name
		c.Type = confcType
		c.DefaultValue = reflect.Zero(val.Field(i).Type()).Interface()
		c.addTags(valType.Field(i).Tag)
		toJsonSlice = append(toJsonSlice, c)
	}
	return json.Marshal(toJsonSlice)
}
