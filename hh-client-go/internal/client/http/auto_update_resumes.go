package httpclient

import (
	"context"
	"log/slog"
	"time"
)

const intervalAutoUpdateResumes = 4 * time.Hour

func (c *Client) AutoUpdateResumes(ctx context.Context, slResumeID []string, chDone chan<- struct{}) {
	slog.Info("Start auto update resumes")

	c.updateResumes(ctx, slResumeID)
	ticker := time.NewTicker(intervalAutoUpdateResumes)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			c.updateResumes(ctx, slResumeID)
		case <-ctx.Done():
			slog.Info("Stop auto update resumes")
			chDone <- struct{}{}
			return
		}
	}
}

func (c *Client) updateResumes(ctx context.Context, slResumeID []string) {
	for _, resumeID := range slResumeID {
		if err := c.PublishResume(ctx, resumeID); err != nil {
			slog.Error("Update resume with id %q: %v", resumeID, err)
		}
	}
}
