package v1

import (
	"encoding/json"
	"net/http"

	"simple-service/internal/notification"
)

type NotificationHandler struct {
	service *notification.Service
}

func NewNotificationHandler(service *notification.Service) *NotificationHandler {
	return &NotificationHandler{service: service}
}

// NotificationRequest представляет запрос на отправку уведомления
type NotificationRequest struct {
	Title string `json:"title" example:"Новое сообщение"`
	Body  string `json:"body" example:"Привет, как дела?"`
}

// @Summary     Отправить уведомление на устройство
// @Description Отправляет push-уведомление на конкретное устройство
// @Tags        notifications
// @Accept      json
// @Produce     json
// @Param       id path int true "ID устройства"
// @Param       notification body NotificationRequest true "Данные уведомления"
// @Success     200
// @Failure     400 {object} ErrorResponse
// @Failure     404 {object} ErrorResponse
// @Router      /notify/device/{id} [post]
func (h *NotificationHandler) SendToDevice(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	uuid := r.URL.Path[len("/notify/device/"):]
	if uuid == "" {
		http.Error(w, "Invalid device UUID", http.StatusBadRequest)
		return
	}

	var req NotificationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.SendToDevice(r.Context(), uuid, req.Title, req.Body); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// @Summary     Отправить уведомление всем устройствам
// @Description Отправляет push-уведомление всем зарегистрированным устройствам
// @Tags        notifications
// @Accept      json
// @Produce     json
// @Param       body body NotificationRequest true "Параметры уведомления"
// @Success     200
// @Failure     400  {object} ErrorResponse
// @Failure     500  {object} ErrorResponse
// @Router      /notify/all [post]
func (h *NotificationHandler) SendToAll(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req NotificationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.SendToAll(r.Context(), req.Title, req.Body); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
