package usecase

import (
	"fmt"
	"helprepet/config"
	"helprepet/internal/common"
	"helprepet/internal/models/user"
	"helprepet/internal/repository"
	"helprepet/pkg/secure"
	"helprepet/pkg/storage"
	"net/smtp"
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
	userData, err := u.userRepo.GetUserByUsername(in.Username)
	if err != nil {
		return nil, common.ErrUnknown.Wrap(err)
	}
	if userData == nil {
		return nil, common.ErrUnauthorized
	}
	if userData.Password == secure.CalcHash(u.auth.Salt, in.Password) {
		access := secure.CalcHash(u.auth.Salt, fmt.Sprintf("access:%s:%v", in.Username, time.Now()))
		refresh := secure.CalcHash(u.auth.Salt, fmt.Sprintf("refresh:%s:%v", in.Username, time.Now()))
		err = u.redis.SetKey(fmt.Sprintf("access_token:%v", access), userData.Username,
			time.Hour*time.Duration(u.auth.AccessTokenLifetime))
		if err != nil {
			return nil, err
		}
		err = u.redis.SetKey(fmt.Sprintf("refresh_token:%v", refresh), userData.Username,
			time.Hour*time.Duration(u.auth.RefreshTokenLifetime))
		if err != nil {
			return nil, err
		}
		return &user.AuthorizeOutput{
			AccessToken:  access,
			RefreshToken: refresh,
		}, nil
	} else {
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
	passwordHash := secure.CalcHash(u.auth.Salt, in.Password)
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
	auth := smtp.PlainAuth("", u.mail.Email, u.mail.Password, u.mail.SMTPHost)
	err = smtp.SendMail(fmt.Sprintf("%s:%s", u.mail.SMTPHost, u.mail.SMTPPort), auth, u.mail.Email, []string{in.Username}, []byte("Subject: Sign up on helprepet\n\nBody of mail"))
	if err != nil {
		fmt.Println(err, auth, *u.mail, fmt.Sprintf("%s:%s", u.mail.SMTPHost, u.mail.SMTPPort))
		return nil, err
	}
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
