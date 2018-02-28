package main

import (
	"flag"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/AntoineAugusti/avurnav"
	"github.com/cloudflare/service"
	"github.com/cloudflare/service/render"
	"github.com/codegangsta/negroni"
	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	"github.com/ulule/limiter"
	"github.com/ulule/limiter/drivers/middleware/stdlib"
	"github.com/ulule/limiter/drivers/store/memory"
)

var (
	buildTag  = "dev"
	buildDate = "0001-01-01T00:00:00Z"
)

func AVURNAVsController(storage *avurnav.Storage) service.WebController {
	wc := service.NewWebController("/avurnavs/regions/{region}")
	validator := avurnavController{
		storage: storage,
	}

	store := memory.NewStore()
	rate := limiter.Rate{
		Period: 1 * time.Minute,
		Limit:  10,
	}

	middleware := stdlib.NewMiddleware(limiter.New(store, rate))

	n := negroni.New()
	n.UseHandler(middleware.Handler(http.HandlerFunc(validator.AVURNAVsRegionController)))

	wc.AddMethodHandler(service.Get, n.ServeHTTP)

	return wc
}

type avurnavController struct {
	storage *avurnav.Storage
}

func (c *avurnavController) AVURNAVsRegionController(w http.ResponseWriter, req *http.Request) {
	res := c.storage.AVURNAVsForRegion(mux.Vars(req)["region"])

	render.JSON(w, http.StatusOK, res)
}

func refreshAVURNAVs(fetcher *avurnav.AVURNAVFetcher, storage avurnav.Storage) {
	avurnavs, _, _ := fetcher.List()
	var res avurnav.AVURNAVs
	for _, avurnav := range avurnavs {
		avurnav, _, _ = fetcher.Get(avurnav)
		res = append(res, avurnav)
	}
	storage.RegisterAVURNAVs(res)
}

func envWithFallback(key, fallback string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	return val
}

func main() {
	service.BuildTag = buildTag
	service.BuildDate = buildDate
	service.VersionRoute = "/version"
	service.HeartbeatRoute = "/heartbeat"

	port := ":" + strings.TrimLeft(envWithFallback("PORT", "8080"), ":")

	address := flag.String("a", port, "HTTP address to listen to")
	redisURL := flag.String("redis-url", envWithFallback("REDIS_URL", ":6379"), "Redis URL to connect to")
	flag.Parse()

	storage := avurnav.NewStorage(redis.NewClient(&redis.Options{
		Addr: *redisURL,
	}))

	go func() {
		client := http.Client{
			Timeout: time.Duration(2 * time.Second),
		}
		for _ = range time.NewTicker(60 * time.Second).C {
			for _, fetcher := range avurnav.NewClient(&client).Fetchers {
				go refreshAVURNAVs(fetcher, storage)
			}
		}
	}()

	ws := service.NewWebService()
	ws.AddWebController(AVURNAVsController(&storage))
	ws.Run(*address)
}
