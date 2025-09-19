package routes

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/Dishank-Sen/Discipline-OS/interfaces"
	"github.com/Dishank-Sen/Discipline-OS/service/auth"
	"github.com/Dishank-Sen/Discipline-OS/types/payload"
	errorhandler "github.com/Dishank-Sen/Discipline-OS/utils/errorHandler"
	payloadhandler "github.com/Dishank-Sen/Discipline-OS/utils/payloadHandler"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Handler struct{
	store interfaces.UserStore
	UserCollection *mongo.Collection
	TempUserCollection *mongo.Collection
}

func NewHandler(userStore interfaces.UserStore, userCollection *mongo.Collection, tempUserCollection *mongo.Collection) *Handler{
	return &Handler{
		store: userStore,
		UserCollection: userCollection,
		TempUserCollection: tempUserCollection,
	}
}

func (h *Handler) RegisterRoutes(router *mux.Router){
	router.HandleFunc("/signup", h.handleSignup)
	router.HandleFunc("/signup/email", h.handleEmail)
	router.HandleFunc("/signup/password", h.handlePassword)
	router.HandleFunc("/signup/verify-otp", h.handleOTP)
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

	
}

func (h *Handler) handleEmail(w http.ResponseWriter, r *http.Request){
	var emailPayload payload.EmailPayload

	err := payloadhandler.ParseJSON(r, &emailPayload)
	if err != nil{
		errors := err.(validator.ValidationErrors)
		errorhandler.WriteError(w, http.StatusInternalServerError, fmt.Errorf("error in parsing: %w", errors))
		return
	}

	if err := errorhandler.Validate.Struct(emailPayload); err != nil{
		errorhandler.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload error: %w",err))
		return
	}

	// insert email to temporary database collection
	signupToken, err := h.store.InsertEmail(emailPayload.Email, h.TempUserCollection)

	// if user already exists
	if (err != nil) && (signupToken == ""){
		errorhandler.WriteError(w, http.StatusConflict, err)
		return
	}

	// if some internal server error
	if (err != nil) && (signupToken != ""){
		// delete the user with this signup token
		filter := bson.M{
			"signupToken": signupToken,
		}
		_, deleteErr := h.store.DeleteRecord(context.Background(), h.TempUserCollection, filter)
		log.Fatal(deleteErr)

		errorhandler.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	errorhandler.WriteJSON(w, http.StatusAccepted, map[string]string{"signupToken": signupToken, "message": "email saved"})
}

func (h *Handler) handlePassword(w http.ResponseWriter, r *http.Request){
	var passwordPayload payload.PasswordPayload

	err := payloadhandler.ParseJSON(r, &passwordPayload)
	if err != nil{
		errors := err.(validator.ValidationErrors)
		errorhandler.WriteError(w, http.StatusInternalServerError, fmt.Errorf("error in parsing: %w", errors))
		return
	}

	if err := errorhandler.Validate.Struct(passwordPayload); err != nil{
		
		// delete the user with this signup token
		filter := bson.M{
			"signupToken": passwordPayload.SignupToken,
		}
		_, deleteErr := h.store.DeleteRecord(context.Background(), h.TempUserCollection, filter)
		log.Fatal(deleteErr)
		
		errorhandler.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload error: %w",err))
		return
	}

	hashedPassword, err := auth.HashPassword(passwordPayload.Password)
	if err != nil{
		// delete the user with this signup token
		filter := bson.M{
			"signupToken": passwordPayload.SignupToken,
		}
		_, deleteErr := h.store.DeleteRecord(context.Background(), h.TempUserCollection, filter)
		log.Fatal(deleteErr)

		errorhandler.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	err = h.store.InsertPassword(hashedPassword, passwordPayload.SignupToken, h.TempUserCollection)

	if err != nil{
		errorhandler.WriteError(w, http.StatusBadRequest, err)
		return
	}

	errorhandler.WriteJSON(w, http.StatusAccepted, map[string]string{"message": "password saved"})
}

func (h *Handler) handleOTP(w http.ResponseWriter, r *http.Request){
	var otpPayload payload.OTPPayload
	
	err := payloadhandler.ParseJSON(r, &otpPayload)
	if err != nil{
		errors := err.(validator.ValidationErrors)
		errorhandler.WriteError(w, http.StatusInternalServerError, fmt.Errorf("error in parsing: %w", errors))
		return
	}

	if err := errorhandler.Validate.Struct(otpPayload); err != nil{
		errorhandler.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload error: %w",err))
		return
	}

	
}