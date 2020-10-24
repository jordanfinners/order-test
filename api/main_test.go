package api

import (
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

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
		expectedStatus: http.StatusMethodNotAllowed,
		expectedBody:   "",
	},
	"Put Request": {
		method:         "PUT",
		body:           `items=501`,
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
		body:           `items=501`,
		expectedStatus: http.StatusCreated,
		expectedBody:   `{"items":501,"packs":[{"quantity":500},{"quantity":250}]}`,
	},
}

func TestHandleOrders(t *testing.T) {
	for name, test := range handleOrdersTests {
		t.Run(name, func(t *testing.T) {
			req := httptest.NewRequest(test.method, "/", strings.NewReader(test.body))
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

			rr := httptest.NewRecorder()
			HandleOrders(rr, req)

			body := rr.Body.String()

			require.Equal(t, test.expectedStatus, rr.Result().StatusCode)
			if test.expectedBody != "" {
				require.Equal(t, test.expectedBody, body)
			}
		})
	}
}

func TestHandleOrdersInvalidForm(t *testing.T) {
	badForm := "items="
	req := httptest.NewRequest("POST", "/", strings.NewReader(badForm))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	HandleOrders(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Result().StatusCode)
}

func TestHandleOrdersInvalidItemsOrdered(t *testing.T) {
	badForm := "items=Seven"
	req := httptest.NewRequest("POST", "/", strings.NewReader(badForm))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	HandleOrders(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Result().StatusCode)
}

type handleWebsiteTest struct {
	method         string
	body           string
	expectedStatus int
	expectedBody   string
}

var handleWebsiteTests = map[string]handleWebsiteTest{
	"Get Request": {
		method:         "GET",
		body:           "",
		expectedStatus: http.StatusOK,
		expectedBody:   "<!doctype html>\n<html lang=\"en\">\n\n<head>\n    <meta charset=\"utf-8\">\n    <meta name=\"generator\" content=\"Jordans Order Test\">\n    <meta name=\"viewport\" content=\"width=device-width, minimum-scale=1, initial-scale=1, user-scalable=yes\">\n\n    <title>Jordans Order Test</title>\n    <meta name=\"description\" content=\"Jordans Order Test\">\n    <meta name=\"author\" content=\"https://github.com/jordanfinners\">\n\n    <meta property=\"og:title\" content=\"Jordans Order Test\">\n    <meta property=\"og:description\" content=\"Jordans Order Test\">\n    <meta property=\"og:image\"\n        content=\"https://avatars2.githubusercontent.com/u/17813098?s=460&u=f8f61170c39933eff8aaf52f87bf6939ecdee7a6&v=4\">\n</head>\n\n<body>\n    <form name=\"order\" method=\"post\" action=\"http://localhost:8080/orders\">\n        <label>How many items do you wish to order?\n            <input type=\"number\" name=\"items\" placeholder=\"e.g. 501\" min=\"1\" required>\n        </label>\n        <button type=\"submit\">Submit Order</button>\n        <button type=\"reset\">Start Over</button>\n    </form>\n</body>\n\n</html>\n",
	},
	"Put Request": {
		method:         "PUT",
		body:           `{"items": 501}`,
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
		body:           `{"items": 501}`,
		expectedStatus: http.StatusMethodNotAllowed,
		expectedBody:   "",
	},
}

func TestHandleWebsite(t *testing.T) {

	err := os.Setenv("WEBSITE_FILE_PATH", "./static/index.html")
	if err != nil {
		log.Printf("Failed to set the WEBSITE_FILE_PATH: %v", err)
	}

	err = os.Setenv("ORDERS_API", "http://localhost:8080/orders")
	if err != nil {
		log.Printf("Failed to set the ORDERS_API: %v", err)
	}

	for name, test := range handleWebsiteTests {
		t.Run(name, func(t *testing.T) {
			req := httptest.NewRequest(test.method, "/", strings.NewReader(test.body))
			req.Header.Add("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			HandleWebsite(rr, req)

			body := rr.Body.String()

			require.Equal(t, test.expectedStatus, rr.Result().StatusCode)
			if test.expectedBody != "" {
				require.Equal(t, test.expectedBody, body)
			}
		})
	}
}
