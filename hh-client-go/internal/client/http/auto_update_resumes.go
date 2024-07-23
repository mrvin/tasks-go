package httpclient

import (
	"context"
	"log/slog"
	"time"

	"github.com/mrvin/tasks-go/hh-client-go/pkg/retry"
)

const intervalAutoUpdateResumes = 4 * time.Hour
const retriesUpdate = 3

func (c *Client) AutoUpdateResumes(ctx context.Context, slResumeID []string, chDone chan<- struct{}) {
	slog.Info("Start auto update resumes")

	c.updateResumes(ctx, slResumeID)
	ticker := time.NewTicker(intervalAutoUpdateResumes)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if time.Now().Hour() != 2 {
				c.updateResumes(ctx, slResumeID)
			}
		case <-ctx.Done():
			slog.Info("Stop auto update resumes")
			chDone <- struct{}{}
			return
		}
	}
}

func (c *Client) retryUpdateResume(ctx context.Context, resumeID string, retries int) {
	retryUpdater := retry.Retry(c.PublishResume, retries)
	if err := retryUpdater(ctx, resumeID); err != nil {
		slog.Error("Update resume: "+err.Error(), slog.String("id", resumeID))
	}
}

func (c *Client) updateResumes(ctx context.Context, slResumeID []string) {
	for _, resumeID := range slResumeID {
		go func() {
			c.retryUpdateResume(ctx, resumeID, retriesUpdate)
		}()
	}
}
