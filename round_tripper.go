package srt

import (
	"fmt"
	fhttp "github.com/bogdanfinn/fhttp"
	tlsclient "github.com/bogdanfinn/tls-client"
	"net/http"
)

type SpoofedRoundTripper struct {
	Client tlsclient.HttpClient
}

func NewSpoofedRoundTripper(httpClientOption ...tlsclient.HttpClientOption) (*SpoofedRoundTripper, error) {
	c, err := tlsclient.NewHttpClient(tlsclient.NewNoopLogger(), httpClientOption...)
	if err != nil {
		return nil, fmt.Errorf("error creating client: %w", err)
	}
	return &SpoofedRoundTripper{
		Client: c,
	}, nil
}

func (s SpoofedRoundTripper) RoundTrip(hReq *http.Request) (*http.Response, error) {
	fReq, err := fhttp.NewRequest(hReq.Method, hReq.URL.String(), hReq.Body)
	if err != nil {
		return nil, err
	}
	fReq.Header = fhttp.Header(hReq.Header)
	fReq.Trailer = fhttp.Header(hReq.Trailer)
	fReq.Form = hReq.Form
	fReq.MultipartForm = hReq.MultipartForm
	fReq.PostForm = hReq.PostForm

	fResp, err := s.Client.Do(fReq)
	if err != nil {
		return nil, fmt.Errorf("error fetching response: %w", err)
	}
	return &http.Response{
		Status:           fResp.Status,
		StatusCode:       fResp.StatusCode,
		Proto:            fResp.Proto,
		ProtoMajor:       fResp.ProtoMajor,
		ProtoMinor:       fResp.ProtoMinor,
		Header:           http.Header(fResp.Header),
		Body:             fResp.Body,
		ContentLength:    fResp.ContentLength,
		TransferEncoding: fResp.TransferEncoding,
		Close:            fResp.Close,
		Uncompressed:     fResp.Uncompressed,
		Trailer:          http.Header(fResp.Trailer),
		Request:          hReq,
		TLS:              nil,
	}, nil
}
