package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/thats-insane/awt-quiz3/internal/data"
	"github.com/thats-insane/awt-quiz3/internal/validator"
)

func (a *appDependencies) createUserHandler(w http.ResponseWriter, r *http.Request) {
	var incomingData struct {
		Name  string `json:"fullname"`
		Email string `json:"email"`
	}

	err := a.readJSON(w, r, &incomingData)
	if err != nil {
		a.badRequestResponse(w, r, err)
		return
	}

	user := &data.User{
		Name:  incomingData.Name,
		Email: incomingData.Email,
	}

	v := validator.New()

	data.Validate(v, user)

	if !v.IsEmpty() {
		a.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = a.userModel.Insert(user)
	if err != nil {
		a.serverErrResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/signup/%d", user.ID))

	data := envelope{
		"user": user,
	}

	err = a.writeJSON(w, http.StatusCreated, data, headers)

	if err != nil {
		a.serverErrResponse(w, r, err)
		return
	}

	fmt.Fprintf(w, "%+v\n", incomingData)
}

func (a *appDependencies) displayUserHandler(w http.ResponseWriter, r *http.Request) {
	id, err := a.readIDParam(r)

	if err != nil {
		a.notFoundResponse(w, r)
		return
	}

	user, err := a.userModel.Get(id)

	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			a.notFoundResponse(w, r)
		default:
			a.serverErrResponse(w, r, err)
		}

		return
	}

	data := envelope{
		"user": user,
	}

	err = a.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		a.serverErrResponse(w, r, err)
		return
	}
}

func (a *appDependencies) updateUserHandler(w http.ResponseWriter, r *http.Request) {
	id, err := a.readIDParam(r)

	if err != nil {
		a.notFoundResponse(w, r)
		return
	}

	user, err := a.userModel.Get(id)

	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			a.notFoundResponse(w, r)
		default:
			a.serverErrResponse(w, r, err)
		}

		return
	}

	var incomingData struct {
		Name  *string `json:"fullname"`
		Email *string `json:"email"`
	}

	err = a.readJSON(w, r, &incomingData)

	if err != nil {
		a.badRequestResponse(w, r, err)
		return
	}

	if incomingData.Name != nil {
		user.Name = *incomingData.Name
	}

	err = a.userModel.Update(user)

	if err != nil {
		a.serverErrResponse(w, r, err)
		return
	}

	data := envelope{
		"user": user,
	}

	err = a.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		a.serverErrResponse(w, r, err)
		return
	}
}

func (a *appDependencies) deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	id, err := a.readIDParam(r)

	if err != nil {
		a.notFoundResponse(w, r)
		return
	}

	err = a.userModel.Delete(id)

	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			a.notFoundResponse(w, r)
		default:
			a.serverErrResponse(w, r, err)
		}

		return
	}

	data := envelope{
		"message": "user successfully deleted",
	}

	err = a.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		a.serverErrResponse(w, r, err)
	}
}
