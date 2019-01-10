package jtl

import (
	"testing"

	"github.com/stretchr/testify/assert"
	yaml "gopkg.in/yaml.v2"
)

func TestCondition(t *testing.T) {

	testCases := []struct {
		Source    interface{}
		Condition Condition
		Expected  bool
	}{
		{
			Source: true,
			Condition: Condition{
				Comparator: OPERATOR_EQ,
				Value:      true,
			},
			Expected: true,
		},
		{
			Source: 10,
			Condition: Condition{
				Comparator: OPERATOR_GTE,
				Value:      9,
			},
			Expected: true,
		},
		{
			Source: 10,
			Condition: Condition{
				Comparator: OPERATOR_GTE,
				Value:      11,
			},
			Expected: false,
		},
		{
			Source: map[string]interface{}{
				"test": true,
			},
			Condition: Condition{
				SourcePath: "test",
				Comparator: OPERATOR_EQ,
				Value:      true,
			},
			Expected: true,
		},
		{
			Source: map[string]interface{}{
				"test": true,
			},
			Condition: Condition{
				SourcePath: "not-found",
				Comparator: OPERATOR_EQ,
				Value:      false,
			},
			Expected: true,
		},
	}

	for i, tc := range testCases {
		out := tc.Condition.Evaluate(tc.Source)
		assert.Equal(t, tc.Expected, out, "TC %v", i+1)
	}

}

func TestConditionUnmarshalYaml(t *testing.T) {

	testCases := []struct {
		Source   string
		Expected Condition
	}{
		{
			Source: `---\n comparator: "="`,
			Expected: Condition{
				Comparator: OPERATOR_EQ,
			},
		},
	}

	for i, tc := range testCases {
		out := Condition{}
		err := yaml.Unmarshal([]byte(tc.Source), &out)
		assert.NoError(t, err)
		assert.Equal(t, tc.Expected, out, "TC %v", i+1)
	}

}
