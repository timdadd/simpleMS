package common_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"lib/common"
	"path"
	"testing"
)

func TestIsJsonEqual(t *testing.T) {
	// A database add test - the created json has uuid etc.
	createJson, err := readTestFile(t, "create.json")
	if !assert.Nil(t, err) {
		t.Log(err)
		t.FailNow()
	}
	createdJson, err := readTestFile(t, "created.json")
	if !assert.Nil(t, err) {
		t.Log(err)
		t.FailNow()
	}
	assert.True(t, common.IsJsonEqual(t, "MemberEqual", createJson, createdJson))
	assert.True(t, common.IsJsonEqual(t, "MemberEqual", createdJson, createJson))
	assert.False(t, common.IsJsonEqual(t, "DeepEqual", createJson, createdJson))
	assert.True(t, common.IsJsonEqual(t, "DeepEqual", createJson, createJson))
	assert.True(t, common.IsJsonEqual(t, "DeepEqual", createdJson, createdJson))
	assert.False(t, common.IsJsonEqual(t, "KeyEqual", createJson, createdJson))

	updatedJson, err := readTestFile(t, "updated.json")
	if !assert.Nil(t, err) {
		t.Log(err)
		t.FailNow()
	}
	assert.False(t, common.IsJsonEqual(t, "DeepEqual", updatedJson, createdJson))
	assert.True(t, common.IsJsonEqual(t, "KeyEqual", updatedJson, createdJson))
}

// Read file from the test file path
func readTestFile(t *testing.T, testFile string) (string, error) {
	tdPath := "test_data" //test_data.GetTestCaseFolder()
	tFile := path.Join(tdPath, testFile)
	t.Logf("Loading test file : %v", tFile)
	res, err := ioutil.ReadFile(tFile)
	if err != nil {
		return "", fmt.Errorf("Error loading test file %s:%w", tFile, err)
	}
	return string(res), nil
}
