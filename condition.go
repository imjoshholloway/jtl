package jtl

import (
	"errors"
)

type Operator int64

const (
	// OPERATOR_EQ is the equality comparator
	OPERATOR_EQ Operator = iota
	// OPERATOR_GT is the greater than comparator
	OPERATOR_GT
	// OPERATOR_GTE is the greater or equal to than comparator
	OPERATOR_GTE
	// OPERATOR_LT is the less than comparator
	OPERATOR_LT
	// OPERATOR_LTE is the less than or equal to than comparator
	OPERATOR_LTE
)

var (
	comparatorStrings map[string]Operator = map[string]Operator{
		"=":  OPERATOR_EQ,
		">":  OPERATOR_GT,
		">=": OPERATOR_GTE,
		"<":  OPERATOR_LT,
		"<=": OPERATOR_LTE,
	}
)

func (c *Operator) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var raw string
	if err := unmarshal(&raw); err != nil {
		return err
	}

	found, ok := comparatorStrings[raw]
	if !ok {
		return errors.New("invalid comparator")
	}

	*c = found

	return nil
}

type Condition struct {
	SourcePath string      `yaml:"sourcePath"`
	Comparator Operator    `yaml:"comparator"`
	Value      interface{} `yaml:"value"`
}

func (c *Condition) Evaluate(in interface{}) bool {

	var val interface{}

	val = in

	// if the value is a map then we should look in it for the SourcePath and
	// use the value as our value to compare
	if mapped, ok := in.(map[string]interface{}); ok {
		if c.SourcePath == "" {
			return false
		}

		val, ok = mapped[c.SourcePath]
		if !ok {
			val = new(interface{})
		}
	}

	switch c.Comparator {
	case OPERATOR_EQ:
		return c.isEqual(val)

	case OPERATOR_GT:
		return c.isGreaterThan(val)

	case OPERATOR_GTE:
		return c.isGreaterThanOrEqualTo(val)

	case OPERATOR_LT:
		return c.isLessThan(val)

	case OPERATOR_LTE:
		return c.isLessThanOrEqualTo(val)
	}

	return false
}

func (c *Condition) isEqual(val interface{}) bool {

	switch c.Value.(type) {
	case string:
		if newVal, ok := val.(string); ok {
			return c.Value.(string) == newVal
		}
		return false
	case int:
		if newVal, ok := val.(int); ok {
			return c.Value.(int) == newVal
		}
	case bool:
		if newVal, ok := val.(bool); ok {
			return c.Value.(bool) == newVal
		}

		return !c.Value.(bool)
	}

	return false
}

func (c *Condition) isLessThanOrEqualTo(val interface{}) bool {
	valInt, ok := val.(int)
	if !ok {
		return false
	}

	compareTo, ok := c.Value.(int)
	if !ok {
		return false
	}

	return valInt <= compareTo
}

func (c *Condition) isLessThan(val interface{}) bool {
	valInt, ok := val.(int)
	if !ok {
		return false
	}

	compareTo, ok := c.Value.(int)
	if !ok {
		return false
	}

	return valInt < compareTo
}

func (c *Condition) isGreaterThanOrEqualTo(val interface{}) bool {
	valInt, ok := val.(int)
	if !ok {
		return false
	}

	compareTo, ok := c.Value.(int)
	if !ok {
		return false
	}

	return valInt >= compareTo
}

func (c *Condition) isGreaterThan(val interface{}) bool {
	valInt, ok := val.(int)
	if !ok {
		return false
	}

	compareTo, ok := c.Value.(int)
	if !ok {
		return false
	}

	return valInt > compareTo
}
