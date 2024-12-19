package state

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// type Domain string
type AuthDomains struct {
	Merchant Domain
	System   Domain
	User     Domain
}

type Domain struct {
	Name string
	ID   string
}

func (a AuthDomains) Validate() error {
	return validation.Validate([]string{
		a.Merchant.ID, a.System.ID, a.Merchant.Name, a.System.Name, a.User.Name,
	}, validation.Each(validation.Required))
}

// OpaConfigs represents the configuration structure for an Open Policy Agent (OPA) service.
type OpaConfigs struct {
	// ServiceID is the unique identifier for the OPA service.
	ServiceID string
	// Password is the authentication password for accessing the OPA service.
	Password string
	// Host is the host address where the OPA service is running.
	Host string
	// Admin is the username or identifier for the administrative account of the OPA service.
	// This is only used for the mock
	Admin string
}

func (o OpaConfigs) Validate() error {
	return validation.ValidateStruct(&o,
		validation.Field(&o.ServiceID, validation.Required),
		validation.Field(&o.Password, validation.Required),
		validation.Field(&o.Host, validation.Required),
		validation.Field(&o.Admin, validation.Required),
	)
}

type AccountingParams struct {
	// Asset 1000
	TelebirrCashInAccount      string
	TelebirrCashOutAccount     string
	AirtimeCashInAccount       string
	AirtimeCashOutAccount      string
	EnatBankCashInAccount      string
	EnatBankCashOutAccount     string
	CBEBirrCashInAccount       string
	CBEMobileCashInAccount     string
	MpesaCashInAccount         string
	AmharaBankCashInAccount    string
	ZemenBankBankCashInAccount string
	EtSwitchCashInAccount      string

	// Liability 2000
	MerchantsAccountParentID string
	VatAccount               string

	// Equity 3000
	RetainedBalance string

	//Revenue 4000
	ServiceCharge       string
	ServiceChargeCustom string
	// receivable Account
	ReceivableAccount string

	// Expense 5000
	OperationalCost string
	Reimbursement   string
	Gift            string
}

func (a AccountingParams) Validate() error {
	return validation.Validate([]string{
		a.MpesaCashInAccount,
		a.TelebirrCashInAccount,
		a.TelebirrCashOutAccount,
		a.AirtimeCashInAccount,
		a.AirtimeCashOutAccount,
		a.EnatBankCashInAccount,
		a.EnatBankCashOutAccount,
		a.CBEBirrCashInAccount,
		a.CBEMobileCashInAccount,
		a.MerchantsAccountParentID,
		a.RetainedBalance,
		a.ServiceCharge,
		a.OperationalCost,
		a.Reimbursement,
		a.Gift,
		a.AmharaBankCashInAccount,
		a.ZemenBankBankCashInAccount,
		a.EtSwitchCashInAccount,
		a.VatAccount,
		a.ReceivableAccount,
		a.ServiceChargeCustom,
	}, validation.Each(validation.Required))
}

type RedisConfig struct {
	EventType string `json:"event_type,omitempty"`
	RedisURL  string `json:"redis_url,omitempty"`
}

func (r RedisConfig) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.EventType, validation.Required.Error("event_type is required")),
		validation.Field(&r.RedisURL, validation.Required.Error("reddis_url is required")),
	)
}

type UploadParams struct {
	FileTypes []FileType
}

type FileType struct {
	Name            string
	Types           []string
	MaxSize         int64
	AllowCustomName bool
}

func (f *FileType) SetValues(values map[string]any) {
	name, ok := values["name"].(string)
	if ok {
		f.Name = name
	}
	types, ok := values["types"].([]any)
	if ok {
		for _, v := range types {
			typeString, ok := v.(string)
			if ok {
				f.Types = append(f.Types, typeString)
			}
		}
	}

	size, ok := values["max_size"].(int)
	if ok {
		f.MaxSize = int64(size)
	}

	allowCustomName, ok := values["allow_custom_name"].(bool)
	if ok {
		f.AllowCustomName = allowCustomName
	}
}

type HTTPTransport struct {
	MaxIdleConnsPerHost int
	MaxIdleConns        int
	MaxConnsPerHost     int
	Timeout             time.Duration
}

func (a HTTPTransport) Validate() error {
	return validation.Validate([]int{
		a.MaxConnsPerHost, a.MaxIdleConns, a.MaxIdleConnsPerHost, int(a.Timeout),
	}, validation.Each(validation.Required))
}
