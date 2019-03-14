package cart

import (
	"encoding/json"
	"flamingo.me/flamingo-commerce/v3/price/domain"
	"time"

	"flamingo.me/flamingo/v3/framework/web"
	"github.com/pkg/errors"
)

type (
	// Provider should be used to create the cart Value objects
	Provider func() *Cart

	// Cart Value Object (immutable data - because the cartservice is responsible to return a cart).
	Cart struct {
		//ID is the main identifier of the cart
		ID string
		//EntityID is a second identifier that may be used by some backends
		EntityID string

		// ReservedOrderID is an ID already known by the Cart of the future order ID
		ReservedOrderID string

		//CartTotals - the cart totals (contain summary costs and discounts etc)
		CartTotals Totals

		//BillingAdress - the main billing address (relevant for all payments/invoices)
		BillingAdress Address

		//Purchaser - additional infos for the legal contact person in this order
		Purchaser Person

		//Deliveries - list of desired Deliverys (or Shippments) involved in this cart
		Deliveries []Delivery

		//AdditionalData   can be used for Custom attributes
		AdditionalData AdditionalData

		//BelongsToAuthenticatedUser - false = Guest Cart true = cart from the authenticated user
		BelongsToAuthenticatedUser bool
		AuthenticatedUserID        string

		AppliedCouponCodes []CouponCode
	}

	// Teaser - represents some teaser infos for cart
	Teaser struct {
		ProductCount  int
		ItemCount     int
		DeliveryCodes []string
	}

	// CouponCode value object
	CouponCode struct {
		Code string
	}

	// Person value object
	Person struct {
		Address         *Address
		PersonalDetails PersonalDetails
		//ExistingCustomerData if the current purchaser is an existing customer - this contains infos about existing customer
		ExistingCustomerData *ExistingCustomerData
	}

	// ExistingCustomerData value object
	ExistingCustomerData struct {
		//ID of the customer
		ID string
	}

	// PersonalDetails value object
	PersonalDetails struct {
		DateOfBirth     string
		PassportCountry string
		PassportNumber  string
		Nationality     string
	}

	// Delivery - represents the DeliveryInfo and the assigned Items
	Delivery struct {
		//DeliveryInfo - The details for this delivery - normaly completed during checkout
		DeliveryInfo DeliveryInfo
		//Cartitems - list of cartitems
		Cartitems      []Item
		//DeliveryTotals - Totals with the intent to use them to display the customer summary costs for this delivery
		DeliveryTotals DeliveryTotals
		//ShippingItem	- The Shipping Costs that may be involved in this delivery
		ShippingItem   ShippingItem
	}

	// DeliveryInfo - represents the Delivery
	DeliveryInfo struct {
		// Code - is a project specific idendifier for the Delivery - you need it for the AddToCart Request for example
		// The code can follow the convention in the Readme: Type_Method_LocationType_LocationCode
		Code string
		//Type - The Type of the Delivery - e.g. delivery or pickup - this might trigger different workflows
		Workflow string
		//Method - The shippingmethod something that is project specific and that can mean different delivery qualities with different deliverycosts
		Method string
		//Carrier - Optional the name of the Carrier that should be responsible for executing the delivery
		Carrier                 string
		//DeliveryLocation The target Location for the delivery
		DeliveryLocation        DeliveryLocation
		//DesiredTime - Optional - the desired time of the delivery
		DesiredTime             time.Time
		//AdditionalData  - Possibility for key value based information on the delivery - can be used flexible by each project
		AdditionalData          map[string]string
		//AdditionalDeliveryInfos - similar to AdditionalData this can be used to store "any" other object on a delivery encoded as json.RawMessage
		AdditionalDeliveryInfos map[string]json.RawMessage
	}

	//AdditionalDeliverInfo is an interface that allows to store "any" additional objects on the cart
	// see DeliveryInfoUpdateCommand
	AdditionalDeliverInfo interface {
		Marshal() (json.RawMessage, error)
		Unmarshal(json.RawMessage) error
	}

	// DeliveryLocation value object
	DeliveryLocation struct {
		Type string
		//Address - only set for type adress
		Address *Address
		//Code - optional idendifier of this location/destination - is used in special destination Types
		Code string
	}

	// Totals value object
	Totals struct {
		Totalitems        []Totalitem
		TotalShippingItem ShippingItem
		//Final sum that need to be payed: GrandTotal = SubTotal + TaxAmount - DiscountAmount + SOME of Totalitems = (Sum of Items RowTotalWithDiscountInclTax) + SOME of Totalitems
		GrandTotal domain.Price
		//SubTotal = SUM of Item RowTotal
		SubTotal domain.Price
		//SubTotalInclTax = SUM of Item RowTotalInclTax
		SubTotalInclTax domain.Price
		//SubTotalWithDiscounts = SubTotal - Sum of Item ItemRelatedDiscountAmount
		SubTotalWithDiscounts domain.Price
		//SubTotalWithDiscountsAndTax= Sum of RowTotalWithItemRelatedDiscountInclTax
		SubTotalWithDiscountsAndTax domain.Price

		//TotalDiscountAmount = SUM of Item TotalDiscountAmount
		TotalDiscountAmount domain.Price
		//TotalNonItemRelatedDiscountAmount= SUM of Item NonItemRelatedDiscountAmount
		TotalNonItemRelatedDiscountAmount domain.Price
		//TaxAmount = Sum of Item TaxAmount
		TaxAmount domain.Price
	}

	// DeliveryTotals value object
	DeliveryTotals struct {
		//SubTotal = SUM of Item RowTotal
		SubTotal domain.Price
		//SubTotalInclTax = SUM of Item RowTotalInclTax
		SubTotalInclTax domain.Price
		//SubTotalWithDiscounts = SubTotal - Sum of Item ItemRelatedDiscountAmount
		SubTotalWithDiscounts domain.Price
		//SubTotalWithDiscountsAndTax= Sum of RowTotalWithItemRelatedDiscountInclTax
		SubTotalWithDiscountsAndTax domain.Price

		//TotalDiscountAmount = SUM of Item TotalDiscountAmount
		TotalDiscountAmount domain.Price
		//TotalNonItemRelatedDiscountAmount= SUM of Item NonItemRelatedDiscountAmount
		TotalNonItemRelatedDiscountAmount domain.Price
	}

	// Item for Cart
	Item struct {
		//ID of the item - need to be unique under a delivery
		ID string
		//
		UniqueID        string
		MarketplaceCode string
		//VariantMarketPlaceCode is used for Configurable products
		VariantMarketPlaceCode string
		ProductName            string

		// Source Id of where the items should be initial picked - This is set by the SourcingLogic
		SourceID string

		Qty int

		AdditionalData map[string]string

		//brutto for single item
		SinglePrice domain.Price
		//netto for single item
		SinglePriceInclTax domain.Price
		//RowTotal = SinglePrice * Qty
		RowTotal domain.Price
		//TaxAmount=Qty * (SinglePriceInclTax-SinglePrice)
		TaxAmount domain.Price
		//RowTotalInclTax= RowTotal + TaxAmount
		RowTotalInclTax domain.Price
		//AppliedDiscounts contains the details about the discounts applied to this item - they can be "itemrelated" or not
		AppliedDiscounts []ItemDiscount
		// TotalDiscountAmount = Sum of AppliedDiscounts = ItemRelatedDiscountAmount +NonItemRelatedDiscountAmount
		TotalDiscountAmount domain.Price
		// ItemRelatedDiscountAmount = Sum of AppliedDiscounts where IsItemRelated = True
		ItemRelatedDiscountAmount domain.Price
		//NonItemRelatedDiscountAmount = Sum of AppliedDiscounts where IsItemRelated = false
		NonItemRelatedDiscountAmount domain.Price
		//RowTotalWithItemRelatedDiscountInclTax=RowTotal-ItemRelatedDiscountAmount
		RowTotalWithItemRelatedDiscount domain.Price
		//RowTotalWithItemRelatedDiscountInclTax=RowTotalInclTax-ItemRelatedDiscountAmount
		RowTotalWithItemRelatedDiscountInclTax domain.Price
		//This is the price the customer finaly need to pay for this item:  RowTotalWithDiscountInclTax = RowTotalInclTax-TotalDiscountAmount
		RowTotalWithDiscountInclTax domain.Price
	}

	// ItemCartReference - value object that can be used to reference a Item in a Cart
	//@todo - Use in ServicePort methods...
	ItemCartReference struct {
		ItemID       string
		DeliveryCode string
	}

	// ItemDiscount value object
	ItemDiscount struct {
		Code  string
		Title string
		Price domain.Price
		//IsItemRelated is a flag indicating if the discount should be displayed in the item or if it the result of a cart discount
		IsItemRelated bool
	}

	// Totalitem for totals
	Totalitem struct {
		Code  string
		Title string
		Price domain.Price
		Type  string
	}

	// ShippingItem value object
	ShippingItem struct {
		Title string
		Price domain.Price

		TaxAmount      domain.Price
		DiscountAmount domain.Price

	}

	// InvalidateCartEvent value object
	InvalidateCartEvent struct {
		Session *web.Session
	}

	// AdditionalData defines the supplementary cart data
	AdditionalData struct {
		CustomAttributes map[string]string
		SelectedPayment  SelectedPayment
	}

	// SelectedPayment value object
	SelectedPayment struct {
		Provider string
		Method   string
	}

	// PlacedOrderInfos represents a slice of PlacedOrderInfo
	PlacedOrderInfos []PlacedOrderInfo

	// PlacedOrderInfo defines the additional info struct for placed orders
	PlacedOrderInfo struct {
		OrderNumber  string
		DeliveryCode string
	}
)

var (
	// ErrAdditionalInfosNotFound is returned if the additional infos are not set
	ErrAdditionalInfosNotFound = errors.New("additional infos not found")
)

// Key constants
const (
	DeliveryWorkflowPickup      = "pickup"
	DeliveryWorkflowDelivery    = "delivery"
	DeliveryWorkflowUnspecified = "unspecified"

	DeliverylocationTypeUnspecified = "unspecified"
	DeliverylocationTypeCollectionpoint = "collection-point"
	DeliverylocationTypeStore           = "store"
	DeliverylocationTypeAddress         = "address"
	DeliverylocationTypeFreightstation  = "freight-station"

	TotalsTypeDiscount      = "totals_type_discount"
	TotalsTypeVoucher       = "totals_type_voucher"
	TotalsTypeTax           = "totals_type_tax"
	TotalsTypeLoyaltypoints = "totals_loyaltypoints"
	TotalsTypeShipping      = "totals_type_shipping"
)

// GetMainShippingEMail returns the main shipping address email, empty string if not available
func (Cart Cart) GetMainShippingEMail() string {
	for _, deliveries := range Cart.Deliveries {
		if deliveries.DeliveryInfo.DeliveryLocation.Address != nil {
			if deliveries.DeliveryInfo.DeliveryLocation.Address.Email != "" {
				return deliveries.DeliveryInfo.DeliveryLocation.Address.Email
			}
		}
	}

	return ""
}

// GetDeliveryByCode gets a delivery by code
func (Cart Cart) GetDeliveryByCode(deliveryCode string) (*Delivery, bool) {
	for _, delivery := range Cart.Deliveries {
		if delivery.DeliveryInfo.Code == deliveryCode {
			return &delivery, true
		}
	}

	return nil, false
}

// HasDeliveryForCode checks if a delivery with the given code exists in the cart
func (Cart Cart) HasDeliveryForCode(deliveryCode string) bool {
	_, found := Cart.GetDeliveryByCode(deliveryCode)

	return found == true
}

// GetDeliveryCodes returns a slice of all delivery codes in cart that have at least one cart item
func (Cart Cart) GetDeliveryCodes() []string {
	var deliveryCodes []string

	for _, delivery := range Cart.Deliveries {
		if len(delivery.Cartitems) > 0 {
			deliveryCodes = append(deliveryCodes, delivery.DeliveryInfo.Code)
		}
	}

	return deliveryCodes
}

// GetByItemID gets an item by its id
func (Cart Cart) GetByItemID(itemID string, deliveryCode string) (*Item, error) {
	delivery, found := Cart.GetDeliveryByCode(deliveryCode)
	if found != true {
		return nil, errors.Errorf("Delivery for code %v not found", deliveryCode)
	}
	for _, currentItem := range delivery.Cartitems {
		if currentItem.ID == itemID {
			return &currentItem, nil
		}
	}

	return nil, errors.Errorf("itemId %v in cart not existing", itemID)
}

func inStruct(value string, list []string) bool {
	for _, item := range list {
		if item == value {
			return true
		}
	}

	return false
}

// ItemCount - returns amount of Cartitems
func (Cart Cart) ItemCount() int {
	count := 0
	for _, delivery := range Cart.Deliveries {
		for _, item := range delivery.Cartitems {
			count += item.Qty
		}
	}

	return count
}

// ProductCount - returns amount of different products
func (Cart Cart) ProductCount() int {
	count := 0
	for _, delivery := range Cart.Deliveries {
		count += len(delivery.Cartitems)
	}

	return count
}

// GetItemCartReferences returns a slice of all ItemCartReferences
func (Cart Cart) GetItemCartReferences() []ItemCartReference {
	var ids []ItemCartReference
	for _, delivery := range Cart.Deliveries {
		for _, item := range delivery.Cartitems {
			ids = append(ids, ItemCartReference{
				ItemID:       item.ID,
				DeliveryCode: delivery.DeliveryInfo.Code,
			})
		}
	}

	return ids
}

// GetVoucherSavings returns the savings of all vouchers
func (Cart Cart) GetVoucherSavings() domain.Price {
	price := domain.Price{}
	for _, item := range Cart.CartTotals.Totalitems {
		if item.Type == TotalsTypeVoucher {
			price, err := price.Add(item.Price)
			if err != nil {
				return price
			}
		}
	}
	if price.IsNegative() {
		return domain.Price{}
	}
	return price
}

// GetSavings retuns the total of all savings
func (Cart Cart) GetSavings() domain.Price {
	price := domain.Price{}
	for _, item := range Cart.CartTotals.Totalitems {
		if item.Type == TotalsTypeDiscount {
			price, err := price.Add(item.Price)
			if err != nil {
				return price
			}
		}
	}

	if price.IsNegative() {
		return domain.Price{}
	}
	return price
}

// HasAppliedCouponCode checks if a coupon code is applied to the cart
func (Cart Cart) HasAppliedCouponCode() bool {
	return len(Cart.AppliedCouponCodes) > 0
}

// GetCartTeaser returns the teaser
func (Cart Cart) GetCartTeaser() *Teaser {
	return &Teaser{
		DeliveryCodes: Cart.GetDeliveryCodes(),
		ItemCount:     Cart.ItemCount(),
		ProductCount:  Cart.ProductCount(),
	}
}

// GetTotalItemsByType gets a slice of all Totalitems by typeCode
func (ct Totals) GetTotalItemsByType(typeCode string) []Totalitem {
	var totalitems []Totalitem
	for _, item := range ct.Totalitems {
		if item.Type == typeCode {
			totalitems = append(totalitems, item)
		}
	}

	return totalitems
}

// GetSavingsByItem gets the savings by item
func (item Item) GetSavingsByItem() domain.Price {
	price := domain.Price{}
	for _, discount := range item.AppliedDiscounts {
		price, err := price.Add(discount.Price)
		if err != nil {
			return price
		}
	}

	if price.IsNegative() {
		return domain.Price{}
	}
	return price
}

// GetOrderNumberForDeliveryCode returns the order number for a delivery code
func (poi PlacedOrderInfos) GetOrderNumberForDeliveryCode(deliveryCode string) string {
	for _, v := range poi {
		if v.DeliveryCode == deliveryCode {
			return v.OrderNumber
		}
	}
	return ""
}

//LoadAdditionalInfo - returns the additional Data
func (d *DeliveryInfo) LoadAdditionalInfo(key string, info AdditionalDeliverInfo) error {
	if d.AdditionalDeliveryInfos == nil {
		return ErrAdditionalInfosNotFound
	}
	if val, ok := d.AdditionalDeliveryInfos[key]; ok {
		return info.Unmarshal(val)
	}
	return ErrAdditionalInfosNotFound
}
