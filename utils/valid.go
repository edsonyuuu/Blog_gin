package utils

import (
	"github.com/go-playground/validator/v10"
	"reflect"
)

func GetValidMsg(err error, obj any) string {
	//注意传obj的指针
	getObj := reflect.TypeOf(obj)
	//
	if errs, ok := err.(validator.ValidationErrors); ok {
		//
		for _, e := range errs {
			//根据错误，获取结构体的具体字段
			if f, exits := getObj.Elem().FieldByName(e.Field()); exits {
				msg := f.Tag.Get("msg")
				return msg
			}
		}
	}
	return err.Error()
}
