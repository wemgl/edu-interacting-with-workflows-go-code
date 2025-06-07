package async

import (
	"context"
	"encoding/base64"
	"time"

	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/workflow"
)

func Workflow(ctx workflow.Context, input string) (string, error) {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 300 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	logger := workflow.GetLogger(ctx)
	logger.Info("Async workflow started", "input", input)

	var result string

	err := workflow.ExecuteActivity(ctx, Activity, input).Get(ctx, &result)
	if err != nil {
		logger.Error("Activity failed.", "Error", err)
		return "", err
	}

	logger.Info("Async workflow completed.", "result", result)
	return result, nil
}

func Activity(ctx context.Context, input string) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Activity", "input", input)

	activityInfo := activity.GetInfo(ctx)
	// Encode as hex…
	//taskToken := hex.EncodeToString(activityInfo.TaskToken)
	// Or Base64
	taskToken := base64.StdEncoding.EncodeToString(activityInfo.TaskToken)
	logger.Info("Async. Activity Completion", "taskToken", taskToken)

	return "", activity.ErrResultPending
}
