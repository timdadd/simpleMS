package common_test

import (
	"encoding/json"
	"testing"
)

// TestLevel is DeepEqual, MemberEqual, KeyEqual
func IsJsonEqual(t *testing.T, testLevel, s1, s2 string) bool {
	var o1 interface{}
	var o2 interface{}
	var err error
	t.Logf("%s test of:\n%s\n%s", testLevel, s1, s2)
	err = json.Unmarshal([]byte(s1), &o1)
	if err != nil {
		t.Logf("Error mashalling string 1 :: %s", err.Error())
		return false
	}
	err = json.Unmarshal([]byte(s2), &o2)
	if err != nil {
		t.Logf("Error mashalling string 2 :: %s", err.Error())
		return false
	}

	// Go through all members in o1 and see if the same named member in o2 is equal
	return compareMembers(t, testLevel, o1, o2)
}

func compareMembers(t *testing.T, testLevel string, in1, in2 interface{}) bool {
	if m1, ok := in1.(map[string]interface{}); ok {
		if m2, ok := in2.(map[string]interface{}); ok {
			return compareMaps(t, testLevel, m1, m2)
		}
		return false
	}
	if a1, ok := in1.([]interface{}); ok {
		if a2, ok := in2.([]interface{}); ok {
			return compareArrays(t, testLevel, a1, a2)
		}
		return false
	}
	t.Logf("%s member test of:'%v' == '%v'", testLevel, in1, in2)
	if _, ok := in2.(map[string]interface{}); ok {
		return false
	}
	if _, ok := in2.([]interface{}); ok {
		return false
	}
	if testLevel != "KeyEqual" && in1 != in2 {
		return false
	}

	t.Logf("Test:'%v'=='%v' OK", in1, in2)

	return true
}

func compareMaps(t *testing.T, testLevel string, in1, in2 map[string]interface{}) bool {
	t.Logf("%s map test of:\n%v\n%v", testLevel, in1, in2)
	if testLevel != "MemberEqual" && len(in1) != len(in2) {
		return false
	}
	for k, v := range in1 {
		if v2, ok := in2[k]; ok {
			t.Logf("Checking key '%v'", k)
			if !compareMembers(t, testLevel, v, v2) {
				return false
			}
		} else if testLevel == "MemberEqual" {
			t.Logf("Key '%v' missing, skipping for member equality", k)
		} else {
			return false
		}
	}
	return true
}

// Arrays have to be the same size, but number of members can be different
func compareArrays(t *testing.T, testLevel string, in1, in2 []interface{}) bool {
	t.Logf("%s array test of:\n%v\n%v", testLevel, in1, in2)
	for i := range in1 {
		if !compareMembers(t, testLevel, in1[i], in2[i]) {
			return false
		}
	}
	return true
}
