package response

import (
	"encoding/json"
	"log"
	"net/http"
)

type CoreResponse struct {
	Status  http.ConnState `json:"status"`
	Message string         `json:"message"`
}

func SendJSON(w http.ResponseWriter, i interface{}) { // 200, success

	js, err := json.Marshal(i)

	if err != nil {
		http.Error(w, "JSON error:"+err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	w.Write(js)

}

func SendJSONMessage(w http.ResponseWriter, status http.ConnState, message string) {

	i := &CoreResponse{
		Status:  status,
		Message: message}

	js, err := json.Marshal(i)
	if err != nil {
		http.Error(w, "JSON error:"+err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(int(status))
	w.Write(js)
}
