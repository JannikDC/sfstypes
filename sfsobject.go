package sfstypes

import (
	"fmt"
)

type SFSObject struct {
	dataHolder map[string]sfsDataWrapper
}

func NewSFSObject() *SFSObject {
	return &SFSObject{
		dataHolder: make(map[string]sfsDataWrapper),
	}
}

func NewSFSObjectFromBinaryData(data []byte) (*SFSObject, error) {
	return newSFSObjectfromBinary(data)
}

func NewSFSObjectFromJsonData(jsonString string) (*SFSObject, error) {
	return sfsObjectFromJson(jsonString)
}

/*
// TODO
func NewSFSObjectFromResultSet() (*SFSObject, error) {
	return nil, nil
}
*/

func (sfsobject *SFSObject) ToBinary() []byte {
	return encodeSFSObject(sfsobject)
}

func (sfsobject *SFSObject) ToJson() string {
	return sfsObjectToJson(sfsobject)
}

func (sfsobject *SFSObject) Size() int {
	return len(sfsobject.dataHolder)
}

func (sfsobject *SFSObject) GetHexDump() string {
	bytes := sfsobject.ToBinary()
	hexString := fmt.Sprintf("binary size: %d\n", len(bytes))
	for i, b := range bytes {
		if i%16 == 0 && i != 0 {
			hexString += "\n"
		}
		hexString += fmt.Sprintf("%02x ", b)
	}
	return hexString
}

func (sfsobject *SFSObject) GetKeys() []string {
	keys := make([]string, 0, len(sfsobject.dataHolder))
	for k := range sfsobject.dataHolder {
		keys = append(keys, k)
	}
	return keys
}

func (sfsobject *SFSObject) ContainsKey(key string) bool {
	if _, exists := sfsobject.dataHolder[key]; !exists {
		return false
	}
	return true
}

func (sfsobject *SFSObject) RemoveElement(key string) error {
	if _, exists := sfsobject.dataHolder[key]; !exists {
		return &ErrKeyNotFound{key: key}
	}
	delete(sfsobject.dataHolder, key)
	return nil
}

func (sfsobject *SFSObject) IsNull(key string) (bool, error) {
	value, err := sfsobject.getData(key)
	if err == nil {
		if value == nil {
			return true, err
		}
		return false, err
	}
	return false, err
}

func (sfsobject *SFSObject) getWrapper(key string) (*sfsDataWrapper, error) {
	if value, exists := sfsobject.dataHolder[key]; exists {
		return &value, nil
	}
	return nil, &ErrKeyNotFound{key: key}
}

func (sfsobject *SFSObject) getData(key string) (interface{}, error) {
	if value, exists := sfsobject.dataHolder[key]; exists {
		return value.data, nil
	}
	return nil, &ErrKeyNotFound{key: key}
}

func (sfsobject *SFSObject) Get(key string) (interface{}, error) {
	value, err := sfsobject.getWrapper(key)
	if err != nil {
		return nil, err
	}
	switch value.typeId {
	case type_NULL:
		return nil, nil
	case type_BOOL:
		return value.data.(bool), nil
	case type_BYTE:
		return value.data.(int8), nil
	case type_SHORT:
		return value.data.(int16), nil
	case type_INT:
		return value.data.(int32), nil
	case type_LONG:
		return value.data.(int64), nil
	case type_FLOAT:
		return value.data.(float32), nil
	case type_DOUBLE:
		return value.data.(float64), nil
	case type_UTF_STRING:
		return value.data.(string), nil
	case type_BOOL_ARRAY:
		return value.data.([]bool), nil
	case type_BYTE_ARRAY:
		return value.data.([]int8), nil
	case type_SHORT_ARRAY:
		return value.data.([]int16), nil
	case type_INT_ARRAY:
		return value.data.([]int32), nil
	case type_LONG_ARRAY:
		return value.data.([]int64), nil
	case type_FLOAT_ARRAY:
		return value.data.([]float32), nil
	case type_DOUBLE_ARRAY:
		return value.data.([]float64), nil
	case type_UTF_STRING_ARRAY:
		return value.data.([]string), nil
	case type_SFS_ARRAY:
		return value.data.(SFSArray), nil
	case type_SFS_OBJECT:
		return value.data.(SFSObject), nil
	case type_TEXT:
		return value.data.(string), nil
	default:
		return nil, &ErrUnsupportedType{value: value}
	}
}

func (sfsobject *SFSObject) GetBool(key string) (bool, error) {
	value, err := sfsobject.getWrapper(key)
	if err == nil {
		if value.typeId == type_BOOL {
			return value.data.(bool), nil
		}
		return false, &ErrWrongType{actualType: value.typeId, wantedType: type_BOOL}
	}
	return false, err
}

func (sfsobject *SFSObject) GetBoolArray(key string) ([]bool, error) {
	value, err := sfsobject.getWrapper(key)
	if err == nil {
		if value.typeId == type_BOOL_ARRAY {
			return value.data.([]bool), nil
		}
		return nil, &ErrWrongType{actualType: value.typeId, wantedType: type_BOOL_ARRAY}
	}
	return nil, err
}

func (sfsobject *SFSObject) GetByte(key string) (int8, error) {
	value, err := sfsobject.getWrapper(key)
	if err == nil {
		if value.typeId == type_BYTE {
			return value.data.(int8), nil
		}
		return 0, &ErrWrongType{actualType: value.typeId, wantedType: type_BYTE}
	}
	return 0, err
}

func (sfsobject *SFSObject) GetByteArray(key string) ([]int8, error) {
	value, err := sfsobject.getWrapper(key)
	if err == nil {
		if value.typeId == type_BYTE_ARRAY {
			return value.data.([]int8), nil
		}
		return nil, &ErrWrongType{actualType: value.typeId, wantedType: type_BYTE_ARRAY}
	}
	return nil, err
}

func (sfsobject *SFSObject) GetDouble(key string) (float64, error) {
	value, err := sfsobject.getWrapper(key)
	if err == nil {
		if value.typeId == type_DOUBLE {
			return value.data.(float64), nil
		}
		return 0, &ErrWrongType{actualType: value.typeId, wantedType: type_DOUBLE}
	}
	return 0, err
}

func (sfsobject *SFSObject) GetDoubleArray(key string) ([]float64, error) {
	value, err := sfsobject.getWrapper(key)
	if err == nil {
		if value.typeId == type_DOUBLE_ARRAY {
			return value.data.([]float64), nil
		}
		return nil, &ErrWrongType{actualType: value.typeId, wantedType: type_DOUBLE_ARRAY}
	}
	return nil, err
}

func (sfsobject *SFSObject) GetFloat(key string) (float32, error) {
	value, err := sfsobject.getWrapper(key)
	if err == nil {
		if value.typeId == type_FLOAT {
			return value.data.(float32), nil
		}
		return 0, &ErrWrongType{actualType: value.typeId, wantedType: type_FLOAT}
	}
	return 0, err
}

func (sfsobject *SFSObject) GetFloatArray(key string) ([]float32, error) {
	value, err := sfsobject.getWrapper(key)
	if err == nil {
		if value.typeId == type_FLOAT_ARRAY {
			return value.data.([]float32), nil
		}
		return nil, &ErrWrongType{actualType: value.typeId, wantedType: type_FLOAT_ARRAY}
	}
	return nil, err
}

func (sfsobject *SFSObject) GetInt(key string) (int32, error) {
	value, err := sfsobject.getWrapper(key)
	if err == nil {
		if value.typeId == type_INT {
			return value.data.(int32), nil
		}
		return 0, &ErrWrongType{actualType: value.typeId, wantedType: type_INT}
	}
	return 0, err
}

func (sfsobject *SFSObject) GetIntArray(key string) ([]int32, error) {
	value, err := sfsobject.getWrapper(key)
	if err == nil {
		if value.typeId == type_INT_ARRAY {
			return value.data.([]int32), nil
		}
		return nil, &ErrWrongType{actualType: value.typeId, wantedType: type_INT_ARRAY}
	}
	return nil, err
}

func (sfsobject *SFSObject) GetLong(key string) (int64, error) {
	value, err := sfsobject.getWrapper(key)
	if err == nil {
		if value.typeId == type_LONG {
			return value.data.(int64), nil
		}
		return 0, &ErrWrongType{actualType: value.typeId, wantedType: type_LONG}
	}
	return 0, err
}

func (sfsobject *SFSObject) GetLongArray(key string) ([]int64, error) {
	value, err := sfsobject.getWrapper(key)
	if err == nil {
		if value.typeId == type_LONG_ARRAY {
			return value.data.([]int64), nil
		}
		return nil, &ErrWrongType{actualType: value.typeId, wantedType: type_LONG_ARRAY}
	}
	return nil, err
}

func (sfsobject *SFSObject) GetSFSArray(key string) (*SFSArray, error) {
	value, err := sfsobject.getWrapper(key)
	if err == nil {
		if value.typeId == type_SFS_ARRAY {
			return value.data.(*SFSArray), nil
		}
		return NewSFSArray(), &ErrWrongType{actualType: value.typeId, wantedType: type_SFS_ARRAY}
	}
	return NewSFSArray(), err
}

func (sfsobject *SFSObject) GetSFSObject(key string) (*SFSObject, error) {
	value, err := sfsobject.getWrapper(key)
	if err == nil {
		if value.typeId == type_SFS_OBJECT {
			return value.data.(*SFSObject), nil
		}
		return NewSFSObject(), &ErrWrongType{actualType: value.typeId, wantedType: type_SFS_OBJECT}
	}
	return NewSFSObject(), err
}

func (sfsobject *SFSObject) GetShort(key string) (int16, error) {
	value, err := sfsobject.getWrapper(key)
	if err == nil {
		if value.typeId == type_SHORT {
			return value.data.(int16), nil
		}
		return 0, &ErrWrongType{actualType: value.typeId, wantedType: type_SHORT}
	}
	return 0, err
}

func (sfsobject *SFSObject) GetShortArray(key string) ([]int16, error) {
	value, err := sfsobject.getWrapper(key)
	if err == nil {
		if value.typeId == type_SHORT_ARRAY {
			return value.data.([]int16), nil
		}
		return nil, &ErrWrongType{actualType: value.typeId, wantedType: type_SHORT_ARRAY}
	}
	return nil, err
}

func (sfsobject *SFSObject) GetUnsignedByte(key string) (uint8, error) {
	value, err := sfsobject.getWrapper(key)
	if err == nil {
		if value.typeId == type_BYTE {
			return value.data.(uint8), nil
		}
		return 0, &ErrWrongType{actualType: value.typeId, wantedType: type_BYTE}
	}
	return 0, err
}

func (sfsobject *SFSObject) GetUnsignedByteArray(key string) ([]uint8, error) {
	value, err := sfsobject.getWrapper(key)
	if err == nil {
		if value.typeId == type_BYTE_ARRAY {
			return value.data.([]uint8), nil
		}
		return nil, &ErrWrongType{actualType: value.typeId, wantedType: type_BYTE_ARRAY}
	}
	return nil, err
}

func (sfsobject *SFSObject) GetUtfString(key string) (string, error) {
	value, err := sfsobject.getWrapper(key)
	if err == nil {
		if value.typeId == type_UTF_STRING {
			return value.data.(string), nil
		}
		return "", &ErrWrongType{actualType: value.typeId, wantedType: type_UTF_STRING}
	}
	return "", err
}

func (sfsobject *SFSObject) GetText(key string) (string, error) {
	value, err := sfsobject.getWrapper(key)
	if err == nil {
		if value.typeId == type_TEXT {
			return value.data.(string), nil
		}
		return "", &ErrWrongType{actualType: value.typeId, wantedType: type_TEXT}
	}
	return "", err
}

func (sfsobject *SFSObject) GetUtfStringArray(key string) ([]string, error) {
	value, err := sfsobject.getWrapper(key)
	if err == nil {
		if value.typeId == type_UTF_STRING_ARRAY {
			return value.data.([]string), nil
		}
		return nil, &ErrWrongType{actualType: value.typeId, wantedType: type_UTF_STRING_ARRAY}
	}
	return nil, err
}

func (sfsobject *SFSObject) putsfsDataWrapper(key string, wrapper *sfsDataWrapper) error {
	if key == "" {
		return ErrKeyEmpty
	} else if len(key) > 255 {
		return &ErrInvalidKeySize{key: key, length: len(key)}
	}
	sfsobject.dataHolder[key] = *wrapper
	return nil
}

func (sfsobject *SFSObject) putData(key string, value interface{}, typeId sfsDataType) error {
	if typeId != type_NULL && value == nil {
		return ErrDataNull
	}
	return sfsobject.putsfsDataWrapper(key, newsfsDataWrapper(typeId, value))
}

func (sfsobject *SFSObject) Put(key string, value interface{}) error {
	switch v := value.(type) {
	case bool:
		return sfsobject.PutBool(key, v)
	case []bool:
		return sfsobject.PutBoolArray(key, v)
	case int8:
		return sfsobject.PutByte(key, v)
	case []int8:
		return sfsobject.PutByteArray(key, v)
	case float32:
		return sfsobject.PutFloat(key, v)
	case []float32:
		return sfsobject.PutFloatArray(key, v)
	case float64:
		return sfsobject.PutDouble(key, v)
	case []float64:
		return sfsobject.PutDoubleArray(key, v)
	case int:
		return sfsobject.PutInt(key, int32(v))
	case []int:
		arr := make([]int32, len(v))
		for i, element := range v {
			arr[i] = int32(element)
		}
		return sfsobject.PutIntArray(key, arr)
	case int32:
		return sfsobject.PutInt(key, v)
	case []int32:
		return sfsobject.PutIntArray(key, v)
	case int64:
		return sfsobject.PutLong(key, v)
	case []int64:
		return sfsobject.PutLongArray(key, v)
	case nil:
		return sfsobject.PutNull(key)
	case SFSArray:
		return sfsobject.PutSFSArray(key, &v)
	case SFSObject:
		return sfsobject.PutSFSObject(key, &v)
	case *SFSArray:
		return sfsobject.PutSFSArray(key, v)
	case *SFSObject:
		return sfsobject.PutSFSObject(key, v)
	case int16:
		return sfsobject.PutShort(key, v)
	case []int16:
		return sfsobject.PutShortArray(key, v)
	case string:
		return sfsobject.PutUtfString(key, v)
	case []string:
		return sfsobject.PutUtfStringArray(key, v)
	}
	return &ErrUnsupportedType{value: value}
}

func (sfsobject *SFSObject) PutBool(key string, value bool) error {
	return sfsobject.putData(key, value, type_BOOL)
}

func (sfsobject *SFSObject) PutBoolArray(key string, value []bool) error {
	return sfsobject.putData(key, value, type_BOOL_ARRAY)
}

func (sfsobject *SFSObject) PutByte(key string, value int8) error {
	return sfsobject.putData(key, value, type_BYTE)
}

func (sfsobject *SFSObject) PutByteArray(key string, value []int8) error {
	return sfsobject.putData(key, value, type_BYTE_ARRAY)
}

func (sfsobject *SFSObject) PutDouble(key string, value float64) error {
	return sfsobject.putData(key, value, type_DOUBLE)
}

func (sfsobject *SFSObject) PutDoubleArray(key string, value []float64) error {
	return sfsobject.putData(key, value, type_DOUBLE_ARRAY)
}

func (sfsobject *SFSObject) PutFloat(key string, value float32) error {
	return sfsobject.putData(key, value, type_FLOAT)
}

func (sfsobject *SFSObject) PutFloatArray(key string, value []float32) error {
	return sfsobject.putData(key, value, type_FLOAT_ARRAY)
}

func (sfsobject *SFSObject) PutInt(key string, value int32) error {
	return sfsobject.putData(key, value, type_INT)
}

func (sfsobject *SFSObject) PutIntArray(key string, value []int32) error {
	return sfsobject.putData(key, value, type_INT_ARRAY)
}

func (sfsobject *SFSObject) PutLong(key string, value int64) error {
	return sfsobject.putData(key, value, type_LONG)
}

func (sfsobject *SFSObject) PutLongArray(key string, value []int64) error {
	return sfsobject.putData(key, value, type_LONG_ARRAY)
}

func (sfsobject *SFSObject) PutNull(key string) error {
	return sfsobject.putData(key, nil, type_NULL)
}

func (sfsobject *SFSObject) PutSFSArray(key string, value *SFSArray) error {
	return sfsobject.putData(key, *value, type_SFS_ARRAY)
}

func (sfsobject *SFSObject) PutSFSObject(key string, value *SFSObject) error {
	return sfsobject.putData(key, *value, type_SFS_OBJECT)
}

func (sfsobject *SFSObject) PutShort(key string, value int16) error {
	return sfsobject.putData(key, value, type_SHORT)
}

func (sfsobject *SFSObject) PutShortArray(key string, value []int16) error {
	return sfsobject.putData(key, value, type_SHORT_ARRAY)
}

func (sfsobject *SFSObject) PutUtfString(key string, value string) error {
	return sfsobject.putData(key, value, type_UTF_STRING)
}

func (sfsobject *SFSObject) PutText(key string, value string) error {
	return sfsobject.putData(key, value, type_TEXT)
}

func (sfsobject *SFSObject) PutUtfStringArray(key string, value []string) error {
	return sfsobject.putData(key, value, type_UTF_STRING_ARRAY)
}
