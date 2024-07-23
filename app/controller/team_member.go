package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	help "github.com/adamnasrudin03/go-helpers"
	response_mapper "github.com/adamnasrudin03/go-helpers/response-mapper/v1"
	"github.com/adamnasrudin03/go-skeleton-chi/app/dto"
	"github.com/adamnasrudin03/go-skeleton-chi/app/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

type TeamMemberController interface {
	Create(w http.ResponseWriter, r *http.Request)
	GetDetail(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	GetList(w http.ResponseWriter, r *http.Request)
}

type TeamMemberHandler struct {
	Service  service.TeamMemberService
	Logger   *logrus.Logger
	Validate *validator.Validate
}

func NewTeamMemberDelivery(
	srv service.TeamMemberService,
	logger *logrus.Logger,
	validator *validator.Validate,
) TeamMemberController {
	return &TeamMemberHandler{
		Service:  srv,
		Logger:   logger,
		Validate: validator,
	}
}

func (c *TeamMemberHandler) Create(w http.ResponseWriter, r *http.Request) {
	var (
		opName = "TeamMemberController-Create"
		input  dto.TeamMemberCreateReq
		err    error
	)

	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		c.Logger.Errorf("%v error bind json: %v ", opName, err)
		response_mapper.RenderJSON(w, http.StatusBadRequest, response_mapper.ErrGetRequest())
		return
	}

	// validation input user
	err = c.Validate.Struct(input)
	if err != nil {
		response_mapper.RenderJSON(w, http.StatusBadRequest, response_mapper.FormatValidationError(err))
		return
	}

	res, err := c.Service.Create(r.Context(), input)
	if err != nil {
		response_mapper.RenderJSON(w, http.StatusInternalServerError, err)
		return
	}

	response_mapper.RenderJSON(w, http.StatusCreated, res)
}

func (c *TeamMemberHandler) GetDetail(w http.ResponseWriter, r *http.Request) {
	var (
		opName  = "TeamMemberController-GetDetail"
		idParam = strings.TrimSpace(chi.URLParam(r, "id"))
		err     error
	)

	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.Logger.Errorf("%v error parse param: %v ", opName, err)
		response_mapper.RenderJSON(w, http.StatusBadRequest, response_mapper.ErrInvalid("ID Anggota team", "Team Member ID"))
		return
	}

	res, err := c.Service.GetByID(r.Context(), id)
	if err != nil {
		response_mapper.RenderJSON(w, http.StatusInternalServerError, err)
		return
	}

	response_mapper.RenderJSON(w, http.StatusOK, res)
}

func (c *TeamMemberHandler) Delete(w http.ResponseWriter, r *http.Request) {
	var (
		opName  = "TeamMemberController-Delete"
		idParam = strings.TrimSpace(chi.URLParam(r, "id"))
		err     error
	)

	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.Logger.Errorf("%v error parse param: %v ", opName, err)
		response_mapper.RenderJSON(w, http.StatusBadRequest, response_mapper.ErrInvalid("ID Anggota team", "Team Member ID"))
		return
	}

	err = c.Service.DeleteByID(r.Context(), id)
	if err != nil {
		response_mapper.RenderJSON(w, http.StatusInternalServerError, err)
		return
	}

	response_mapper.RenderJSON(w, http.StatusOK, response_mapper.MultiLanguages{
		ID: "Anggota Tim Berhasil Dihapus",
		EN: "Team Member Deleted Successfully",
	})
}

func (c *TeamMemberHandler) Update(w http.ResponseWriter, r *http.Request) {
	var (
		opName  = "TeamMemberController-Update"
		idParam = strings.TrimSpace(chi.URLParam(r, "id"))
		input   dto.TeamMemberUpdateReq
		err     error
	)

	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.Logger.Errorf("%v error parse param: %v ", opName, err)
		response_mapper.RenderJSON(w, http.StatusBadRequest, response_mapper.ErrInvalid("ID Anggota team", "Team Member ID"))
		return
	}

	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		c.Logger.Errorf("%v error bind json: %v ", opName, err)
		response_mapper.RenderJSON(w, http.StatusBadRequest, response_mapper.ErrGetRequest())
		return
	}
	input.ID = id
	// validation input user
	err = c.Validate.Struct(input)
	if err != nil {
		response_mapper.RenderJSON(w, http.StatusBadRequest, response_mapper.FormatValidationError(err))
		return
	}

	err = c.Service.Update(r.Context(), input)
	if err != nil {
		response_mapper.RenderJSON(w, http.StatusInternalServerError, err)
		return
	}

	response_mapper.RenderJSON(w, http.StatusOK, response_mapper.MultiLanguages{
		ID: "Anggota Tim Berhasil Diperbarui",
		EN: "Team Member Updated Successfully",
	})
}

func (c *TeamMemberHandler) GetList(w http.ResponseWriter, r *http.Request) {
	var (
		opName = "TeamMemberController-GetList"
		query  = r.URL.Query()
		input  dto.TeamMemberListReq
		err    error
	)

	err = help.NewQueryDecoder(query).Decode(&input)
	if err != nil {
		c.Logger.Errorf("%v error bind json: %v ", opName, err)
		response_mapper.RenderJSON(w, http.StatusBadRequest, response_mapper.ErrGetRequest())
		return
	}

	fmt.Println(input)
	res, err := c.Service.GetList(r.Context(), input)
	if err != nil {
		c.Logger.Errorf("%v error: %v ", opName, err)
		response_mapper.RenderJSON(w, http.StatusInternalServerError, err)
		return
	}

	response_mapper.RenderJSON(w, http.StatusOK, res)
}
