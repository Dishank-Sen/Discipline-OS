package errorhandler

import (
	"encoding/json"
	"log"
	"net/http"
	"github.com/go-playground/validator/v10"
)

var Validate = validator.New()

func HandleError(err error, cmt string){
	if err != nil{
		log.Fatal(cmt,": ", err)
	}
}

func WriteJSON(w http.ResponseWriter, status int, v any) error{
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(v)
}

func WriteError(w http.ResponseWriter, status int, err error){
	WriteJSON(w, status, map[string]string{"error": err.Error()})
}