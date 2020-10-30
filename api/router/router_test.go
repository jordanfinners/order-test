package router

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"jordanfinners/api/model"
	"jordanfinners/api/storage"
)

func seedPacks(testClient storage.Client) {
	testClient.SavePack(context.TODO(), model.Pack{Quantity: 250})
	testClient.SavePack(context.TODO(), model.Pack{Quantity: 500})
	testClient.SavePack(context.TODO(), model.Pack{Quantity: 1000})
	testClient.SavePack(context.TODO(), model.Pack{Quantity: 2000})
	testClient.SavePack(context.TODO(), model.Pack{Quantity: 5000})
}

func TestMain(m *testing.M) {
	testClient := storage.StartTestDB()
	seedPacks(testClient)
	os.Exit(m.Run())
}

func createAzureFunctionRequestBody(method, body string) string {
	return fmt.Sprintf(`{"Data":{"req":{"Url":"http://localhost:7071/api/orders","Method":"%s","Query":"{}","Headers":{"Content-Type":["application/json"]},"Params":{},"Body":"%s"}},"Metadata":{}}`, method, body)
}

type handleOrdersTest struct {
	method         string
	body           string
	expectedStatus int
	expectedBody   string
}

var handleOrdersTests = map[string]handleOrdersTest{
	"Put Request": {
		method:         "PUT",
		body:           `{\"items\":501}`,
		expectedStatus: http.StatusMethodNotAllowed,
		expectedBody:   "",
	},
	"Delete Request": {
		method:         "DELETE",
		body:           "",
		expectedStatus: http.StatusMethodNotAllowed,
		expectedBody:   "",
	},
}

func TestHandleNotAllowedMethodOrdersRequests(t *testing.T) {
	for name, test := range handleOrdersTests {
		t.Run(name, func(t *testing.T) {

			functionHostRequest := createAzureFunctionRequestBody(test.method, test.body)
			req := httptest.NewRequest("POST", "/", strings.NewReader(functionHostRequest))
			req.Header.Add("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			HandleOrdersRequest(rr, req)

			require.Equal(t, test.expectedStatus, rr.Result().StatusCode)

			var body InvokeResponse
			err := json.Unmarshal(rr.Body.Bytes(), &body)
			require.NoError(t, err)
		})
	}
}

func TestHandlePOSTOrdersRequest(t *testing.T) {
	functionHostRequest := createAzureFunctionRequestBody("POST", `{\"items\":501}`)
	req := httptest.NewRequest("POST", "/", strings.NewReader(functionHostRequest))
	req.Header.Add("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	HandleOrdersRequest(rr, req)

	require.Equal(t, 201, rr.Result().StatusCode)

	var body InvokeResponse
	err := json.Unmarshal(rr.Body.Bytes(), &body)
	require.NoError(t, err)

	output := body.Outputs["res"].(map[string]interface{})
	responseBody := output["body"].(string)

	var orderDoc model.OrderDocument
	err = json.Unmarshal([]byte(responseBody), &orderDoc)
	require.NoError(t, err)

	require.Equal(t, 501, orderDoc.Items)

	expectedPacks := []model.Pack{{Quantity: 500}, {Quantity: 250}}
	require.Equal(t, expectedPacks, orderDoc.Packs)
}

func TestHandleGETOrdersRequest(t *testing.T) {
	functionHostRequest := createAzureFunctionRequestBody("GET", "")
	req := httptest.NewRequest("POST", "/", strings.NewReader(functionHostRequest))
	req.Header.Add("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	HandleOrdersRequest(rr, req)

	require.Equal(t, 200, rr.Result().StatusCode)

	var body InvokeResponse
	err := json.Unmarshal(rr.Body.Bytes(), &body)
	require.NoError(t, err)

	output := body.Outputs["res"].(map[string]interface{})
	responseBody := output["body"].(string)

	var orderDocs []model.OrderDocument
	err = json.Unmarshal([]byte(responseBody), &orderDocs)
	require.NoError(t, err)
}
