package service

import (
	"context"
	hass_ws "github.com/kjbreil/hass-ws"
	"github.com/kjbreil/hass-ws/services"
	"log/slog"
	"strings"
	"time"
)

func Run(ctx context.Context, ws *hass_ws.Client, logger *slog.Logger, sChan Chan) {
	chanBuffer := make(map[string]services.Service)
	ticker := time.NewTicker(time.Millisecond * 100)
	go func() {
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
