// Code generated by mockery v2.10.0. DO NOT EDIT.

package mocks

import (
	context "context"

	cart "flamingo.me/flamingo-commerce/v3/cart/domain/cart"

	mock "github.com/stretchr/testify/mock"
)

// GiftCardAndVoucherBehaviour is an autogenerated mock type for the GiftCardAndVoucherBehaviour type
type GiftCardAndVoucherBehaviour struct {
	mock.Mock
}

// ApplyAny provides a mock function with given fields: ctx, _a1, anyCode
func (_m *GiftCardAndVoucherBehaviour) ApplyAny(ctx context.Context, _a1 *cart.Cart, anyCode string) (*cart.Cart, cart.DeferEvents, error) {
	ret := _m.Called(ctx, _a1, anyCode)

	var r0 *cart.Cart
	if rf, ok := ret.Get(0).(func(context.Context, *cart.Cart, string) *cart.Cart); ok {
		r0 = rf(ctx, _a1, anyCode)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*cart.Cart)
		}
	}

	var r1 cart.DeferEvents
	if rf, ok := ret.Get(1).(func(context.Context, *cart.Cart, string) cart.DeferEvents); ok {
		r1 = rf(ctx, _a1, anyCode)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(cart.DeferEvents)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context, *cart.Cart, string) error); ok {
		r2 = rf(ctx, _a1, anyCode)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}
