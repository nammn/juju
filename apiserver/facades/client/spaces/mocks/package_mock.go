// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/juju/juju/apiserver/facades/client/spaces (interfaces: Backing,BlockChecker,Machine,RenameSpace,RenameSpaceState,Settings,OpFactory,RemoveSpace,RemoveSpaceState)

// Package mocks is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	set "github.com/juju/collections/set"
	networkingcommon "github.com/juju/juju/apiserver/common/networkingcommon"
	spaces "github.com/juju/juju/apiserver/facades/client/spaces"
	controller "github.com/juju/juju/controller"
	network "github.com/juju/juju/core/network"
	settings "github.com/juju/juju/core/settings"
	environs "github.com/juju/juju/environs"
	config "github.com/juju/juju/environs/config"
	state "github.com/juju/juju/state"
	names_v3 "gopkg.in/juju/names.v3"
	txn "gopkg.in/mgo.v2/txn"
	reflect "reflect"
)

// MockBacking is a mock of Backing interface
type MockBacking struct {
	ctrl     *gomock.Controller
	recorder *MockBackingMockRecorder
}

// MockBackingMockRecorder is the mock recorder for MockBacking
type MockBackingMockRecorder struct {
	mock *MockBacking
}

// NewMockBacking creates a new mock instance
func NewMockBacking(ctrl *gomock.Controller) *MockBacking {
	mock := &MockBacking{ctrl: ctrl}
	mock.recorder = &MockBackingMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockBacking) EXPECT() *MockBackingMockRecorder {
	return m.recorder
}

// AddSpace mocks base method
func (m *MockBacking) AddSpace(arg0 string, arg1 network.Id, arg2 []string, arg3 bool) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddSpace", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddSpace indicates an expected call of AddSpace
func (mr *MockBackingMockRecorder) AddSpace(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddSpace", reflect.TypeOf((*MockBacking)(nil).AddSpace), arg0, arg1, arg2, arg3)
}

// AllEndpointBindings mocks base method
func (m *MockBacking) AllEndpointBindings() ([]spaces.ApplicationEndpointBindingsShim, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AllEndpointBindings")
	ret0, _ := ret[0].([]spaces.ApplicationEndpointBindingsShim)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AllEndpointBindings indicates an expected call of AllEndpointBindings
func (mr *MockBackingMockRecorder) AllEndpointBindings() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AllEndpointBindings", reflect.TypeOf((*MockBacking)(nil).AllEndpointBindings))
}

// AllMachines mocks base method
func (m *MockBacking) AllMachines() ([]spaces.Machine, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AllMachines")
	ret0, _ := ret[0].([]spaces.Machine)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AllMachines indicates an expected call of AllMachines
func (mr *MockBackingMockRecorder) AllMachines() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AllMachines", reflect.TypeOf((*MockBacking)(nil).AllMachines))
}

// AllSpaces mocks base method
func (m *MockBacking) AllSpaces() ([]networkingcommon.BackingSpace, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AllSpaces")
	ret0, _ := ret[0].([]networkingcommon.BackingSpace)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AllSpaces indicates an expected call of AllSpaces
func (mr *MockBackingMockRecorder) AllSpaces() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AllSpaces", reflect.TypeOf((*MockBacking)(nil).AllSpaces))
}

// ApplyOperation mocks base method
func (m *MockBacking) ApplyOperation(arg0 state.ModelOperation) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ApplyOperation", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// ApplyOperation indicates an expected call of ApplyOperation
func (mr *MockBackingMockRecorder) ApplyOperation(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ApplyOperation", reflect.TypeOf((*MockBacking)(nil).ApplyOperation), arg0)
}

// CloudSpec mocks base method
func (m *MockBacking) CloudSpec() (environs.CloudSpec, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CloudSpec")
	ret0, _ := ret[0].(environs.CloudSpec)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CloudSpec indicates an expected call of CloudSpec
func (mr *MockBackingMockRecorder) CloudSpec() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CloudSpec", reflect.TypeOf((*MockBacking)(nil).CloudSpec))
}

// ConstraintsTagForSpaceName mocks base method
func (m *MockBacking) ConstraintsTagForSpaceName(arg0 string) ([]names_v3.Tag, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ConstraintsTagForSpaceName", arg0)
	ret0, _ := ret[0].([]names_v3.Tag)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ConstraintsTagForSpaceName indicates an expected call of ConstraintsTagForSpaceName
func (mr *MockBackingMockRecorder) ConstraintsTagForSpaceName(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ConstraintsTagForSpaceName", reflect.TypeOf((*MockBacking)(nil).ConstraintsTagForSpaceName), arg0)
}

// ControllerConfig mocks base method
func (m *MockBacking) ControllerConfig() (controller.Config, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ControllerConfig")
	ret0, _ := ret[0].(controller.Config)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ControllerConfig indicates an expected call of ControllerConfig
func (mr *MockBackingMockRecorder) ControllerConfig() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ControllerConfig", reflect.TypeOf((*MockBacking)(nil).ControllerConfig))
}

// IsControllerModel mocks base method
func (m *MockBacking) IsControllerModel() (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsControllerModel")
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsControllerModel indicates an expected call of IsControllerModel
func (mr *MockBackingMockRecorder) IsControllerModel() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsControllerModel", reflect.TypeOf((*MockBacking)(nil).IsControllerModel))
}

// ModelConfig mocks base method
func (m *MockBacking) ModelConfig() (*config.Config, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ModelConfig")
	ret0, _ := ret[0].(*config.Config)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ModelConfig indicates an expected call of ModelConfig
func (mr *MockBackingMockRecorder) ModelConfig() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ModelConfig", reflect.TypeOf((*MockBacking)(nil).ModelConfig))
}

// ModelTag mocks base method
func (m *MockBacking) ModelTag() names_v3.ModelTag {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ModelTag")
	ret0, _ := ret[0].(names_v3.ModelTag)
	return ret0
}

// ModelTag indicates an expected call of ModelTag
func (mr *MockBackingMockRecorder) ModelTag() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ModelTag", reflect.TypeOf((*MockBacking)(nil).ModelTag))
}

// ReloadSpaces mocks base method
func (m *MockBacking) ReloadSpaces(arg0 environs.BootstrapEnviron) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReloadSpaces", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// ReloadSpaces indicates an expected call of ReloadSpaces
func (mr *MockBackingMockRecorder) ReloadSpaces(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReloadSpaces", reflect.TypeOf((*MockBacking)(nil).ReloadSpaces), arg0)
}

// SpaceByName mocks base method
func (m *MockBacking) SpaceByName(arg0 string) (networkingcommon.BackingSpace, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SpaceByName", arg0)
	ret0, _ := ret[0].(networkingcommon.BackingSpace)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SpaceByName indicates an expected call of SpaceByName
func (mr *MockBackingMockRecorder) SpaceByName(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SpaceByName", reflect.TypeOf((*MockBacking)(nil).SpaceByName), arg0)
}

// SubnetByCIDR mocks base method
func (m *MockBacking) SubnetByCIDR(arg0 string) (networkingcommon.BackingSubnet, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SubnetByCIDR", arg0)
	ret0, _ := ret[0].(networkingcommon.BackingSubnet)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SubnetByCIDR indicates an expected call of SubnetByCIDR
func (mr *MockBackingMockRecorder) SubnetByCIDR(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SubnetByCIDR", reflect.TypeOf((*MockBacking)(nil).SubnetByCIDR), arg0)
}

// MockBlockChecker is a mock of BlockChecker interface
type MockBlockChecker struct {
	ctrl     *gomock.Controller
	recorder *MockBlockCheckerMockRecorder
}

// MockBlockCheckerMockRecorder is the mock recorder for MockBlockChecker
type MockBlockCheckerMockRecorder struct {
	mock *MockBlockChecker
}

// NewMockBlockChecker creates a new mock instance
func NewMockBlockChecker(ctrl *gomock.Controller) *MockBlockChecker {
	mock := &MockBlockChecker{ctrl: ctrl}
	mock.recorder = &MockBlockCheckerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockBlockChecker) EXPECT() *MockBlockCheckerMockRecorder {
	return m.recorder
}

// ChangeAllowed mocks base method
func (m *MockBlockChecker) ChangeAllowed() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ChangeAllowed")
	ret0, _ := ret[0].(error)
	return ret0
}

// ChangeAllowed indicates an expected call of ChangeAllowed
func (mr *MockBlockCheckerMockRecorder) ChangeAllowed() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChangeAllowed", reflect.TypeOf((*MockBlockChecker)(nil).ChangeAllowed))
}

// RemoveAllowed mocks base method
func (m *MockBlockChecker) RemoveAllowed() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveAllowed")
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveAllowed indicates an expected call of RemoveAllowed
func (mr *MockBlockCheckerMockRecorder) RemoveAllowed() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveAllowed", reflect.TypeOf((*MockBlockChecker)(nil).RemoveAllowed))
}

// MockMachine is a mock of Machine interface
type MockMachine struct {
	ctrl     *gomock.Controller
	recorder *MockMachineMockRecorder
}

// MockMachineMockRecorder is the mock recorder for MockMachine
type MockMachineMockRecorder struct {
	mock *MockMachine
}

// NewMockMachine creates a new mock instance
func NewMockMachine(ctrl *gomock.Controller) *MockMachine {
	mock := &MockMachine{ctrl: ctrl}
	mock.recorder = &MockMachineMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockMachine) EXPECT() *MockMachineMockRecorder {
	return m.recorder
}

// AllSpaces mocks base method
func (m *MockMachine) AllSpaces() (set.Strings, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AllSpaces")
	ret0, _ := ret[0].(set.Strings)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AllSpaces indicates an expected call of AllSpaces
func (mr *MockMachineMockRecorder) AllSpaces() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AllSpaces", reflect.TypeOf((*MockMachine)(nil).AllSpaces))
}

// MockRenameSpace is a mock of RenameSpace interface
type MockRenameSpace struct {
	ctrl     *gomock.Controller
	recorder *MockRenameSpaceMockRecorder
}

// MockRenameSpaceMockRecorder is the mock recorder for MockRenameSpace
type MockRenameSpaceMockRecorder struct {
	mock *MockRenameSpace
}

// NewMockRenameSpace creates a new mock instance
func NewMockRenameSpace(ctrl *gomock.Controller) *MockRenameSpace {
	mock := &MockRenameSpace{ctrl: ctrl}
	mock.recorder = &MockRenameSpaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockRenameSpace) EXPECT() *MockRenameSpaceMockRecorder {
	return m.recorder
}

// Id mocks base method
func (m *MockRenameSpace) Id() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Id")
	ret0, _ := ret[0].(string)
	return ret0
}

// Id indicates an expected call of Id
func (mr *MockRenameSpaceMockRecorder) Id() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Id", reflect.TypeOf((*MockRenameSpace)(nil).Id))
}

// Name mocks base method
func (m *MockRenameSpace) Name() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Name")
	ret0, _ := ret[0].(string)
	return ret0
}

// Name indicates an expected call of Name
func (mr *MockRenameSpaceMockRecorder) Name() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Name", reflect.TypeOf((*MockRenameSpace)(nil).Name))
}

// Refresh mocks base method
func (m *MockRenameSpace) Refresh() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Refresh")
	ret0, _ := ret[0].(error)
	return ret0
}

// Refresh indicates an expected call of Refresh
func (mr *MockRenameSpaceMockRecorder) Refresh() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Refresh", reflect.TypeOf((*MockRenameSpace)(nil).Refresh))
}

// RenameSpaceOps mocks base method
func (m *MockRenameSpace) RenameSpaceOps(arg0 string) []txn.Op {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RenameSpaceOps", arg0)
	ret0, _ := ret[0].([]txn.Op)
	return ret0
}

// RenameSpaceOps indicates an expected call of RenameSpaceOps
func (mr *MockRenameSpaceMockRecorder) RenameSpaceOps(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RenameSpaceOps", reflect.TypeOf((*MockRenameSpace)(nil).RenameSpaceOps), arg0)
}

// MockRenameSpaceState is a mock of RenameSpaceState interface
type MockRenameSpaceState struct {
	ctrl     *gomock.Controller
	recorder *MockRenameSpaceStateMockRecorder
}

// MockRenameSpaceStateMockRecorder is the mock recorder for MockRenameSpaceState
type MockRenameSpaceStateMockRecorder struct {
	mock *MockRenameSpaceState
}

// NewMockRenameSpaceState creates a new mock instance
func NewMockRenameSpaceState(ctrl *gomock.Controller) *MockRenameSpaceState {
	mock := &MockRenameSpaceState{ctrl: ctrl}
	mock.recorder = &MockRenameSpaceStateMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockRenameSpaceState) EXPECT() *MockRenameSpaceStateMockRecorder {
	return m.recorder
}

// ConstraintsOpsForSpaceNameChange mocks base method
func (m *MockRenameSpaceState) ConstraintsOpsForSpaceNameChange(arg0, arg1 string) ([]txn.Op, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ConstraintsOpsForSpaceNameChange", arg0, arg1)
	ret0, _ := ret[0].([]txn.Op)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ConstraintsOpsForSpaceNameChange indicates an expected call of ConstraintsOpsForSpaceNameChange
func (mr *MockRenameSpaceStateMockRecorder) ConstraintsOpsForSpaceNameChange(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ConstraintsOpsForSpaceNameChange", reflect.TypeOf((*MockRenameSpaceState)(nil).ConstraintsOpsForSpaceNameChange), arg0, arg1)
}

// ControllerConfig mocks base method
func (m *MockRenameSpaceState) ControllerConfig() (controller.Config, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ControllerConfig")
	ret0, _ := ret[0].(controller.Config)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ControllerConfig indicates an expected call of ControllerConfig
func (mr *MockRenameSpaceStateMockRecorder) ControllerConfig() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ControllerConfig", reflect.TypeOf((*MockRenameSpaceState)(nil).ControllerConfig))
}

// MockSettings is a mock of Settings interface
type MockSettings struct {
	ctrl     *gomock.Controller
	recorder *MockSettingsMockRecorder
}

// MockSettingsMockRecorder is the mock recorder for MockSettings
type MockSettingsMockRecorder struct {
	mock *MockSettings
}

// NewMockSettings creates a new mock instance
func NewMockSettings(ctrl *gomock.Controller) *MockSettings {
	mock := &MockSettings{ctrl: ctrl}
	mock.recorder = &MockSettingsMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockSettings) EXPECT() *MockSettingsMockRecorder {
	return m.recorder
}

// DeltaOps mocks base method
func (m *MockSettings) DeltaOps(arg0 string, arg1 settings.ItemChanges) ([]txn.Op, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeltaOps", arg0, arg1)
	ret0, _ := ret[0].([]txn.Op)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeltaOps indicates an expected call of DeltaOps
func (mr *MockSettingsMockRecorder) DeltaOps(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeltaOps", reflect.TypeOf((*MockSettings)(nil).DeltaOps), arg0, arg1)
}

// MockOpFactory is a mock of OpFactory interface
type MockOpFactory struct {
	ctrl     *gomock.Controller
	recorder *MockOpFactoryMockRecorder
}

// MockOpFactoryMockRecorder is the mock recorder for MockOpFactory
type MockOpFactoryMockRecorder struct {
	mock *MockOpFactory
}

// NewMockOpFactory creates a new mock instance
func NewMockOpFactory(ctrl *gomock.Controller) *MockOpFactory {
	mock := &MockOpFactory{ctrl: ctrl}
	mock.recorder = &MockOpFactoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockOpFactory) EXPECT() *MockOpFactoryMockRecorder {
	return m.recorder
}

// NewRemoveSpaceModelOp mocks base method
func (m *MockOpFactory) NewRemoveSpaceModelOp(arg0 string) (state.ModelOperation, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewRemoveSpaceModelOp", arg0)
	ret0, _ := ret[0].(state.ModelOperation)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NewRemoveSpaceModelOp indicates an expected call of NewRemoveSpaceModelOp
func (mr *MockOpFactoryMockRecorder) NewRemoveSpaceModelOp(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewRemoveSpaceModelOp", reflect.TypeOf((*MockOpFactory)(nil).NewRemoveSpaceModelOp), arg0)
}

// NewRenameSpaceModelOp mocks base method
func (m *MockOpFactory) NewRenameSpaceModelOp(arg0, arg1 string) (state.ModelOperation, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewRenameSpaceModelOp", arg0, arg1)
	ret0, _ := ret[0].(state.ModelOperation)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NewRenameSpaceModelOp indicates an expected call of NewRenameSpaceModelOp
func (mr *MockOpFactoryMockRecorder) NewRenameSpaceModelOp(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewRenameSpaceModelOp", reflect.TypeOf((*MockOpFactory)(nil).NewRenameSpaceModelOp), arg0, arg1)
}

// MockRemoveSpace is a mock of RemoveSpace interface
type MockRemoveSpace struct {
	ctrl     *gomock.Controller
	recorder *MockRemoveSpaceMockRecorder
}

// MockRemoveSpaceMockRecorder is the mock recorder for MockRemoveSpace
type MockRemoveSpaceMockRecorder struct {
	mock *MockRemoveSpace
}

// NewMockRemoveSpace creates a new mock instance
func NewMockRemoveSpace(ctrl *gomock.Controller) *MockRemoveSpace {
	mock := &MockRemoveSpace{ctrl: ctrl}
	mock.recorder = &MockRemoveSpaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockRemoveSpace) EXPECT() *MockRemoveSpaceMockRecorder {
	return m.recorder
}

// Id mocks base method
func (m *MockRemoveSpace) Id() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Id")
	ret0, _ := ret[0].(string)
	return ret0
}

// Id indicates an expected call of Id
func (mr *MockRemoveSpaceMockRecorder) Id() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Id", reflect.TypeOf((*MockRemoveSpace)(nil).Id))
}

// Name mocks base method
func (m *MockRemoveSpace) Name() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Name")
	ret0, _ := ret[0].(string)
	return ret0
}

// Name indicates an expected call of Name
func (mr *MockRemoveSpaceMockRecorder) Name() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Name", reflect.TypeOf((*MockRemoveSpace)(nil).Name))
}

// Refresh mocks base method
func (m *MockRemoveSpace) Refresh() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Refresh")
	ret0, _ := ret[0].(error)
	return ret0
}

// Refresh indicates an expected call of Refresh
func (mr *MockRemoveSpaceMockRecorder) Refresh() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Refresh", reflect.TypeOf((*MockRemoveSpace)(nil).Refresh))
}

// RemoveSpaceOps mocks base method
func (m *MockRemoveSpace) RemoveSpaceOps() []txn.Op {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveSpaceOps")
	ret0, _ := ret[0].([]txn.Op)
	return ret0
}

// RemoveSpaceOps indicates an expected call of RemoveSpaceOps
func (mr *MockRemoveSpaceMockRecorder) RemoveSpaceOps() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveSpaceOps", reflect.TypeOf((*MockRemoveSpace)(nil).RemoveSpaceOps))
}

// MockRemoveSpaceState is a mock of RemoveSpaceState interface
type MockRemoveSpaceState struct {
	ctrl     *gomock.Controller
	recorder *MockRemoveSpaceStateMockRecorder
}

// MockRemoveSpaceStateMockRecorder is the mock recorder for MockRemoveSpaceState
type MockRemoveSpaceStateMockRecorder struct {
	mock *MockRemoveSpaceState
}

// NewMockRemoveSpaceState creates a new mock instance
func NewMockRemoveSpaceState(ctrl *gomock.Controller) *MockRemoveSpaceState {
	mock := &MockRemoveSpaceState{ctrl: ctrl}
	mock.recorder = &MockRemoveSpaceStateMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockRemoveSpaceState) EXPECT() *MockRemoveSpaceStateMockRecorder {
	return m.recorder
}
