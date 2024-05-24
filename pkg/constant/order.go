package constant

type OrderStatusType string
type PaymentStatusType string

var (
	OrderStatusPaymentPending OrderStatusType = "payment_pending"
	OrderStatusPaid           OrderStatusType = "paid"

	PaymentStatusPending    PaymentStatusType = "pending"
	PaymentStatusSettlement PaymentStatusType = "settlement"
)
