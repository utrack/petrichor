package confc

// SettingDesc fully describes some setting of an app.
type SettingDesc struct {
	// Tags are used to map this setting to business features' collection.
	Tags []string
	// Module is an ID or a name of the app's module that uses this setting.
	Module string
	// Description is this setting's description.
	Description string
}

type NamedSettingDesc struct {
	SettingDesc
	Name string
}

type TypedSettingDesc struct {
	NamedSettingDesc

	Type         SettingDescType
	DefaultValue interface{}

	// TODO need value list for an enum
}

func newTypedDesc(
	name string,
	def interface{}, vType SettingDescType,
	desc SettingDesc,
) TypedSettingDesc {
	return TypedSettingDesc{
		Type:         vType,
		DefaultValue: def,
		NamedSettingDesc: NamedSettingDesc{
			Name:        name,
			SettingDesc: desc,
		},
	}
}

// SettingDescType describes a type of a setting.
type SettingDescType uint

const (
	TypeString SettingDescType = iota
	TypeBoolean
	TypeInteger
	TypeDuration
	TypeEnum
)
