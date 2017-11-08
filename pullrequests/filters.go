package pullrequests

import (
	"errors"
	"fmt"
	"strconv"
	"time"
)

func DoFilter(operator string, f Filter) (bool, error) {
	switch operator {
	case "eq":
		return f.Eq()
	case "neq":
		return f.Neq()
	}

	return false, errors.New("Unrecognised operator '%s'")
}

func NewColumnFilter(columnName, columnValue interface{}, compareValue string) (Filter, error) {
	switch columnName {
	case "Author", "RepoName", "URL":
		return &StringFilter{ColumnValue: columnValue.(string), CompareValue: compareValue}, nil

	case "TotalReviews":
		v, err := strconv.Atoi(compareValue)
		if err != nil {
			return nil, err
		}

		return &IntFilter{ColumnValue: columnValue.(int), CompareValue: v}, nil

	case "Approved", "ChangesRequested":
		v, err := strconv.ParseBool(compareValue)
		if err != nil {
			return nil, err
		}

		return &BoolFilter{ColumnValue: columnValue.(bool), CompareValue: v}, nil
	}

	return nil, fmt.Errorf("Column '%s' not supported", columnName)
}

type Filter interface {
	Eq() (bool, error)
	Neq() (bool, error)
}

type BoolFilter struct {
	Operator     string
	ColumnValue  bool
	CompareValue bool
}

func (f *BoolFilter) Eq() (bool, error) {
	if f.ColumnValue == f.CompareValue {
		return true, nil
	}
	return false, nil
}

func (f *BoolFilter) Neq() (bool, error) {
	if f.ColumnValue != f.CompareValue {
		return true, nil
	}
	return false, nil
}

type IntFilter struct {
	Operator     string
	ColumnValue  int
	CompareValue int
}

func (f *IntFilter) Eq() (bool, error) {
	if f.ColumnValue == f.CompareValue {
		return true, nil
	}
	return false, nil
}

func (f *IntFilter) Neq() (bool, error) {
	if f.ColumnValue != f.CompareValue {
		return true, nil
	}
	return false, nil
}

type StringFilter struct {
	Operator     string
	ColumnValue  string
	CompareValue string
}

func (f *StringFilter) Eq() (bool, error) {
	if f.ColumnValue == f.CompareValue {
		return true, nil
	}
	return false, nil
}

func (f *StringFilter) Neq() (bool, error) {
	if f.ColumnValue != f.CompareValue {
		return true, nil
	}
	return false, nil
}

type TimeFilter struct {
	Operator     string
	ColumnValue  time.Time
	CompareValue time.Time
}
