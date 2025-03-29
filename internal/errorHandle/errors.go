package errorHandle

import "errors"

var (
	ItemNameExists     = errors.New("Item with this name alredy exists")
	ItemIdExists       = errors.New("Item with this ID alredy exists")
	ServerError        = errors.New("Error in server")
	ErrorFormatJson    = errors.New("Invalid format JSON")
	EmptyFile          = errors.New("Menu has not items")
	NotFoundID         = errors.New("Item with this ID does not exists ")
	ChangeID           = errors.New("You can't change ID item")
	PriceLessZero      = errors.New("Price is less than zero")
	QuantityLessZero   = errors.New("Quantity is less than zero")
	IdOrder            = errors.New("ID already exists")
	EmptyFileInventory = errors.New("Inventory has not items")
	OrderID            = errors.New("ID already exists")
	Ingred             = errors.New("Ingredients are missing")
	ChangeName         = errors.New("You can't update the name of customer")
	StatusExists       = errors.New("Status already close")
	DeleteOrder        = errors.New("You can't delete the order")
)

func CheckErrors(e error) int {
	if e == IdOrder || e == ItemNameExists || e == ItemIdExists || e == ErrorFormatJson || e == ChangeID || e == PriceLessZero || e == QuantityLessZero {
		return 400
	}
	if e == ServerError {
		return 500
	}
	if e == NotFoundID || e == EmptyFile || e == EmptyFileInventory {
		return 404
	}
	return 500
}
