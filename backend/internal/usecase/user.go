package usecase

import (
	"backend/config"
	"backend/internal/common"
	"backend/internal/models/user"
	"backend/internal/repository"
	"backend/pkg/secure"
	"backend/pkg/storage"
	"context"
	"fmt"
	"strings"
	"time"
)

type UserUC struct {
	userRepo *repository.UserRepo
	redis    *storage.Redis
	auth     *config.AuthConfig
	mail     *config.MailConfig
}

func NewUserUC(userRepo *repository.UserRepo, redis *storage.Redis, auth *config.AuthConfig, mail *config.MailConfig) *UserUC {
	return &UserUC{
		userRepo: userRepo,
		redis:    redis,
		auth:     auth,
		mail:     mail,
	}
}

func (u *UserUC) AuthorizeUser(in *user.AuthorizeInput) (*user.AuthorizeOutput, error) {
	fmt.Printf("AuthorizeUser: using salt: %s\n", u.auth.Salt)
	userData, err := u.userRepo.GetUserByUsername(in.Username)
	if err != nil {
		fmt.Printf("Error getting user by username: %v\n", err)
		return nil, common.ErrUnknown.Wrap(err)
	}
	if userData == nil {
		fmt.Printf("User not found: %s\n", in.Username)
		return nil, common.ErrUnauthorized
	}

	providedHash := secure.CalcHash(u.auth.Salt, in.Password)
	fmt.Printf("Provided password hash: %s\n", providedHash)
	fmt.Printf("Stored password hash: %s\n", userData.Password)

	if userData.Password == providedHash {
		access := secure.CalcHash(u.auth.Salt, fmt.Sprintf("access:%s:%v", in.Username, time.Now()))
		refresh := secure.CalcHash(u.auth.Salt, fmt.Sprintf("refresh:%s:%v", in.Username, time.Now()))
		err = u.redis.SetKey(fmt.Sprintf("access_token:%v", access), userData.Username,
			time.Hour*time.Duration(u.auth.AccessTokenLifetime))
		if err != nil {
			fmt.Printf("Error setting access token: %v\n", err)
			return nil, err
		}
		err = u.redis.SetKey(fmt.Sprintf("refresh_token:%v", refresh), userData.Username,
			time.Hour*time.Duration(u.auth.RefreshTokenLifetime))
		if err != nil {
			fmt.Printf("Error setting refresh token: %v\n", err)
			return nil, err
		}
		return &user.AuthorizeOutput{
			AccessToken:  access,
			RefreshToken: refresh,
		}, nil
	} else {
		fmt.Printf("Password mismatch for user: %s\n", in.Username)
		return nil, common.ErrUnauthorized
	}
}

func (u *UserUC) CheckToken(token string) (*string, error) {
	key, err := u.redis.GetKey(fmt.Sprintf("access_token:%v", token))
	if err != nil {
		return nil, err
	}
	if key == nil {
		return nil, nil
	}
	userData, err := u.userRepo.GetUserByUsername(*key)
	if err == nil && userData != nil {
		return &userData.Username, nil
	}
	return nil, err
}

func (u *UserUC) SignUp(in *user.SignUpInput) (*user.SignUpOutput, error) {
	fmt.Printf("SignUp: using salt: %s\n", u.auth.Salt)
	passwordHash := secure.CalcHash(u.auth.Salt, in.Password)
	fmt.Printf("SignUp: password hash: %s\n", passwordHash)
	role := user.RoleStudent
	if in.IsTeacher > 0 {
		role = user.RoleTeacher
	}
	err := u.userRepo.InsertUser(&user.User{
		Username:   in.Username,
		Password:   passwordHash,
		Name:       in.Name,
		InviteCode: strings.ToUpper(secure.CalcHash(u.auth.Salt, in.Username)[0:6]),
		Role:       role,
	})
	if err != nil {
		return nil, err
	}
	// auth := smtp.PlainAuth("", u.mail.Email, u.mail.Password, u.mail.SMTPHost)
	// err = smtp.SendMail(fmt.Sprintf("%s:%s", u.mail.SMTPHost, u.mail.SMTPPort), auth, u.mail.Email, []string{in.Username}, []byte("Subject: Sign up on backend\n\nBody of mail"))
	// if err != nil {
	// 	fmt.Println(err, auth, *u.mail, fmt.Sprintf("%s:%s", u.mail.SMTPHost, u.mail.SMTPPort))
	// 	return nil, err
	// }
	if in.IsTeacher == 0 && in.InviteCode != nil {
		teacher, err := u.userRepo.GetUserIDByInvite(*in.InviteCode)
		if err != nil {
			return nil, err
		}
		student, err := u.GetMe(in.Username)
		if err != nil {
			return nil, err
		}
		err = u.userRepo.InsertTeacherStudent(teacher.ID, student.ID)
	}
	return &user.SignUpOutput{Success: true}, nil
}

func (u *UserUC) GetMe(username string) (*user.User, error) {
	return u.userRepo.GetUserByUsername(username)
}

func (u *UserUC) GetStudents(username string) ([]user.User, error) {
	usr, err := u.userRepo.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}
	res, err := u.userRepo.GetStudents(usr.ID)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return []user.User{}, nil
	}
	return res, nil
}

func (u *UserUC) GetStudentStats(username string, days int) ([]user.UserStats, error) {
	usr, err := u.userRepo.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}
	dt := time.Unix(0, 0)
	if days > 0 {
		dt = time.Now().Add(-time.Duration(days*24) * time.Hour)
	}
	return u.userRepo.GetStudentStats(usr.ID, dt)
}

func (u *UserUC) GetTeachers(username string) ([]user.User, error) {
	usr, err := u.userRepo.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}
	res, err := u.userRepo.GetTeachers(usr.ID)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return []user.User{}, nil
	}
	return res, nil
}

func (u *UserUC) GetTeacherStats(username string, days int) ([]user.UserStats, error) {
	usr, err := u.userRepo.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}
	dt := time.Unix(0, 0)
	if days > 0 {
		dt = time.Now().Add(-time.Duration(days*24) * time.Hour)
	}
	return u.userRepo.GetTeacherStats(usr.ID, dt)
}

func (u *UserUC) CreateTeacherStudentRelation(teacherID uint, studentID uint) error {
	return u.userRepo.InsertTeacherStudent(int64(teacherID), int64(studentID))
}

func (uc *UserUC) GetProfile(username string) (*user.User, error) {
	return uc.userRepo.GetByUsername(username)
}

func (uc *UserUC) UpdateProfile(ctx context.Context, userID int64, input user.UserInput) error {
	u, err := uc.userRepo.GetByUsername(input.Email)
	if err != nil {
		return fmt.Errorf("error getting user: %w", err)
	}
	if u == nil {
		return common.ErrNotFound
	}

	u.Name = input.Name
	u.WorkingHours = input.WorkingHours

	if err := uc.userRepo.Update(u); err != nil {
		return fmt.Errorf("error updating user: %w", err)
	}

	return nil
}
