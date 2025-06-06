package proxy

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/vysotskaya-a/balancer/internal/cache"
	"go.uber.org/zap"
)

type ProxyHandler struct {
	balancer *Balancer
	cache    *cache.RedisCache
	logger   *zap.Logger
}

func NewProxyHandler(b *Balancer, c *cache.RedisCache, logger *zap.Logger) *ProxyHandler {
	return &ProxyHandler{
		balancer: b,
		cache:    c,
		logger:   logger,
	}
}

func (h *ProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	key := "cache:" + r.URL.Path

	// Пробуем достать из кеша
	cached, err := h.cache.Get(key)
	if err == nil {
		h.logger.Info("Ответ из кеша", zap.String("path", r.URL.Path))
		w.Header().Set("X-Cache", "HIT")
		w.Write([]byte(cached))
		return
	}

	target := h.balancer.NextTarget()

	proxy := httputil.NewSingleHostReverseProxy(target)
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		req.Host = target.Host
	}

	proxy.ModifyResponse = func(resp *http.Response) error {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		// Кешируем на 30 секунд
		err = h.cache.Set(key, string(bodyBytes), 30*time.Second)
		if err != nil {
			h.logger.Error("Ошибка записи в кеш", zap.Error(err))
		}
		resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		return nil
	}

	h.logger.Info("Проксируем запрос", zap.String("to", target.String()), zap.String("path", r.URL.Path))
	w.Header().Set("X-Cache", "MISS")
	proxy.ServeHTTP(w, r)
}
