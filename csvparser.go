//Package provides simple and noobie parsing of each csv line into struct
//Supports only strings because I didn't need others :)
//See simple example in _example directory
package csvparser

import (
	"encoding/csv"
	"os"
	"reflect"
)

type Reader struct {
	file        *os.File
	headers     []string
	type_struct reflect.Type
	reader      *csv.Reader
}

func (r *Reader) getHeaders() (err error) {
	r.headers, err = r.reader.Read()
	if err != nil {
		return err
	}
	return
}

// Parses each line of csv file and passing interface{} to callback function
func ParseEach(file string, v interface{}, callback func(diag interface{})) error {
	var err error
	r := new(Reader)
	dataType := reflect.TypeOf(v)
	newData := reflect.New(dataType).Elem()
	if r.file, err = os.Open(file); err != nil {
		return
	}
	defer r.file.Close()
	r.reader = csv.NewReader(r.file)
	r.reader.FieldsPerRecord = -1
	r.headers = make([]string, 0)
	err = r.getHeaders()
	if err != nil {
		return
	}

	for {
		row, err := r.reader.Read()
		if err != nil {
			break
		}
		for i := 0; i < dataType.NumField(); i++ {
			f := dataType.Field(i)
			index := 0
			field_name := f.Tag.Get("csv")
			for k, v := range r.headers {
				if v == field_name {
					index = k
				}
			}
			new_field := newData.FieldByName(f.Name)
			if new_field.IsValid() {
				if new_field.CanSet() {
					new_field.Set(reflect.ValueOf(row[index]))
				}
			}
		}
		callback(newData.Interface())
	}
	return nil
}
