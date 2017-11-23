package util

import (
	"bytes"
	"encoding/json"
	"github.com/golang/glog"
)

func ConvertByte2Json(in []byte) string {
	var out bytes.Buffer
	err := json.Indent(&out, in, "", "  ")
	if err != nil {
		glog.Errorf("json.Marshal error: err(%s)", err.Error())
		return ""
	}
	return out.String()
}

func ConvertString2Json(in string) string {
	var out bytes.Buffer
	err := json.Indent(&out, []byte(in), "", "\t")
	if err != nil {
		glog.Errorf("json.Marshal error: err(%s)", err.Error())
		return ""
	}
	return out.String()
}

func ConvertStruct2Json(in interface{}) string {
	b, err := json.Marshal(in)
	if err != nil {
		glog.Errorf("json.Marshal error: err(%s)", err.Error())
		return ""
	}

	var out bytes.Buffer
	err = json.Indent(&out, []byte(b), "", "\t")
	if err != nil {
		return ""
	}
	return out.String()
}
