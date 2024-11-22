package admin

import (
	"path"
	"regexp"
	"time"

	"github.com/TheDoctor028/bazar/internal/config/db"
	"github.com/qor/admin"
	"github.com/qor/notification"
	"github.com/qor/notification/channels/database"
)

// SetupNotification add notification
func SetupNotification(Admin *admin.Admin) {
	Notification := notification.New(&notification.Config{})
	Notification.RegisterChannel(database.New(&database.Config{DB: db.DB}))
	Notification.Action(&notification.Action{
		Name:         "Check it out",
		MessageTypes: []string{"order_paid_cancelled", "order_processed", "order_returned"},
		URL: func(data *notification.QorNotification, context *admin.Context) string {
			return path.Join("/admin/orders/", regexp.MustCompile(`#(\d+)`).FindStringSubmatch(data.Body)[1])
		},
	})
	Notification.Action(&notification.Action{
		Name:         "Dismiss",
		MessageTypes: []string{"order_paid_cancelled", "info", "order_processed", "order_returned"},
		Visible: func(data *notification.QorNotification, context *admin.Context) bool {
			return data.ResolvedAt == nil
		},
		Handler: func(argument *notification.ActionArgument) error {
			return argument.Context.GetDB().Model(argument.Message).Update("resolved_at", time.Now()).Error
		},
		Undo: func(argument *notification.ActionArgument) error {
			return argument.Context.GetDB().Model(argument.Message).Update("resolved_at", nil).Error
		},
	})
	Admin.NewResource(Notification)
}
