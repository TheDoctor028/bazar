package account

import (
	"net/http"

	"github.com/TheDoctor028/bazar/models/users"
	"github.com/TheDoctor028/bazar/utils"
	"github.com/qor/render"
)

// Controller products controller
type Controller struct {
	View *render.Render
}

// Profile profile show page
func (ctrl Controller) Profile(w http.ResponseWriter, req *http.Request) {
	var (
		currentUser                     = utils.GetCurrentUser(req)
		tx                              = utils.GetDB(req)
		billingAddress, shippingAddress users.Address
	)

	// TODO refactor
	tx.Model(currentUser).Related(&currentUser.Addresses, "Addresses")
	tx.Model(currentUser).Related(&billingAddress, "DefaultBillingAddress")
	tx.Model(currentUser).Related(&shippingAddress, "DefaultShippingAddress")

	ctrl.View.Execute("profile", map[string]interface{}{
		"CurrentUser": currentUser, "DefaultBillingAddress": billingAddress, "DefaultShippingAddress": shippingAddress,
	}, req, w)
}

// Update update profile page
func (ctrl Controller) Update(w http.ResponseWriter, req *http.Request) {
	// FIXME
}

// AddCredit add credit
func (ctrl Controller) AddCredit(w http.ResponseWriter, req *http.Request) {
	// FIXME
}
