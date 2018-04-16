package convertor

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/utrack/petrichor/client/confc"
	"fmt"
)

var goTypeToConfc = map[string]confc.SettingDescType{
	"bool":          confc.TypeBoolean,
	"string":        confc.TypeString,
	"int":           confc.TypeInteger,
	"float64":       confc.TypeFloat,
	"time.Duration": confc.TypeDuration,
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
			m[valType.Field(i).Name] = kind
		} else {
			parseTypedSettingsDesc(val.Field(i).Interface(), m)
		}
	}
	return m
}

var tagsMap map[string]reflect.Kind

func init() {
	tagsMap = parseTypedSettingsDesc(Json, make(map[string]reflect.Kind))
	tagsMap["json"] = reflect.String
}

func (j toJson) Convert(obj interface{}) ([]byte, error) {
	return nil, nil
}

func (j toJson) C(t *testing.T, obj interface{}) ([]byte, error) {
	t.Log(tagsMap)
	toJsonSlice := make([]toJson, 0)
	val := reflect.Indirect(reflect.ValueOf(obj))
	valType := val.Type()
	fieldsCount := val.NumField()
	for i := 0; i < fieldsCount; i++ {
		c := toJson{}
		t.Log(valType.Field(i).Name, valType.Field(i).Type, valType.Field(i).Tag, val.Field(i).Kind())
		// ==== > to be continued
		//v, ok := valType.Field(i).Tag.Lookup("tags")
		//t.Log(v, ok)
		confcType, ok := goTypeToConfc[valType.Field(i).Type.String()]
		if !ok {
			return nil, fmt.Errorf("unsupported type: %v", valType.Field(i).Type.String())
		}
		c.Name = valType.Field(i).Name
		c.Type = confcType
		c.DefaultValue = reflect.Zero(val.Field(i).Type()).Interface()
		toJsonSlice = append(toJsonSlice, c)
	}
	return json.Marshal(toJsonSlice)
}