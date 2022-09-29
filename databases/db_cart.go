package databases

import "errors"

var (
	ErrCantFindProduct    = errors.New("can't find the product")
	ErrCantDecodeProducts = errors.New("can't find the product")
	ErrUserIDIsNotValid   = errors.New("this user is invalid")
	ErrCantUpdateUser     = errors.New("can't update user")
	ErrCantRemoveItemCart = errors.New("cannot remove item cart")
	ErrCantGetItems       = errors.New("can't get item from the cart")
	ErrCantBuyCartItem    = errors.New("can't update the purchase")
)

func AddProductToChart() {

}

func RemoveCartItem() {

}

func BuyItemFromCart() {

}

func InstantBuyer() {

}
