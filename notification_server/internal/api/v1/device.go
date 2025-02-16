package v1

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"simple-service/internal/device"
)

// DeviceHandler обрабатывает HTTP-запросы для управления устройствами
type DeviceHandler struct {
	store *device.PostgresRepository
}

// NewDeviceHandler создает новый обработчик с заданным репозиторием
func NewDeviceHandler(store *device.PostgresRepository) *DeviceHandler {
	return &DeviceHandler{store: store}
}

// DeviceRequest представляет запрос на создание устройства
type DeviceRequest struct {
	Token string `json:"token" example:"fcm-token-123"`
	Model string `json:"model" example:"iPhone 12"`
}

// @Summary     Обработка запросов к устройствам
// @Description Обрабатывает GET и POST запросы для получения списка устройств и создания новых устройств
// @Tags        devices
// @Accept      json
// @Produce     json
// @Success     200 {array}  device.Device
// @Router      /devices [get]
func (h *DeviceHandler) HandleDevices(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetDevices(w, r)
	case http.MethodPost:
		h.CreateDevice(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// @Summary     Получить список всех устройств
// @Description Возвращает список всех зарегистрированных устройств
// @Tags        devices
// @Produce     json
// @Success     200 {array}  device.Device
// @Failure     500 {object} ErrorResponse
// @Router      /devices [get]
func (h *DeviceHandler) GetDevices(w http.ResponseWriter, r *http.Request) {
	devices, err := h.store.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(devices)
}

// @Summary     Создать новое устройство
// @Description Регистрирует новое устройство в системе
// @Tags        devices
// @Accept      json
// @Produce     json
// @Param       device body DeviceRequest true "Информация об устройстве"
// @Success     200 {object} device.Device
// @Failure     400 {object} ErrorResponse
// @Router      /devices [post]
func (h *DeviceHandler) CreateDevice(w http.ResponseWriter, r *http.Request) {
	var req DeviceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	device := &device.Device{
		Token: req.Token,
		Model: req.Model,
	}

	if err := h.store.Create(device); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(device)
}

// @Summary     Обработка запросов к конкретному устройству
// @Description Обрабатывает GET и DELETE запросы для получения информации об устройстве и его удаления
// @Tags        devices
// @Param       id path int true "ID устройства"
// @Success     200 {object} device.Device
// @Failure     400 {object} ErrorResponse
// @Failure     404 {object} ErrorResponse
// @Router      /devices/{id} [get]
func (h *DeviceHandler) HandleDevice(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/devices/")
	if idStr == "" {
		http.Error(w, "Device ID is required", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.GetDevice(w, r, idStr)
	case http.MethodDelete:
		h.DeleteDevice(w, r, idStr)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// @Summary     Получить информацию об устройстве
// @Description Возвращает детальную информацию об устройстве по его ID
// @Tags        devices
// @Produce     json
// @Param       id path int true "ID устройства"
// @Success     200 {object} device.Device
// @Failure     400 {object} ErrorResponse
// @Failure     404 {object} ErrorResponse
// @Router      /devices/{id} [get]
func (h *DeviceHandler) GetDevice(w http.ResponseWriter, r *http.Request, idStr string) {
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid device ID", http.StatusBadRequest)
		return
	}

	device, err := h.store.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(device)
}

// @Summary     Удалить устройство
// @Description Удаляет устройство из системы по его ID
// @Tags        devices
// @Param       id path int true "ID устройства"
// @Success     204 "No Content"
// @Failure     400 {object} ErrorResponse
// @Failure     404 {object} ErrorResponse
// @Router      /devices/{id} [delete]
func (h *DeviceHandler) DeleteDevice(w http.ResponseWriter, r *http.Request, idStr string) {
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid device ID", http.StatusBadRequest)
		return
	}

	if err := h.store.Delete(id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ErrorResponse представляет структуру ответа с ошибкой
type ErrorResponse struct {
	Error string `json:"error" example:"error message"`
}
