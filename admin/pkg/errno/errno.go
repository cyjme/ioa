package errno

import (
	"fmt"
)

var (
	// common errors
	OK                  	= &Errno{Code: 0, Msg: "OK"}
	InternalServerError 	= &Errno{Code: 10001, Msg: "Internal server error"}
	ErrBind             	= &Errno{Code: 10002, Msg: "Error occurred while binding the request body to the struct."}
	BadRequestError     	= &Errno{Code: 10003, Msg: "Bad request Error."}

	// admin operate route
	AddRouteError       	= &Errno{Code: 11001, Msg: "Add route operation failed."}
	UpdateRouteError    	= &Errno{Code: 11002, Msg: "Update route operation failed."}
	DeleteRouteError    	= &Errno{Code: 11003, Msg: "Delete route operation failed."}
	ListRouteError      	= &Errno{Code: 11004, Msg: "List route operation failed."}
	NotExistsRouteError 	= &Errno{Code: 11005, Msg: "The route not exists, get failed."}
	
	// admin operate service
	AddServiceError     	= &Errno{Code: 12001, Msg: "Add service operation failed."}
	UpdateServiceError  	= &Errno{Code: 12002, Msg: "Update service operation failed."}
	DeleteServiceError  	= &Errno{Code: 12003, Msg: "Delete service operation failed."}
	ListServiceError    	= &Errno{Code: 12004, Msg: "List service operation failed."}
	NotExistsServiceError 	= &Errno{Code: 12005, Msg: "The service not exists, get failed."}

	// admin login logout
	AdminLoginError     	= &Errno{Code: 13001, Msg: "admin login operation failed."}
	AdminLogoutError    	= &Errno{Code: 13002, Msg: "admin logout operation failed."}
	AdminNoLoginError   	= &Errno{Code: 13003, Msg: "admin user not login."}
)

type Errno struct {
	Code int
	Msg  string
	Err  error
}

func (e *Errno) Error() string {
	return fmt.Sprintf("Err - code: %d, message: %s, error: %s", e.Code, e.Msg, e.Err)
}

func DecodeErr(err error) (int, string) {
	if err == nil {
		return OK.Code, OK.Msg
	}

	switch typed := err.(type) {
	case *Errno:
		return typed.Code, typed.Msg
	}

	return InternalServerError.Code, err.Error()
}
