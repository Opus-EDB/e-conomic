package economic

import "testing"

func TestFilter(t *testing.T) {
	// Simple test to check if the filter string is created correctly
	f := &Filter{}
	f.AndCondition("name", FilterOperatorEquals, "test")
	f.AndCondition("age", FilterOperatorGreaterThan, 10)
	expected := `name$eq:test$and:age$gt:10`
	if f.String() != expected {
		t.Errorf("Expected %s, got %s", expected, f)
	}
	// Test in/not int arrays
	fArr := &Filter{}
	fArr.AndCondition("name", FilterOperatorIn, []string{"test", "test2"})
	fArr.AndCondition("age", FilterOperatorNotIn, []int{10, 20})
	expectedArr := `name$in:[test,test2]$and:age$nin:[10,20]`
	if fArr.String() != expectedArr {
		t.Errorf("Expected %s, got %s", expectedArr, fArr)
	}
	//TODO: Nesting using $and:()
}
