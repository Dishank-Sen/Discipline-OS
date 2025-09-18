package routes

import (
	"fmt"
	"net/http"

	"github.com/Dishank-Sen/Discipline-OS/interfaces"
	"github.com/Dishank-Sen/Discipline-OS/types/payload"
	errorhandler "github.com/Dishank-Sen/Discipline-OS/utils/errorHandler"
	payloadhandler "github.com/Dishank-Sen/Discipline-OS/utils/payloadHandler"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type Handler struct{
	store interfaces.UserStore
}

func NewHandler(userStore interfaces.UserStore) *Handler{
	return &Handler{
		store: userStore,
	}
}

func (h *Handler) RegisterRoutes(router *mux.Router){
	router.HandleFunc("/signup", h.handleSignup)
	router.HandleFunc("/login",nil)
}

func (h *Handler) handleSignup(w http.ResponseWriter, r *http.Request){
	var userPayload payload.SignupPayload

	err := payloadhandler.ParseJSON(r, &userPayload)
	if err != nil{
		errors := err.(validator.ValidationErrors)
		errorhandler.WriteError(w, http.StatusInternalServerError, fmt.Errorf("error in parsing: %w", errors))
		return
	}

	if err := errorhandler.Validate.Struct(userPayload); err != nil{
		errorhandler.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload error: %w",err))
		return
	}

	// user, _ := h.store.GetUserByEmail("aklhfds")
	// collection := client.Database("testing").Collection("numbers")
}