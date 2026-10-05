/*
 _____               _____
|_   _| __ _   _  __|_   _|__ _ __ ___  _ __ ___
  | || '__| | | |/ _ \| |/ _ \ '_ ` _ \| '_ ` _ \
  | || |  | |_| |  __/| |  __/ | | | | | | | | | |
  |_||_|   \__,_|\___||_|\___|_| |_| |_|_| |_| |_|

 _____ __  __       ____                               __ _
|_   _|  \/  |     |  _ \ _ __ __ _  __ _  ___  _ __  / _| |_   _
  | | | |\/| |_____| | | | '__/ _` |/ _` |/ _ \| '_ \| |_| | | | |
  | | | |  | |_____| |_| | | | (_| | (_| | (_) | | | |  _| | |_| |
  |_| |_|  |_|     |____/|_|  \__,_|\__, |\___/|_| |_|_| |_|\__, |
                                    |___/                   |___/

@author TrueTemm
@link   https://github.com/TrueTemm
TM-Dragonfly Project
*/

package server

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const (
	authCacheAge = 30 * 24 * time.Hour

	authTimeout = 45 * time.Second
)

var authClient = &http.Client{
	Timeout: authTimeout,
	Transport: cachingTransport{base: &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		DialContext:           dialIPv4First,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          10,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   15 * time.Second,
		ExpectContinueTimeout: time.Second,
	}},
}

func dialIPv4First(ctx context.Context, network, addr string) (net.Conn, error) {
	d := &net.Dialer{Timeout: 15 * time.Second, KeepAlive: 30 * time.Second}
	if network != "tcp" {
		return d.DialContext(ctx, network, addr)
	}
	conn, err := d.DialContext(ctx, "tcp4", addr)
	if err == nil {
		return conn, nil
	}
	conn6, err6 := d.DialContext(ctx, "tcp6", addr)
	if err6 != nil {
		return nil, err
	}
	return conn6, nil
}

type cachingTransport struct {
	base http.RoundTripper
}

func (t cachingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if req.Method != http.MethodGet {
		return t.base.RoundTrip(req)
	}
	key := req.URL.String()
	resp, err := t.base.RoundTrip(req)
	if err == nil && resp.StatusCode == http.StatusOK {
		body, readErr := io.ReadAll(resp.Body)
		_ = resp.Body.Close()
		if readErr != nil {
			return nil, readErr
		}
		storeAuthAnswer(key, body)
		resp.Body = io.NopCloser(bytes.NewReader(body))
		return resp, nil
	}
	if resp != nil {
		_ = resp.Body.Close()
	}
	body, ok := loadAuthAnswer(key)
	if !ok {
		if err != nil {
			return nil, err
		}
		return t.base.RoundTrip(req)
	}
	return &http.Response{
		Status: "200 OK", StatusCode: http.StatusOK,
		Proto: "HTTP/1.1", ProtoMajor: 1, ProtoMinor: 1,
		Header:        http.Header{"Content-Type": []string{"application/json"}},
		Body:          io.NopCloser(bytes.NewReader(body)),
		ContentLength: int64(len(body)),
		Request:       req,
	}, nil
}

func authAnswerPath(key string) (string, bool) {
	dir, err := os.UserCacheDir()
	if err != nil {
		dir = os.TempDir()
	}
	dir = filepath.Join(dir, "tm-dragonfly", "auth")
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return "", false
	}
	sum := sha256.Sum256([]byte(key))
	return filepath.Join(dir, hex.EncodeToString(sum[:8])+".json"), true
}

func storeAuthAnswer(key string, body []byte) {
	path, ok := authAnswerPath(key)
	if !ok {
		return
	}
	_ = os.WriteFile(path, body, 0o600)
}

func loadAuthAnswer(key string) ([]byte, bool) {
	path, ok := authAnswerPath(key)
	if !ok {
		return nil, false
	}
	info, err := os.Stat(path)
	if err != nil || time.Since(info.ModTime()) > authCacheAge {
		return nil, false
	}
	body, err := os.ReadFile(path)
	if err != nil || len(body) == 0 {
		return nil, false
	}
	return body, true
}
