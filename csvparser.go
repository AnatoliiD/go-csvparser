//Package csvparser provides simple and noobie parsing of each csv line into struct
//Supports only strings because I didn't need others :)
//See simple example in _example directory
package csvparser

import (
	"encoding/csv"
	"os"
	"reflect"
)

type reader struct {
	file    *os.File
	headers []string
	reader  *csv.Reader
}

func (r *reader) getHeaders() (err error) {
	r.headers, err = r.reader.Read()
	if err != nil {
		return err
	}
	return
}

// ParseEach parses each line of csv file and passing interface{} to callback function
func ParseEach(file string, v interface{}, callback func(result interface{})) error {
	var err error
	r := new(reader)
	dataType := reflect.TypeOf(v)
	newData := reflect.New(dataType).Elem()
	if r.file, err = os.Open(file); err != nil {
		return err
	}
	defer r.file.Close()
	r.reader = csv.NewReader(r.file)
	r.reader.FieldsPerRecord = -1
	r.headers = make([]string, 0)
	err = r.getHeaders()
	if err != nil {
		return err
	}

	for {
		row, err := r.reader.Read()
		if err != nil {
			break
		}
		for i := 0; i < dataType.NumField(); i++ {
			f := dataType.Field(i)
			index := 0
			fieldName := f.Tag.Get("csv")
			for k, v := range r.headers {
				if v == fieldName {
					index = k
				}
			}
			newField := newData.FieldByName(f.Name)
			if newField.IsValid() {
				if newField.CanSet() {
					newField.Set(reflect.ValueOf(row[index]))
				}
			}
		}
		callback(newData.Interface())
	}
	return nil
}
