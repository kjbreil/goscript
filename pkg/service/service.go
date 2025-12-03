package service

import (
	"context"
	"log/slog"
	"strings"
	"time"

	hass_ws "github.com/kjbreil/hass-ws/pkg/hass"
	"github.com/kjbreil/hass-ws/services"
)

const (
	serviceTickInterval = 100 * time.Millisecond
)

func Run(ctx context.Context, ws *hass_ws.Client, logger *slog.Logger, sChan Chan) {
	chanBuffer := make(map[string]services.Service)
	ticker := time.NewTicker(serviceTickInterval)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case s := <-sChan:
				// logger.Info(fmt.Sprintf("%s", s.JSON()))
				entitiesString := strings.Join(s.Targets(), ",")
				chanBuffer[entitiesString] = s
			case <-ticker.C:
				for k, v := range chanBuffer {
					ws.CallService(v)
					delete(chanBuffer, k)
				}
			}
		}
	}()
}
