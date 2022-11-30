package input

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"strings"
)

func formatTitle(w io.Writer, s, field any, entry SchemaEntry) error {
	fmt.Fprintf(w, "%s\n", field)
	return nil
}

func formatOutList(w io.Writer, s, field any, entry SchemaEntry) error {
	fmt.Fprintf(w, "%12s    %-14s - %s\n", "", entry.Field, entry.Desc)
	for _, v := range field.([]string) {
		fmt.Fprintf(w, "%-12s\n", `"`+v+`"`)
	}
	w.Write([]byte("END of input file (\"END\" must appear " +
		"in the first 3 columns of this last OutList line)\n"))
	w.Write(bytes.Repeat([]byte("-"), 80))
	w.Write([]byte("\n"))
	return nil
}

func formatText(s any, schema Schema) ([]byte, error) {

	w := &bytes.Buffer{}
	sVal := reflect.ValueOf(s).Elem()
	defaults := sVal.FieldByName("Defaults").Interface().(map[string]struct{})

	// Loop through entries in schema
	for _, entry := range schema {

		fieldVal := sVal.FieldByName(entry.Field)

		// Write heading
		if entry.Heading == "-" {
			w.WriteString("\n")
			continue
		} else if entry.Heading != "" {
			w.WriteString(strings.Repeat("-", 15) + " " + entry.Heading + " " +
				strings.Repeat("-", 132-15-len(entry.Heading)) + "\n")
			continue
		}

		// If entry has format function
		if entry.Format != nil {
			if err := entry.Format(w, s, fieldVal.Interface(), entry); err != nil {
				return nil, err
			}
			continue
		}

		// Build description from desc, units, and options
		desc := entry.Desc
		if entry.Unit != "" {
			desc += " (" + entry.Unit + ")"
		}
		if entry.Show != nil {
			desc += " [used if"
			for _, cond := range entry.Show {
				desc += fmt.Sprintf(" %s %s %v and", cond.Field, cond.Relation, cond.Value)
			}
			desc = desc[:len(desc)-4] + "]"
		}
		if entry.Options != nil {
			desc += " {"
			for _, opt := range entry.Options {
				if opt.Value == opt.Text {
					desc += opt.Text + "; "
				} else {
					desc += fmt.Sprintf("%v=%s; ", opt.Value, opt.Text)
				}
			}
			desc = desc[:len(desc)-2] + "}"
		}

		// If optional field and value is nil, print default
		if entry.CanBeDefault {
			if _, ok := defaults[entry.Keyword]; ok {
				fmt.Fprintf(w, "%12s    %-14s - %s\n", `"default"`, entry.Keyword, desc)
				continue
			}
		}

		switch entry.Type {
		case String:
			fmt.Fprintf(w, "%12s    %-14s - %s\n", `"`+fieldVal.String()+`"`, entry.Keyword, desc)
		// case FieldFloat:
		// 	fmt.Fprintf(w, "%12g    %-14s - %s\n", fieldVal.Interface(), entry.Keyword, desc)
		// case FieldInt:
		// 	fmt.Fprintf(w, "%12d    %-14s - %s\n", fieldVal.Interface(), entry.Keyword, desc)
		// case FieldBool:
		// 	fmt.Fprintf(w, "%12t    %-14s - %s\n", fieldVal.Interface(), entry.Keyword, desc)
		// case FieldSlice:
		default:
			switch entry.Dims {
			case 0:
				fmt.Fprintf(w, "%12v    %-14s - %s\n", fieldVal.Interface(), entry.Keyword, desc)
			case 1:
				tmp := fmt.Sprintf("%v", fieldVal.Interface())
				fmt.Fprintf(w, "%12s    %-14s - %s\n", tmp[1:len(tmp)-1], entry.Keyword, desc)
			}
			// return nil, fmt.Errorf("formatting for %s not implemented", entry.Keyword)
		}
	}

	return w.Bytes(), nil
}
