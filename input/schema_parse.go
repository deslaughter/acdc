package input

import (
	"bufio"
	"bytes"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func splitLines(text []byte) ([]string, error) {
	lines := make([]string, 0, 256)
	scanner := bufio.NewScanner(bytes.NewReader(text))
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error scanning text: %w", err)
	}
	return lines, nil
}

func findKeywordLine(keyword string, lines []string) (string, []string, error) {
	line := ""
	keywordLower := strings.ToLower(keyword)
	for {
		if len(lines) == 0 {
			return "", []string{}, fmt.Errorf("label %s not found", keyword)
		}
		line, lines = lines[0], lines[1:]
		if strings.Contains(strings.ToLower(line), keywordLower) {
			break
		}
	}
	return line, lines, nil
}

func parseText(s any, schema Schema, lines []string) error {

	sVal := reflect.ValueOf(s).Elem()
	defaults := sVal.FieldByName("Defaults").Interface().(map[string]struct{})

	// Loop through entries in schema
	for _, entry := range schema {

		// If entry is heading, continue
		if entry.Heading != "" {
			continue
		}

		// Get field value
		fieldVal := sVal.FieldByName(entry.Field)

		// If function given for parsing this keyword, use it instead
		if entry.Parse != nil {
			if ls, err := entry.Parse(s, fieldVal, lines); err != nil {
				return err
			} else {
				lines = ls
			}
			continue
		}

		keyword := entry.Keyword
		keywordLower := strings.ToLower(keyword)

		// Loop through lines until field keyword is found
		line := ""
		for {
			line, lines = lines[0], lines[1:]
			line, _, _ = strings.Cut(line, "- ")
			ind := strings.Index(strings.ToLower(line), keywordLower)
			if ind != -1 {
				line = strings.TrimSpace(line[:ind])
				break
			}
			if len(lines) == 0 {
				return fmt.Errorf("keyword '%s' not found", keyword)
			}
		}

		// If entry is optional and value is "default"
		if entry.CanBeDefault && strings.Contains(strings.ToLower(line), "default") {
			defaults[entry.Keyword] = struct{}{}
			continue
		}

		// Split line in elements
		elems := strings.FieldsFunc(line, func(r rune) bool {
			return r == ' ' || r == ','
		})

		// Switch based on field kind
		switch entry.Type {

		case String:
			fieldVal.SetString(strings.Trim(line, `"`))

		case Bool:
			values := make([]bool, 0, len(elems))
			for _, elem := range elems {
				v, err := strconv.ParseBool(elem)
				if err != nil {
					return fmt.Errorf("invalid value '%s' for boolean field '%s'", line, keyword)
				} else {
					values = append(values, v)
				}
			}
			if entry.Dims == 0 {
				fieldVal.SetBool(values[0])
			} else {
				fieldVal.Set(reflect.ValueOf(values))
			}

		case Int:
			values := make([]int, 0, len(elems))
			for _, elem := range elems {
				if v, err := strconv.ParseInt(elem, 0, 0); err != nil {
					return fmt.Errorf("invalid value '%s' for int field '%s'", line, keyword)
				} else {
					values = append(values, int(v))
				}
			}
			if entry.Dims == 0 {
				fieldVal.SetInt(int64(values[0]))
			} else {
				fieldVal.Set(reflect.ValueOf(values))
			}

		case Float:
			values := make([]float64, 0, len(elems))
			for _, elem := range elems {
				if v, err := strconv.ParseFloat(elem, 64); err != nil {
					return fmt.Errorf("invalid value '%s' for float64 field '%s'", line, keyword)
				} else {
					values = append(values, v)
				}
			}
			if entry.Dims == 0 {
				fieldVal.SetFloat(values[0])
			} else {
				fieldVal.Set(reflect.ValueOf(values))
			}

		default:
			return fmt.Errorf("parsing of field '%s' not implemented", entry.Keyword)
		}
	}

	return nil
}

func parseTitle(s any, v reflect.Value, lines []string) ([]string, error) {
	v.SetString(lines[1])
	return lines[2:], nil
}

func parseOutList(s any, v reflect.Value, lines []string) ([]string, error) {
	line, lines, err := findKeywordLine("OutList", lines)
	if err != nil {
		return nil, err
	}
	_ = line
	vars := []string{}
	for {
		line, lines = lines[0], lines[1:]
		if strings.HasPrefix(strings.ToLower(line), "end") {
			break
		}
		if i := strings.LastIndex(line, `"`); i != -1 {
			line = line[:i]
		}
		tmp := strings.FieldsFunc(line, func(r rune) bool {
			return r == '"' || r == ',' || r == ' ' || r == '\t'
		})
		if len(tmp) == 0 {
			continue
		}
		vars = append(vars, tmp...)
	}
	v.Set(reflect.ValueOf(vars))
	return lines, nil
}
