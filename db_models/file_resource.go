package db_models

type BucketType string

const (
	TempBucket BucketType = "temp"
	ProdBucket BucketType = "production"
)

type ResourceType string

func (ResourceType) EnumList() []string {
	return []string{
		"product_resources",
		"transaction_resources",
		"profile_picture_resources",
		"warehouse_resources",
		"invoice_resources",
		"withdrawal_resources",
	}
}

const (
	ProductResource        ResourceType = "product_resources"
	TransactionResource    ResourceType = "transaction_resources"
	ProfilePictureResource ResourceType = "profile_picture_resources"
	WarehouseResource      ResourceType = "warehouse_resources"
	InvoiceResource        ResourceType = "invoice_resources"
	WithdrawalResource     ResourceType = "withdrawal_resources"
)
