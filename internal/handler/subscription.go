package handler

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "net/http"
    "strings"

    "effmob-subscriptions/internal/model"
    "github.com/google/uuid"
    "github.com/sirupsen/logrus"
)

// CreateSubscription создаёт новую подписку.
// @Summary      Создать подписку
// @Description  Создаёт запись о подписке пользователя.
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        body body model.Subscription true "Данные подписки"
// @Success      201 {object} model.Subscription
// @Failure      400 {object} map[string]string
// @Router       /subscriptions [post]
func (h *SubscriptionHandler) Create(w http.ResponseWriter, r *http.Request) {
    var req model.Subscription
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, `{"error":"invalid JSON"}`, http.StatusBadRequest)
        return
    }
    startDate, err := ParseMonthYear(req.StartDate)
    if err != nil {
        http.Error(w, fmt.Sprintf(`{"error":"%v"}`, err), http.StatusBadRequest)
        return
    }
    id := uuid.New().String()
    var endDate *string
    if req.EndDate != nil && *req.EndDate != "" {
        endDateStr, err := ParseMonthYear(*req.EndDate)
        if err != nil {
            http.Error(w, fmt.Sprintf(`{"error":"%v"}`, err), http.StatusBadRequest)
            return
        }
        endDate = &endDateStr
    }
    _, err = h.db.Exec(`INSERT INTO subscriptions (id, service_name, price, user_id, start_date, end_date)
        VALUES ($1, $2, $3, $4, $5, $6)`, id, req.ServiceName, req.Price, req.UserID, startDate, endDate)
    if err != nil {
        h.log.Errorf("insert subscription: %v", err)
        http.Error(w, `{"error":"internal error"}`, http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(req) // или вернуть созданный объект с id
}

// Aggregate получает суммарную стоимость подписок за период.
// @Summary      Суммарная стоимость
// @Description  Подсчёт суммарной стоимости подписок с фильтрацией.
// @Tags         subscriptions
// @Produce      json
// @Param        user_id       query string false "ID пользователя"
// @Param        service_name  query string false "Название сервиса"
// @Param        period_start  query string false "Начало периода (MM-YYYY)"
// @Param        period_end    query string false "Конец периода (MM-YYYY)"
// @Success      200 {object} map[string]int
// @Router       /subscriptions/aggregate [get]
func (h *SubscriptionHandler) Aggregate(w http.ResponseWriter, r *http.Request) {
    userID := r.URL.Query().Get("user_id")
    serviceName := r.URL.Query().Get("service_name")
    periodStart := r.URL.Query().Get("period_start")
    periodEnd := r.URL.Query().Get("period_end")

    query := `SELECT COALESCE(SUM(price), 0) FROM subscriptions WHERE 1=1`
    args := []interface{}{}
    argIdx := 1

    if userID != "" {
        query += fmt.Sprintf(" AND user_id = $%d", argIdx)
        args = append(args, userID)
        argIdx++
    }
    if serviceName != "" {
        query += fmt.Sprintf(" AND service_name ILIKE $%d", argIdx)
        args = append(args, "%"+serviceName+"%")
        argIdx++
    }
    if periodStart != "" {
        startDate, err := ParseMonthYear(periodStart)
        if err == nil {
            query += fmt.Sprintf(" AND (end_date IS NULL OR end_date >= $%d)", argIdx)
            args = append(args, startDate)
            argIdx++
        }
    }
    if periodEnd != "" {
        endDate, err := ParseMonthYear(periodEnd)
        if err == nil {
            // последний день месяца
            endOfMonth := endDate[:8] + fmt.Sprintf("%02d", daysInMonth(endDate))
            query += fmt.Sprintf(" AND start_date <= $%d", argIdx)
            args = append(args, endOfMonth)
            argIdx++
        }
    }

    var total int
    err := h.db.QueryRow(query, args...).Scan(&total)
    if err != nil {
        h.log.Errorf("aggregate query: %v", err)
        http.Error(w, `{"error":"internal error"}`, http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]int{"total": total})
}

