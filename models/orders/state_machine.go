package orders

import (
	"errors"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/qor/transition"
)

var (
	// OrderState order's state machine
	OrderState = transition.New(&Order{})

	// ItemState order item's state machine
	ItemState = transition.New(&OrderItem{})
)

var (
	// DraftState draft state
	DraftState = "draft"
)

func init() {
	// Define Order's States
	OrderState.Initial("draft")

	OrderState.State("pending").Enter(func(value interface{}, tx *gorm.DB) (err error) {
		order := value.(*Order)
		tx.Model(order).Association("OrderItems").Find(&order.OrderItems)

		if err != nil {
			order.PaymentLog += "\n" + err.Error()
		} else {
			for idx, orderItem := range order.OrderItems {
				order.OrderItems[idx].Price = orderItem.SellingPrice()
			}
			order.PaymentAmount = order.Amount()
			order.PaymentTotal = order.Total()
		}
		return err
	})

	OrderState.State("processing").Enter(func(value interface{}, tx *gorm.DB) (err error) {
		order := value.(*Order)

		switch order.PaymentMethod {
		case COD:
		default:
			err = errors.New("unsupported pay method")
		}

		if err != nil {
			order.PaymentLog += "\n" + err.Error()
		}

		return
	})

	OrderState.State("cancelled").Enter(func(value interface{}, tx *gorm.DB) (err error) {
		order := value.(*Order)
		method := ""

		switch order.PaymentMethod {
		case COD:
		default:
			err = errors.New("unsupported pay method")
		}

		order.PaymentLog += "\n\n" + method + "\n" + fmt.Sprintf("Order cancelled at %#v", time.Now().String())

		if err != nil {
			order.PaymentLog += fmt.Sprintf("with error %v", err.Error())
		} else {
			now := time.Now()
			order.CancelledAt = &now
		}
		return
	})

	OrderState.State("shipped").Enter(func(value interface{}, tx *gorm.DB) (err error) {
		order := value.(*Order)

		switch order.PaymentMethod {
		case COD:
		default:
			err = errors.New("unsupported pay method")
		}

		if err != nil {
			order.PaymentLog += "\n" + err.Error()
		} else {
			now := time.Now()
			order.ShippedAt = &now
		}
		return
	})

	OrderState.State("paid_cancelled").Enter(func(value interface{}, tx *gorm.DB) (err error) {
		order := value.(*Order)

		switch order.PaymentMethod {
		case COD:
		default:
			err = errors.New("unsupported pay method")
		}

		order.PaymentLog += "\n\nRefund\n" + fmt.Sprintf("Order paid cancelled at %#v", time.Now().String())

		if err != nil {
			order.PaymentLog += fmt.Sprintf("with error %v", err.Error())
		} else {
			now := time.Now()
			order.CancelledAt = &now
		}
		return
	})

	OrderState.State("returned").Enter(func(value interface{}, tx *gorm.DB) error {
		order := value.(*Order)

		// check returned or not
		now := time.Now()
		order.ReturnedAt = &now
		return nil
	})

	OrderState.Event("checkout").To("pending").From("draft").After(func(value interface{}, tx *gorm.DB) (err error) {
		order := value.(*Order)
		for _, item := range order.OrderItems {
			ItemState.Trigger("checkout", &item, tx)
		}
		return nil
	})

	OrderState.Event("process").To("processing").From("pending").After(func(value interface{}, tx *gorm.DB) (err error) {
		order := value.(*Order)
		tx.Model(order).Association("OrderItems").Find(&order.OrderItems)

		for _, item := range order.OrderItems {
			ItemState.Trigger("process", &item, tx)
		}
		return nil
	})

	cancelEvent := OrderState.Event("cancel")
	cancelEvent.To("cancelled").From("draft", "pending").After(func(value interface{}, tx *gorm.DB) (err error) {
		order := value.(*Order)
		tx.Model(order).Association("OrderItems").Find(&order.OrderItems)

		for _, item := range order.OrderItems {
			ItemState.Trigger("cancel", &item, tx)
		}
		return nil
	})

	cancelEvent.To("paid_cancelled").From("processing", "shipped").After(func(value interface{}, tx *gorm.DB) (err error) {
		order := value.(*Order)
		tx.Model(order).Association("OrderItems").Find(&order.OrderItems)

		for _, item := range order.OrderItems {
			ItemState.Trigger("cancel", &item, tx)
		}
		return nil
	})

	OrderState.Event("ship").To("shipped").From("processing").After(func(value interface{}, tx *gorm.DB) (err error) {
		order := value.(*Order)
		tx.Model(order).Association("OrderItems").Find(&order.OrderItems)

		for _, item := range order.OrderItems {
			ItemState.Trigger("ship", &item, tx)
		}
		return nil
	})

	OrderState.Event("return").To("returned").From("shipped").After(func(value interface{}, tx *gorm.DB) (err error) {
		order := value.(*Order)
		tx.Model(order).Association("OrderItems").Find(&order.OrderItems)

		for _, item := range order.OrderItems {
			ItemState.Trigger("return", &item, tx)
		}
		return nil
	})

	// Define ItemItem's States
	ItemState.Initial("draft")
	ItemState.State("pending").Enter(func(value interface{}, tx *gorm.DB) error {
		// freeze stock, update order state
		return nil
	})
	ItemState.State("cancelled").Enter(func(value interface{}, tx *gorm.DB) error {
		// release stock, upate order state
		return nil
	})
	ItemState.State("processing")
	ItemState.State("shipped")
	ItemState.State("paid_cancelled").Enter(func(value interface{}, tx *gorm.DB) error {
		// do refund, release stock, update order state
		return nil
	})
	ItemState.State("returned")

	ItemState.Event("checkout").To("pending").From("draft")
	ItemState.Event("process").To("processing").From("checkout")
	cancelItemEvent := ItemState.Event("cancel")
	cancelItemEvent.To("cancelled").From("checkout")
	cancelItemEvent.To("paid_cancelled").From("paid")
	ItemState.Event("process").To("processing").From("paid")
	ItemState.Event("ship").To("shipped").From("processing")
	ItemState.Event("return").To("returned").From("shipped")
}
