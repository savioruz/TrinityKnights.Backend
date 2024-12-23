package venue_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/TrinityKnights/Backend/internal/delivery/http/handler/venue"
	"github.com/TrinityKnights/Backend/internal/domain/model"
	domainErrors "github.com/TrinityKnights/Backend/pkg/errors"
	mockVenue "github.com/TrinityKnights/Backend/test/mock/service/venue"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func setupTest(t *testing.T) (*venue.VenueHandlerImpl, *mockVenue.MockVenueService, *echo.Echo) {
	ctrl := gomock.NewController(t)
	mockVenueService := mockVenue.NewMockVenueService(ctrl)
	logger := logrus.New()
	handler := venue.NewVenueHandler(logger, mockVenueService).(*venue.VenueHandlerImpl)
	e := echo.New()
	return handler, mockVenueService, e
}

func TestVenueHandler_CreateVenue(t *testing.T) {
	handler, mockVenueService, e := setupTest(t)

	tests := []struct {
		name           string
		requestBody    string
		setupMock      func()
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Success",
			requestBody: `{
				"name": "Test Venue",
				"address": "123 Test St",
				"capacity": 1000,
				"city": "Test City",
				"state": "TS",
				"zip": "12345"
			}`,
			setupMock: func() {
				mockVenueService.EXPECT().
					CreateVenue(gomock.Any(), &model.CreateVenueRequest{
						Name:     "Test Venue",
						Address:  "123 Test St",
						Capacity: 1000,
						City:     "Test City",
						State:    "TS",
						Zip:      "12345",
					}).
					Return(&model.VenueResponse{
						ID:       1,
						Name:     "Test Venue",
						Address:  "123 Test St",
						Capacity: 1000,
						City:     "Test City",
						State:    "TS",
						Zip:      "12345",
					}, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   `{"data":{"id":1,"name":"Test Venue","address":"123 Test St","capacity":1000,"city":"Test City","state":"TS","zip":"12345"}}`,
		},
		{
			name:           "Invalid JSON",
			requestBody:    `{"invalid json`,
			setupMock:      func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":{"code":400,"message":"code=400, message=unexpected EOF, internal=unexpected EOF"}}`,
		},
		{
			name: "Empty Fields",
			requestBody: `{
				"name": "",
				"address": "",
				"city": "",
				"state": "",
				"zip": ""
			}`,
			setupMock: func() {
				mockVenueService.EXPECT().
					CreateVenue(gomock.Any(), gomock.Any()).
					Return(nil, domainErrors.ErrBadRequest)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":{"code":400,"message":"bad request"}}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/venues", strings.NewReader(tc.requestBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			tc.setupMock()

			err := handler.CreateVenue(c)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedStatus, rec.Code)

			var actualBody, expectedBody map[string]interface{}
			json.Unmarshal(rec.Body.Bytes(), &actualBody)
			json.Unmarshal([]byte(tc.expectedBody), &expectedBody)
			assert.Equal(t, expectedBody, actualBody)
		})
	}
}

func TestVenueHandler_GetVenueByID(t *testing.T) {
	handler, mockVenueService, e := setupTest(t)

	tests := []struct {
		name           string
		venueID        string
		setupMock      func()
		expectedStatus int
		expectedBody   string
	}{
		{
			name:    "Success",
			venueID: "1",
			setupMock: func() {
				mockVenueService.EXPECT().
					GetVenueByID(gomock.Any(), &model.GetVenueRequest{
						ID: 1,
					}).
					Return(&model.VenueResponse{
						ID:       1,
						Name:     "Test Venue",
						Address:  "123 Test St",
						Capacity: 1000,
						City:     "Test City",
						State:    "TS",
						Zip:      "12345",
					}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"data":{"id":1,"name":"Test Venue","address":"123 Test St","capacity":1000,"city":"Test City","state":"TS","zip":"12345"}}`,
		},
		{
			name:    "Venue Not Found",
			venueID: "999",
			setupMock: func() {
				mockVenueService.EXPECT().
					GetVenueByID(gomock.Any(), &model.GetVenueRequest{
						ID: 999,
					}).
					Return(nil, domainErrors.ErrNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"error":{"code":404,"message":"not found"}}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/venues/"+tc.venueID, http.NoBody)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(tc.venueID)

			tc.setupMock()

			err := handler.GetVenueByID(c)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedStatus, rec.Code)

			var actualBody, expectedBody map[string]interface{}
			json.Unmarshal(rec.Body.Bytes(), &actualBody)
			json.Unmarshal([]byte(tc.expectedBody), &expectedBody)
			assert.Equal(t, expectedBody, actualBody)
		})
	}
}

func TestVenueHandler_UpdateVenue(t *testing.T) {
	handler, mockVenueService, e := setupTest(t)

	tests := []struct {
		name           string
		venueID        string
		requestBody    string
		setupMock      func()
		expectedStatus int
		expectedBody   string
	}{
		{
			name:    "Success",
			venueID: "1",
			requestBody: `{
				"name": "Updated Venue",
				"address": "456 Update St",
				"capacity": 2000,
				"city": "Update City",
				"state": "UP",
				"zip": "54321"
			}`,
			setupMock: func() {
				mockVenueService.EXPECT().
					UpdateVenue(gomock.Any(), &model.UpdateVenueRequest{
						ID:       1,
						Name:     "Updated Venue",
						Address:  "456 Update St",
						Capacity: 2000,
						City:     "Update City",
						State:    "UP",
						Zip:      "54321",
					}).
					Return(&model.VenueResponse{
						ID:       1,
						Name:     "Updated Venue",
						Address:  "456 Update St",
						Capacity: 2000,
						City:     "Update City",
						State:    "UP",
						Zip:      "54321",
					}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"data":{"id":1,"name":"Updated Venue","address":"456 Update St","capacity":2000,"city":"Update City","state":"UP","zip":"54321"}}`,
		},
		{
			name:    "Venue Not Found",
			venueID: "999",
			requestBody: `{
				"name": "Updated Venue"
			}`,
			setupMock: func() {
				mockVenueService.EXPECT().
					UpdateVenue(gomock.Any(), gomock.Any()).
					Return(nil, domainErrors.ErrNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"error":{"code":404,"message":"not found"}}`,
		},
		{
			name:           "Invalid JSON",
			venueID:        "1",
			requestBody:    `{"invalid json`,
			setupMock:      func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":{"code":400,"message":"code=400, message=unexpected EOF, internal=unexpected EOF"}}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPut, "/venues/"+tc.venueID, strings.NewReader(tc.requestBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(tc.venueID)

			tc.setupMock()

			err := handler.UpdateVenue(c)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedStatus, rec.Code)

			var actualBody, expectedBody map[string]interface{}
			json.Unmarshal(rec.Body.Bytes(), &actualBody)
			json.Unmarshal([]byte(tc.expectedBody), &expectedBody)
			assert.Equal(t, expectedBody, actualBody)
		})
	}
}
