package router

import (
	"encoding/json"
	"net/http"

	"jordanfinners/api/handlers"
	"jordanfinners/api/model"
)

// InvokeRequest represents an Azure Inbound object to a handler
type InvokeRequest struct {
	Data     map[string]json.RawMessage
	Metadata map[string]interface{}
}

// InvokeResponse represents the expected output from a handler
type InvokeResponse struct {
	Outputs     map[string]interface{}
	Logs        []string
	ReturnValue interface{}
}

func handleInvokeRequest(r *http.Request) model.Request {
	var invokeRequest InvokeRequest

	d := json.NewDecoder(r.Body)
	d.Decode(&invokeRequest)

	var reqData map[string]interface{}
	json.Unmarshal(invokeRequest.Data["req"], &reqData)

	method := reqData["Method"].(string)

	body := ""
	if reqData["Body"] != nil {
		body = reqData["Body"].(string)
	}
	queryParams := ""
	if reqData["Query"] != nil {
		queryParams = reqData["Query"].(string)
	}

	return model.Request{
		Method:      method,
		Body:        body,
		QueryParams: queryParams,
	}
}

func handleInvokeResponse(w http.ResponseWriter, response model.Response) {
	outputs := make(map[string]interface{})
	res := make(map[string]interface{})

	headers := make(map[string]interface{})
	headers["content-type"] = "application/json"
	res["headers"] = headers

	res["body"] = response.Body
	res["status"] = response.Status
	outputs["res"] = res

	responseJSON, _ := json.Marshal(
		InvokeResponse{
			Outputs:     outputs,
			Logs:        nil,
			ReturnValue: nil,
		})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.Status)
	w.Write(responseJSON)
}

// HandleOrdersRequest deals with inbound requests to any path /orders
func HandleOrdersRequest(w http.ResponseWriter, r *http.Request) {
	request := handleInvokeRequest(r)
	switch request.Method {
	case "GET":
		response := handlers.GetOrders(request)
		handleInvokeResponse(w, response)
	case "POST":
		response := handlers.PostOrders(request)
		handleInvokeResponse(w, response)
	default:
		response := model.Response{
			Status: http.StatusMethodNotAllowed,
			Body:   "",
		}
		handleInvokeResponse(w, response)
		return
	}
}
