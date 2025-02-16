package device

import (
	"time"
)

// DeviceRequest представляет запрос на создание/обновление устройства
// @Description Запрос на создание или обновление устройства
type DeviceRequest struct {
	UUID  string `json:"uuid" example:"device-uuid-123" binding:"required"`
	Token string `json:"token" example:"fcm-token-123" binding:"required"`
	Model string `json:"model" example:"iPhone 12"`
}

// Device представляет информацию об устройстве
// @Description Информация об устройстве
type Device struct {
	UUID        string    `json:"uuid" gorm:"primaryKey;type:string"`
	Token       string    `json:"token" gorm:"unique"`
	Model       string    `json:"model"`
	FirstSeenAt time.Time `json:"first_seen_at"`
	LastSeenAt  time.Time `json:"last_seen_at"`
}

type DeviceRepository interface {
	Create(device *Device) error
	Update(device *Device) error
	Delete(uuid string) error
	GetByUUID(uuid string) (*Device, error)
	GetByToken(token string) (*Device, error)
	GetAll() ([]*Device, error)
}
