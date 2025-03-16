package errors

import (
	"net/http"

	"github.com/joomcode/errorx"
)

// list of error namespaces
var (
	databaseError           = errorx.NewNamespace("database error").ApplyModifiers(errorx.TypeModifierOmitStackTrace)
	invalidInput            = errorx.NewNamespace("validation error").ApplyModifiers(errorx.TypeModifierOmitStackTrace)
	resourceNotFound        = errorx.NewNamespace("not found").ApplyModifiers(errorx.TypeModifierOmitStackTrace)
	AccessDenied            = errorx.RegisterTrait("You are not authorized to perform the action")
	Ineligible              = errorx.RegisterTrait("You are not eligible to perform the action")
	serverError             = errorx.NewNamespace("server error")
	httpError               = errorx.NewNamespace("http error")
	badRequest              = errorx.NewNamespace("bad request error")
	Unauthenticated         = errorx.NewNamespace("user authentication failed")
	mpesaServiceError       = errorx.NewNamespace("mpesa service error")
	unauthorized            = errorx.NewNamespace("unauthorized").ApplyModifiers(errorx.TypeModifierOmitStackTrace)
	authoriztionClientError = errorx.NewNamespace("authorization client error")
	zemenBankServiceError   = errorx.NewNamespace("zemen bank service error")
	etSwitchServiceError    = errorx.NewNamespace("etswitch service error")
	redisServerError        = errorx.NewNamespace("redis service error")
)

var (
	ErrAuthClient               = errorx.NewType(authoriztionClientError, "authorization client error")
	ErrAcessError               = errorx.NewType(unauthorized, "Unauthorized", AccessDenied)
	ErrInvalidUserInput         = errorx.NewType(invalidInput, "invalid user input")
	ErrUnableToGet              = errorx.NewType(databaseError, "unable to get")
	ErrInternalServerError      = errorx.NewType(serverError, "internal server error")
	ErrUnableToUpdate           = errorx.NewType(databaseError, "unable to update")
	ErrUnableToCreate           = errorx.NewType(databaseError, "unable to create")
	ErrDBDelError               = errorx.NewType(databaseError, "could not delete record")
	ErrNoRecordFound            = errorx.NewType(resourceNotFound, "no record found")
	ErrHTTPRequestPrepareFailed = errorx.NewType(httpError, "couldn't prepare http request")
	ErrBadRequest               = errorx.NewType(badRequest, "bad request error")
	ErrSSOAuthenticationFailed  = errorx.NewType(Unauthenticated, "user authentication failed")
	ErrSSOError                 = errorx.NewType(serverError, "sso communication failed")
	ErrTelebirrInvalidRequest   = errorx.NewType(badRequest, "parsing service failed")
	ErrAirtimeInvalidRequest    = errorx.NewType(badRequest, "parsing service failed")
	ErrSMSInvalidRequest        = errorx.NewType(badRequest, "parsing service failed")
	ErrEnatBankError            = errorx.NewType(serverError, "error from enat bank")
	ErrEnatBankInvalidRequest   = errorx.NewType(badRequest, "parsing service failed")
	ErrTelebirrError            = errorx.NewType(serverError, "error from telebirr")
	ErrAirtimeError             = errorx.NewType(serverError, "error from airtime")

	ErrInvalidAccessToken = errorx.NewType(Unauthenticated, "invalid token").
				ApplyModifiers(errorx.TypeModifierOmitStackTrace)
	ErrAccountingClientError         = errorx.NewType(serverError, "accounting client error")
	ErrMpesaConn                     = errorx.NewType(mpesaServiceError, "mpesa service failed")
	ErrAuthError                     = errorx.NewType(unauthorized, "you are not authorized.")
	ErrHashError                     = errorx.NewType(serverError, "error generating hash")
	ErrWebHookSendError              = errorx.NewType(httpError, "can't send webhook")
	ErrWebHookRejectedByMerchant     = errorx.NewType(httpError, "merchant rejected the webhook request") // Fix typo
	ErrInvalidUserSubscriptionStatus = errorx.NewType(badRequest, "invalid user subscription status")
	ErrInsufficientBalance           = errorx.NewType(badRequest, "account insufficient balance error")
	ErrAmharaBankClientError         = errorx.NewType(serverError, "amhara bank client error")
	ErrZemenBankError                = errorx.NewType(zemenBankServiceError, "zemen bank service failed")
	ErrEtSwitchClient                = errorx.NewType(etSwitchServiceError, "etswitch client error")
	ErrRedisPubSubUnableToGetRoute   = errorx.NewType(redisServerError, "redis pubsub error")
	ErrReadError                     = errorx.NewType(redisServerError, "redis read message error")
	ErrRedisPubSubBadEvent           = errorx.NewType(redisServerError, "redis read message error")
	ErrUnsupportedPublicKeyFormat    = errorx.NewType(invalidInput, "unsupported public key format")
	ErrCBEMobileError                = errorx.NewType(invalidInput, "cbe mobile error")
)

var ErrorMap = map[*errorx.Type]int{
	ErrAuthClient:                    http.StatusBadRequest,
	ErrReadError:                     http.StatusBadRequest,
	ErrRedisPubSubBadEvent:           http.StatusBadRequest,
	ErrRedisPubSubUnableToGetRoute:   http.StatusInternalServerError,
	ErrAcessError:                    http.StatusForbidden,
	ErrInvalidUserInput:              http.StatusBadRequest,
	ErrMpesaConn:                     http.StatusBadRequest,
	ErrInvalidUserSubscriptionStatus: http.StatusBadRequest,
	ErrWebHookSendError:              http.StatusBadRequest,
	ErrWebHookRejectedByMerchant:     http.StatusBadRequest,
	ErrHashError:                     http.StatusInternalServerError,
	ErrInternalServerError:           http.StatusInternalServerError,
	ErrUnableToGet:                   http.StatusInternalServerError,
	ErrNoRecordFound:                 http.StatusNotFound,
	ErrDBDelError:                    http.StatusInternalServerError,
	ErrUnableToUpdate:                http.StatusInternalServerError,
	ErrUnableToCreate:                http.StatusInternalServerError,
	ErrHTTPRequestPrepareFailed:      http.StatusInternalServerError,
	ErrBadRequest:                    http.StatusBadRequest,
	ErrSSOAuthenticationFailed:       http.StatusUnauthorized,
	ErrSSOError:                      http.StatusInternalServerError,
	ErrInvalidAccessToken:            http.StatusUnauthorized,
	ErrTelebirrInvalidRequest:        http.StatusBadRequest,
	ErrSMSInvalidRequest:             http.StatusBadRequest,
	ErrAirtimeInvalidRequest:         http.StatusBadRequest,
	ErrEtSwitchClient:                http.StatusBadRequest,
	ErrEnatBankError:                 http.StatusInternalServerError,
	ErrAuthError:                     http.StatusUnauthorized,
	ErrTelebirrError:                 http.StatusInternalServerError,
	ErrAirtimeError:                  http.StatusInternalServerError,
	ErrEnatBankInvalidRequest:        http.StatusInternalServerError,
	ErrAccountingClientError:         http.StatusInternalServerError,
	ErrAmharaBankClientError:         http.StatusBadRequest,
	ErrInsufficientBalance:           http.StatusBadRequest,
	ErrZemenBankError:                http.StatusBadRequest,
	ErrUnsupportedPublicKeyFormat:    http.StatusBadRequest,
}
