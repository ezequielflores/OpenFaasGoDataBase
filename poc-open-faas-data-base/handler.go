package function

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"handler/function/internal/adapter/chache"
	cacheModel "handler/function/internal/adapter/chache/model"
	"handler/function/internal/adapter/controller"
	"handler/function/internal/adapter/respository"
	"handler/function/internal/application/usecase/starwar"
	"log"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"syscall"
	"time"
)

var routes = map[string]func(http.ResponseWriter, *http.Request){}

func init() {
	go func() {

		databaseUrl := os.Getenv("database_url")
		dataBasePool, err := pgxpool.New(context.Background(), databaseUrl)
		if err != nil {
			log.Printf("Data base connection error %s\n", err.Error())
			return
		}

		repositoryAdapter, dataBaseErr := respository.NewStarwarRepositoryAdapter(dataBasePool)
		if dataBaseErr != nil {
			log.Println("Error initializing repository adapter")
			return
		}

		cacheUrl := os.Getenv("cache_url")
		cachePassword := os.Getenv("cache_password")
		ttl := os.Getenv("cache_ttl")

		ttlMillisecond, ttlErr := strconv.Atoi(ttl)

		if ttlErr != nil {
			log.Printf("Error parsing cache ttl %s\n", ttlErr.Error())
			return
		}

		cacheOptions := &cacheModel.CacheOptions{
			Ttl: time.Duration(ttlMillisecond) * time.Millisecond,
		}

		redisClient := redis.NewClient(&redis.Options{
			Addr:      cacheUrl,
			Password:  cachePassword,
			TLSConfig: &tls.Config{},
		})

		redisAdapter, cacheErr := chache.NewStarwarRedisAdapter(redisClient, cacheOptions)

		if cacheErr != nil {
			log.Println("Error initializing cache adapter")
			return
		}

		createCharacter := starwar.NewCreateCharacterUseCase(repositoryAdapter)
		findCharacter := starwar.NewFindCharacterUseCase(repositoryAdapter, redisAdapter)

		starwarsController := controller.NewStarWarController(createCharacter, findCharacter)

		routes["/characters"] = starwarsController.CreateStarWarCharacter
		routes["/characters/"] = starwarsController.FindStarWarCharacter

		channel := make(chan os.Signal, 2)
		signal.Notify(channel, syscall.SIGINT, syscall.SIGTERM)
		log.Println("Server running")

		<-channel
		dataBasePool.Close()

		log.Println("Server shutdown")
	}()
}

func Handle(w http.ResponseWriter, r *http.Request) {

	isGetCharacterDetail, _ := regexp.MatchString("^/api/v1/starwar/characters/[0-9]+$", r.URL.Path)
	if r.Method == http.MethodGet && isGetCharacterDetail {
		routes["/characters/"](w, r)
		return
	}

	isCreateCharacter, _ := regexp.MatchString("^/api/v1/starwar/characters$", r.URL.Path)

	if r.Method == http.MethodPost && isCreateCharacter {
		routes["/characters"](w, r)
		return
	}

	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(fmt.Sprintf("Url: %s %s", r.URL.Path, http.StatusText(http.StatusNotFound))))
}
