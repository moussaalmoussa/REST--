package model

// model/subscription.go
type Subscription struct {
    ID          string  `json:"id" db:"id"`
    ServiceName string  `json:"service_name" db:"service_name"`
    Price       int     `json:"price" db:"price"`
    UserID      string  `json:"user_id" db:"user_id"`
    StartDate   string  `json:"start_date" db:"start_date"`   // формат "MM-YYYY"
    EndDate     *string `json:"end_date,omitempty" db:"end_date"`
}
