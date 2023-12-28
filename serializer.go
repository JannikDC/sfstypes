package sfstypes

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
)

func encodeSFSObject(object *SFSObject) []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.BigEndian, byte(type_SFS_OBJECT))
	binary.Write(&buf, binary.BigEndian, int16(object.Size()))

	keys := object.GetKeys()
	for _, key := range keys {
		encodeSFSObjectKey(&buf, key)
		wrapper, _ := object.getWrapper(key)
		dataObj := wrapper.data
		encodeData(&buf, wrapper.typeId, dataObj)
	}
	return buf.Bytes()
}

func encodeSFSObjectKey(buf *bytes.Buffer, value string) {
	binary.Write(buf, binary.BigEndian, int16(len(value)))
	binary.Write(buf, binary.BigEndian, []byte(value))
}

func encodeSFSArray(array *SFSArray) []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.BigEndian, byte(type_SFS_ARRAY))
	binary.Write(&buf, binary.BigEndian, int16(array.Size()))

	for i := 0; i < array.Size(); i++ {
		wrapper, _ := array.getWrapper(i)
		encodeData(&buf, wrapper.typeId, wrapper.data)
	}
	return buf.Bytes()
}

func encodeData(buf *bytes.Buffer, typeId sfsDataType, object interface{}) {
	binary.Write(buf, binary.BigEndian, typeId)
	switch typeId {
	case type_NULL:
		{

		}
	case type_BOOL:
		if object.(bool) {
			binary.Write(buf, binary.BigEndian, byte(1))
		} else {
			binary.Write(buf, binary.BigEndian, byte(0))
		}
	case type_BYTE:
		binary.Write(buf, binary.BigEndian, object.(int8))
	case type_SHORT:
		binary.Write(buf, binary.BigEndian, object.(int16))
	case type_INT:
		binary.Write(buf, binary.BigEndian, object.(int32))
	case type_LONG:
		binary.Write(buf, binary.BigEndian, object.(int64))
	case type_FLOAT:
		binary.Write(buf, binary.BigEndian, object.(float32))
	case type_DOUBLE:
		binary.Write(buf, binary.BigEndian, object.(float64))
	case type_UTF_STRING:
		str := object.(string)
		binary.Write(buf, binary.BigEndian, int16(len(str)))
		binary.Write(buf, binary.BigEndian, []byte(str))
	case type_TEXT:
		str := object.(string)
		binary.Write(buf, binary.BigEndian, int32(len(str)))
		binary.Write(buf, binary.BigEndian, []byte(str))
	case type_BOOL_ARRAY:
		array := object.([]bool)
		binary.Write(buf, binary.BigEndian, int16(len(array)))
		for _, element := range array {
			if element {
				binary.Write(buf, binary.BigEndian, byte(1))
			} else {
				binary.Write(buf, binary.BigEndian, byte(0))
			}
		}
	case type_BYTE_ARRAY:
		array := object.([]int8)
		binary.Write(buf, binary.BigEndian, int32(len(array)))
		for _, element := range array {
			binary.Write(buf, binary.BigEndian, element)
		}
	case type_SHORT_ARRAY:
		array := object.([]int16)
		binary.Write(buf, binary.BigEndian, int16(len(array)))
		for _, element := range array {
			binary.Write(buf, binary.BigEndian, element)
		}
	case type_INT_ARRAY:
		array := object.([]int32)
		binary.Write(buf, binary.BigEndian, int16(len(array)))
		for _, element := range array {
			binary.Write(buf, binary.BigEndian, element)
		}
	case type_LONG_ARRAY:
		array := object.([]int64)
		binary.Write(buf, binary.BigEndian, int16(len(array)))
		for _, element := range array {
			binary.Write(buf, binary.BigEndian, element)
		}
	case type_FLOAT_ARRAY:
		array := object.([]float32)
		binary.Write(buf, binary.BigEndian, int16(len(array)))
		for _, element := range array {
			binary.Write(buf, binary.BigEndian, element)
		}
	case type_DOUBLE_ARRAY:
		array := object.([]float64)
		binary.Write(buf, binary.BigEndian, int16(len(array)))
		for _, element := range array {
			binary.Write(buf, binary.BigEndian, element)
		}
	case type_UTF_STRING_ARRAY:
		array := object.([]string)
		binary.Write(buf, binary.BigEndian, int16(len(array)))
		for _, element := range array {
			binary.Write(buf, binary.BigEndian, int16(len(element)))
			binary.Write(buf, binary.BigEndian, []byte(element))
		}
	case type_SFS_ARRAY:
		buf.Truncate(buf.Len() - 1)
		arr := object.(SFSArray)
		binary.Write(buf, binary.BigEndian, encodeSFSArray(&arr))
	case type_SFS_OBJECT:
		buf.Truncate(buf.Len() - 1)
		obj := object.(SFSObject)
		binary.Write(buf, binary.BigEndian, encodeSFSObject(&obj))
		//case CLASS:
		//	addData(buffer, object2binary(pojo2sfs(object)));
	}
}

func sfsObjectToJson(sfsobject *SFSObject) string {
	return dataToJson(convertSFSObjectToMap(sfsobject))
}

func sfsArrayToJson(sfsarray *SFSArray) string {
	return dataToJson(convertSFSArrayToSlice(sfsarray))
}

func sfsObjectFromJson(input string) (*SFSObject, error) {
	dataMap, err := convertJsonToMap(input)
	if err != nil {
		return nil, err
	}
	return convertMapToSFSObject(dataMap)
}

func sfsArrayFromJson(input string) (*SFSArray, error) {
	dataArray, err := convertJsonToSlice(input)
	if err != nil {
		return nil, err
	}
	return convertSliceToSFSArray(dataArray)
}

func dataToJson(input interface{}) string {
	jsonData, _ := json.MarshalIndent(input, "", "    ")
	return string(jsonData)
}

func convertJsonToMap(jsonStr string) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	err := json.Unmarshal([]byte(jsonStr), &result)
	return result, err
}

func convertJsonToSlice(jsonStr string) ([]interface{}, error) {
	result := make([]interface{}, 0)
	err := json.Unmarshal([]byte(jsonStr), &result)
	return result, err
}

func convertSFSObjectToMap(sfsobject *SFSObject) map[string]interface{} {
	result := make(map[string]interface{})
	keys := sfsobject.GetKeys()
	for _, key := range keys {
		currentData := sfsobject.dataHolder[key]
		switch currentData.typeId {
		case type_SFS_OBJECT:
			obj := currentData.data.(SFSObject)
			result[key] = convertSFSObjectToMap(&obj)
		case type_SFS_ARRAY:
			arr := currentData.data.(SFSArray)
			result[key] = convertSFSArrayToSlice(&arr)
		default:
			result[key] = currentData.data
		}
	}
	return result
}

func convertSFSArrayToSlice(sfsarray *SFSArray) []interface{} {
	result := make([]interface{}, 0)
	for index := range sfsarray.dataHolder {
		currentData := sfsarray.dataHolder[index]

		switch currentData.typeId {
		case type_SFS_OBJECT:
			obj := currentData.data.(SFSObject)
			result = append(result, convertSFSObjectToMap(&obj))
		case type_SFS_ARRAY:
			arr := currentData.data.(SFSArray)
			result = append(result, convertSFSArrayToSlice(&arr))
		default:
			result = append(result, currentData.data)
		}
	}
	return result
}

func convertMapToSFSObject(input map[string]interface{}) (*SFSObject, error) {
	sfsObject := NewSFSObject()
	for key, element := range input {
		switch el := element.(type) {
		case map[string]interface{}:
			obj, err1 := convertMapToSFSObject(el)
			if err1 != nil {
				return nil, err1
			}
			err2 := sfsObject.Put(key, obj)
			if err2 != nil {
				return nil, err2
			}
		case []interface{}:
			arr, err1 := convertSliceToSFSArray(el)
			if err1 != nil {
				return nil, err1
			}
			err2 := sfsObject.Put(key, arr)
			if err2 != nil {
				return nil, err2
			}
		default:
			err := sfsObject.Put(key, el)
			if err != nil {
				return nil, err
			}
		}
	}
	return sfsObject, nil
}

func convertSliceToSFSArray(input []interface{}) (*SFSArray, error) {
	sfsArray := NewSFSArray()
	for _, element := range input {
		switch el := element.(type) {
		case map[string]interface{}:
			obj, err1 := convertMapToSFSObject(el)
			if err1 != nil {
				return nil, err1
			}
			err2 := sfsArray.Add(obj)
			if err2 != nil {
				return nil, err2
			}
		case []interface{}:
			arr, err1 := convertSliceToSFSArray(el)
			if err1 != nil {
				return nil, err1
			}
			err2 := sfsArray.Add(arr)
			if err2 != nil {
				return nil, err2
			}
		default:
			err := sfsArray.Add(el)
			if err != nil {
				return nil, err
			}
		}
	}
	return sfsArray, nil
}

func newSFSObjectfromBinary(data []byte) (*SFSObject, error) {
	if size := len(data); size < 3 {
		return nil, &ErrInsufficientByteData{sfsType: type_SFS_OBJECT, size: len(data)}
	}
	buf := bytes.NewBuffer(data)
	return decodeSFSObject(buf)
}

func newSFSArrayFromBinaryData(data []byte) (*SFSArray, error) {
	if size := len(data); size < 3 {
		return nil, &ErrInsufficientByteData{sfsType: type_SFS_ARRAY, size: len(data)}
	}
	buf := bytes.NewBuffer(data)
	return decodeSFSArray(buf)
}

func decodeSFSObject(buf *bytes.Buffer) (*SFSObject, error) {
	sfsObject := NewSFSObject()

	var header byte
	if err := binary.Read(buf, binary.BigEndian, &header); err != nil {
		return nil, &ErrReadingData{TypeToRead: "value header", Len: buf.Len(), Cap: buf.Cap(), IoErr: err}
	} else if header != byte(type_SFS_OBJECT) {
		return nil, &ErrWrongType{actualType: sfsDataType(header), wantedType: type_SFS_OBJECT}
	}

	var size uint16
	if err := binary.Read(buf, binary.BigEndian, &size); err != nil {
		return nil, &ErrReadingData{TypeToRead: "SFSObject size", Len: buf.Len(), Cap: buf.Cap(), IoErr: err}
	}
	for i := uint16(0); i < size; i++ {

		var keySize uint16
		if err := binary.Read(buf, binary.BigEndian, &keySize); err != nil {
			return nil, &ErrReadingData{TypeToRead: "key size", Len: buf.Len(), Cap: buf.Cap(), IoErr: err}
		}
		if keySize > 255 {
			return nil, &ErrInvalidKeySize{key: "", length: int(keySize)}
		}
		keyStringBytes := make([]byte, keySize)
		if err := binary.Read(buf, binary.BigEndian, &keyStringBytes); err != nil {
			return nil, err
		}
		key := string(keyStringBytes)

		decodedObject, decodeError := decodeData(buf)
		if decodeError != nil {
			return nil, decodeError
		}
		sfsObject.putsfsDataWrapper(key, decodedObject)
	}
	return sfsObject, nil
}

func decodeSFSArray(buf *bytes.Buffer) (*SFSArray, error) {
	sfsArray := NewSFSArray()

	var header byte
	if err := binary.Read(buf, binary.BigEndian, &header); err != nil {
		return nil, &ErrReadingData{TypeToRead: "SFSArry header", Len: buf.Len(), Cap: buf.Cap(), IoErr: err}
	} else if header != byte(type_SFS_ARRAY) {
		return nil, &ErrWrongType{actualType: sfsDataType(header), wantedType: type_SFS_ARRAY}
	}

	var size uint16
	if err := binary.Read(buf, binary.BigEndian, &size); err != nil {
		return nil, &ErrReadingData{TypeToRead: "SFSObject size", Len: buf.Len(), Cap: buf.Cap(), IoErr: err}
	}

	for i := uint16(0); i < size; i++ {
		wrapper, err := decodeData(buf)
		if err != nil {
			return nil, err
		}
		sfsArray.addsfsDataWrapper(*wrapper)
	}

	return sfsArray, nil
}

func decodeData(buf *bytes.Buffer) (*sfsDataWrapper, error) {
	var header byte
	if err := binary.Read(buf, binary.BigEndian, &header); err != nil {
		return nil, &ErrReadingData{TypeToRead: "value header", Len: buf.Len(), Cap: buf.Cap(), IoErr: err}
	}
	switch sfsDataType(header) {
	case type_NULL:
		return newsfsDataWrapper(type_NULL, nil), nil
	case type_BOOL:
		var input byte
		if err := binary.Read(buf, binary.BigEndian, &input); err != nil {
			return nil, &ErrReadingData{TypeToRead: "bool", Len: buf.Len(), Cap: buf.Cap(), IoErr: err}
		}
		switch input {
		case 0:
			return newsfsDataWrapper(type_BOOL, false), nil
		case 1:
			return newsfsDataWrapper(type_BOOL, true), nil
		default:
			return nil, ErrDecodingBool
		}
	case type_BYTE:
		var input int8
		if err := binary.Read(buf, binary.BigEndian, &input); err != nil {
			return nil, &ErrReadingData{TypeToRead: "byte", Len: buf.Len(), Cap: buf.Cap(), IoErr: err}
		}
		return newsfsDataWrapper(type_BYTE, input), nil
	case type_SHORT:
		var input int16
		if err := binary.Read(buf, binary.BigEndian, &input); err != nil {
			return nil, &ErrReadingData{TypeToRead: "short", Len: buf.Len(), Cap: buf.Cap(), IoErr: err}
		}
		return newsfsDataWrapper(type_SHORT, input), nil
	case type_INT:
		var input int32
		if err := binary.Read(buf, binary.BigEndian, &input); err != nil {
			return nil, &ErrReadingData{TypeToRead: "int", Len: buf.Len(), Cap: buf.Cap(), IoErr: err}
		}
		return newsfsDataWrapper(type_INT, input), nil
	case type_LONG:
		var input int64
		if err := binary.Read(buf, binary.BigEndian, &input); err != nil {
			return nil, &ErrReadingData{TypeToRead: "long", Len: buf.Len(), Cap: buf.Cap(), IoErr: err}
		}
		return newsfsDataWrapper(type_LONG, input), nil
	case type_FLOAT:
		var input float32
		if err := binary.Read(buf, binary.BigEndian, &input); err != nil {
			return nil, &ErrReadingData{TypeToRead: "float", Len: buf.Len(), Cap: buf.Cap(), IoErr: err}
		}
		return newsfsDataWrapper(type_FLOAT, input), nil
	case type_DOUBLE:
		var input float64
		if err := binary.Read(buf, binary.BigEndian, &input); err != nil {
			return nil, &ErrReadingData{TypeToRead: "double", Len: buf.Len(), Cap: buf.Cap(), IoErr: err}
		}
		return newsfsDataWrapper(type_DOUBLE, input), nil
	case type_UTF_STRING:
		var len uint16
		if err := binary.Read(buf, binary.BigEndian, &len); err != nil {
			return nil, &ErrReadingData{TypeToRead: "string length", Len: buf.Len(), Cap: buf.Cap(), IoErr: err}
		}
		stringBytes := make([]byte, len)
		if err := binary.Read(buf, binary.BigEndian, &stringBytes); err != nil {
			return nil, &ErrReadingData{TypeToRead: "string bytes", Len: buf.Len(), Cap: buf.Cap(), IoErr: err}
		}
		decodedString := string(stringBytes)
		return newsfsDataWrapper(type_UTF_STRING, decodedString), nil
	case type_TEXT:
		var len uint32
		if err := binary.Read(buf, binary.BigEndian, &len); err != nil {
			return nil, &ErrReadingData{TypeToRead: "text length", Len: buf.Len(), Cap: buf.Cap(), IoErr: err}
		}
		stringBytes := make([]byte, len)
		if err := binary.Read(buf, binary.BigEndian, &stringBytes); err != nil {
			return nil, &ErrReadingData{TypeToRead: "text bytes", Len: buf.Len(), Cap: buf.Cap(), IoErr: err}
		}
		decodedString := string(stringBytes)
		return newsfsDataWrapper(type_TEXT, decodedString), nil
	case type_BOOL_ARRAY:
		var length uint16
		if err := binary.Read(buf, binary.BigEndian, &length); err != nil {
			return nil, &ErrReadingData{TypeToRead: "bool array length", Len: buf.Len(), Cap: buf.Cap(), IoErr: err}
		}

		results := make([]bool, length)
		for i := uint16(0); i < length; i++ {
			var tempValue byte
			if err := binary.Read(buf, binary.BigEndian, &tempValue); err != nil {
				return nil, &ErrReadingData{TypeToRead: "bool array element", Len: buf.Len(), Cap: buf.Cap(), IoErr: err}
			}
			switch tempValue {
			case 0:
				results[i] = false
			case 1:
				results[i] = true
			default:
				return nil, ErrDecodingBool
			}
		}
		return newsfsDataWrapper(type_BOOL_ARRAY, results), nil
	case type_BYTE_ARRAY:
		var len uint32
		if err := binary.Read(buf, binary.BigEndian, &len); err != nil {
			return nil, &ErrReadingData{TypeToRead: "byte array length", Len: buf.Len(), Cap: buf.Cap(), IoErr: err}
		}

		results := make([]int8, len)
		if err := binary.Read(buf, binary.BigEndian, &results); err != nil {
			return nil, &ErrReadingData{TypeToRead: "byte array", Len: buf.Len(), Cap: buf.Cap(), IoErr: err}
		}
		return newsfsDataWrapper(type_BYTE_ARRAY, results), nil
	case type_SHORT_ARRAY:
		var len uint16
		if err := binary.Read(buf, binary.BigEndian, &len); err != nil {
			return nil, &ErrReadingData{TypeToRead: "short array size", Len: buf.Len(), Cap: buf.Cap(), IoErr: err}
		}

		results := make([]int16, len)
		if err := binary.Read(buf, binary.BigEndian, &results); err != nil {
			return nil, &ErrReadingData{TypeToRead: "short array", Len: buf.Len(), Cap: buf.Cap(), IoErr: err}
		}
		return newsfsDataWrapper(type_SHORT_ARRAY, results), nil
	case type_INT_ARRAY:
		var len uint16
		if err := binary.Read(buf, binary.BigEndian, &len); err != nil {
			return nil, &ErrReadingData{TypeToRead: "int array size", Len: buf.Len(), Cap: buf.Cap(), IoErr: err}
		}

		results := make([]int32, len)
		if err := binary.Read(buf, binary.BigEndian, &results); err != nil {
			return nil, &ErrReadingData{TypeToRead: "int array", Len: buf.Len(), Cap: buf.Cap(), IoErr: err}
		}
		return newsfsDataWrapper(type_INT_ARRAY, results), nil
	case type_LONG_ARRAY:
		var len uint16
		if err := binary.Read(buf, binary.BigEndian, &len); err != nil {
			return nil, &ErrReadingData{TypeToRead: "long array size", Len: buf.Len(), Cap: buf.Cap(), IoErr: err}
		}

		results := make([]int64, len)
		if err := binary.Read(buf, binary.BigEndian, &results); err != nil {
			return nil, &ErrReadingData{TypeToRead: "long array", Len: buf.Len(), Cap: buf.Cap(), IoErr: err}
		}
		return newsfsDataWrapper(type_LONG_ARRAY, results), nil
	case type_FLOAT_ARRAY:
		var len uint16
		if err := binary.Read(buf, binary.BigEndian, &len); err != nil {
			return nil, &ErrReadingData{TypeToRead: "float array size", Len: buf.Len(), Cap: buf.Cap(), IoErr: err}
		}

		results := make([]float32, len)
		if err := binary.Read(buf, binary.BigEndian, &results); err != nil {
			return nil, &ErrReadingData{TypeToRead: "float array", Len: buf.Len(), Cap: buf.Cap(), IoErr: err}
		}
		return newsfsDataWrapper(type_FLOAT_ARRAY, results), nil
	case type_DOUBLE_ARRAY:
		var len uint16
		if err := binary.Read(buf, binary.BigEndian, &len); err != nil {
			return nil, &ErrReadingData{TypeToRead: "double array size", Len: buf.Len(), Cap: buf.Cap(), IoErr: err}
		}

		results := make([]float64, len)
		if err := binary.Read(buf, binary.BigEndian, &results); err != nil {
			return nil, &ErrReadingData{TypeToRead: "double array", Len: buf.Len(), Cap: buf.Cap(), IoErr: err}
		}
		return newsfsDataWrapper(type_DOUBLE_ARRAY, results), nil
	case type_UTF_STRING_ARRAY:
		var len uint16
		if err := binary.Read(buf, binary.BigEndian, &len); err != nil {
			return nil, &ErrReadingData{TypeToRead: "string array length", Len: buf.Len(), Cap: buf.Cap(), IoErr: err}
		}

		results := make([]string, 0)
		for i := uint16(0); i < len; i++ {
			var strLen uint16
			if err := binary.Read(buf, binary.BigEndian, &strLen); err != nil {
				return nil, &ErrReadingData{TypeToRead: "string array element length", Len: buf.Len(), Cap: buf.Cap(), IoErr: err}
			}
			stringBytes := make([]byte, strLen)
			if err := binary.Read(buf, binary.BigEndian, &stringBytes); err != nil {
				return nil, &ErrReadingData{TypeToRead: "string array element bytes", Len: buf.Len(), Cap: buf.Cap(), IoErr: err}
			}
			results = append(results, string(stringBytes))
		}
		return newsfsDataWrapper(type_UTF_STRING_ARRAY, results), nil
	case type_SFS_ARRAY:
		buf.UnreadByte()
		obj, err := decodeSFSArray(buf)
		if err != nil {
			return nil, err
		}
		return newsfsDataWrapper(type_SFS_ARRAY, *obj), nil
	case type_SFS_OBJECT:
		buf.UnreadByte()
		obj, err := decodeSFSObject(buf)
		if err != nil {
			return nil, err
		}
		return newsfsDataWrapper(type_SFS_OBJECT, *obj), nil
	default: // sfsDataType(header)
		return nil, &ErrDecodingUnsupportedType{sfsType: sfsDataType(header), Len: buf.Len(), Cap: buf.Cap()}
	}
}
