package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/k20ku/see/entity"
	"github.com/k20ku/see/store"
	"github.com/k20ku/see/testutil"
)

func TestAddItem(t *testing.T) {
	t.Parallel()

	datadir := "testdata/add_item/"
	type want struct {
		status  int
		rspFile string
	}

	tests := map[string]struct {
		reqFile string
		want    want
	}{
		"ok": {
			reqFile: datadir + "ok_req.json.golden",
			want: want{
				status:  http.StatusOK,
				rspFile: datadir + "ok_rsp.json.golden",
			},
		},
		"badRequest": {
			reqFile: datadir + "bad_req.json.golden",
			want: want{
				status:  http.StatusBadRequest,
				rspFile: datadir + "bad_rsp.json.golden",
			},
		},
		"invalidJson": {
			reqFile: datadir + "invalid_json_req.json.golden",
			want: want{
				status:  http.StatusBadRequest,
				rspFile: datadir + "invalid_json_rsp.json.golden",
			},
		},
	}

	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			w := httptest.NewRecorder()
			r := httptest.NewRequest(
				http.MethodPost,
				"/task",
				bytes.NewReader(testutil.LoadFile(t, tt.reqFile)),
			)

			validate := validator.New()
			validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
				name, _, _ := strings.Cut(fld.Tag.Get("json"), ",")
				// skip if tag key says it should be ignored
				if name == "-" {
					return ""
				}
				return name
			})

			sut := AddItem{
				Store: &store.ItemStore{
					Items: map[entity.ItemId]*entity.Item{},
				},
				Validate: validate,
			}

			sut.ServeHTTP(w, r)

			resp := w.Result()
			testutil.AssertResponse(t,
				resp, tt.want.status, testutil.LoadFile(t, tt.want.rspFile),
			)
		})
	}
}
