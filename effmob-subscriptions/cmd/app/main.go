package main

import (
    _ "effmob-subscriptions/docs" 
    httpSwagger "github.com/swaggo/http-swagger"
    "effmob-subscriptions/internal/config"
    "effmob-subscriptions/internal/database"
    "effmob-subscriptions/internal/handler"
    "effmob-subscriptions/internal/logger"
    "github.com/gorilla/mux"
    "net/http"
    "os"
)

// @title          Subscription Aggregator API
// @version        1.0
// @description    Сервис управления подписками пользователей.

// @host           localhost:8080
// @BasePath       /

func main() {
    log := logger.New()
    cfg, err := config.Load()
    if err != nil {
        log.Fatalf("config error: %v", err)
    }
    db, err := database.Connect(cfg.DB)
    if err != nil {
        log.Fatalf("db connection: %v", err)
    }
    if err := database.RunMigrations(db); err != nil {
        log.Fatalf("migrations: %v", err)
    }
    r := mux.NewRouter()
    handler.RegisterHandlers(r, db, log)
    log.Infof("starting server on %s", cfg.ServerAddr)
    if err := http.ListenAndServe(cfg.ServerAddr, r); err != nil {
        log.Fatal(err)
    }
}
