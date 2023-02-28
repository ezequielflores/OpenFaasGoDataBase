package controller

import (
	"encoding/json"
	"fmt"
	controllerModel "handler/function/internal/adapter/controller/model"
	"handler/function/internal/application/model"
	"handler/function/internal/application/port/in"
	"handler/function/pkg"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type StarWarController struct {
	createCharacter in.CreateCharacter
	findCharacter   in.FindCharacter
}

func NewStarWarController(
	createCharacter in.CreateCharacter,
	findCharacter in.FindCharacter,
) *StarWarController {
	return &StarWarController{
		createCharacter: createCharacter,
		findCharacter:   findCharacter,
	}
}

func GetRequestBody(r *http.Request) (*controllerModel.CreaterCharacterRequest, *pkg.GenericException) {

	decoder := json.NewDecoder(r.Body)
	requestBody := &controllerModel.CreaterCharacterRequest{}
	err := decoder.Decode(requestBody)

	if err != nil {
		log.Printf("Error reading request body: %s", err.Error())
		return nil, &pkg.GenericException{
			Msj:        fmt.Sprintf("Invalid request body: %s", err.Error()),
			StatusCode: http.StatusBadRequest,
		}
	}

	return requestBody, nil
}

func (c *StarWarController) CreateStarWarCharacter(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	requestBody, requestError := GetRequestBody(r)

	if requestError != nil {
		w.WriteHeader(requestError.StatusCode)
		w.Write([]byte(requestError.Msj))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	searchResult, err := c.createCharacter.CreateCharacter(requestBody.ToDomain(), ctx)

	if err != nil {
		statusCode, errorMsj := pkg.GetErrorDetail(err)
		w.WriteHeader(statusCode)
		w.Write([]byte(errorMsj))
		return
	}

	jsonResult, jsonError := json.Marshal(controllerModel.CreateResponseFromDomain(searchResult))
	if jsonError != nil {
		log.Printf("Error converting to json %s", jsonError.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(jsonError.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonResult)
}

func (c *StarWarController) FindStarWarCharacter(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	split := strings.Split(r.URL.Path, "/")
	count := len(split)

	pathParam, pathError := strconv.Atoi(split[count-1])

	if pathError != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid character identifier"))
		return
	}

	characterIdentifier := &model.CharacterIdentifier{Id: pathParam}

	w.Header().Set("Content-Type", "application/json")
	findResult, err := c.findCharacter.FindCharacter(characterIdentifier, ctx)

	if err != nil {
		statusCode, errorMsj := pkg.GetErrorDetail(err)
		w.WriteHeader(statusCode)
		w.Write([]byte(errorMsj))
		return
	}

	jsonResult, jsonError := json.Marshal(controllerModel.FindResponseFromDomain(findResult))
	if jsonError != nil {
		log.Printf("Error converting to json %s", jsonError.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(jsonError.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonResult)
}
