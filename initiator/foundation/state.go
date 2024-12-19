package foundation

import (
	"Eccomerce-website/internal/constant"
	"Eccomerce-website/internal/constant/errors"
	"Eccomerce-website/internal/constant/state"
	"Eccomerce-website/platform/asset"
	"context"
	"time"

	"github.com/eangle9/log"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type State struct {
	AuthDomains      state.AuthDomains
	HTTPConfig       state.HTTPTransport
	TokenConfig      state.TokenKey
	AccountingParams state.AccountingParams
	OpaConfig        state.OpaConfigs
	UploadParams     state.UploadParams
	RedisConfig      state.RedisConfig
}

func InitState(logger log.Logger) State {
	authDomains := state.AuthDomains{
		Merchant: state.Domain{
			ID:   viper.GetString("service.authorization.domain.merchant"),
			Name: string(constant.Merchant),
		},
		System: state.Domain{
			ID:   viper.GetString("service.authorization.domain.system"),
			Name: string(constant.System),
		},
		User: state.Domain{
			Name: string(constant.User),
			ID:   viper.GetString("service.authorization.domain.user"),
		},
	}
	if err := authDomains.Validate(); err != nil {
		logger.Fatal(context.Background(), "invalid domain input",
			zap.Error(err),
		)
	}

	httpconfig := state.HTTPTransport{
		MaxIdleConnsPerHost: viper.GetInt("http.max_idle_conns_per_host"),
		MaxIdleConns:        viper.GetInt("http.max_idle_conns"),
		MaxConnsPerHost:     viper.GetInt("http.max_conns_per_host"),
		Timeout:             viper.GetDuration("http.timeout"),
	}
	if err := httpconfig.Validate(); err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid http configuration")
		logger.Fatal(context.Background(), "all http fields are required", zap.Error(err))
	}

	tokenConfig := state.TokenKey{
		SymmetricKey: viper.GetString("security_credential.symmetric_key"),
		Issuer:       viper.GetString("security_credential.issuer"),
		Footer:       viper.GetString("security_credential.footer"),
		KeyLength:    viper.GetInt("security_credential.key_length"),
		ExpireAt: time.Minute * time.Duration(
			viper.GetInt64("security_credential.expire_at")),
	}
	if err := tokenConfig.Validate(); err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid token config configuration")
		logger.Fatal(context.Background(), "all tokenKey fields are required", zap.Error(err))
	}

	coa := state.AccountingParams{
		VatAccount:             viper.GetString("service.accounting.chart_of_account.vat_account"),
		TelebirrCashInAccount:  viper.GetString("service.accounting.chart_of_account.bank.telebirr_account.in"),
		TelebirrCashOutAccount: viper.GetString("service.accounting.chart_of_account.bank.telebirr_account.out"),
		AirtimeCashInAccount:   viper.GetString("service.accounting.chart_of_account.bank.airtime_account.in"),
		AirtimeCashOutAccount:  viper.GetString("service.accounting.chart_of_account.bank.airtime_account.out"),
		EnatBankCashInAccount:  viper.GetString("service.accounting.chart_of_account.bank.enat_bank_account.in"),
		ReceivableAccount: viper.GetString(
			"service.accounting.chart_of_account.receivable_account",
		),
		EnatBankCashOutAccount: viper.GetString(
			"service.accounting.chart_of_account.bank.enat_bank_account.out",
		),
		CBEBirrCashInAccount:   viper.GetString("service.accounting.chart_of_account.bank.cbe_account.cbe_birr_cash_in"),
		CBEMobileCashInAccount: viper.GetString("service.accounting.chart_of_account.bank.cbe_account.cbe_mobile_cash_in"),
		MpesaCashInAccount:     viper.GetString("service.accounting.chart_of_account.bank.mpesa_account.in"),
		MerchantsAccountParentID: viper.GetString(
			"service.accounting.chart_of_account.merchant_account",
		),
		RetainedBalance: viper.GetString(
			"service.accounting.chart_of_account.retained_earnings_account",
		),
		OperationalCost: viper.GetString(
			"service.accounting.chart_of_account.operational_cost_account",
		),
		Gift: viper.GetString(
			"service.accounting.chart_of_account.gift_account",
		),
		Reimbursement: viper.GetString(
			"service.accounting.chart_of_account.reimbursement_account",
		),

		ServiceCharge: viper.GetString(
			"service.accounting.chart_of_account.service_charge_account",
		),
		AmharaBankCashInAccount: viper.GetString(
			"service.accounting.chart_of_account.bank.amhara_bank_account.in",
		),
		ZemenBankBankCashInAccount: viper.GetString(
			"service.accounting.chart_of_account.bank.zemen_bank_account.in",
		),
		EtSwitchCashInAccount: viper.GetString(
			"service.accounting.chart_of_account.bank.etswitch_account.in",
		),

		ServiceChargeCustom: viper.GetString(
			"service.accounting.chart_of_account.service_charge_account_custom",
		),
	}
	if err := coa.Validate(); err != nil {
		logger.Fatal(context.Background(), "missing some account values",
			zap.Error(err),
		)
	}

	opa := state.OpaConfigs{
		ServiceID: viper.GetString("service.authorization.service_id"),
		Password:  viper.GetString("service.authorization.password"),
		Host:      viper.GetString("service.authorization.host"),
		Admin:     viper.GetString("service.authorization.admin"),
	}
	if err := opa.Validate(); err != nil {
		logger.Fatal(context.Background(), "missing some opa configs",
			zap.Error(err),
		)
	}

	assets := GetMapSlice("assets")
	fileTypes := make([]state.FileType, 0, len(assets))

	for _, v := range assets {
		var fileType state.FileType

		fileType.SetValues(v)
		fileTypes = append(fileTypes, fileType)
	}
	return State{
		AuthDomains:      authDomains,
		HTTPConfig:       httpconfig,
		TokenConfig:      tokenConfig,
		AccountingParams: coa,
		OpaConfig:        opa,
		UploadParams: asset.SetParams(logger, state.UploadParams{
			FileTypes: fileTypes,
		}),
	}
}

func GetMapSlice(path string) []map[string]any {
	value := viper.Get(path)
	mapInterfaceSlice, ok := value.([]any)
	if !ok {
		return nil
	}

	var mapStringAny []map[string]any
	for _, v := range mapInterfaceSlice {
		v, ok := v.(map[string]any)
		if ok {
			mapStringAny = append(mapStringAny, v)
		}
	}

	return mapStringAny
}
