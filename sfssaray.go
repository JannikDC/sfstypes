package sfstypes

import (
	"fmt"
)

type SFSArray struct {
	dataHolder []sfsDataWrapper
}

func NewSFSArray() *SFSArray {
	return &SFSArray{
		dataHolder: make([]sfsDataWrapper, 0),
	}
}

func NewSFSArrayFromBinaryData(data []byte) (*SFSArray, error) {
	return newSFSArrayFromBinaryData(data)
}

func NewSFSArrayFromJsonData(jsonStr string) (*SFSArray, error) {
	return sfsArrayFromJson(jsonStr)
}

/*
// TODO
func NewSFSArrayFromResultSet() (*SFSArray, error) {

}
*/

func (sfsarray *SFSArray) GetHexDump() string {
	bytes := sfsarray.ToBinary()
	hexString := fmt.Sprintf("binary size: %d\n", len(bytes))
	for i, b := range bytes {
		if i%16 == 0 && i != 0 {
			hexString += "\n"
		}
		hexString += fmt.Sprintf("%02x ", b)
	}
	return hexString
}

func (sfsarray *SFSArray) ToBinary() []byte {
	return encodeSFSArray(sfsarray)
}

func (sfsarray *SFSArray) ToJson() string {
	return sfsArrayToJson(sfsarray)
}

func (sfsarray *SFSArray) IsNull(index int) (bool, error) {
	test, err := sfsarray.getWrapper(index)
	if err != nil {
		return false, err
	}
	if test == nil {
		return true, nil
	}
	return false, nil
}

func (sfsarray *SFSArray) Contains(value interface{}) bool {
	for _, v := range sfsarray.dataHolder {
		if v == value {
			return true
		}
	}
	return false
}

func (sfsarray *SFSArray) GetElementAt(index int) (interface{}, error) {
	if index >= len(sfsarray.dataHolder) || index < 0 {
		return nil, &ErrIndexNotInRange{index: index}
	}
	return sfsarray.dataHolder[index], nil
}

func (sfsarray *SFSArray) RemoveElementAt(index int) error {
	if index >= len(sfsarray.dataHolder) || index < 0 {
		return &ErrIndexNotInRange{index: index}
	}
	sfsarray.dataHolder = append(sfsarray.dataHolder[:index], sfsarray.dataHolder[index+1:]...)
	return nil
}

func (sfsarray *SFSArray) Size() int {
	return len(sfsarray.dataHolder)
}

func (sfsarray *SFSArray) getWrapper(index int) (*sfsDataWrapper, error) {
	if index >= 0 && index <= sfsarray.Size() {
		return &sfsarray.dataHolder[index], nil
	}
	return nil, &ErrIndexNotInRange{index: index}
}

func (sfsarray *SFSArray) Get(index int) (interface{}, error) {
	if index < 0 && index >= len(sfsarray.dataHolder) {
		return nil, &ErrIndexNotInRange{index: index}
	}
	value, _ := sfsarray.getWrapper(index)
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

func (sfsarray *SFSArray) GetBool(index int) (bool, error) {
	value, err := sfsarray.getWrapper(index)
	if err == nil {
		if value.typeId == type_BOOL {
			return value.data.(bool), nil
		}
		return false, &ErrWrongType{actualType: value.typeId, wantedType: type_BOOL}
	}
	return false, err
}

func (sfsarray *SFSArray) GetBoolArray(index int) ([]bool, error) {
	value, err := sfsarray.getWrapper(index)
	if err == nil {
		if value.typeId == type_BOOL_ARRAY {
			return value.data.([]bool), nil
		}
		return nil, &ErrWrongType{actualType: value.typeId, wantedType: type_BOOL_ARRAY}
	}
	return nil, err
}

func (sfsarray *SFSArray) GetByte(index int) (int8, error) {
	value, err := sfsarray.getWrapper(index)
	if err == nil {
		if value.typeId == type_BYTE {
			return value.data.(int8), nil
		}
		return 0, &ErrWrongType{actualType: value.typeId, wantedType: type_BYTE}
	}
	return 0, err
}

func (sfsarray *SFSArray) GetByteArray(index int) ([]int8, error) {
	value, err := sfsarray.getWrapper(index)
	if err == nil {
		if value.typeId == type_BYTE_ARRAY {
			return value.data.([]int8), nil
		}
		return nil, &ErrWrongType{actualType: value.typeId, wantedType: type_BYTE_ARRAY}
	}
	return nil, err
}

func (sfsarray *SFSArray) GetDouble(index int) (float64, error) {
	value, err := sfsarray.getWrapper(index)
	if err == nil {
		if value.typeId == type_DOUBLE {
			return value.data.(float64), nil
		}
		return 0, &ErrWrongType{actualType: value.typeId, wantedType: type_DOUBLE}
	}
	return 0, err
}

func (sfsarray *SFSArray) GetDoubleArray(index int) ([]float64, error) {
	value, err := sfsarray.getWrapper(index)
	if err == nil {
		if value.typeId == type_DOUBLE_ARRAY {
			return value.data.([]float64), nil
		}
		return nil, &ErrWrongType{actualType: value.typeId, wantedType: type_DOUBLE_ARRAY}
	}
	return nil, err
}

func (sfsarray *SFSArray) GetFloat(index int) (float32, error) {
	value, err := sfsarray.getWrapper(index)
	if err == nil {
		if value.typeId == type_FLOAT {
			return value.data.(float32), nil
		}
		return 0, &ErrWrongType{actualType: value.typeId, wantedType: type_FLOAT}
	}
	return 0, err
}

func (sfsarray *SFSArray) GetFloatArray(index int) ([]float32, error) {
	value, err := sfsarray.getWrapper(index)
	if err == nil {
		if value.typeId == type_FLOAT_ARRAY {
			return value.data.([]float32), nil
		}
		return nil, &ErrWrongType{actualType: value.typeId, wantedType: type_FLOAT_ARRAY}
	}
	return nil, err
}

func (sfsarray *SFSArray) GetInt(index int) (int32, error) {
	value, err := sfsarray.getWrapper(index)
	if err == nil {
		if value.typeId == type_INT {
			return value.data.(int32), nil
		}
		return 0, &ErrWrongType{actualType: value.typeId, wantedType: type_INT}
	}
	return 0, err
}

func (sfsarray *SFSArray) GetIntArray(index int) ([]int32, error) {
	value, err := sfsarray.getWrapper(index)
	if err == nil {
		if value.typeId == type_INT_ARRAY {
			return value.data.([]int32), nil
		}
		return nil, &ErrWrongType{actualType: value.typeId, wantedType: type_INT_ARRAY}
	}
	return nil, err
}

func (sfsarray *SFSArray) GetLong(index int) (int64, error) {
	value, err := sfsarray.getWrapper(index)
	if err == nil {
		if value.typeId == type_LONG {
			return value.data.(int64), nil
		}
		return 0, &ErrWrongType{actualType: value.typeId, wantedType: type_LONG}
	}
	return 0, err
}

func (sfsarray *SFSArray) GetLongArray(index int) ([]int64, error) {
	value, err := sfsarray.getWrapper(index)
	if err == nil {
		if value.typeId == type_LONG_ARRAY {
			return value.data.([]int64), nil
		}
		return nil, &ErrWrongType{actualType: value.typeId, wantedType: type_LONG_ARRAY}
	}
	return nil, err
}

func (sfsarray *SFSArray) GetSFSObject(index int) (SFSObject, error) {
	value, err := sfsarray.getWrapper(index)
	if err == nil {
		if value.typeId == type_SFS_OBJECT {
			return value.data.(SFSObject), nil
		}
		return *NewSFSObject(), &ErrWrongType{actualType: value.typeId, wantedType: type_SFS_OBJECT}
	}
	return *NewSFSObject(), err
}

func (sfsarray *SFSArray) GetSFSArray(index int) (SFSArray, error) {
	value, err := sfsarray.getWrapper(index)
	if err == nil {
		if value.typeId == type_SFS_ARRAY {
			return *value.data.(*SFSArray), nil
		}
		return *NewSFSArray(), &ErrWrongType{actualType: value.typeId, wantedType: type_SFS_ARRAY}
	}
	return *NewSFSArray(), err
}

func (sfsarray *SFSArray) GetShort(index int) (int16, error) {
	value, err := sfsarray.getWrapper(index)
	if err == nil {
		if value.typeId == type_SHORT {
			return value.data.(int16), nil
		}
		return 0, &ErrWrongType{actualType: value.typeId, wantedType: type_SHORT}
	}
	return 0, err
}

func (sfsarray *SFSArray) GetShortArray(index int) ([]int16, error) {
	value, err := sfsarray.getWrapper(index)
	if err == nil {
		if value.typeId == type_SHORT_ARRAY {
			return value.data.([]int16), nil
		}
		return nil, &ErrWrongType{actualType: value.typeId, wantedType: type_SHORT_ARRAY}
	}
	return nil, err
}

func (sfsarray *SFSArray) GetUnsignedByte(index int) (uint8, error) {
	value, err := sfsarray.getWrapper(index)
	if err == nil {
		if value.typeId == type_BYTE {
			return value.data.(uint8), nil
		}
		return 0, &ErrWrongType{actualType: value.typeId, wantedType: type_BYTE}
	}
	return 0, err
}

func (sfsarray *SFSArray) GetUnsignedByteArray(index int) ([]uint8, error) {
	value, err := sfsarray.getWrapper(index)
	if err == nil {
		if value.typeId == type_BYTE_ARRAY {
			return value.data.([]uint8), nil
		}
		return nil, &ErrWrongType{actualType: value.typeId, wantedType: type_BYTE_ARRAY}
	}
	return nil, err
}

func (sfsarray *SFSArray) GetUtfString(index int) (string, error) {
	value, err := sfsarray.getWrapper(index)
	if err == nil {
		if value.typeId == type_UTF_STRING {
			return value.data.(string), nil
		}
		return "", &ErrWrongType{actualType: value.typeId, wantedType: type_UTF_STRING}
	}
	return "", err
}

func (sfsarray *SFSArray) GetText(index int) (string, error) {
	value, err := sfsarray.getWrapper(index)
	if err == nil {
		if value.typeId == type_TEXT {
			return value.data.(string), nil
		}
		return "", &ErrWrongType{actualType: value.typeId, wantedType: type_TEXT}
	}
	return "", err
}

func (sfsarray *SFSArray) GetUtfStringArray(index int) ([]string, error) {
	value, err := sfsarray.getWrapper(index)
	if err == nil {
		if value.typeId == type_UTF_STRING_ARRAY {
			return value.data.([]string), nil
		}
		return nil, &ErrWrongType{actualType: value.typeId, wantedType: type_UTF_STRING_ARRAY}
	}
	return nil, err
}

func (sfsarray *SFSArray) addData(value interface{}, typeId sfsDataType) {
	sfsarray.addsfsDataWrapper(*newsfsDataWrapper(typeId, value))
}

func (sfsarray *SFSArray) addsfsDataWrapper(value sfsDataWrapper) {
	sfsarray.dataHolder = append(sfsarray.dataHolder, value)
}

func (sfsarray *SFSArray) Add(value interface{}) error {
	switch v := value.(type) {
	case bool:
		sfsarray.AddBool(v)
		return nil
	case []bool:
		sfsarray.AddBoolArray(v)
		return nil
	case int8:
		sfsarray.AddByte(v)
		return nil
	case []int8:
		sfsarray.AddByteArray(v)
		return nil
	case float32:
		sfsarray.AddFloat(v)
		return nil
	case []float32:
		sfsarray.AddFloatArray(v)
		return nil
	case float64:
		sfsarray.AddDouble(v)
		return nil
	case []float64:
		sfsarray.AddDoubleArray(v)
		return nil
	case int:
		sfsarray.AddInt(int32(v))
		return nil
	case int32:
		sfsarray.AddInt(v)
		return nil
	case []int:
		// TODO
	case []int32:
		sfsarray.AddIntArray(v)
		return nil
	case int64:
		sfsarray.AddLong(v)
		return nil
	case []int64:
		sfsarray.AddLongArray(v)
		return nil
	case nil:
		sfsarray.AddNull()
		return nil
	case *SFSArray:
		sfsarray.AddSFSArray(v)
		return nil
	case SFSArray:
		sfsarray.AddSFSArray(&v)
		return nil
	case *SFSObject:
		sfsarray.AddSFSObject(v)
		return nil
	case SFSObject:
		sfsarray.AddSFSObject(&v)
		return nil
	case int16:
		sfsarray.AddShort(v)
		return nil
	case []int16:
		sfsarray.AddShortArray(v)
		return nil
	case string:
		sfsarray.AddUtfString(v)
		return nil
	case []string:
		sfsarray.AddUtfStringArray(v)
		return nil
	}
	return &ErrUnsupportedType{value: value}
}

func (sfsarray *SFSArray) AddBool(value bool) {
	sfsarray.addData(value, type_BOOL)
}

func (sfsarray *SFSArray) AddBoolArray(value []bool) {
	sfsarray.addData(value, type_BOOL_ARRAY)
}

func (sfsarray *SFSArray) AddByte(value int8) {
	sfsarray.addData(value, type_BYTE)
}

func (sfsarray *SFSArray) AddByteArray(value []int8) {
	sfsarray.addData(value, type_BYTE_ARRAY)
}

func (sfsarray *SFSArray) AddDouble(value float64) {
	sfsarray.addData(value, type_DOUBLE)
}

func (sfsarray *SFSArray) AddDoubleArray(value []float64) {
	sfsarray.addData(value, type_DOUBLE_ARRAY)
}

func (sfsarray *SFSArray) AddFloat(value float32) {
	sfsarray.addData(value, type_FLOAT)
}

func (sfsarray *SFSArray) AddFloatArray(value []float32) {
	sfsarray.addData(value, type_FLOAT_ARRAY)
}

func (sfsarray *SFSArray) AddInt(value int32) {
	sfsarray.addData(value, type_INT)
}

func (sfsarray *SFSArray) AddIntArray(value []int32) {
	sfsarray.addData(value, type_INT_ARRAY)
}

func (sfsarray *SFSArray) AddLong(value int64) {
	sfsarray.addData(value, type_LONG)
}

func (sfsarray *SFSArray) AddLongArray(value []int64) {
	sfsarray.addData(value, type_LONG_ARRAY)
}

func (sfsarray *SFSArray) AddNull() {
	sfsarray.addData(nil, type_NULL)
}

func (sfsarray *SFSArray) AddSFSArray(value *SFSArray) {
	sfsarray.addData(*value, type_SFS_ARRAY)
}

func (sfsarray *SFSArray) AddSFSObject(value *SFSObject) {
	sfsarray.addData(*value, type_SFS_OBJECT)
}

func (sfsarray *SFSArray) AddShort(value int16) {
	sfsarray.addData(value, type_SHORT)
}

func (sfsarray *SFSArray) AddShortArray(value []int16) {
	sfsarray.addData(value, type_SHORT_ARRAY)
}

func (sfsarray *SFSArray) AddUtfString(value string) {
	sfsarray.addData(value, type_UTF_STRING)
}

func (sfsarray *SFSArray) AddText(value string) {
	sfsarray.addData(value, type_TEXT)
}

func (sfsarray *SFSArray) AddUtfStringArray(value []string) {
	sfsarray.addData(value, type_UTF_STRING_ARRAY)
}
