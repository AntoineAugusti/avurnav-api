package main

import (
	"crypto/tls"
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
)

var (
	buildTag  = "dev"
	buildDate = "0001-01-01T00:00:00Z"
)

// AVURNAVsController constructs a web controller to deal with AVURNAVs
func AVURNAVsController(storage *avurnav.Storage) service.WebController {
	wc := service.NewWebController("/avurnavs/regions/{region}")
	validator := avurnavController{
		storage: storage,
	}

	n := negroni.New()
	n.UseHandler(http.HandlerFunc(validator.AVURNAVsRegionController))

	wc.AddMethodHandler(service.Get, n.ServeHTTP)

	return wc
}

type avurnavController struct {
	storage *avurnav.Storage
}

// AVURNAVsRegionController lists AVURNAVs for a specific region
func (c *avurnavController) AVURNAVsRegionController(w http.ResponseWriter, req *http.Request) {
	res := c.storage.AVURNAVsForRegion(mux.Vars(req)["region"])
	w.Header().Set("Access-Control-Allow-Origin", "*")

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

// NewRedis builds a Redis client from a Redis URL
func NewRedis(url string) *redis.Client {
	redisOptions, err := redis.ParseURL(url)
	if err != nil {
		panic(err)
	}

	redisClient := redis.NewClient(redisOptions)
	if err := redisClient.Ping().Err(); err != nil {
		panic(err)
	}
	return redisClient
}

func main() {
	service.BuildTag = buildTag
	service.BuildDate = buildDate
	service.VersionRoute = "/version"
	service.HeartbeatRoute = "/heartbeat"

	defaultHTTPPort := ":" + strings.TrimLeft(envWithFallback("PORT", "8080"), ":")
	defaultRedisURL := envWithFallback("REDIS_URL", "redis://localhost:6379")

	address := flag.String("a", defaultHTTPPort, "HTTP address to listen to")
	redisURL := flag.String("redis-url", defaultRedisURL, "Redis URL to connect to")
	flag.Parse()

	storage := avurnav.NewStorage(NewRedis(*redisURL))

	go func() {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client := http.Client{
			Timeout:   time.Duration(2 * time.Second),
			Transport: tr,
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
