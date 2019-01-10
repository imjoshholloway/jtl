package jtl

import (
	"strings"
)

type Spec struct {
	// Type is the type of spec
	Type string `yaml:"type"`

	// SourcePath is the json-esque path for where in the source tree the
	// transformations should take place. Dot-notation is used to describe
	// paths to nested objects
	SourcePath string `yaml:"sourcePath"`

	// TargetPath is the json-esque path for where the transformations should
	// be applied to. Dot-notation is used to describe paths to nested objects
	TargetPath string `yaml:"targetPath"`

	UseKeyValueAsKey string `yaml:"useKeyValueAsKey"`

	// Condition
	Condition *Condition `yaml:"condition"`

	// Spec is a slice of specs
	Specs []Spec `yaml:"specs"`
}

// Process
func (spec *Spec) Process(in interface{}) interface{} {

	if in == nil {
		return nil
	}

	// If we've got a base path to start from lets get the data from that
	extracted := extractValue(in, strings.Split(spec.SourcePath, ".")...)
	if extracted == nil {
		return nil
	}

	// the data we found is an array so we need to process each item in the
	// array
	if set, ok := extracted.([]interface{}); ok {
		return spec.processArraySet(set)
	}

	// if we have a condition and it evaluates to false don't return the value
	if spec.Condition != nil {
		if !spec.Condition.Evaluate(extracted) {
			return nil
		}
	}

	// if the value we found isn't a map, let's return it
	if _, ok := extracted.(map[string]interface{}); !ok {
		return extracted
	}

	out := make(map[string]interface{})

	for _, mapping := range spec.Specs {

		res := mapping.Process(extracted)
		if res == nil {
			continue
		}

		mapped, ok := res.(map[string]interface{})
		if !ok {
			res := storeAtPath(res, strings.Split(mapping.TargetPath, ".")...)

			mapped, ok = res.(map[string]interface{})
			if !ok {
				return res
				//continue
			}
		}

		out = mergeMap(out, mapped)
	}

	return storeAtPath(out, strings.Split(spec.TargetPath, ".")...)
}

// processArraySet takes a set of interfaces as processes each item in the set
func (spec *Spec) processArraySet(set []interface{}) interface{} {

	// since items that are in an array can make use of the
	// UseKeyValueAsKey option, we need to keep a track of a map and a
	// slice of the results. This allows us to simply pass back whatever
	// was expected at the end, rather than doing multiple ifs and type
	// casts
	keyedResults := make(map[string]interface{})
	res := make([]interface{}, 0)

	// create a copy of the spec so that don't look inherit our SourcePath
	// and TargetPath - this stops us incorrectly trying to look down the
	// tree, when we've already looked
	tmpSpec := Spec{
		Specs:     spec.Specs,
		Condition: spec.Condition,
	}

	for _, item := range set {

		out := tmpSpec.Process(item)
		if out == nil {
			continue
		}

		key, ok := extractValue(item, strings.Split(spec.UseKeyValueAsKey, ".")...).(string)
		if ok {
			keyedResults[key] = out
		}

		res = append(res, out)
	}

	var out interface{}
	out = res

	if len(spec.Specs) == 1 && len(res) == 1 {
		out = res[0]
	}

	if spec.UseKeyValueAsKey != "" {
		out = keyedResults
	}

	return storeAtPath(out, strings.Split(spec.TargetPath, ".")...)
}

// extractValue traverses the provided map structure returning the value at the
// specified key. If the key contains json-esque dot notation it is split and
// each part is used as a key
func extractValue(in interface{}, path ...string) interface{} {

	if in == nil {
		return nil
	}

	source, ok := in.(map[string]interface{})
	if !ok {
		return in
	}

	if len(path) <= 0 || path[0] == "" {
		return in
	}

	// didn't find anything so bail out early
	data, ok := source[path[0]]
	if !ok {
		return nil
	}

	if mapped, ok := data.(map[string]interface{}); ok {
		return extractValue(mapped, path[1:]...)
	}

	if mapped, ok := data.([]map[string]interface{}); ok {
		out := make([]interface{}, 0)

		for _, val := range mapped {
			res := extractValue(val, path[1:]...)
			if res != nil {
				out = append(out, res)
			}
		}

		return out
	}

	return data
}

func storeAtPath(value interface{}, path ...string) interface{} {

	if value == nil {
		return nil
	}

	out := make(map[string]interface{})

	current := 0
	if len(path) <= 0 || path[current] == "" {
		return value
	}

	next := 1

	if len(path) > 1 {
		value = storeAtPath(value, path[next:]...)
	}

	out[path[current]] = value
	return out
}

func mergeMap(maps ...map[string]interface{}) map[string]interface{} {

	out := make(map[string]interface{})

	if len(maps) == 0 {
		return out
	}

	// for each of the maos
	for _, m := range maps {

		for k := range m {

			if _, exists := out[k]; !exists {
				out[k] = m[k]
				continue
			}

			if mapped, ok := out[k].(map[string]interface{}); ok {

				newMap, ok := m[k].(map[string]interface{})
				if !ok {
					newMap = make(map[string]interface{})
					newMap[k] = m[k]
				}

				out[k] = mergeMap(mapped, newMap)
				continue
			}
		}
	}

	return out
}
