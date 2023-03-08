package http_encoder

import (
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/transport/http"
	httpNet "net/http"
	"time"
)

type response struct {
	Code   int         `json:"code"`
	Data   interface{} `json:"data"`
	Reason string      `json:"reason,omitempty"`
	Msg    string      `json:"msg"`
	Ts     string      `json:"ts"`
}

func PeroResponseEncoder(w httpNet.ResponseWriter, r *httpNet.Request, v interface{}) error {
	if v == nil {
		return nil
	}
	if rd, ok := v.(http.Redirector); ok {
		url, code := rd.Redirect()
		httpNet.Redirect(w, r, url, code)
		return nil
	}
	reply := &response{
		Code: 0,
		Data: v,
		Msg:  "ok",
		Ts:   time.Now().String(),
	}
	codec, _ := http.CodecForRequest(r, "Accept")
	data, err := codec.Marshal(reply)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func PeroErrorResponseEncoder(w httpNet.ResponseWriter, r *httpNet.Request, err error) {
	se := errors.FromError(err)
	codec, _ := http.CodecForRequest(r, "Accept")
	reply := &response{
		Code:   int(se.Code),
		Reason: se.Reason,
		Msg:    se.Message,
		Ts:     time.Now().String(),
	}
	body, err := codec.Marshal(reply)
	if err != nil {
		w.WriteHeader(httpNet.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if 0 < int(se.Code) && int(se.Code) <= 600 {
		w.WriteHeader(int(se.Code))
	} else {
		w.WriteHeader(200)
	}

	_, _ = w.Write(body)
}
