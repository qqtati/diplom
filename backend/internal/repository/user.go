package repository

import (
	"backend/internal/models/user"
	"backend/pkg/storage"

	"database/sql"
	"fmt"
	"time"
)

type UserRepo struct {
	db storage.Postgres
}

func NewUserRepo(db storage.Postgres) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

func (r *UserRepo) InsertUser(user *user.User) error {
	_, err := r.db.Exec(`INSERT INTO "user"(username, password, role, invite_code, name) VALUES ($1, $2, $3, $4, $5)`,
		user.Username, user.Password, user.Role, user.InviteCode, user.Name)
	return err
}

func (r *UserRepo) GetUserByUsername(username string) (*user.User, error) {
	var data []user.User
	err := r.db.Select(&data, `SELECT * FROM "user" WHERE username=$1`, username)
	if err != nil {
		return nil, fmt.Errorf("error selecting user: %w", err)
	}
	if len(data) == 0 {
		return nil, nil
	}
	return &data[0], nil
}

func (r *UserRepo) GetUserIDByInvite(inviteCode string) (*user.User, error) {
	var data []user.User
	err := r.db.Select(&data, `SELECT * FROM "user" WHERE invite_code=$1`, inviteCode)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, nil
	}
	return &data[0], nil
}

func (r *UserRepo) InsertTeacherStudent(teacherID int64, studentID int64) error {
	fmt.Printf("Inserting teacher-student relation: teacher_id=%d, student_id=%d\n", teacherID, studentID)
	_, err := r.db.Exec(`INSERT INTO teacher_student(teacher_id, student_id) VALUES($1, $2)`, teacherID, studentID)
	if err != nil {
		fmt.Printf("Error inserting teacher-student relation: %v\n", err)
		return err
	}
	fmt.Printf("Successfully inserted teacher-student relation\n")
	return nil
}

func (r *UserRepo) GetStudents(userID int64) ([]user.User, error) {
	var data []user.User
	err := r.db.Select(&data, `SELECT u.* FROM "user" u left join teacher_student ts on ts.student_id=u.id WHERE ts.teacher_id=$1`, userID)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (r *UserRepo) GetTeachers(userID int64) ([]user.User, error) {
	var data []user.User
	err := r.db.Select(&data, `SELECT u.* FROM "user" u left join teacher_student ts on ts.teacher_id=u.id WHERE ts.student_id=$1`, userID)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (r *UserRepo) GetStudentStats(userID int64, dt time.Time) ([]user.UserStats, error) {
	var data []user.UserStats
	err := r.db.Select(&data,
		`SELECT u.name as name, coalesce(sum(ev.price) filter ( where ev.skipped=false ), 0) as total_income, coalesce(count(*) filter ( where ev.skipped=false ), 0) as event_count, coalesce(avg(ev.duration) filter ( where ev.skipped=false ), 0) as avg_duration, coalesce(count(*) filter ( where ev.skipped=true ), 0) as skip_count, coalesce(avg(ev.rating) filter ( where ev.rating > 0 or ev.skipped=true ), 4) as rating FROM teacher_student ts 
         join "event" ev on ev.teacher_id=ts.teacher_id and ev.student_id=ts.student_id	
		 join "user" u on ts.student_id = u.id
         WHERE ts.teacher_id=$1 and ev.start_time < now() and ev.start_time >= $2 group by ts.student_id, u.name`, userID, dt)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (r *UserRepo) GetTeacherStats(userID int64, dt time.Time) ([]user.UserStats, error) {
	var data []user.UserStats
	err := r.db.Select(&data,
		`SELECT u.name as name, coalesce(sum(ev.price) filter ( where ev.skipped=false ), 0) as total_income, coalesce(count(*) filter ( where ev.skipped=false ), 0) as event_count, coalesce(avg(ev.duration) filter ( where ev.skipped=false ), 0) as avg_duration, coalesce(count(*) filter ( where ev.skipped=true ), 0) as skip_count, coalesce(avg(ev.rating) filter ( where ev.rating > 0 or ev.skipped=true ), 4) as rating FROM teacher_student ts 
         join "event" ev on ev.teacher_id=ts.teacher_id and ev.student_id=ts.student_id	
		 join "user" u on ts.teacher_id = u.id
         WHERE ts.student_id=$1 and ev.start_time < now() and ev.start_time >= $2 group by ts.student_id, u.name`, userID, dt)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (r *UserRepo) GetByUsername(username string) (*user.User, error) {
	var u user.User
	err := r.db.QueryRow(
		`SELECT id, username, password, role, invite_code, name, working_hours, created_at, updated_at 
		FROM "user" WHERE username = $1`,
		username,
	).Scan(
		&u.ID,
		&u.Username,
		&u.Password,
		&u.Role,
		&u.InviteCode,
		&u.Name,
		&u.WorkingHours,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error getting user: %w", err)
	}
	return &u, nil
}

func (r *UserRepo) Update(u *user.User) error {
	_, err := r.db.Exec(
		`UPDATE "user" 
		SET name = $1, working_hours = $2, updated_at = CURRENT_TIMESTAMP 
		WHERE id = $3`,
		u.Name,
		u.WorkingHours,
		u.ID,
	)
	if err != nil {
		return fmt.Errorf("error updating user: %w", err)
	}
	return nil
}
