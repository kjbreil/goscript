package service

import (
	"context"
	"fmt"
	"github.com/go-logr/logr"
	hass_ws "github.com/kjbreil/hass-ws"
	"github.com/kjbreil/hass-ws/services"
	"strings"
	"time"
)

func Run(ctx context.Context, ws *hass_ws.Client, logger logr.Logger, sChan Chan) {
	chanBuffer := make(map[string]services.Service)
	ticker := time.NewTicker(time.Millisecond * 100)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case s := <-sChan:
				logger.V(4).Info(fmt.Sprintf("%s", s.JSON()))
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
