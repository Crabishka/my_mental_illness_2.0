package device

import "time"

// DeviceRequest представляет запрос на создание/обновление устройства
// @Description Запрос на создание или обновление устройства
type DeviceRequest struct {
    Token string `json:"token" example:"fcm-token-123" binding:"required"`
    Model string `json:"model" example:"iPhone 12"`
}

// Device представляет информацию об устройстве
// @Description Информация об устройстве
type Device struct {
    ID          int64     `json:"id"`
    Token       string    `json:"token"`
    Model       string    `json:"model"`
    FirstSeenAt time.Time `json:"first_seen_at"`
    LastSeenAt  time.Time `json:"last_seen_at"`
}

type DeviceRepository interface {
    Create(device *Device) error
    Update(device *Device) error
    Delete(id int64) error
    GetByID(id int64) (*Device, error)
    GetAll() ([]*Device, error)
} 