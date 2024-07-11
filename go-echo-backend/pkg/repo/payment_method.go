package repo

import (
	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/logger"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/stripehelper"
	"github.com/samber/lo"
	"github.com/stripe/stripe-go/v74"
)

type PaymentMethodRepo struct {
	db     *db.DB
	logger *logger.Logger
}

func NewPaymentMethodRepo(db *db.DB) *PaymentMethodRepo {
	return &PaymentMethodRepo{
		db:     db,
		logger: logger.New("repo/PaymentMethod"),
	}
}

func (r *PaymentMethodRepo) AttachPaymentMethod(userID string, form models.UserPaymentMethodCreateForm) (*stripe.PaymentMethod, error) {
	var user models.User
	var err = r.db.Select("ID", "StripeCustomerID").First(&user, "id = ?", userID).Error
	if err != nil {
		return nil, err
	}

	if user.StripeCustomerID == "" {
		err = user.CreateStripeCustomer(r.db, &form.PaymentMethodID, false)
		if err != nil {
			return nil, err
		}
	}

	var pm *stripe.PaymentMethod

	if form.IsDefault {
		pm, err = stripehelper.GetInstance().SetAsDefaultPaymentMethod(user.StripeCustomerID, form.PaymentMethodID)

	} else {
		pm, err = stripehelper.GetInstance().AddPaymentMethod(user.StripeCustomerID, form.PaymentMethodID)
	}

	return pm, err
}

type DetactPaymentMethodParams struct {
	UserID string `json:"-"`

	PaymentMethodID string `json:"payment_method_id" param:"payment_method_id" query:"payment_method_id" validate:"required"`
}

func (r *PaymentMethodRepo) DetachPaymentMethod(params DetactPaymentMethodParams) (*stripe.PaymentMethod, error) {
	var user models.User
	var err = r.db.Select("ID", "StripeCustomerID").First(&user, "id = ?", params.UserID).Error
	if err != nil {
		return nil, err
	}
	existingPm, err := stripehelper.GetInstance().GetPaymentMethod(params.PaymentMethodID)
	if err != nil {
		return nil, err
	}

	var isDefault = existingPm.Customer != nil && existingPm.Customer.InvoiceSettings != nil && existingPm.Customer.InvoiceSettings.DefaultPaymentMethod != nil && existingPm.Customer.InvoiceSettings.DefaultPaymentMethod.ID == params.PaymentMethodID
	pm, err := stripehelper.GetInstance().DetachPaymentMethod(params.PaymentMethodID)
	if err != nil {
		return nil, err
	}

	if isDefault {
		var list = stripehelper.GetInstance().GetPaymentMethods(user.StripeCustomerID)
		if len(list) > 0 {
			_, err = stripehelper.GetInstance().SetAsDefaultPaymentMethod(user.StripeCustomerID, list[0].ID)
			if err != nil {
				return nil, err
			}
		}

	}
	return pm, err
}

type GetPaymentMethodsParams struct {
	IsDefault *bool  `json:"is_default" query:"is_default" param:"is_default"`
	UserID    string `json:"-"`
}

func (r *PaymentMethodRepo) GetPaymentMethods(params GetPaymentMethodsParams) ([]*models.UserPaymentMethod, error) {
	var user models.User
	var err = r.db.Select("ID", "StripeCustomerID").First(&user, "id = ?", params.UserID).Error
	if err != nil {
		return nil, err
	}

	if user.StripeCustomerID == "" {
		err = user.CreateStripeCustomer(r.db, nil, false)
		if err != nil {
			return nil, err
		}
	}

	var list = stripehelper.GetInstance().GetPaymentMethods(user.StripeCustomerID)

	if params.IsDefault != nil && *params.IsDefault {
		list = lo.Filter(list, func(item *stripehelper.PaymentMethodWithDefault, index int) bool {
			return item.IsDefault
		})
	}

	var result = lo.Map(list, func(item *stripehelper.PaymentMethodWithDefault, index int) *models.UserPaymentMethod {
		var card = &models.UserPaymentMethod{
			ID:        item.ID,
			Type:      "card",
			CreatedAt: item.Created,
		}

		if item.Customer != nil {
			card.Name = item.Customer.Name
		}

		if item.Card != nil {
			card.ExpMonth = item.Card.ExpMonth
			card.ExpYear = item.Card.ExpYear
			card.Last4 = item.Card.Last4
			card.Brand = string(item.Card.Brand)
			card.Type = "card"
		}

		if item.USBankAccount != nil {
			card.Name = item.USBankAccount.BankName
			card.Last4 = item.USBankAccount.Last4
			card.Type = "us_bank_account"
		}

		return card
	})

	return result, nil

}

type MarkDefaultPaymentMethodParams struct {
	UserID string `json:"-"`

	PaymentMethodID string `json:"payment_method_id" param:"payment_method_id" query:"payment_method_id" validate:"required"`
}

func (r *PaymentMethodRepo) MarkDefaultPaymentMethod(params MarkDefaultPaymentMethodParams) (*stripe.PaymentMethod, error) {
	var user models.User
	var err = r.db.Select("ID", "StripeCustomerID").First(&user, "id = ?", params.UserID).Error
	if err != nil {
		return nil, err
	}

	pm, err := stripehelper.GetInstance().SetAsDefaultPaymentMethod(user.StripeCustomerID, params.PaymentMethodID)

	return pm, err
}

type GetPaymentMethodParams struct {
	UserID string `json:"-"`

	PaymentMethodID string `json:"payment_method_id" param:"payment_method_id" query:"payment_method_id" validate:"required"`
}

func (r *PaymentMethodRepo) GetPaymentMethod(params GetPaymentMethodParams) (*stripe.PaymentMethod, error) {
	var user models.User
	var err = r.db.Select("ID", "StripeCustomerID").First(&user, "id = ?", params.UserID).Error
	if err != nil {
		return nil, err
	}

	pm, err := stripehelper.GetInstance().GetPaymentMethod(params.PaymentMethodID)
	if err != nil {
		return nil, err
	}

	return pm, err
}
