package handler

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "net/http"
    // другие импорты (uuid, ваш логгер, model...)
)

type SubscriptionHandler struct {
    db  *sql.DB
    log *logrus.Logger
}

func NewSubscriptionHandler(db *sql.DB, log *logrus.Logger) *SubscriptionHandler {
    return &SubscriptionHandler{db: db, log: log}
}

// ===== CRUDL =====
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
    func (h *SubscriptionHandler) Create(w http.ResponseWriter, r *http.Request) {
    var req model.Subscription
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        // ...
    }
    // валидация даты, преобразование "MM-YYYY" -> date "YYYY-MM-01"
    startDate, err := ParseMonthYear(req.StartDate)
    // ...
    id := uuid.New().String()
    _, err = h.db.Exec(`INSERT INTO subscriptions (id, service_name, price, user_id, start_date, end_date)
        VALUES ($1, $2, $3, $4, $5, $6)`, id, req.ServiceName, req.Price, req.UserID, startDate, endDate)
    // ...
}
}

func (h *SubscriptionHandler) Get(w http.ResponseWriter, r *http.Request) { ... }
func (h *SubscriptionHandler) Update(w http.ResponseWriter, r *http.Request) { ... }
func (h *SubscriptionHandler) Delete(w http.ResponseWriter, r *http.Request) { ... }
func (h *SubscriptionHandler) List(w http.ResponseWriter, r *http.Request) { ... }

// ===== Агрегация =====

func (h *SubscriptionHandler) Aggregate(w http.ResponseWriter, r *http.Request) {
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
        startDate := parseMonthYear(periodStart) // YYYY-MM-01
        query += fmt.Sprintf(" AND (end_date IS NULL OR end_date >= $%d)", argIdx)
        args = append(args, startDate)
        argIdx++
    }
    if periodEnd != "" {
        endDate := parseMonthYear(periodEnd)
        // если период — по месяц включительно, конец периода = последний день месяца
        endOfMonth := endDate.AddDate(0, 1, -1)
        query += fmt.Sprintf(" AND start_date <= $%d", argIdx)
        args = append(args, endOfMonth)
        argIdx++
    }

    var total int
    err := h.db.QueryRow(query, args...).Scan(&total)
    // ...
}
}
