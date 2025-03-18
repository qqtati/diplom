package user

const (
	RoleTeacher = int64(0)
	RoleStudent = int64(1)
)

type User struct {
	ID         int64  `json:"id" db:"id"`
	Username   string `json:"username" db:"username"`
	Password   string `json:"-" db:"password"`
	Role       int64  `json:"role" db:"role"`
	InviteCode string `json:"invite_code" db:"invite_code"`
	Name       string `json:"name" db:"name"`
}

type UserStats struct {
	Name        string  `json:"name" db:"name"`
	EventCount  int64   `json:"event_count" db:"event_count"`
	TotalIncome int64   `json:"total_income" db:"total_income"`
	AvgDuration float64 `json:"avg_duration" db:"avg_duration"`
	SkipCount   int64   `json:"skip_count" db:"skip_count"`
	Rating      float64 `json:"rating" db:"rating"`
}
