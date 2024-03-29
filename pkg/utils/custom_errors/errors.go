package customerrors

import "errors"

var (
	ErrNotFound                     = errors.New("not found")
	ErrEmailAlredyExist             = errors.New("email is used")
	ErrBadRequestBody               = errors.New("bad request body")
	ErrInvalidId                    = errors.New("invalid id")
	ErrNotAllowedDeleteDefaultAdmin = errors.New("default admin user not allowed to delete")
	ErrFailedDeleteUser             = errors.New("failed delete user")
	ErrInvalidPassword              = errors.New("invalid password")
	ErrPermission                   = errors.New("not have permission to access")
	ErrInvalidParam                 = errors.New("invalid param")
	ErrDuplicateData                = errors.New("data duplicate")
	ErrTimeLocation                 = errors.New("time location error")
	ErrQtyOrder                     = errors.New("order qty exceeds stock or less that 1")
	ErrOrderCode                    = errors.New("invalid order code")
	ErrCheckpointNotCovered         = errors.New("not found checkpoint in your location")
	ErrUpdateStatusOrder            = errors.New("cant update status order")
	ErrGenerateQR                   = errors.New("error when generate qrcode")
	ErrCodeUsed                     = errors.New("code is used")
	ErrWrongCheckpoint              = errors.New("cant pick up this order at this checkpoint")
)
