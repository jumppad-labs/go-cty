package gocty

import (
	"fmt"
	"math/big"
	"reflect"
	"testing"

	"github.com/zclconf/go-cty/cty"
)

type Embedded struct {
	ID string `hcl:"id"`
}

type Main struct {
	Embedded `hcl:",remain"`
	Name     string `hcl:"name"`
}

func TestFromCtyValueWithEmbeddedTypes(t *testing.T) {
	m := Main{}

	ctyT := map[string]cty.Value{
		"id":   cty.StringVal("abc"),
		"name": cty.StringVal("Erik"),
	}

	err := FromCtyValue(cty.ObjectVal(ctyT), &m)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(m)
}

func TestToCtyValueWithEmbeddedTypes(t *testing.T) {
	m := Main{
		Embedded: Embedded{ID: "123"},
		Name:     "Nic",
	}

	typ, err := ImpliedType(m)
	if err != nil {
		t.Fatal(err)
	}

	ctyVal, err := ToCtyValue(m, typ)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(ctyVal)
}

func TestOut(t *testing.T) {
	capsuleANative := &capsuleType1Native{"capsuleA"}

	type (
		stringAlias  string
		boolAlias    bool
		intAlias     int
		float32Alias float32
		float64Alias float64
		bigIntAlias  big.Int
		listIntAlias []int
		mapIntAlias  map[string]int
	)

	tests := []struct {
		CtyValue   cty.Value
		TargetType reflect.Type
		Want       interface{}
	}{

		// Bool
		{
			CtyValue:   cty.True,
			TargetType: reflect.TypeOf(false),
			Want:       true,
		},
		{
			CtyValue:   cty.False,
			TargetType: reflect.TypeOf(false),
			Want:       false,
		},
		{
			CtyValue:   cty.True,
			TargetType: reflect.PtrTo(reflect.TypeOf(false)),
			Want:       testOutAssertPtrVal(true),
		},
		{
			CtyValue:   cty.NullVal(cty.Bool),
			TargetType: reflect.PtrTo(reflect.TypeOf(false)),
			Want:       (*bool)(nil),
		},
		{
			CtyValue:   cty.True,
			TargetType: reflect.TypeOf((*boolAlias)(nil)).Elem(),
			Want:       boolAlias(true),
		},

		// String
		{
			CtyValue:   cty.StringVal("hello"),
			TargetType: reflect.TypeOf(""),
			Want:       "hello",
		},
		{
			CtyValue:   cty.StringVal(""),
			TargetType: reflect.TypeOf(""),
			Want:       "",
		},
		{
			CtyValue:   cty.StringVal("hello"),
			TargetType: reflect.PtrTo(reflect.TypeOf("")),
			Want:       testOutAssertPtrVal("hello"),
		},
		{
			CtyValue:   cty.NullVal(cty.String),
			TargetType: reflect.PtrTo(reflect.TypeOf("")),
			Want:       (*string)(nil),
		},
		{
			CtyValue:   cty.StringVal("hello"),
			TargetType: reflect.TypeOf((*stringAlias)(nil)).Elem(),
			Want:       stringAlias("hello"),
		},

		// Number
		{
			CtyValue:   cty.NumberIntVal(5),
			TargetType: reflect.TypeOf(int(0)),
			Want:       int(5),
		},
		{
			CtyValue:   cty.NumberIntVal(5),
			TargetType: reflect.TypeOf(int8(0)),
			Want:       int8(5),
		},
		{
			CtyValue:   cty.NumberIntVal(5),
			TargetType: reflect.TypeOf(int16(0)),
			Want:       int16(5),
		},
		{
			CtyValue:   cty.NumberIntVal(5),
			TargetType: reflect.TypeOf(int32(0)),
			Want:       int32(5),
		},
		{
			CtyValue:   cty.NumberIntVal(5),
			TargetType: reflect.TypeOf(int64(0)),
			Want:       int64(5),
		},
		{
			CtyValue:   cty.NumberIntVal(5),
			TargetType: reflect.TypeOf(uint(0)),
			Want:       uint(5),
		},
		{
			CtyValue:   cty.NumberIntVal(5),
			TargetType: reflect.TypeOf(uint8(0)),
			Want:       uint8(5),
		},
		{
			CtyValue:   cty.NumberIntVal(5),
			TargetType: reflect.TypeOf(uint16(0)),
			Want:       uint16(5),
		},
		{
			CtyValue:   cty.NumberIntVal(5),
			TargetType: reflect.TypeOf(uint32(0)),
			Want:       uint32(5),
		},
		{
			CtyValue:   cty.NumberIntVal(5),
			TargetType: reflect.TypeOf(uint64(0)),
			Want:       uint64(5),
		},
		{
			CtyValue:   cty.NumberFloatVal(1.5),
			TargetType: reflect.TypeOf(float32(0)),
			Want:       float32(1.5),
		},
		{
			CtyValue:   cty.NumberFloatVal(1.5),
			TargetType: reflect.TypeOf(float64(0)),
			Want:       float64(1.5),
		},
		{
			CtyValue:   cty.NumberFloatVal(1.5),
			TargetType: reflect.PtrTo(bigFloatType),
			Want:       big.NewFloat(1.5),
		},
		{
			CtyValue:   cty.NumberIntVal(5),
			TargetType: reflect.PtrTo(bigIntType),
			Want:       big.NewInt(5),
		},
		{
			CtyValue:   cty.NumberIntVal(5),
			TargetType: reflect.TypeOf((*intAlias)(nil)).Elem(),
			Want:       intAlias(5),
		},
		{
			CtyValue:   cty.NumberFloatVal(1.5),
			TargetType: reflect.TypeOf((*float32Alias)(nil)).Elem(),
			Want:       float32Alias(1.5),
		},
		{
			CtyValue:   cty.NumberFloatVal(1.5),
			TargetType: reflect.TypeOf((*float64Alias)(nil)).Elem(),
			Want:       float64Alias(1.5),
		},
		{
			CtyValue:   cty.NumberIntVal(5),
			TargetType: reflect.TypeOf((*bigIntAlias)(nil)),
			Want:       (*bigIntAlias)(big.NewInt(5)),
		},

		// Lists
		{
			CtyValue:   cty.ListValEmpty(cty.Number),
			TargetType: reflect.TypeOf(([]int)(nil)),
			Want:       []int{},
		},
		{
			CtyValue:   cty.ListVal([]cty.Value{cty.NumberIntVal(1), cty.NumberIntVal(5)}),
			TargetType: reflect.TypeOf(([]int)(nil)),
			Want:       []int{1, 5},
		},
		{
			CtyValue:   cty.NullVal(cty.List(cty.Number)),
			TargetType: reflect.TypeOf(([]int)(nil)),
			Want:       ([]int)(nil),
		},
		{
			CtyValue:   cty.ListVal([]cty.Value{cty.NumberIntVal(1), cty.NumberIntVal(5)}),
			TargetType: reflect.ArrayOf(2, reflect.TypeOf(0)),
			Want:       [2]int{1, 5},
		},
		{
			CtyValue:   cty.ListValEmpty(cty.Number),
			TargetType: reflect.ArrayOf(0, reflect.TypeOf(0)),
			Want:       [0]int{},
		},
		{
			CtyValue:   cty.ListValEmpty(cty.Number),
			TargetType: reflect.PtrTo(reflect.ArrayOf(0, reflect.TypeOf(0))),
			Want:       testOutAssertPtrVal([0]int{}),
		},
		{
			CtyValue:   cty.ListVal([]cty.Value{cty.NumberIntVal(1), cty.NumberIntVal(5)}),
			TargetType: reflect.TypeOf((listIntAlias)(nil)),
			Want:       listIntAlias{1, 5},
		},

		// Maps
		{
			CtyValue:   cty.MapValEmpty(cty.Number),
			TargetType: reflect.TypeOf((map[string]int)(nil)),
			Want:       map[string]int{},
		},
		{
			CtyValue: cty.MapVal(map[string]cty.Value{
				"one":  cty.NumberIntVal(1),
				"five": cty.NumberIntVal(5),
			}),
			TargetType: reflect.TypeOf(map[string]int{}),
			Want: map[string]int{
				"one":  1,
				"five": 5,
			},
		},
		{
			CtyValue:   cty.NullVal(cty.Map(cty.Number)),
			TargetType: reflect.TypeOf((map[string]int)(nil)),
			Want:       (map[string]int)(nil),
		},
		{
			CtyValue: cty.MapVal(map[string]cty.Value{
				"one":  cty.NumberIntVal(1),
				"five": cty.NumberIntVal(5),
			}),
			TargetType: reflect.TypeOf(mapIntAlias(nil)),
			Want: mapIntAlias{
				"one":  1,
				"five": 5,
			},
		},

		// Sets
		{
			CtyValue:   cty.SetValEmpty(cty.Number),
			TargetType: reflect.TypeOf(([]int)(nil)),
			Want:       []int{},
		},
		{
			CtyValue:   cty.SetVal([]cty.Value{cty.NumberIntVal(1), cty.NumberIntVal(5)}),
			TargetType: reflect.TypeOf(([]int)(nil)),
			Want:       []int{1, 5},
		},
		{
			CtyValue:   cty.SetVal([]cty.Value{cty.NumberIntVal(1), cty.NumberIntVal(5)}),
			TargetType: reflect.TypeOf([2]int{}),
			Want:       [2]int{1, 5},
		},

		// Objects
		{
			CtyValue:   cty.EmptyObjectVal,
			TargetType: reflect.TypeOf(struct{}{}),
			Want:       struct{}{},
		},
		{
			CtyValue: cty.ObjectVal(map[string]cty.Value{
				"name": cty.StringVal("Stephen"),
			}),
			TargetType: reflect.TypeOf(testStruct{}),
			Want: testStruct{
				Name: "Stephen",
			},
		},
		{
			CtyValue: cty.ObjectVal(map[string]cty.Value{
				"name":   cty.StringVal("Stephen"),
				"number": cty.NumberIntVal(12),
			}),
			TargetType: reflect.TypeOf(testStruct{}),
			Want: testStruct{
				Name:   "Stephen",
				Number: ptrToInt(12),
			},
		},

		// Tuples
		{
			CtyValue:   cty.EmptyTupleVal,
			TargetType: reflect.TypeOf(struct{}{}),
			Want:       struct{}{},
		},
		{
			CtyValue: cty.TupleVal([]cty.Value{
				cty.StringVal("Stephen"),
				cty.NumberIntVal(5),
			}),
			TargetType: reflect.TypeOf(testTupleStruct{}),
			Want:       testTupleStruct{"Stephen", 5},
		},

		// Capsules
		{
			CtyValue:   cty.CapsuleVal(capsuleType1, capsuleANative),
			TargetType: reflect.TypeOf(capsuleType1Native{}),
			Want:       capsuleType1Native{"capsuleA"},
		},
		{
			CtyValue:   cty.CapsuleVal(capsuleType1, capsuleANative),
			TargetType: reflect.PtrTo(reflect.TypeOf(capsuleType1Native{})),
			Want:       capsuleANative, // should recover original pointer
		},

		// Passthrough
		{
			CtyValue:   cty.NumberIntVal(2),
			TargetType: valueType,
			Want:       cty.NumberIntVal(2),
		},
		{
			CtyValue:   cty.UnknownVal(cty.Bool),
			TargetType: valueType,
			Want:       cty.UnknownVal(cty.Bool),
		},
		{
			CtyValue:   cty.NullVal(cty.Bool),
			TargetType: valueType,
			Want:       cty.NullVal(cty.Bool),
		},
		{
			CtyValue:   cty.DynamicVal,
			TargetType: valueType,
			Want:       cty.DynamicVal,
		},
		{
			CtyValue:   cty.NullVal(cty.DynamicPseudoType),
			TargetType: valueType,
			Want:       cty.NullVal(cty.DynamicPseudoType),
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%#v into %s", test.CtyValue, test.TargetType), func(t *testing.T) {
			target := reflect.New(test.TargetType)
			err := FromCtyValue(test.CtyValue, target.Interface())
			if err != nil {
				t.Fatalf("FromCtyValue returned error: %s", err)
			}

			got := target.Elem().Interface()

			if assertFunc, ok := test.Want.(testOutAssertFunc); ok {
				assertFunc(test.CtyValue, test.TargetType, got, t)
			} else if wantV, ok := test.Want.(cty.Value); ok {
				if gotV, ok := got.(cty.Value); ok {
					if !gotV.RawEquals(wantV) {
						testOutWrongResult(test.CtyValue, test.TargetType, got, test.Want, t)
					}
				} else {
					testOutWrongResult(test.CtyValue, test.TargetType, got, test.Want, t)
				}
			} else {
				if !reflect.DeepEqual(got, test.Want) {
					testOutWrongResult(test.CtyValue, test.TargetType, got, test.Want, t)
				}
			}
		})
	}
}

type testOutAssertFunc func(cty.Value, reflect.Type, interface{}, *testing.T)

func testOutAssertPtrVal(want interface{}) testOutAssertFunc {
	return func(ctyValue cty.Value, targetType reflect.Type, gotPtr interface{}, t *testing.T) {
		wantVal := reflect.ValueOf(want)
		gotVal := reflect.ValueOf(gotPtr)

		if gotVal.Kind() != reflect.Ptr {
			t.Fatalf("wrong type %s; want pointer to %T", gotVal.Type(), want)
		}
		gotVal = gotVal.Elem()

		want := wantVal.Interface()
		got := gotVal.Interface()
		if got != want {
			testOutWrongResult(
				ctyValue,
				targetType,
				got,
				want,
				t,
			)
		}
	}
}

func testOutWrongResult(ctyValue cty.Value, targetType reflect.Type, got interface{}, want interface{}, t *testing.T) {
	t.Errorf("wrong result\ninput:       %#v\ntarget type: %s\ngot:         %#v\nwant:        %#v", ctyValue, targetType, got, want)
}

type testStruct struct {
	Name   string `cty:"name"`
	Number *int   `cty:"number"`
}

type testTupleStruct struct {
	Name   string
	Number int
}
