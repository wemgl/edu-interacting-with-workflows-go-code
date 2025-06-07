package queries

import (
	"context"
	"errors"
	"time"

	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/workflow"
)

type FulfillOrderSignal struct {
	Fulfilled bool
}

func Workflow(ctx workflow.Context, input string) (string, error) {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
		// Example of supporting Workflow Cancellation. This setting allow all in progress Activities to either fail,
		// complete, or accept the Cancellation.
		WaitForCancellation: true,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	logger := workflow.GetLogger(ctx)
	logger.Info("Query workflow started", "input", input)

	// Custom search attributes can be added to from the workflow like this.
	// Upserting a custom search attribute like this overwrites any that were sent by the client.
	err := workflow.UpsertTypedSearchAttributes(ctx, CustomerIDSearchAttribute.ValueSet("3"))
	if err != nil {
		return "", err
	}

	// To remove a custom search attribute that was previously set:
	//err = workflow.UpsertTypedSearchAttributes(ctx, CustomerIDSearchAttribute.ValueUnset())
	//if err != nil {
	//	return "", err
	//}

	// Example of supporting Workflow Cancellation
	defer func() {
		// This logic ensures cleanup only happens if there is a Cancellation error
		if !errors.Is(ctx.Err(), workflow.ErrCanceled) {
			return
		}
		logger.Info("Workflow cancelled. Starting cleanup Activity now…")
		// For the Workflow to execute an Activity after it receives a Cancellation Request
		// It has to get a new disconnected context.
		disconnectedCtx, _ := workflow.NewDisconnectedContext(ctx)
		err := workflow.ExecuteActivity(disconnectedCtx, CleanupActivity).Get(disconnectedCtx, nil)
		if err != nil {
			logger.Error("CleanupActivity failed", "Error", err)
		}
	}()

	currentState := "started"
	queryType := "current_state"
	err = workflow.SetQueryHandler(ctx, queryType, func() (string, error) {
		return currentState, nil
	})
	if err != nil {
		currentState = "failed to register query handler"
		return "", err
	}

	currentState = "waiting for signal"

	var signal FulfillOrderSignal
	var result string

	signalChan := workflow.GetSignalChannel(ctx, "fulfill-order-signal")
	signalChan.Receive(ctx, &signal)
	if signal.Fulfilled == true {
		err := workflow.ExecuteActivity(ctx, Activity, input).Get(ctx, &result)
		if err != nil {
			logger.Error("Activity failed.", "Error", err)
			return "", err
		}
		currentState = "workflow completed"
		logger.Info("Signal workflow completed.", "result", result)
	}

	return result, nil
}

func Activity(ctx context.Context, input string) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Activity", "input", input)

	return "Received " + input, nil
}

func CleanupActivity(ctx context.Context) error {
	return nil
}
