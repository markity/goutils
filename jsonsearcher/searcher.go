package jsonsearcher

import (
	"errors"

	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// New a json searcher. Return error when the json data is invalid
func New(data []byte) (*searcher, error) {
	obj := make(map[string]interface{})
	if err := json.Unmarshal(data, &obj); err != nil {
		return nil, err
	}
	return &searcher{obj: obj}, nil
}

// ResultType is the type of json field
type resultType int

const (
	_ resultType = iota
	TypeNumber
	TypeBool
	TypeString
	TypeArray
	TypeObject
	TypeNull
)

func (rt resultType) String() string {
	switch rt {
	// The filed is undefined
	case resultType(0):
		return "UndefinedType"
	case TypeNumber:
		return "NumberType"
	case TypeBool:
		return "BoolType"
	case TypeString:
		return "StringType"
	case TypeArray:
		return "ArrayType"
	case TypeObject:
		return "ObjectType"
	case TypeNull:
		return "NullType"
	// Unaccessable
	default:
		return "InvalidType"
	}
}

type searcher struct {
	obj map[string]interface{}
}

// Query specific json field. Args' type must be int or string(if not, the function will panic)
func (s *searcher) Query(args ...interface{}) *Result {
	result := &Result{}

	v := interface{}(s.obj)
	for _, path := range args {
		switch p := path.(type) {
		case int:
			value, ok := v.([]interface{})
			if !ok {
				return result
			}
			if len(value) < p+1 || p < 0 {
				return result
			}
			v = value[p]
		case string:
			value, ok := v.(map[string]interface{})
			if !ok {
				return result
			}
			v, ok = value[p]
			if !ok {
				return result
			}
		default:
			panic(errors.New("unexpected type"))
		}
	}

	result.exists = true
	switch v.(type) {
	case float64:
		result.resType = TypeNumber
	case bool:
		result.resType = TypeBool
	case string:
		result.resType = TypeString
	case []interface{}:
		result.resType = TypeArray
	case map[string]interface{}:
		result.resType = TypeObject
	case nil:
		result.resType = TypeNull
	}
	result.value = v

	return result
}

type Result struct {
	resType resultType
	exists  bool
	value   interface{}
}

func (r *Result) Type() resultType {
	return r.resType
}

func (r *Result) Exists() bool {
	return r.exists
}

func (r *Result) GetValue() interface{} {
	return r.value
}

func (r *Result) GetInt64() int64 {
	if r.Type() != TypeNumber {
		panic(errors.New("mistaken type of the result"))
	}
	return int64(r.value.(float64))
}

func (r *Result) GetUint64() uint64 {
	if r.Type() != TypeNumber {
		panic(errors.New("mistaken type of the result"))
	}
	return uint64(r.value.(float64))
}

func (r *Result) GetFloat64() float64 {
	if r.Type() != TypeNumber {
		panic(errors.New("mistaken type of the result"))
	}
	return r.value.(float64)
}

func (r *Result) GetBool() bool {
	if r.Type() != TypeBool {
		panic(errors.New("mistaken type of the result"))
	}
	return r.value.(bool)
}

func (r *Result) GetString() string {
	if r.Type() != TypeString {
		panic(errors.New("mistaken type of the result"))
	}
	return r.value.(string)
}

func (r *Result) GetObject() map[string]interface{} {
	if r.Type() != TypeObject {
		panic(errors.New("mistaken type of the result"))
	}
	return r.value.(map[string]interface{})
}

func (r *Result) GetArray() []interface{} {
	if r.Type() != TypeArray {
		panic(errors.New("mistaken type of the result"))
	}
	return r.value.([]interface{})
}
