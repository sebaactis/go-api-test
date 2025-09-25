package users

import (
	"encoding/json"
	"net/http"
	"sebaactis/go-api-simple/httpResponses"
	"sebaactis/go-api-simple/models"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type Handler struct {
	DB *gorm.DB
}

func NewHandler(db *gorm.DB) *Handler {
	return &Handler{DB: db}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {

	type int struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	var in int

	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		httpResponses.WriteError(w, http.StatusBadRequest, "JSON Inválido")
		return
	}

	if strings.TrimSpace(in.Name) == "" || !strings.Contains(in.Email, "@") {
		httpResponses.WriteError(w, http.StatusBadRequest, "Datos inválidos")
		return
	}

	u := models.User{
		Name:  in.Name,
		Email: in.Email,
	}

	if err := h.DB.Create(&u).Error; err != nil {
		// ej: unique violation en email
		httpResponses.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	httpResponses.WriteJSON(w, http.StatusBadRequest, u)
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {

	idStr := mux.Vars(r)["id"]

	id, err := strconv.ParseUint(idStr, 10, 64)

	if err != nil {
		httpResponses.WriteError(w, http.StatusBadRequest, "Id inválido")
		return
	}

	var u models.User

	if err := h.DB.First(&u, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			httpResponses.WriteError(w, http.StatusNotFound, "No encontrado")
			return
		}

		httpResponses.WriteError(w, http.StatusInternalServerError, err.Error())
		return

	}

	httpResponses.WriteJSON(w, http.StatusOK, u)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {

	var out []models.User

	if err := h.DB.Find(&out).Error; err != nil {
		httpResponses.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	httpResponses.WriteJSON(w, http.StatusOK, out)
}
