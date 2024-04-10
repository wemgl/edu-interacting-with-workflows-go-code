package pizza

import (
	"errors"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type FulfillOrderSignal struct {
	Fulfilled bool
}

var signal = FulfillOrderSignal{
	Fulfilled: true,
}

func PizzaWorkflow(ctx workflow.Context, order PizzaOrder) (OrderConfirmation, error) {
	retrypolicy := &temporal.RetryPolicy{
		MaximumInterval: time.Second * 10,
	}

	options := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 5,
		RetryPolicy:         retrypolicy,
	}

	ctx = workflow.WithActivityOptions(ctx, options)

	logger := workflow.GetLogger(ctx)

	var totalPrice int
	for _, pizza := range order.Items {
		totalPrice += pizza.Price
	}

	var distance Distance
	err := workflow.ExecuteActivity(ctx, GetDistance, order.Address).Get(ctx, &distance)
	if err != nil {
		logger.Error("Unable get distance", "Error", err)
		return OrderConfirmation{}, err
	}

	if order.IsDelivery && distance.Kilometers > 25 {
		return OrderConfirmation{}, errors.New("customer lives too far away for delivery")
	}

	signalChan := workflow.GetSignalChannel(ctx, "fulfill-order-signal")
	signalChan.Receive(ctx, &signal)
	if signal.Fulfilled == true {
		bill := Bill{
			CustomerID:  order.Customer.CustomerID,
			OrderNumber: order.OrderNumber,
			Amount:      totalPrice,
			Description: "Pizza",
		}

		var confirmation OrderConfirmation
		err = workflow.ExecuteActivity(ctx, SendBill, bill).Get(ctx, &confirmation)
		if err != nil {
			logger.Error("Unable to bill customer", "Error", err)
			return OrderConfirmation{}, err
		}
	}

	return OrderConfirmation{}, nil
}

func FulFillOrderWorkflow(ctx workflow.Context, order PizzaOrder) error {
	retrypolicy := &temporal.RetryPolicy{
		MaximumInterval: time.Second * 10,
	}

	options := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 5,
		RetryPolicy:         retrypolicy,
	}

	ctx = workflow.WithActivityOptions(ctx, options)

	logger := workflow.GetLogger(ctx)

	err := workflow.ExecuteActivity(ctx, MakePizzas, order)
	if err != nil {
		logger.Error("Unable to make pizzas", "Error", err)
	}

	err = workflow.ExecuteActivity(ctx, DeliverPizzas, order)
	if err != nil {
		logger.Error("Unable to deliver pizzas", "Error", err)
	}

	err = workflow.SignalExternalWorkflow(ctx, "signals", "", "fulfill-order-signal", signal)
	if err != nil {
		logger.Error("Error sending the Signal", "Error", err)
	}

	return nil
}
