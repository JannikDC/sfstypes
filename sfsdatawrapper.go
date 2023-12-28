package sfstypes

type sfsDataType byte

const (
	type_NULL             sfsDataType = 0
	type_BOOL             sfsDataType = 1
	type_BYTE             sfsDataType = 2
	type_SHORT            sfsDataType = 3
	type_INT              sfsDataType = 4
	type_LONG             sfsDataType = 5
	type_FLOAT            sfsDataType = 6
	type_DOUBLE           sfsDataType = 7
	type_UTF_STRING       sfsDataType = 8
	type_BOOL_ARRAY       sfsDataType = 9
	type_BYTE_ARRAY       sfsDataType = 10
	type_SHORT_ARRAY      sfsDataType = 11
	type_INT_ARRAY        sfsDataType = 12
	type_LONG_ARRAY       sfsDataType = 13
	type_FLOAT_ARRAY      sfsDataType = 14
	type_DOUBLE_ARRAY     sfsDataType = 15
	type_UTF_STRING_ARRAY sfsDataType = 16
	type_SFS_ARRAY        sfsDataType = 17
	type_SFS_OBJECT       sfsDataType = 18
	type_CLASS            sfsDataType = 19
	type_TEXT             sfsDataType = 20
)

type sfsDataWrapper struct {
	typeId sfsDataType
	data   interface{}
}

func newsfsDataWrapper(typeId sfsDataType, data interface{}) *sfsDataWrapper {
	return &sfsDataWrapper{
		typeId: typeId,
		data:   data,
	}
}

func sfsTypeToString(sfsType sfsDataType) string {
	switch sfsType {
	case type_NULL:
		return "NULL/nil"
	case type_BOOL:
		return "BOOL/bool"
	case type_BYTE:
		return "BYTE/int8"
	case type_SHORT:
		return "SHORT/int16"
	case type_INT:
		return "INT/int32"
	case type_LONG:
		return "LONG/int64"
	case type_FLOAT:
		return "FLOAT/float32"
	case type_DOUBLE:
		return "DOUBLE/float64"
	case type_UTF_STRING:
		return "UTF_STRING/string"
	case type_BOOL_ARRAY:
		return "BOOL_ARRAY/[]bool"
	case type_BYTE_ARRAY:
		return "BYTE_ARRAY/[]int8"
	case type_SHORT_ARRAY:
		return "SHORT_ARRAY/[]int16"
	case type_INT_ARRAY:
		return "INT_ARRAY/[]int32"
	case type_LONG_ARRAY:
		return "LONG_ARRAY/[]int64"
	case type_FLOAT_ARRAY:
		return "FLOAT_ARRAY/[]float32"
	case type_DOUBLE_ARRAY:
		return "DOUBLE_ARRAY/[]float64"
	case type_UTF_STRING_ARRAY:
		return "UTF_STRING_ARRAY/[]string"
	case type_SFS_ARRAY:
		return "SFS_ARRAY"
	case type_SFS_OBJECT:
		return "SFS_OBJECT"
	case type_CLASS:
		return "CLASS (unsupported)"
	case type_TEXT:
		return "TEXT/string"
	default:
		return "unknown type"
	}
}
