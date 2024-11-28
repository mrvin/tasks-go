package httpclient

import (
	"context"
	"fmt"
	"log/slog"
	"time"
)

func (c *Client) refreshToken() {
	userAuth, err := c.getToken(context.TODO(), "", "", "")
	if err != nil {
		slog.Error("Refresh token: " + err.Error())
	}
	c.mutexUserAuth.Lock()
	c.userAuth = userAuth
	c.mutexUserAuth.Unlock()

	c.mutexUserAuth.RLock()
	fmt.Println(c.userAuth)
	time.AfterFunc(c.userAuth.expiresIn, c.refreshToken)
	slog.Info("Refresh token will start", slog.String("duration", c.userAuth.expiresIn.String()))
	c.mutexUserAuth.RUnlock()
}
