// This is a work in progress.

package economic

import (
	"fmt"
	"strings"
)

type FilterOperator string

// “$” is escaped with “$$”
// “(” is escaped with “$(”
// “)” is escaped with “$)”
// “*” is escaped with “$*”
// “,” is escaped with “$,”
// “[” is escaped with “$[”
// “]” is escaped with “$]”

const (
	FilterOperatorEquals             = FilterOperator("$eq")
	FilterOperatorNotEquals          = FilterOperator("$ne")
	FilterOperatorGreaterThan        = FilterOperator("$gt")
	FilterOperatorGreaterThanOrEqual = FilterOperator("$gte")
	FilterOperatorLessThan           = FilterOperator("$lt")
	FilterOperatorLessThanOrEqual    = FilterOperator("$lte")
	FilterOperatorSubstringMatch     = FilterOperator("$like")
	FilterOperatorAndAlso            = FilterOperator("$and")
	FilterOperatorOrElse             = FilterOperator("$or")
	FilterOperatorIn                 = FilterOperator("$in")
	FilterOperatorNotIn              = FilterOperator("$nin")
)

type Filter struct {
	filterStr string
}

func (f *Filter) AndCondition(field string, operator FilterOperator, value any) {
	f.condition(field, operator, value, true)
}

func (f *Filter) OrCondition(field string, operator FilterOperator, value any) {
	f.condition(field, operator, value, false)
}

func (f *Filter) condition(field string, operator FilterOperator, value any, and bool) {
	value = fmt.Sprint(value)
	if operator == FilterOperatorIn || operator == FilterOperatorNotIn {
		value = strings.ReplaceAll(fmt.Sprintf("%v", value), " ", ",")
	}
	if len(f.filterStr) > 0 {
		if and {
			f.filterStr += "$and:"
		} else {
			f.filterStr += "$or:"
		}
	}
	f.filterStr += fmt.Sprintf("%s%s:%s", field, string(operator), fmt.Sprint(value))
}

func (f *Filter) String() string {
	return f.filterStr
}
