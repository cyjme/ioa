package cache

import (
	"errors"
	"ioa/context"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type filter struct {
	Path string
	TTL  int64
}

func New(arg string) (*filter, error) {
	filter := filter{}
	ttl, err := strconv.ParseInt(arg, 10, 64)
	if err != nil {
		return &filter, err
	}
	if ttl < 0 {
		return &filter, errors.New("ttl must max 0")
	}
	filter.TTL = ttl
	return &filter, nil
}

func (f *filter) Name() string {
	return "Cache"
}

func (f *filter) Request(ctx *context.Context) error {
	if ctx.Request.Method != "GET" {
		// cache only valid for GET
		return nil
	}
	f.Path = ctx.RouteUri + ctx.Request.RequestURI
	if f.TTL > 0 {
		if _, ok := hub.caches[f.Path]; ok {
			c := hub.caches[f.Path]
			if time.Now().Unix()-c.CreateAt.Unix() < f.TTL {
				ctx.Response = c.Response
				ctx.ResponseBody = c.ResponseBody
			}
		}
	}
	return nil
}

func (f *filter) Response(ctx *context.Context) error {
	if ctx.Request.Method != "GET" {
		return nil
	}

	if f.TTL > 0 {
		hub.mux.Lock()
		if _, ok := hub.caches[f.Path]; ok {
			c := hub.caches[f.Path]
			if time.Now().Unix()-c.CreateAt.Unix() > f.TTL {
				//cache is invalid and update the cache
				c.CreateAt = time.Now()
				c.Response = ctx.Response
				c.ResponseBody = ctx.ResponseBody
				hub.caches[f.Path] = c
			}

		} else {
			// init cache
			c := cache{
				CreateAt:     time.Now(),
				Response:     ctx.Response,
				ResponseBody: ctx.ResponseBody,
			}
			hub.caches[f.Path] = c
		}
		hub.mux.Unlock()

	}
	return nil
}

type cache struct {
	ResponseBody []byte
	Response     *http.Response
	CreateAt     time.Time
}

type cacheHub struct {
	caches map[string]cache
	mux    sync.RWMutex
}

var hub = cacheHub{}

func init() {
	hub.caches = make(map[string]cache)
}
