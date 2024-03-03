// Code generated by mockery v2.41.0. DO NOT EDIT.

package jetstream

import (
	context "context"

	jetstream "github.com/nats-io/nats.go/jetstream"
	mock "github.com/stretchr/testify/mock"

	nats "github.com/nats-io/nats.go"

	time "time"
)

// Msg is an autogenerated mock type for the Msg type
type Msg struct {
	mock.Mock
}

type Msg_Expecter struct {
	mock *mock.Mock
}

func (_m *Msg) EXPECT() *Msg_Expecter {
	return &Msg_Expecter{mock: &_m.Mock}
}

// Ack provides a mock function with given fields:
func (_m *Msg) Ack() error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Ack")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Msg_Ack_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Ack'
type Msg_Ack_Call struct {
	*mock.Call
}

// Ack is a helper method to define mock.On call
func (_e *Msg_Expecter) Ack() *Msg_Ack_Call {
	return &Msg_Ack_Call{Call: _e.mock.On("Ack")}
}

func (_c *Msg_Ack_Call) Run(run func()) *Msg_Ack_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Msg_Ack_Call) Return(_a0 error) *Msg_Ack_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Msg_Ack_Call) RunAndReturn(run func() error) *Msg_Ack_Call {
	_c.Call.Return(run)
	return _c
}

// Data provides a mock function with given fields:
func (_m *Msg) Data() []byte {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Data")
	}

	var r0 []byte
	if rf, ok := ret.Get(0).(func() []byte); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	return r0
}

// Msg_Data_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Data'
type Msg_Data_Call struct {
	*mock.Call
}

// Data is a helper method to define mock.On call
func (_e *Msg_Expecter) Data() *Msg_Data_Call {
	return &Msg_Data_Call{Call: _e.mock.On("Data")}
}

func (_c *Msg_Data_Call) Run(run func()) *Msg_Data_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Msg_Data_Call) Return(_a0 []byte) *Msg_Data_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Msg_Data_Call) RunAndReturn(run func() []byte) *Msg_Data_Call {
	_c.Call.Return(run)
	return _c
}

// DoubleAck provides a mock function with given fields: _a0
func (_m *Msg) DoubleAck(_a0 context.Context) error {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for DoubleAck")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Msg_DoubleAck_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DoubleAck'
type Msg_DoubleAck_Call struct {
	*mock.Call
}

// DoubleAck is a helper method to define mock.On call
//   - _a0 context.Context
func (_e *Msg_Expecter) DoubleAck(_a0 interface{}) *Msg_DoubleAck_Call {
	return &Msg_DoubleAck_Call{Call: _e.mock.On("DoubleAck", _a0)}
}

func (_c *Msg_DoubleAck_Call) Run(run func(_a0 context.Context)) *Msg_DoubleAck_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *Msg_DoubleAck_Call) Return(_a0 error) *Msg_DoubleAck_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Msg_DoubleAck_Call) RunAndReturn(run func(context.Context) error) *Msg_DoubleAck_Call {
	_c.Call.Return(run)
	return _c
}

// Headers provides a mock function with given fields:
func (_m *Msg) Headers() nats.Header {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Headers")
	}

	var r0 nats.Header
	if rf, ok := ret.Get(0).(func() nats.Header); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(nats.Header)
		}
	}

	return r0
}

// Msg_Headers_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Headers'
type Msg_Headers_Call struct {
	*mock.Call
}

// Headers is a helper method to define mock.On call
func (_e *Msg_Expecter) Headers() *Msg_Headers_Call {
	return &Msg_Headers_Call{Call: _e.mock.On("Headers")}
}

func (_c *Msg_Headers_Call) Run(run func()) *Msg_Headers_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Msg_Headers_Call) Return(_a0 nats.Header) *Msg_Headers_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Msg_Headers_Call) RunAndReturn(run func() nats.Header) *Msg_Headers_Call {
	_c.Call.Return(run)
	return _c
}

// InProgress provides a mock function with given fields:
func (_m *Msg) InProgress() error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for InProgress")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Msg_InProgress_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'InProgress'
type Msg_InProgress_Call struct {
	*mock.Call
}

// InProgress is a helper method to define mock.On call
func (_e *Msg_Expecter) InProgress() *Msg_InProgress_Call {
	return &Msg_InProgress_Call{Call: _e.mock.On("InProgress")}
}

func (_c *Msg_InProgress_Call) Run(run func()) *Msg_InProgress_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Msg_InProgress_Call) Return(_a0 error) *Msg_InProgress_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Msg_InProgress_Call) RunAndReturn(run func() error) *Msg_InProgress_Call {
	_c.Call.Return(run)
	return _c
}

// Metadata provides a mock function with given fields:
func (_m *Msg) Metadata() (*jetstream.MsgMetadata, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Metadata")
	}

	var r0 *jetstream.MsgMetadata
	var r1 error
	if rf, ok := ret.Get(0).(func() (*jetstream.MsgMetadata, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() *jetstream.MsgMetadata); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*jetstream.MsgMetadata)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Msg_Metadata_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Metadata'
type Msg_Metadata_Call struct {
	*mock.Call
}

// Metadata is a helper method to define mock.On call
func (_e *Msg_Expecter) Metadata() *Msg_Metadata_Call {
	return &Msg_Metadata_Call{Call: _e.mock.On("Metadata")}
}

func (_c *Msg_Metadata_Call) Run(run func()) *Msg_Metadata_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Msg_Metadata_Call) Return(_a0 *jetstream.MsgMetadata, _a1 error) *Msg_Metadata_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Msg_Metadata_Call) RunAndReturn(run func() (*jetstream.MsgMetadata, error)) *Msg_Metadata_Call {
	_c.Call.Return(run)
	return _c
}

// Nak provides a mock function with given fields:
func (_m *Msg) Nak() error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Nak")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Msg_Nak_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Nak'
type Msg_Nak_Call struct {
	*mock.Call
}

// Nak is a helper method to define mock.On call
func (_e *Msg_Expecter) Nak() *Msg_Nak_Call {
	return &Msg_Nak_Call{Call: _e.mock.On("Nak")}
}

func (_c *Msg_Nak_Call) Run(run func()) *Msg_Nak_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Msg_Nak_Call) Return(_a0 error) *Msg_Nak_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Msg_Nak_Call) RunAndReturn(run func() error) *Msg_Nak_Call {
	_c.Call.Return(run)
	return _c
}

// NakWithDelay provides a mock function with given fields: delay
func (_m *Msg) NakWithDelay(delay time.Duration) error {
	ret := _m.Called(delay)

	if len(ret) == 0 {
		panic("no return value specified for NakWithDelay")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(time.Duration) error); ok {
		r0 = rf(delay)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Msg_NakWithDelay_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'NakWithDelay'
type Msg_NakWithDelay_Call struct {
	*mock.Call
}

// NakWithDelay is a helper method to define mock.On call
//   - delay time.Duration
func (_e *Msg_Expecter) NakWithDelay(delay interface{}) *Msg_NakWithDelay_Call {
	return &Msg_NakWithDelay_Call{Call: _e.mock.On("NakWithDelay", delay)}
}

func (_c *Msg_NakWithDelay_Call) Run(run func(delay time.Duration)) *Msg_NakWithDelay_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(time.Duration))
	})
	return _c
}

func (_c *Msg_NakWithDelay_Call) Return(_a0 error) *Msg_NakWithDelay_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Msg_NakWithDelay_Call) RunAndReturn(run func(time.Duration) error) *Msg_NakWithDelay_Call {
	_c.Call.Return(run)
	return _c
}

// Reply provides a mock function with given fields:
func (_m *Msg) Reply() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Reply")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// Msg_Reply_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Reply'
type Msg_Reply_Call struct {
	*mock.Call
}

// Reply is a helper method to define mock.On call
func (_e *Msg_Expecter) Reply() *Msg_Reply_Call {
	return &Msg_Reply_Call{Call: _e.mock.On("Reply")}
}

func (_c *Msg_Reply_Call) Run(run func()) *Msg_Reply_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Msg_Reply_Call) Return(_a0 string) *Msg_Reply_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Msg_Reply_Call) RunAndReturn(run func() string) *Msg_Reply_Call {
	_c.Call.Return(run)
	return _c
}

// Subject provides a mock function with given fields:
func (_m *Msg) Subject() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Subject")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// Msg_Subject_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Subject'
type Msg_Subject_Call struct {
	*mock.Call
}

// Subject is a helper method to define mock.On call
func (_e *Msg_Expecter) Subject() *Msg_Subject_Call {
	return &Msg_Subject_Call{Call: _e.mock.On("Subject")}
}

func (_c *Msg_Subject_Call) Run(run func()) *Msg_Subject_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Msg_Subject_Call) Return(_a0 string) *Msg_Subject_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Msg_Subject_Call) RunAndReturn(run func() string) *Msg_Subject_Call {
	_c.Call.Return(run)
	return _c
}

// Term provides a mock function with given fields:
func (_m *Msg) Term() error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Term")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Msg_Term_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Term'
type Msg_Term_Call struct {
	*mock.Call
}

// Term is a helper method to define mock.On call
func (_e *Msg_Expecter) Term() *Msg_Term_Call {
	return &Msg_Term_Call{Call: _e.mock.On("Term")}
}

func (_c *Msg_Term_Call) Run(run func()) *Msg_Term_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Msg_Term_Call) Return(_a0 error) *Msg_Term_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Msg_Term_Call) RunAndReturn(run func() error) *Msg_Term_Call {
	_c.Call.Return(run)
	return _c
}

// TermWithReason provides a mock function with given fields: reason
func (_m *Msg) TermWithReason(reason string) error {
	ret := _m.Called(reason)

	if len(ret) == 0 {
		panic("no return value specified for TermWithReason")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(reason)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Msg_TermWithReason_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'TermWithReason'
type Msg_TermWithReason_Call struct {
	*mock.Call
}

// TermWithReason is a helper method to define mock.On call
//   - reason string
func (_e *Msg_Expecter) TermWithReason(reason interface{}) *Msg_TermWithReason_Call {
	return &Msg_TermWithReason_Call{Call: _e.mock.On("TermWithReason", reason)}
}

func (_c *Msg_TermWithReason_Call) Run(run func(reason string)) *Msg_TermWithReason_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *Msg_TermWithReason_Call) Return(_a0 error) *Msg_TermWithReason_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Msg_TermWithReason_Call) RunAndReturn(run func(string) error) *Msg_TermWithReason_Call {
	_c.Call.Return(run)
	return _c
}

// NewMsg creates a new instance of Msg. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMsg(t interface {
	mock.TestingT
	Cleanup(func())
}) *Msg {
	mock := &Msg{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
