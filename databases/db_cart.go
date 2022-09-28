package databases

import "errors"

var (
	ErrCantFindProduct    = errors.New("")
	ErrCantDecodeProducts = errors.New("")
	ErrUserIDIsNotValid   = errors.New("")
	ErrCantUpdateUser     = errors.New("")
	ErrCantRemoveItemCart = errors.New("")
	ErrCantGetItems       = errors.New("")
	ErrCantBuyCartItem    = errors.New("")
)
