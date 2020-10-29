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

type handleOrdersTest struct {
	method         string
	body           string
	expectedStatus int
	expectedBody   string
}

var handleOrdersTests = map[string]handleOrdersTest{
	"Get Request": {
		method:         "GET",
		body:           "",
		expectedStatus: http.StatusOK,
		expectedBody:   "",
	},
	"Put Request": {
		method:         "PUT",
		body:           `{"items":501}`,
		expectedStatus: http.StatusMethodNotAllowed,
		expectedBody:   "",
	},
	"Delete Request": {
		method:         "DELETE",
		body:           "",
		expectedStatus: http.StatusMethodNotAllowed,
		expectedBody:   "",
	},
	"Post Request": {
		method:         "POST",
		body:           `{\"items\":501}`,
		expectedStatus: http.StatusCreated,
		expectedBody:   `{"items":501,"packs":[{"quantity":500},{"quantity":250}]}`,
	},
}

func TestHandleOrdersRequest(t *testing.T) {
	for name, test := range handleOrdersTests {
		t.Run(name, func(t *testing.T) {
			functionHostRequest := fmt.Sprintf(`{"Data":{"req":{"Url":"http://localhost:7071/api/orders","Method":"%v","Query":"{}","Headers":{"Content-Type":["application/json"]},"Params":{},"Body":"%v"}},"Metadata":{}}`, test.method, test.body)

			req := httptest.NewRequest(test.method, "/", strings.NewReader(functionHostRequest))
			req.Header.Add("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			HandleOrdersRequest(rr, req)

			require.Equal(t, test.expectedStatus, rr.Result().StatusCode)

			var body InvokeResponse
			err := json.Unmarshal(rr.Body.Bytes(), &body)
			require.NoError(t, err)

			output := body.Outputs["res"].(map[string]interface{})
			responseBody := output["body"].(string)

			if test.expectedBody != "" {
				require.Equal(t, test.expectedBody, responseBody)
			}
		})
	}
}
