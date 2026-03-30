package gocty

import (
	"reflect"

	"github.com/zclconf/go-cty/cty"
)

// maxImpliedTypeDepth prevents infinite recursion in self-referential types.
// Matches Terraform's recursion limit.
const maxImpliedTypeDepth = 10

// ImpliedType takes an arbitrary Go value (as an interface{}) and attempts
// to find a suitable cty.Type instance that could be used for a conversion
// with ToCtyValue.
//
// This allows -- for simple situations at least -- types to be defined just
// once in Go and the cty types derived from the Go types, but in the process
// it makes some assumptions that may be undesirable so applications are
// encouraged to build their cty types directly if exacting control is
// required.
//
// Not all Go types can be represented as cty types, so an error may be
// returned which is usually considered to be a bug in the calling program.
// In particular, ImpliedType will never use capsule types in its returned
// type, because it cannot know the capsule types supported by the calling
// program.
func ImpliedType(gv interface{}) (cty.Type, error) {
	rt := reflect.TypeOf(gv)
	var path cty.Path
	return impliedTypeWithDepth(rt, path, 0)
}

func impliedType(rt reflect.Type, path cty.Path) (cty.Type, error) {
	// Keep existing function for backward compatibility
	return impliedTypeWithDepth(rt, path, 0)
}

func impliedTypeWithDepth(rt reflect.Type, path cty.Path, depth int) (cty.Type, error) {
	// Check depth to prevent infinite recursion
	if depth > maxImpliedTypeDepth {
		// Return a generic map type to break the recursion
		return cty.Map(cty.DynamicPseudoType), nil
	}

	switch rt.Kind() {

	case reflect.Ptr:
		return impliedTypeWithDepth(rt.Elem(), path, depth)

	// Primitive types
	case reflect.Bool:
		return cty.Bool, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return cty.Number, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return cty.Number, nil
	case reflect.Float32, reflect.Float64:
		return cty.Number, nil
	case reflect.String:
		return cty.String, nil

	// Collection types
	case reflect.Slice:
		path := append(path, cty.IndexStep{Key: cty.UnknownVal(cty.Number)})
		ety, err := impliedTypeWithDepth(rt.Elem(), path, depth+1)
		if err != nil {
			return cty.NilType, err
		}
		return cty.List(ety), nil
	case reflect.Map:
		if !stringType.AssignableTo(rt.Key()) {
			return cty.NilType, path.NewErrorf("no cty.Type for %s (must have string keys)", rt)
		}
		path := append(path, cty.IndexStep{Key: cty.UnknownVal(cty.String)})
		ety, err := impliedTypeWithDepth(rt.Elem(), path, depth+1)
		if err != nil {
			return cty.NilType, err
		}
		return cty.Map(ety), nil

	// Structural types
	case reflect.Struct:
		return impliedStructTypeWithDepth(rt, path, depth+1)

	default:
		return cty.NilType, path.NewErrorf("no cty.Type for %s", rt)
	}
}

func impliedStructType(rt reflect.Type, path cty.Path) (cty.Type, error) {
	// Keep existing function for backward compatibility
	return impliedStructTypeWithDepth(rt, path, 0)
}

func impliedStructTypeWithDepth(rt reflect.Type, path cty.Path, depth int) (cty.Type, error) {
	// Check depth to prevent infinite recursion
	if depth > maxImpliedTypeDepth {
		// Return a generic object type to break the recursion
		return cty.Object(map[string]cty.Type{}), nil
	}

	if valueType.AssignableTo(rt) {
		// Special case: cty.Value represents cty.DynamicPseudoType, for
		// type conformance checking.
		return cty.DynamicPseudoType, nil
	}

	fieldIdxs := structTagIndices(rt)
	if len(fieldIdxs) == 0 {
		return cty.NilType, path.NewErrorf("no cty.Type for %s (no hcl field tags)", rt)
	}

	atys := make(map[string]cty.Type, len(fieldIdxs))

	{
		// Temporary extension of path for attributes
		path := append(path, nil)

		for k, fi := range fieldIdxs {
			path[len(path)-1] = cty.GetAttrStep{Name: k}

			ft := rt.Field(fi).Type
			aty, err := impliedTypeWithDepth(ft, path, depth)
			if err != nil {
				return cty.NilType, err
			}

			atys[k] = aty
		}
	}

	// now check any anonymous fields
	ct := rt.NumField()
	for i := 0; i < ct; i++ {
		field := rt.Field(i)
		if field.Anonymous {
			anon, err := impliedStructTypeWithDepth(field.Type, path, depth)
			if err != nil {
				return cty.NilType, err
			}

			for k, v := range anon.AttributeTypes() {
				atys[k] = v
			}
		}
	}

	return cty.Object(atys), nil
}
