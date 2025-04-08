package buffer

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"reflect"

	"github.com/iancoleman/strcase"
	"gopkg.in/yaml.v3"
)

const (
	HeaderSeperator = "--"
	UnknownKind     = "unknown" // special key that persists
)

type Buffer struct {
	Content *bytes.Buffer  // Raw contents of the buffer
	Header  map[string]any // Raw header
}

func NewBuffer() *Buffer {
	return &Buffer{
		Content: new(bytes.Buffer),
		Header:  map[string]any{"kind": UnknownKind},
	}
}

func FromFile(filename string) (*Buffer, error) {
	var f *os.File
	var err error

	if _, err = os.Stat(filename); os.IsNotExist(err) {
		f, err = os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0644)
	} else {
		f, err = os.Open(filename)
	}

	if err != nil {
		return nil, err
	}

	defer f.Close()

	b := NewBuffer()

	if err := b.FromReader(f); err != nil {
		fmt.Printf("err: %v\n", err)

		return nil, err
	}

	return b, nil
}

func (b *Buffer) FromReader(f io.Reader) error {
	var headerBuffer bytes.Buffer
	var contentBuffer bytes.Buffer

	scanner := bufio.NewScanner(f)
	noHeader := false
	inContent := false
	lineNum := 1

	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return err
		}

		text := scanner.Text()

		// does not contain header
		if lineNum == 1 && text != HeaderSeperator {
			noHeader = true
			inContent = true
		} else if lineNum == 1 && text == HeaderSeperator {
			noHeader = false
			inContent = false
			lineNum++
			continue
		} else if lineNum != 1 && text == HeaderSeperator && !noHeader {
			inContent = true
			lineNum++
			continue
		}

		if inContent {
			if _, err := contentBuffer.WriteString(text); err != nil {
				return err
			}
			if err := contentBuffer.WriteByte('\n'); err != nil {
				return err
			}
		} else {
			if _, err := headerBuffer.WriteString(text); err != nil {
				return err
			}
			if err := headerBuffer.WriteByte('\n'); err != nil {
				return err
			}
		}

		lineNum++
	}

	// parse header
	result := make(map[string]any)
	if err := yaml.Unmarshal(headerBuffer.Bytes(), &result); err != nil {
		return err
	}

	for k, v := range result {
		switch v.(type) {
		case map[string]interface{}, map[interface{}]interface{}:
			delete(result, k) // delete all map values
		}
	}

	// set itself
	b.Content = &contentBuffer
	b.Header = result

	// add special key "kind"
	if _, ok := b.Header["kind"]; !ok {
		b.Header["kind"] = UnknownKind
	}

	return nil
}

func ToFile(b *Buffer, filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	return b.ToWriter(f)
}

func (b *Buffer) ToWriter(f io.Writer) error {
	writer := bufio.NewWriter(f)

	// Write first header separator
	if _, err := writer.WriteString(HeaderSeperator + "\n"); err != nil {
		return err
	}

	// Marshal header to YAML and write it
	headerBytes, err := yaml.Marshal(b.Header)
	if err != nil {
		return err
	}
	if _, err := writer.Write(headerBytes); err != nil {
		return err
	}

	// Write second header separator
	if _, err := writer.WriteString(HeaderSeperator + "\n"); err != nil {
		return err
	}

	// Write content
	if _, err := writer.Write(b.Content.Bytes()); err != nil {
		return err
	}

	return writer.Flush()
}

func (b *Buffer) ReadHeader(obj any, preferStruct bool) error {
	v := reflect.ValueOf(obj)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return &reflect.ValueError{Method: "PopulateStructFromMap", Kind: v.Kind()}
	}

	v = v.Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		structField := t.Field(i)

		// Skip unexported fields
		if !field.CanSet() {
			continue
		}

		tag := structField.Tag.Get("buffer")
		if tag == "" {
			tag = strcase.ToSnake(structField.Name) // fallback to field name snake cased
		}

		if val, ok := b.Header[tag]; ok && !preferStruct {
			valValue := reflect.ValueOf(val)

			// If assignable, set it
			if valValue.Type().AssignableTo(field.Type()) {
				field.Set(valValue)
			} else if valValue.Type().ConvertibleTo(field.Type()) {
				field.Set(valValue.Convert(field.Type()))
			} else {
				in, ok := val.([]interface{})
				out := make([]string, 0, len(in))
				if !ok {
					continue // skip non-flat header
				}
				for _, v := range in {
					str, ok := v.(string)
					if !ok {
						continue
					}
					out = append(out, str)
				}
				field.Set(reflect.ValueOf(out))
			}

		}
	}

	return nil
}

func (b *Buffer) WriteHeader(obj any, preferMap bool) error {
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return &reflect.ValueError{Method: "WriteHeader", Kind: v.Kind()}
	}

	t := v.Type()
	if b.Header == nil {
		b.Header = make(map[string]any)
	}

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		structField := t.Field(i)

		// Skip unexported fields
		if structField.PkgPath != "" {
			continue
		}

		tag := structField.Tag.Get("buffer")
		if tag == "" {
			tag = strcase.ToSnake(structField.Name)
		}

		if _, ok := b.Header[tag]; !ok || (ok && !preferMap) { // not on map
			b.Header[tag] = field.Interface()
		}
	}

	return nil
}
