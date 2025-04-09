package client

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"
	"reflect"

	"github.com/iancoleman/strcase"
	"gopkg.in/yaml.v3"

	"github.com/ubombar/soa/api"
	"github.com/ubombar/soa/internal/util"
)

const (
	headerSeperator = "--"
	UnknownKind     = "unknown" // special key that persists
)

var ErrCannotSaveInMemoryBuffer = errors.New("cannot save in memory buffer")

type Buffer struct {
	Content *bytes.Buffer  // Raw contents of the buffer
	Header  map[string]any // Raw header
	Origin  string
}

func (c *BufferClient) NewBuffer() *Buffer {
	return &Buffer{
		Content: new(bytes.Buffer),
		Header:  map[string]any{"kind": UnknownKind},
		Origin:  "", // empty means in memory
	}
}

func (c *BufferClient) NewBufferFromFile(filename string, create bool) (*Buffer, error) {
	var f *os.File
	var err error

	if util.FileExists(filename) {
		f, err = os.Open(filename)
	} else if create {
		f, err = os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0644)
	} else {
		return nil, os.ErrNotExist
	}

	if err != nil {
		return nil, err
	}

	defer f.Close()

	b := c.NewBuffer()

	if err := b.read(f); err != nil {
		return nil, err
	}
	b.Origin = filename // set the origin
	return b, nil
}

func (c *BufferClient) SaveBuffer(b *Buffer) error {
	if b.Origin == "" {
		return ErrCannotSaveInMemoryBuffer
	}

	f, err := os.Create(b.Origin)
	if err != nil {
		return err
	}
	defer f.Close()

	return b.write(f)
}

func (b *Buffer) read(f io.Reader) error {
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
		if lineNum == 1 && text != headerSeperator {
			noHeader = true
			inContent = true
		} else if lineNum == 1 && text == headerSeperator {
			noHeader = false
			inContent = false
			lineNum++
			continue
		} else if lineNum != 1 && text == headerSeperator && !noHeader {
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

func (b *Buffer) write(f io.Writer) error {
	writer := bufio.NewWriter(f)

	// Write first header separator
	if _, err := writer.WriteString(headerSeperator + "\n"); err != nil {
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
	if _, err := writer.WriteString(headerSeperator + "\n"); err != nil {
		return err
	}

	// Write content
	if _, err := writer.Write(b.Content.Bytes()); err != nil {
		return err
	}

	return writer.Flush()
}

func (b *Buffer) readHeader(obj any, preferStruct bool) error {
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

func (b *Buffer) writeHeader(obj any, preferMap bool, kind string) error {
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

	b.Header["kind"] = kind

	return nil
}

func GetHeader[T api.Kinder](b *Buffer) (T, error) {
	var header T
	if err := b.readHeader(&header, false); err != nil {
		return header, err
	}
	return header, nil
}

func SetHeader[T api.Kinder](b *Buffer, header T) error {
	if err := b.writeHeader(&header, false, header.Kind()); err != nil {
		return err
	}
	return nil
}
