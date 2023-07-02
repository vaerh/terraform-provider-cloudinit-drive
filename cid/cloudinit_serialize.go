package cid

import (
	"context"
	"encoding/base64"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func strParseSlice(s string) any {
	if len(s) == 2 {
		// []
		return []any{}
	}

	var res []any

	//[a]
	//[a, b]
	//[[]]
	//[[a, b]]
	//[[a, b], [c, d]]

	s = s[1 : len(s)-1]

	//a; 1; [
	if len(s) == 1 {
		return []any{strParse(s)}
	}

	//aaa; a, b; true, true
	if s[0] != '[' {
		for _, v := range strings.Split(s, ",") {
			res = append(res, strParse(strings.TrimSpace(v)))
		}
		return res
	}

	//[]
	if s == "[]" {
		return []any{[]any{}}
	}

	//[a, b]
	//[a, b], [c, d]
	for start, idx := 0, 1; idx < len(s); idx++ {
		if s[idx] == ']' {
			res = append(res, strParseSlice(s[start:idx+1]))
		}
		if s[idx] == '[' {
			start = idx
		}
	}

	return res
}

// strParse Because of explicit typing, we have no way to operate on value types in the YAML file.
// In general we will have unsuccessfully quoted variable values.
// To avoid this, we wrote a wrapper to get Bool, Int, Slice types from strings.
// In the TF configuration, we specify a quoted value ("true", "10", "[abc]") and we will have
// non-string types in the YAML file. For cases where we need to specify that the number is represented
// in string notation, we need to add a single quote "'...'", for example for file rights: "'0644'" -> string(0644),
// "”quoted string”" -> string('quoted string')
// For binary data we need to use the prefix "bin__"
func strParse(s string) any {
	switch {
	case strings.ToLower(s) == "true":
		return true
	case strings.ToLower(s) == "false":
		return false
	case len(s) >= 2 && s[0] == '\'' && s[len(s)-1] == '\'':
		return s[1 : len(s)-1]
	case len(s) >= 5 && s[:5] == "bin__":
		b, err := base64.StdEncoding.DecodeString(s[5:])
		if err != nil {
			return "strParse() base64 decoding error: " + err.Error()
		}
		return b
	case len(s) >= 2 && s[0] == '[' && s[len(s)-1] == ']':
		return strParseSlice(s)
	default:
		if i, e := strconv.ParseInt(s, 10, 0); e == nil {
			return i
		}
		return s
	}
}

// TerraformToGo Function to convert from TF form to Go form.
// Returns a pointer to the new structure.
func TerraformToGo(v any) any {
	return iterator(v).Addr().Interface()
}

func iterator(value any) reflect.Value {
	var newStruct []reflect.StructField
	var newStructValues []any
	var rv reflect.Value
	var rt reflect.Type

	valueOf := reflect.Indirect(reflect.ValueOf(value))
	typeOf := valueOf.Type()

	if typeOf.Kind() == reflect.Struct {
		for i := 0; i < valueOf.NumField(); i++ {
			field := typeOf.Field(i)
			//value := valueOf.Field(i),

			// Skip the specially tagged fields.
			if tag, ok := field.Tag.Lookup("cid"); ok && strings.Contains(tag, "skip") {
				continue
			}

			var v any

			switch x := valueOf.Field(i).Interface().(type) {
			case basetypes.StringValue: // types.String
				if x.IsNull() || x.IsUnknown() {
					continue
				}

				v = strParse(x.ValueString())
			case basetypes.Int64Value: // types.Int64
				if x.IsNull() || x.IsUnknown() {
					continue
				}

				v = x.ValueInt64()
			case basetypes.BoolValue: // types.Bool
				if x.IsNull() || x.IsUnknown() {
					continue
				}

				v = x.ValueBool()
			case basetypes.SetValue: // types.Set
				if x.IsNull() || x.IsUnknown() {
					continue
				}

				switch x.ElementType(context.TODO()) {
				case basetypes.StringType{}: // Element: types.String
					var s []any
					for _, e := range x.Elements() {
						s = append(s, strParse(e.(basetypes.StringValue).ValueString()))
					}
					v = s
				default:
					panic(fmt.Sprintf("unkknown Set element type: %v, %v",
						field.Name, x.ElementType(context.TODO()).String()))
				}
			case basetypes.ListValue: // types.List
				if x.IsNull() || x.IsUnknown() {
					continue
				}

				switch x.ElementType(context.TODO()) {
				case basetypes.StringType{}: // Element: types.String
					var s []any
					for _, e := range x.Elements() {
						s = append(s, strParse(e.(basetypes.StringValue).ValueString()))
					}
					v = s
				case basetypes.Int64Type{}: // Element: types.Int64
					var s []int64
					for _, e := range x.Elements() {
						s = append(s, e.(basetypes.Int64Value).ValueInt64())
					}
					v = s
				case basetypes.ListType{ElemType: basetypes.StringType{}}:
					var list = make([][]any, len(x.Elements()))
					for i, l := range x.Elements() {
						for _, e := range l.(basetypes.ListValue).Elements() {
							list[i] = append(list[i], strParse(e.(basetypes.StringValue).ValueString()))
						}
					}
					v = list
				default:
					panic(fmt.Sprintf("unkknown List element type: %v, %v",
						field.Name, x.ElementType(context.TODO()).String()))
				}
			case basetypes.MapValue: // types.Map
				if x.IsNull() || x.IsUnknown() {
					continue
				}
				switch x.ElementType(context.TODO()) {
				case basetypes.StringType{}: // Element: types.String
					var m = make(map[string]any, len(x.Elements()))
					for mKey, mVal := range x.Elements() {
						m[mKey] = strParse(mVal.(basetypes.StringValue).ValueString())
					}
					v = m
				case basetypes.ListType{ElemType: basetypes.StringType{}}: // Map[List{String}]
					var m = make(map[string]any, len(x.Elements()))
					for mKey, l := range x.Elements() {
						var list []any

						for _, e := range l.(basetypes.ListValue).Elements() {
							list = append(list, e.(basetypes.StringValue).ValueString())
						}

						m[mKey] = list
					}
					v = m
				default:
					panic(fmt.Sprintf("unkknown Map element type: %v, %v",
						field.Name, x.ElementType(context.TODO()).String()))
				}
			default:
				if valueOf.Field(i).Kind() == reflect.Struct { // type strust{}
					if valueOf.Field(i).IsZero() {
						continue
					}
				} else { // type []struct{}
					// IsNil reports whether its argument v is nil. The argument must be
					// a chan, func, interface, map, pointer, or slice value; if it is
					// not, IsNil panics.
					if valueOf.Field(i).IsNil() || valueOf.Field(i).IsZero() {
						continue
					}
				}

				if field.Type.Kind() == reflect.Ptr && field.Type.Elem().Kind() == reflect.Struct {
					// type *struct{}
					v = iterator(valueOf.Field(i).Interface()).Interface()
				} else if field.Type.Kind() == reflect.Slice && field.Type.Elem().Kind() == reflect.Struct {
					// type []struct{}
					var s []any

					for j := 0; j < valueOf.Field(i).Len(); j++ {
						rv := iterator(valueOf.Field(i).Index(j).Interface())
						s = append(s, rv.Interface())
					}
					v = s
				} else if field.Type.Kind() == reflect.Map && field.Type.Elem().Kind() == reflect.Struct {
					// map[string]struct{}
					var m = make(map[string]any)

					iter := valueOf.Field(i).MapRange()
					for iter.Next() {
						m[iter.Key().String()] = iterator(iter.Value().Interface()).Interface()
					}

					v = m
				} else {
					panic(fmt.Sprintf("unkknown Struct (%v = %v) element (%v)",
						field.Type.Kind(), field.Type.String(), field.Type.Elem().Kind()))
				}
			}

			newStruct = append(newStruct, reflect.StructField{
				Name:      field.Name,
				PkgPath:   field.PkgPath,
				Type:      reflect.TypeOf(v),
				Tag:       makeTag(field.Tag),
				Anonymous: field.Anonymous,
			})
			newStructValues = append(newStructValues, v)
		}

		rt = reflect.StructOf(newStruct)
		rv = reflect.New(rt).Elem()

		for i := 0; i < rv.NumField(); i++ {
			rv.Field(i).Set(reflect.ValueOf(newStructValues[i]))
		}
	}

	return rv
}

func makeTag(tag reflect.StructTag) reflect.StructTag {
	var tags = []string{"json", "yaml"}
	var res = make([]string, len(tags))

	for i, tagType := range tags {
		tagVal, ok := tag.Lookup(tagType)
		if ok {
			// json + :" + tag_name + "
			res[i] = tagType + `:"` + tagVal + `"`
			continue
		}
		res[i] = tagType + `:"` + tag.Get("tfsdk") + `"`
	}
	return reflect.StructTag(strings.Join(res, ","))
}
