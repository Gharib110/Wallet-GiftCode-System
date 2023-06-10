package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"reflect"
)

// dResponseWriter use for writing response to the user
func (c *configAPI) dResponseWriter(w http.ResponseWriter, data interface{}, HStat int) error {
	dataType := reflect.TypeOf(data)
	if dataType.Kind() == reflect.String {
		w.WriteHeader(HStat)
		w.Header().Set("Content-Type", "application/text")

		_, err := w.Write([]byte(data.(string)))
		return err
	} else if reflect.PtrTo(dataType).Kind() == dataType.Kind() {
		w.WriteHeader(HStat)
		w.Header().Set("Content-Type", "application/json")

		outData, err := json.MarshalIndent(data, "", "\t")
		if err != nil {
			c.errorLogger.Println(err.Error())
			w.Write([]byte(err.Error()))
			return err
		}

		_, err = w.Write(outData)
		return err
	} else if reflect.Struct == dataType.Kind() {
		w.WriteHeader(HStat)
		w.Header().Set("Content-Type", "application/json")

		outData, err := json.MarshalIndent(data, "", "\t")
		if err != nil {
			c.errorLogger.Println(err.Error())
			w.Write([]byte(err.Error()))
			return err
		}

		_, err = w.Write(outData)
		return err
	} else if reflect.Slice == dataType.Kind() {
		w.WriteHeader(HStat)
		w.Header().Set("Content-Type", "application/json")

		outData, err := json.MarshalIndent(data, "", "\t")
		if err != nil {
			c.errorLogger.Println(err.Error())
			w.Write([]byte(err.Error()))
			return err
		}

		_, err = w.Write(outData)
		return err
	}

	return errors.New("we could not be able to support data type that you passed")
}
