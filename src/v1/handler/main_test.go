package handler

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/Risuii/invoice/src/app"
	"github.com/Risuii/invoice/src/middleware/response"
	"github.com/nsf/jsondiff"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	os.Chdir("../../../")

	app.Init(context.Background())

	exitVal := m.Run()

	os.Exit(exitVal)
}

func CheckBodyResponse(t *testing.T, actualResponse []byte, expected interface{}) response.Response {
	var body response.Response
	err := json.Unmarshal(actualResponse, &body)

	assert.Nil(t, err, "Error when trying to unmarshal response")
	assert.NotNil(t, body.Data, "Your response data is should not be nil or empty")

	actualBytes, err := json.Marshal(body.Data)
	assert.Nil(t, err, "Error when trying to Marshal response data")

	expectedBytes, err := json.Marshal(expected)
	assert.Nil(t, err, "Error when trying to Marshal expected data")

	_, diffStr := jsondiff.Compare(expectedBytes, actualBytes, &jsondiff.Options{
		SkipMatches:      true,
		ChangedSeparator: " expected value is ",
	})

	assert.Empty(t, diffStr, "Your response data and expected data is difference. Check your expected and actual data again!")

	return body
}
