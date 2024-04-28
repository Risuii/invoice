package handler

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Risuii/invoice/src/v1/contract"
	"github.com/go-playground/assert"
	"go.uber.org/mock/gomock"

	errorss "github.com/Risuii/invoice/src/errors"
	mock_handler "github.com/Risuii/invoice/src/v1/handler/mock"
)

func TestHandler_GetListInovice(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	type (
		expected struct {
			statusCode   int
			responseBody string
		}

		given struct {
			param        contract.GetListParam
			svcErrReturn error
		}

		testCase struct {
			name     string
			given    given
			expected expected
		}
	)

	testCases := []testCase{
		{
			name: "err",
			given: given{
				param: contract.GetListParam{
					Page:  1,
					Limit: 10,
				},
				svcErrReturn: errors.New("error"),
			},
			expected: expected{
				statusCode:   500,
				responseBody: `{"data":null,"error":{"code":"err_internal_server","message_title":"Server Error","message":"Failed to process request, please try again in a moment.","message_severity":"error"},"success":false,"metadata":{"request_id":""}}`,
			},
		},
		{
			name: "success",
			given: given{
				param: contract.GetListParam{
					Page:  1,
					Limit: 10,
				},
				svcErrReturn: nil,
			},
			expected: expected{
				statusCode:   200,
				responseBody: `{"data":{"Data":null,"Pagination":null},"error":null,"success":true,"metadata":{"request_id":""}}`,
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "/just/for/testing", nil)
			w := httptest.NewRecorder()

			dataFromService := contract.ListInvoiceResponse{}
			InoviceSvc := mock_handler.NewMockInvoiceService(mockCtrl)

			InoviceSvc.EXPECT().GetList(gomock.Any(), testCase.given.param).
				Return(dataFromService, testCase.given.svcErrReturn).
				Times(1)

			hf := http.HandlerFunc(GetListInvoicesHandler(InoviceSvc))
			hf.ServeHTTP(w, r)

			res := w.Result()
			data, err := io.ReadAll(res.Body)
			if err != nil {
				t.Errorf("expected error to be nil got %v", err)
			}

			assert.Equal(t, testCase.expected.statusCode, res.StatusCode)
			assert.Equal(t, fmt.Sprintf("%s\n", testCase.expected.responseBody), string(data))
		})
	}
}

func TestHandler_GetInovice(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	type (
		expected struct {
			statusCode   int
			responseBody string
		}

		given struct {
			id           string
			svcErrReturn error
		}

		testCase struct {
			name     string
			given    given
			expected expected
		}
	)

	testCases := []testCase{
		{
			name: "err internal server",
			given: given{
				id:           "",
				svcErrReturn: errors.New("error internal server"),
			},
			expected: expected{
				statusCode:   500,
				responseBody: `{"data":null,"error":{"code":"err_internal_server","message_title":"Server Error","message":"Failed to process request, please try again in a moment.","message_severity":"error"},"success":false,"metadata":{"request_id":""}}`,
			},
		},
		{
			name: "err Inovice id not found",
			given: given{
				id:           "",
				svcErrReturn: errorss.ErrInvoiceIdNotFound,
			},
			expected: expected{
				statusCode:   422,
				responseBody: `{"data":null,"error":{"code":"err_invoice_id_not_found","message_title":"err_invoice_id_not_found_title","message":"err_invoice_id_not_found_message","message_severity":"error"},"success":false,"metadata":{"request_id":""}}`,
			},
		},
		{
			name: "err customer id not found",
			given: given{
				id:           "",
				svcErrReturn: errorss.ErrCustomerIdNotFound,
			},
			expected: expected{
				statusCode:   422,
				responseBody: `{"data":null,"error":{"code":"err_customer_id_not_found","message_title":"err_customer_id_not_found_title","message":"err_customer_id_not_found_message","message_severity":"error"},"success":false,"metadata":{"request_id":""}}`,
			},
		},
		{
			name: "success",
			given: given{
				id:           "",
				svcErrReturn: nil,
			},
			expected: expected{
				statusCode:   200,
				responseBody: `{"data":{"invoice_id":"","issue_date":"","subject":"","total_item":0,"item":null,"customer_name":"","due_date":"","status":"","sub_total":0,"tax":0,"grand_total":0},"error":null,"success":true,"metadata":{"request_id":""}}`,
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/just/for/testing/%s", testCase.given.id), nil)
			w := httptest.NewRecorder()

			dataFromService := contract.InvoiceResponse{}
			mockInovice := mock_handler.NewMockInvoiceService(mockCtrl)

			mockInovice.EXPECT().GetDetail(gomock.Any(), testCase.given.id).
				Return(dataFromService, testCase.given.svcErrReturn).
				Times(1)

			hf := http.HandlerFunc(GetDetailInvoicesHandler(mockInovice))
			hf.ServeHTTP(w, r)

			res := w.Result()
			data, err := io.ReadAll(res.Body)
			if err != nil {
				t.Errorf("expected error to be nil got %v", err)
			}

			assert.Equal(t, testCase.expected.statusCode, res.StatusCode)
			assert.Equal(t, fmt.Sprintf("%s\n", testCase.expected.responseBody), string(data))
		})
	}
}

func TestHandler_CreateInvoice(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	type (
		expected struct {
			request      *contract.InvoiceRequest
			statusCode   int
			responseBody string
		}

		given struct {
			payload      string
			svcErrReturn error
		}

		testCase struct {
			name     string
			given    given
			expected expected
		}
	)

	testCases := []testCase{
		{
			name: "err bad request",
			given: given{
				payload: `{
					"subject": "",
					"issue_date": "23-01-2023",
					"due_date": "23-01-2024",
					"sub_total": 300,
					"tax": 10,
					"grand_total": 200,
					"customer_request": {
						"customer_name": "test-customer-name-1",
						"address": "test-address-1"
					},
					"item_request": [
						{
							"name": "test-1",
							"type": "test-type",
							"quantity": 1,
							"unit_price": 1,
							"amount": 1
						},
						{
							"name": "test-2",
							"type": "test-type",
							"quantity": 2,
							"unit_price": 2,
							"amount": 2
						},
						{
							"name": "test-3",
							"type": "test-type",
							"quantity": 3,
							"unit_price": 3,
							"amount": 3
						}
					]
				}`,
			},
			expected: expected{
				request:      nil,
				statusCode:   400,
				responseBody: `{"data":null,"error":{"code":"err_bad_request","message_title":"Bad Request","message":"Invalid request parameters","message_severity":"error"},"success":false,"metadata":{"request_id":""}}`,
			},
		},
		{
			name: "err internal server",
			given: given{
				payload: `{
					"subject": "test-subject-1",
					"issue_date": "23-01-2023",
					"due_date": "23-01-2024",
					"sub_total": 300,
					"tax": 10,
					"grand_total": 200,
					"customer_request": {
						"customer_name": "test-customer-name-1",
						"address": "test-address-1"
					},
					"item_request": [
						{
							"name": "test-1",
							"type": "test-type",
							"quantity": 1,
							"unit_price": 1,
							"amount": 1
						},
						{
							"name": "test-2",
							"type": "test-type",
							"quantity": 2,
							"unit_price": 2,
							"amount": 2
						},
						{
							"name": "test-3",
							"type": "test-type",
							"quantity": 3,
							"unit_price": 3,
							"amount": 3
						}
					]
				}`,
				svcErrReturn: errors.New("error internal server"),
			},
			expected: expected{
				request: &contract.InvoiceRequest{
					Subject:    "test-subject-1",
					IssueDate:  "23-01-2023",
					DueDate:    "23-01-2024",
					SubTotal:   300,
					Tax:        10,
					GrandTotal: 200,
					CustomerRequest: contract.CustomerRequest{
						CustomerName: "test-customer-name-1",
						Address:      "test-address-1",
					},
					ItemRequest: []contract.ItemRequest{
						{
							Name:      "test-1",
							Type:      "test-type",
							Quantity:  1,
							UnitPrice: 1,
							Amount:    1,
						},
						{
							Name:      "test-2",
							Type:      "test-type",
							Quantity:  2,
							UnitPrice: 2,
							Amount:    2,
						},
						{
							Name:      "test-3",
							Type:      "test-type",
							Quantity:  3,
							UnitPrice: 3,
							Amount:    3,
						},
					},
				},
				statusCode:   500,
				responseBody: `{"data":null,"error":{"code":"err_internal_server","message_title":"Server Error","message":"Failed to process request, please try again in a moment.","message_severity":"error"},"success":false,"metadata":{"request_id":""}}`,
			},
		},
		{
			name: "success",
			given: given{
				payload: `{
					"subject": "test-subject-1",
					"issue_date": "23-01-2023",
					"due_date": "23-01-2024",
					"sub_total": 300,
					"tax": 10,
					"grand_total": 200,
					"customer_request": {
						"customer_name": "test-customer-name-1",
						"address": "test-address-1"
					},
					"item_request": [
						{
							"name": "test-1",
							"type": "test-type",
							"quantity": 1,
							"unit_price": 1,
							"amount": 1
						},
						{
							"name": "test-2",
							"type": "test-type",
							"quantity": 2,
							"unit_price": 2,
							"amount": 2
						},
						{
							"name": "test-3",
							"type": "test-type",
							"quantity": 3,
							"unit_price": 3,
							"amount": 3
						}
					]
				}`,
				svcErrReturn: nil,
			},
			expected: expected{
				request: &contract.InvoiceRequest{
					Subject:    "test-subject-1",
					IssueDate:  "23-01-2023",
					DueDate:    "23-01-2024",
					SubTotal:   300,
					Tax:        10,
					GrandTotal: 200,
					CustomerRequest: contract.CustomerRequest{
						CustomerName: "test-customer-name-1",
						Address:      "test-address-1",
					},
					ItemRequest: []contract.ItemRequest{
						{
							Name:      "test-1",
							Type:      "test-type",
							Quantity:  1,
							UnitPrice: 1,
							Amount:    1,
						},
						{
							Name:      "test-2",
							Type:      "test-type",
							Quantity:  2,
							UnitPrice: 2,
							Amount:    2,
						},
						{
							Name:      "test-3",
							Type:      "test-type",
							Quantity:  3,
							UnitPrice: 3,
							Amount:    3,
						},
					},
				},
				statusCode:   200,
				responseBody: `{"data":{"invoice_id":""},"error":null,"success":true,"metadata":{"request_id":""}}`,
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, "/just/for/testing", strings.NewReader(testCase.given.payload))
			w := httptest.NewRecorder()

			dataFromService := contract.InvcResponse{}

			mockInvoiceSvc := mock_handler.NewMockInvoiceService(mockCtrl)

			if testCase.expected.request != nil {
				mockInvoiceSvc.EXPECT().Create(gomock.Any(), *testCase.expected.request).
					Return(dataFromService, testCase.given.svcErrReturn).
					Times(1)
			}

			hf := http.HandlerFunc(CreateInvoiceHandler(mockInvoiceSvc))
			hf.ServeHTTP(w, r)

			res := w.Result()
			data, err := io.ReadAll(res.Body)
			if err != nil {
				t.Errorf("expected error to be nil got %v", err)
			}

			assert.Equal(t, testCase.expected.statusCode, res.StatusCode)
			assert.Equal(t, fmt.Sprintf("%s\n", testCase.expected.responseBody), string(data))
		})
	}
}

func TestHandler_UpdateInvoice(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	type (
		expected struct {
			request      *contract.InvoiceRequest
			statusCode   int
			responseBody string
		}

		given struct {
			id           string
			payload      string
			svcErrReturn error
		}

		testCase struct {
			name     string
			given    given
			expected expected
		}
	)

	request := contract.InvoiceRequest{
		Subject:    "test-subject-2",
		IssueDate:  "23-01-2023",
		DueDate:    "23-01-2024",
		SubTotal:   300,
		Tax:        10,
		GrandTotal: 200,
		CustomerRequest: contract.CustomerRequest{
			CustomerName: "test-customer-name-2",
			Address:      "test-address-2",
		},
		ItemRequest: []contract.ItemRequest{
			{
				Name:      "test-1",
				Type:      "test-type",
				Quantity:  1,
				UnitPrice: 1,
				Amount:    1,
			},
			{
				Name:      "test-2",
				Type:      "test-type",
				Quantity:  2,
				UnitPrice: 2,
				Amount:    2,
			},
		},
	}

	testCases := []testCase{
		{
			name:  "err bad request",
			given: given{},
			expected: expected{
				statusCode:   400,
				responseBody: `{"data":null,"error":{"code":"err_bad_request","message_title":"Bad Request","message":"Invalid request parameters","message_severity":"error"},"success":false,"metadata":{"request_id":""}}`,
			},
		},
		{
			name: "err invoice id not found",
			given: given{
				id: "",
				payload: `{
					"subject": "test-subject-2",
					"issue_date": "23-01-2023",
					"due_date": "23-01-2024",
					"sub_total": 300,
					"tax": 10,
					"grand_total": 200,
					"customer_request": {
						"customer_name": "test-customer-name-2",
						"address": "test-address-2"
					},
					"item_request": [
						{
							"name": "test-1",
							"type": "test-type",
							"quantity": 1,
							"unit_price": 1,
							"amount": 1
						},
						{
							"name": "test-2",
							"type": "test-type",
							"quantity": 2,
							"unit_price": 2,
							"amount": 2
						}
					]
				}`,
				svcErrReturn: errorss.ErrInvoiceIdNotFound,
			},
			expected: expected{
				request:      &request,
				statusCode:   422,
				responseBody: `{"data":null,"error":{"code":"err_invoice_id_not_found","message_title":"err_invoice_id_not_found_title","message":"err_invoice_id_not_found_message","message_severity":"error"},"success":false,"metadata":{"request_id":""}}`,
			},
		},
		{
			name: "err customer id not found",
			given: given{
				id: "",
				payload: `{
					"subject": "test-subject-2",
					"issue_date": "23-01-2023",
					"due_date": "23-01-2024",
					"sub_total": 300,
					"tax": 10,
					"grand_total": 200,
					"customer_request": {
						"customer_name": "test-customer-name-2",
						"address": "test-address-2"
					},
					"item_request": [
						{
							"name": "test-1",
							"type": "test-type",
							"quantity": 1,
							"unit_price": 1,
							"amount": 1
						},
						{
							"name": "test-2",
							"type": "test-type",
							"quantity": 2,
							"unit_price": 2,
							"amount": 2
						}
					]
				}`,
				svcErrReturn: errorss.ErrCustomerIdNotFound,
			},
			expected: expected{
				request:      &request,
				statusCode:   422,
				responseBody: `{"data":null,"error":{"code":"err_customer_id_not_found","message_title":"err_customer_id_not_found_title","message":"err_customer_id_not_found_message","message_severity":"error"},"success":false,"metadata":{"request_id":""}}`,
			},
		},
		{
			name: "err internal server",
			given: given{
				id: "",
				payload: `{
					"subject": "test-subject-2",
					"issue_date": "23-01-2023",
					"due_date": "23-01-2024",
					"sub_total": 300,
					"tax": 10,
					"grand_total": 200,
					"customer_request": {
						"customer_name": "test-customer-name-2",
						"address": "test-address-2"
					},
					"item_request": [
						{
							"name": "test-1",
							"type": "test-type",
							"quantity": 1,
							"unit_price": 1,
							"amount": 1
						},
						{
							"name": "test-2",
							"type": "test-type",
							"quantity": 2,
							"unit_price": 2,
							"amount": 2
						}
					]
				}`,
				svcErrReturn: errors.New("error internal server"),
			},
			expected: expected{
				request:      &request,
				statusCode:   500,
				responseBody: `{"data":null,"error":{"code":"err_internal_server","message_title":"Server Error","message":"Failed to process request, please try again in a moment.","message_severity":"error"},"success":false,"metadata":{"request_id":""}}`,
			},
		},
		{
			name: "err internal server",
			given: given{
				id: "",
				payload: `{
					"subject": "test-subject-2",
					"issue_date": "23-01-2023",
					"due_date": "23-01-2024",
					"sub_total": 300,
					"tax": 10,
					"grand_total": 200,
					"customer_request": {
						"customer_name": "test-customer-name-2",
						"address": "test-address-2"
					},
					"item_request": [
						{
							"name": "test-1",
							"type": "test-type",
							"quantity": 1,
							"unit_price": 1,
							"amount": 1
						},
						{
							"name": "test-2",
							"type": "test-type",
							"quantity": 2,
							"unit_price": 2,
							"amount": 2
						}
					]
				}`,
				svcErrReturn: nil,
			},
			expected: expected{
				request:      &request,
				statusCode:   200,
				responseBody: `{"data":{"invoice_id":""},"error":null,"success":true,"metadata":{"request_id":""}}`,
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, "/just/for/testing", strings.NewReader(testCase.given.payload))
			w := httptest.NewRecorder()

			dataFromService := contract.InvcResponse{}

			mockInvoiceSvc := mock_handler.NewMockInvoiceService(mockCtrl)

			if testCase.expected.request != nil {
				mockInvoiceSvc.EXPECT().Update(gomock.Any(), *testCase.expected.request, testCase.given.id).
					Return(dataFromService, testCase.given.svcErrReturn).
					Times(1)
			}

			hf := http.HandlerFunc(UpdateInvoiceHandler(mockInvoiceSvc))
			hf.ServeHTTP(w, r)

			res := w.Result()
			data, err := io.ReadAll(res.Body)
			if err != nil {
				t.Errorf("expected error to be nil got %v", err)
			}

			assert.Equal(t, testCase.expected.statusCode, res.StatusCode)
			assert.Equal(t, fmt.Sprintf("%s\n", testCase.expected.responseBody), string(data))
		})
	}
}
