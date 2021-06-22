package httpserver

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWriteErrorToClient(t *testing.T){
	testcases := []struct{
		caseDesc string
		msg string
		w *httptest.ResponseRecorder
		httpStatus int
		expectedStatus int
		expectedBody string
	}{
		{
			caseDesc: "invalid status should result in StatusInternalServerError status and valid msg",
			msg:            "hello",
			w: httptest.NewRecorder(),
			httpStatus:     1,
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "hello",
		},
		{
			caseDesc: "valid status should result in same status and valid msg",
			msg:            "foo",
			w: httptest.NewRecorder(),
			httpStatus:     http.StatusOK,
			expectedStatus: http.StatusOK,
			expectedBody:   "foo",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.caseDesc, func(t *testing.T) {
			writeErrorToClient(tc.msg, tc.w, tc.httpStatus)
			if tc.expectedBody != tc.w.Body.String() {
				t.Error("Body expectation mismatch")
			}

			if tc.expectedStatus != tc.w.Code {
				t.Error("Header Status expectation mismatch")
			}
		})

	}

}
