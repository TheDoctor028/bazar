package main

import (
	"fmt"
	"log"
	"time"

	"github.com/TheDoctor028/bazar/app/admin"
	"github.com/TheDoctor028/bazar/internal/config/auth"
	_ "github.com/TheDoctor028/bazar/internal/config/db/migrations"
	"github.com/TheDoctor028/bazar/models/blogs"
	"github.com/TheDoctor028/bazar/models/orders"
	"github.com/TheDoctor028/bazar/models/products"
	adminseo "github.com/TheDoctor028/bazar/models/seo"
	"github.com/TheDoctor028/bazar/models/settings"
	"github.com/TheDoctor028/bazar/models/stores"
	"github.com/TheDoctor028/bazar/models/users"
	"github.com/qor/auth/auth_identity"
	"github.com/qor/auth/providers/password"
	"github.com/qor/banner_editor"
	"github.com/qor/help"
	i18n_database "github.com/qor/i18n/backends/database"
	"github.com/qor/media/asset_manager"
	"github.com/qor/notification"
	"github.com/qor/notification/channels/database"
	"github.com/qor/qor"
)

/* How to run this script
   $ go run db/seeds/main.go db/seeds/seeds.go
*/

/* How to upload file
 * $ brew install s3cmd
 * $ s3cmd --configure (Refer https://github.com/theplant/qor-example)
 * $ s3cmd put local_file_path s3://qor3/
 */

var (
	AdminUser    *users.User
	Notification = notification.New(&notification.Config{})
	Tables       = []interface{}{
		&auth_identity.AuthIdentity{},
		&users.User{}, &users.Address{},
		&products.Category{}, &products.Color{}, &products.Size{}, &products.Material{}, &products.Collection{},
		&products.Product{}, &products.ProductImage{}, &products.ColorVariation{}, &products.SizeVariation{},
		&stores.Store{},
		&orders.Order{}, &orders.OrderItem{},
		&settings.Setting{},
		&adminseo.MySEOSetting{},
		&blogs.Article{},
		&settings.MediaLibrary{},
		&banner_editor.QorBannerEditorSetting{},

		&asset_manager.AssetManager{},
		&i18n_database.Translation{},
		&notification.QorNotification{},
		&admin.QorWidgetSetting{},
		&help.QorHelpEntry{},
	}
)

func main() {
	Notification.RegisterChannel(database.New(&database.Config{}))
	TruncateTables(Tables...)
	createRecords()
}

func createRecords() {
	fmt.Println("Start create sample data...")

	createSetting()
	fmt.Println("--> Created setting.")

	fmt.Println("--> Done!")
}

func createSetting() {
	setting := settings.Setting{}

	setting.ShippingFee = Seeds.Setting.ShippingFee
	setting.GiftWrappingFee = Seeds.Setting.GiftWrappingFee
	setting.CODFee = Seeds.Setting.CODFee
	setting.TaxRate = Seeds.Setting.TaxRate
	setting.Address = Seeds.Setting.Address
	setting.Region = Seeds.Setting.Region
	setting.City = Seeds.Setting.City
	setting.Country = Seeds.Setting.Country
	setting.Zip = Seeds.Setting.Zip
	setting.Latitude = Seeds.Setting.Latitude
	setting.Longitude = Seeds.Setting.Longitude

	if err := DraftDB.Create(&setting).Error; err != nil {
		log.Fatalf("create setting (%v) failure, got err %v", setting, err)
	}
}

func createAdminUsers() {
	AdminUser = &users.User{}
	AdminUser.Email = "dev@getqor.com"
	AdminUser.Confirmed = true
	AdminUser.Name = "QOR Admin"
	AdminUser.Role = "Admin"
	DraftDB.Create(AdminUser)

	provider := auth.Auth.GetProvider("password").(*password.Provider)
	hashedPassword, _ := provider.Encryptor.Digest("testing")
	now := time.Now()

	authIdentity := &auth_identity.AuthIdentity{}
	authIdentity.Provider = "password"
	authIdentity.UID = AdminUser.Email
	authIdentity.EncryptedPassword = hashedPassword
	authIdentity.UserID = fmt.Sprint(AdminUser.ID)
	authIdentity.ConfirmedAt = &now

	DraftDB.Create(authIdentity)

	// Send welcome notification
	Notification.Send(&notification.Message{
		From:        AdminUser,
		To:          AdminUser,
		Title:       "Welcome To QOR Admin",
		Body:        "Welcome To QOR Admin",
		MessageType: "info",
	}, &qor.Context{DB: DraftDB})
}
