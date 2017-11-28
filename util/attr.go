package util

// UNUSED. NOT REVIEW!

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net"
	"regexp"
	"unicode/utf8"
)

const (
	T_String = iota
	T_Int
	T_Bool
	T_Uuid
	T_Struct
	T_Filter
	T_Order
	T_Array
	T_Domain
	T_UrlPath
	T_Text
	T_IP
	T_UuidNull
)

var (
	uuidP     string = `^[a-z0-9]([a-z0-9-])*$`
	textP     string = `^([a-zA-Z0-9\_-]|[\p{Han}])*$`
	domainP   string = `^[a-zA-Z0-9]+([\-\.]{1}[a-z0-9]+)*\.[a-z]{2,6}$`
	urlpathP  string = `^\/[a-zA-Z0-9\/\.\%\?\#\&\=]*$`
	textRe    *regexp.Regexp
	uuidRe    *regexp.Regexp
	domainRe  *regexp.Regexp
	urlpathRe *regexp.Regexp
)

type Validater func(d interface{}) bool

type Attribute struct {
	Type     int
	Val      Validater
	Required bool
	SubAttr  interface{}
}

func validateString(d interface{}) bool {
	_, ok := d.(string)
	return ok
}

func validateBool(d interface{}) bool {
	_, ok := d.(bool)
	return ok
}

func validateInt(d interface{}) bool {
	_, ok := d.(float64)
	return ok
}

func validateUuid(d interface{}) bool {
	if id, ok := d.(string); ok {
		return uuidRe.Match([]byte(id))
	}
	return false
}

func validateUuidNull(d interface{}) bool {
	if id, ok := d.(string); ok {
		return id == "" || uuidRe.Match([]byte(id))
	}
	return false
}

func validateDomain(d interface{}) bool {
	if domain, ok := d.(string); ok {
		return domain == "" || domainRe.Match([]byte(domain))
	}
	return false
}

func validateUrlpath(d interface{}) bool {
	if path, ok := d.(string); ok {
		return urlpathRe.Match([]byte(path))
	}
	return false
}

func validateText(d interface{}) bool {
	if s, ok := d.(string); ok {
		return textRe.Match([]byte(s))
	}
	return false
}

func validateIp(d interface{}) bool {
	if s, ok := d.(string); ok {
		return net.ParseIP(s) != nil
	}
	return false
}

func ValIntRange(s, e int) Validater {
	return func(d interface{}) bool {
		df, ok := d.(float64)
		if !ok {
			return false
		}
		di := int(df)
		return di >= s && di <= e
	}
}

func ValIntEles(eles ...int) Validater {
	em := make(map[int]bool)
	for _, ele := range eles {
		em[ele] = true
	}
	return func(d interface{}) bool {
		df, ok := d.(float64)
		if !ok {
			return false
		}
		di := int(df)
		_, val := em[di]
		return val
	}
}

func ValRegExp(ptn string) Validater {
	return func(d interface{}) bool {
		if s, ok := d.(string); ok {
			matched, err := regexp.MatchString(ptn, s)
			return err == nil && matched
		}
		return false
	}
}

func ValStrLen(le int) Validater {
	return func(d interface{}) bool {
		if s, ok := d.(string); ok {
			return len(s) <= le
		}
		return false
	}
}

func ValTextLen(le int) Validater {
	return func(d interface{}) bool {
		if s, ok := d.(string); ok {
			return utf8.RuneCountInString(s) <= le
		}
		return false
	}
}

func ValStrEq(s string) Validater {
	return func(d interface{}) bool {
		if sd, ok := d.(string); ok {
			return sd == s
		}
		return false
	}
}

func MultiValid(valids ...Validater) Validater {
	return func(d interface{}) bool {
		for _, valid := range valids {
			if !valid(d) {
				return false
			}
		}
		return true
	}
}

func validateArray(a []interface{}, attr *Attribute, r string) (bool, error) {
	if attr.Required && len(a) == 0 {
		return false, nil
	}
	for _, ele := range a {
		if valid, err := validateData(ele, attr, r); !valid {
			return false, err
		}
	}
	return true, nil
}

func validateData(p interface{}, attr *Attribute, r string) (bool, error) {
	var valid bool
	var ce error = nil
	switch attr.Type {
	case T_Int:
		valid = validateInt(p)
	case T_String:
		valid = validateString(p) && (!attr.Required || p.(string) != "")
	case T_Uuid:
		valid = validateUuid(p)
	case T_UuidNull:
		valid = validateUuidNull(p)
	case T_Domain:
		valid = validateDomain(p)
	case T_UrlPath:
		valid = validateUrlpath(p)
	case T_Text:
		valid = validateText(p) && (!attr.Required || p.(string) != "")
	case T_IP:
		valid = validateIp(p)
	case T_Bool:
		valid = validateBool(p)
	case T_Struct:
		if s, ok := p.(map[string]interface{}); ok {
			valid, ce = validateStruct(s, attr.SubAttr.(map[string]*Attribute), r)
		} else {
			valid = false
		}
	case T_Filter:
		if fs, ok := p.([]interface{}); ok {
			valid, ce = validateFilter(fs, attr.SubAttr.(map[string]*Attribute), r)
		} else {
			valid = false
			ce = errors.New("validateData: invalid T_Filter")
		}
	case T_Order:
		if fs, ok := p.([]interface{}); ok {
			valid, ce = validateOrder(fs, attr.SubAttr.(map[string]*Attribute), r)
		} else {
			valid = false
			ce = errors.New("validateData: invalid T_Order")
		}

	case T_Array:
		if a, ok := p.([]interface{}); ok {
			valid, ce = validateArray(a, attr.SubAttr.(*Attribute), r)
		} else {
			valid = false
		}
	default:
		panic("Parameter Format Invalid")
	}
	if valid {
		valid = (attr.Val == nil) || attr.Val(p)
	}
	return valid, ce
}

func validateStruct(m map[string]interface{}, attrs map[string]*Attribute, r string) (bool, error) {
	for k, v := range attrs {
		if p, ok := m[k]; ok && p != nil {
			if valid, ce := validateData(p, v, k); !valid {
				if ce == nil {
					ce = errors.New("validateStruct: invalid " + k)
				}
				return false, ce
			}
		} else {
			if v.Required {
				return false, errors.New("validateStruct: miss " + k)
			}
		}
	}
	return true, nil
}

func ValiAttr(attrs map[string]*Attribute, rr io.Reader, r string, d interface{}) error {
	bt, err := ioutil.ReadAll(rr)
	if err != nil {
		return errors.New("ValiAttr: malformed " + err.Error())
	}
	var f interface{}
	if err := json.Unmarshal(bt, &f); err != nil {
		return errors.New("ValiAttr: malformed " + err.Error())
	}
	m := f.(map[string]interface{})
	if _, err := validateStruct(m, attrs, r); err != nil {
		return err
	}
	if err := json.Unmarshal(bt, d); err != nil {
		return errors.New("ValiAttr: malformed " + err.Error())
	}
	return nil
}

func validateFilter(fs []interface{}, attrs map[string]*Attribute, r string) (bool, error) {
	return validateFo(fs, attrs, r, "field", "value")
}

func validateOrder(fs []interface{}, attrs map[string]*Attribute, r string) (bool, error) {
	return validateFo(fs, attrs, r, "field", "direction")
}

func validateFo(fs []interface{}, attrs map[string]*Attribute, r string, key string, value string) (bool, error) {
	for _, f := range fs {
		if fm, ok := f.(map[string]interface{}); !ok {
			return false, errors.New("validateFo: invalid filter")
		} else {
			field, ok := fm[key]
			if !ok {
				return false, nil
			}
			value, ok := fm[value]
			if !ok {
				return false, nil
			}
			subattr, ok := attrs[field.(string)]
			if !ok {
				return false, nil
			}
			if valid, ce := validateData(value, subattr, r); !valid {
				return valid, ce
			}
		}
	}
	return true, nil
}
