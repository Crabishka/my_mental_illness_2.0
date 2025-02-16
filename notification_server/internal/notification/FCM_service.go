package notification

import (
	"context"
	"simple-service/internal/device"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"google.golang.org/api/option"
)

type Service struct {
	fcmClient  *messaging.Client
	deviceRepo device.DeviceRepository
}

// NewService создает новый экземпляр сервиса уведомлений
func NewService(credentialsFile string, deviceRepo device.DeviceRepository) (*Service, error) {
	opt := option.WithCredentialsFile(credentialsFile)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, err
	}

	fcmClient, err := app.Messaging(context.Background())
	if err != nil {
		return nil, err
	}

	return &Service{
		fcmClient:  fcmClient,
		deviceRepo: deviceRepo,
	}, nil
}

func (s *Service) SendToDevice(ctx context.Context, deviceUUID string, title, body string) error {
	device, err := s.deviceRepo.GetByUUID(deviceUUID)
	if err != nil {
		return err
	}

	message := &messaging.Message{
		Token: device.Token,
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
	}

	_, err = s.fcmClient.Send(ctx, message)
	return err
}

func (s *Service) SendToAll(ctx context.Context, title, body string) error {
	message := &messaging.Message{
		Topic: "all",
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
	}

	_, err := s.fcmClient.Send(ctx, message)
	return err
}
