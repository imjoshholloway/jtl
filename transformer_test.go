package jtl

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProcessSpec(t *testing.T) {

	testCases := []struct {
		Source map[string]interface{}

		Spec Spec

		Expected interface{}
	}{
		{
			Spec: Spec{
				TargetPath: "nested",
				Specs: []Spec{
					{
						SourcePath: "account",
						TargetPath: "account_name",
					},
				},
			},
			Source: map[string]interface{}{
				"account": "some-account",
			},
			Expected: map[string]interface{}{
				"nested": map[string]interface{}{
					"account_name": "some-account",
				},
			},
		},
		// Test nested object mappings with mappings
		{
			Spec: Spec{
				SourcePath:       "nested_source",
				TargetPath:       "nested",
				UseKeyValueAsKey: "name",
				Specs: []Spec{
					{
						SourcePath: "value_field",
						TargetPath: "value",
					},
				},
			},
			Source: map[string]interface{}{
				"nested_source": []map[string]interface{}{
					{
						"name":        "test1",
						"value_field": "test1-value",
					},
					{
						"name":        "test2",
						"value_field": "test2-value",
					},
				},
			},
			Expected: map[string]interface{}{
				"nested": map[string]interface{}{
					"test1": map[string]interface{}{
						"value": "test1-value",
					},
					"test2": map[string]interface{}{
						"value": "test2-value",
					},
				},
			},
		},
		// Nested object with overlapping mappings
		{
			Spec: Spec{
				Specs: []Spec{
					{
						SourcePath: "address",
						TargetPath: "contact",
						Specs: []Spec{
							{
								SourcePath: "line1",
								TargetPath: "address.line1",
							},
						},
					},
					{
						TargetPath: "contact",
						Specs: []Spec{
							{
								SourcePath: "email",
								TargetPath: "email",
							},
						},
					},
				},
			},
			Source: map[string]interface{}{
				"address": map[string]interface{}{
					"line1": "test",
				},
				"email": "test@test.com",
			},
			Expected: map[string]interface{}{
				"contact": map[string]interface{}{
					"email": "test@test.com",
					"address": map[string]interface{}{
						"line1": "test",
					},
				},
			},
		},
		// Test with condition
		{
			Spec: Spec{
				Specs: []Spec{
					{
						SourcePath: "addresses",
						TargetPath: "contact.addresses",
						Specs: []Spec{
							{
								SourcePath: "line1",
								TargetPath: "line1",
							},
							{
								SourcePath: "city",
								TargetPath: "city",
							},
						},
					},
					{
						TargetPath: "contact",
						Specs: []Spec{
							{
								SourcePath: "email",
								TargetPath: "email",
							},
						},
					},
					{
						SourcePath: "addresses",
						TargetPath: "contact.preferred_city",
						Condition: &Condition{
							Comparator: OPERATOR_EQ,
							SourcePath: "preferred",
							Value:      true,
						},
						Specs: []Spec{
							{
								SourcePath: "city",
							},
						},
					},
				},
			},
			Source: map[string]interface{}{
				"addresses": []map[string]interface{}{
					{
						"line1":     "test",
						"preferred": false,
						"city":      "New York",
					},
					{
						"line1":     "test",
						"city":      "London",
						"preferred": true,
					},
				},
				"email": "test@test.com",
			},
			Expected: map[string]interface{}{
				"contact": map[string]interface{}{
					"email": "test@test.com",
					"addresses": []interface{}{
						map[string]interface{}{
							"line1": "test",
							"city":  "New York",
						},
						map[string]interface{}{
							"line1": "test",
							"city":  "London",
						},
					},
					"preferred_city": "London",
				},
			},
		},
	}

	for i, tc := range testCases {
		results := tc.Spec.Process(tc.Source)
		assert.Equal(t, tc.Expected, results, "TC %v", i+1)
	}
}

func TestSpecProcessArraySet(t *testing.T) {

	testCases := []struct {
		Source []interface{}

		Spec Spec

		Expected interface{}
	}{
		{
			Spec: Spec{
				TargetPath: "accounts",
				Specs: []Spec{
					{
						SourcePath: "account",
						TargetPath: "account_name",
					},
				},
			},
			Source: []interface{}{
				map[string]interface{}{
					"account": "some-account",
				},
				map[string]interface{}{
					"account": "some-other-account",
				},
			},
			Expected: map[string]interface{}{
				"accounts": []interface{}{
					map[string]interface{}{
						"account_name": "some-account",
					},
					map[string]interface{}{
						"account_name": "some-other-account",
					},
				},
			},
		},
		{
			Spec: Spec{
				UseKeyValueAsKey: "account",
				TargetPath:       "accounts",
				Specs: []Spec{
					{
						SourcePath: "account",
						TargetPath: "account_name",
					},
				},
			},
			Source: []interface{}{
				map[string]interface{}{
					"account": "some-account",
				},
				map[string]interface{}{
					"account": "some-other-account",
				},
			},
			Expected: map[string]interface{}{
				"accounts": map[string]interface{}{
					"some-account": map[string]interface{}{
						"account_name": "some-account",
					},
					"some-other-account": map[string]interface{}{
						"account_name": "some-other-account",
					},
				},
			},
		},
	}

	for i, tc := range testCases {
		results := tc.Spec.processArraySet(tc.Source)
		assert.Equal(t, tc.Expected, results, "TC %v", i+1)
	}
}

func BenchmarkProcess(b *testing.B) {

	spec := Spec{
		TargetPath: "hierarchy",
		Specs: []Spec{
			{
				SourcePath: "nested.child.grandchild",
				TargetPath: "great_grand_children",
			},
			{
				SourcePath: "name",
				TargetPath: "name",
			},
		},
	}

	source := map[string]interface{}{
		"name": "value",
		"nested": map[string]interface{}{
			"child": map[string]interface{}{
				"grandchild": map[string]interface{}{
					"great_grandchild": map[string]interface{}{
						"kin": "test",
					},
				},
			},
		},
	}

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		spec.Process(source)
	}
}

func TestExtractValue(t *testing.T) {

	testCases := []struct {
		Source   map[string]interface{}
		Path     []string
		Expected interface{}
	}{
		// Test basic value extraction
		{
			Source: map[string]interface{}{
				"basic": "value",
			},
			Path:     []string{"basic"},
			Expected: "value",
		},

		// Test basic nested object extraction
		{
			Source: map[string]interface{}{
				"nested_source": map[string]interface{}{
					"key": "value",
				},
			},
			Path:     []string{"nested_source", "key"},
			Expected: "value",
		},
		// Test object extraction
		{
			Source: map[string]interface{}{
				"nested_source": map[string]interface{}{
					"key": "value",
				},
			},
			Path: []string{"nested_source"},
			Expected: map[string]interface{}{
				"key": "value",
			},
		},
		// Test array extraction
		{
			Source: map[string]interface{}{
				"nested_source": []map[string]interface{}{
					{
						"key": "value",
					},
				},
			},
			Path: []string{"nested_source"},
			Expected: []interface{}{
				map[string]interface{}{
					"key": "value",
				},
			},
		},
		// Test extraction
		{
			Source: map[string]interface{}{
				"nested_source": []map[string]interface{}{
					{
						"key": "value",
					},
				},
			},
			Path: []string{"nested_source"},
			Expected: []interface{}{
				map[string]interface{}{
					"key": "value",
				},
			},
		},
		// Test when no key
		{
			Source: map[string]interface{}{
				"nested_source": []map[string]interface{}{
					{
						"key": "value",
					},
				},
			},
			Path: []string{},
			Expected: map[string]interface{}{
				"nested_source": []map[string]interface{}{
					{
						"key": "value",
					},
				},
			},
		},
	}

	for i, tc := range testCases {
		out := extractValue(tc.Source, tc.Path...)
		assert.Equal(t, tc.Expected, out, "TC %v, Expected to find data at path %v", i+1, strings.Join(tc.Path, "."))
	}
}

func TestStoreAtPath(t *testing.T) {
	testCases := []struct {
		Data     interface{}
		Expected interface{}
		Path     []string
	}{
		{
			Data: map[string]interface{}{
				"some": "data",
			},
			Expected: map[string]interface{}{
				"nested": map[string]interface{}{
					"some": "data",
				},
			},
			Path: []string{"nested"},
		},
		{
			Data: map[string]interface{}{
				"some": "data",
			},
			Expected: map[string]interface{}{
				"deeply": map[string]interface{}{
					"nested": map[string]interface{}{
						"some": "data",
					},
				},
			},
			Path: []string{"deeply", "nested"},
		},
	}

	for i, tc := range testCases {
		out := storeAtPath(tc.Data, tc.Path...)
		assert.Equal(t, tc.Expected, out, "TC %v", i+1)
	}
}

func TestMergeMap(t *testing.T) {

	testCases := []struct {
		SourceMaps []map[string]interface{}
		Expected   map[string]interface{}
	}{
		{
			SourceMaps: []map[string]interface{}{
				{
					"a": map[string]interface{}{
						"b": true,
					},
					"c": "abc",
				},
				{
					"a": map[string]interface{}{
						"c": "nested",
					},
					"d": "abcd",
				},
			},
			Expected: map[string]interface{}{
				"a": map[string]interface{}{
					"b": true,
					"c": "nested",
				},
				"c": "abc",
				"d": "abcd",
			},
		},
	}

	for i, tc := range testCases {

		out := mergeMap(tc.SourceMaps...)
		assert.Equal(t, tc.Expected, out, "TC %v", i+1)
	}

}
