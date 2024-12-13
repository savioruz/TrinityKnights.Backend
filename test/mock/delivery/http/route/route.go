// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/delivery/http/route/route.go

// Package mock_route is a generated GoMock package.
package mock_route

import (
	reflect "reflect"

	route "github.com/TrinityKnights/Backend/pkg/route"
	gomock "github.com/golang/mock/gomock"
)

// MockRouteConfig is a mock of RouteConfig interface.
type MockRouteConfig struct {
	ctrl     *gomock.Controller
	recorder *MockRouteConfigMockRecorder
}

// MockRouteConfigMockRecorder is the mock recorder for MockRouteConfig.
type MockRouteConfigMockRecorder struct {
	mock *MockRouteConfig
}

// NewMockRouteConfig creates a new mock instance.
func NewMockRouteConfig(ctrl *gomock.Controller) *MockRouteConfig {
	mock := &MockRouteConfig{ctrl: ctrl}
	mock.recorder = &MockRouteConfigMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRouteConfig) EXPECT() *MockRouteConfigMockRecorder {
	return m.recorder
}

// PrivateRoute mocks base method.
func (m *MockRouteConfig) PrivateRoute() []route.Route {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PrivateRoute")
	ret0, _ := ret[0].([]route.Route)
	return ret0
}

// PrivateRoute indicates an expected call of PrivateRoute.
func (mr *MockRouteConfigMockRecorder) PrivateRoute() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PrivateRoute", reflect.TypeOf((*MockRouteConfig)(nil).PrivateRoute))
}

// PublicRoute mocks base method.
func (m *MockRouteConfig) PublicRoute() []route.Route {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PublicRoute")
	ret0, _ := ret[0].([]route.Route)
	return ret0
}

// PublicRoute indicates an expected call of PublicRoute.
func (mr *MockRouteConfigMockRecorder) PublicRoute() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PublicRoute", reflect.TypeOf((*MockRouteConfig)(nil).PublicRoute))
}

// SwaggerRoutes mocks base method.
func (m *MockRouteConfig) SwaggerRoutes() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SwaggerRoutes")
}

// SwaggerRoutes indicates an expected call of SwaggerRoutes.
func (mr *MockRouteConfigMockRecorder) SwaggerRoutes() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SwaggerRoutes", reflect.TypeOf((*MockRouteConfig)(nil).SwaggerRoutes))
}
