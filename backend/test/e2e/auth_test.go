package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type SignUpResponse struct {
	Status string `json:"status"`
	Result struct {
		Success bool `json:"success"`
	} `json:"result"`
}

type SignInResponse struct {
	Status string `json:"status"`
	Result struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	} `json:"result"`
}

func TestSignUpAndSignIn(t *testing.T) {
	// Тест регистрации
	username := fmt.Sprintf("testuser_%d", time.Now().Unix())
	password := "password123"
	name := "Test User"
	role := 0 // student

	// Регистрация
	signUpData := map[string]interface{}{
		"username": username,
		"password": password,
		"name":     name,
		"role":     role,
	}

	jsonData, err := json.Marshal(signUpData)
	assert.NoError(t, err)

	// Логируем запрос
	t.Logf("Sign up request: %s", string(jsonData))

	resp, err := http.Post(fmt.Sprintf("%s/user/sign_up", BaseURL), "application/json", bytes.NewBuffer(jsonData))
	assert.NoError(t, err)

	// Логируем ответ
	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	t.Logf("Sign up response status: %d, body: %s", resp.StatusCode, string(body))

	assert.Equal(t, 200, resp.StatusCode)

	var signUpResponse SignUpResponse
	err = json.NewDecoder(bytes.NewReader(body)).Decode(&signUpResponse)
	assert.NoError(t, err)
	assert.True(t, signUpResponse.Result.Success)

	// Ждем немного, чтобы убедиться, что данные сохранились
	time.Sleep(1 * time.Second)

	// Авторизация
	signInData := map[string]interface{}{
		"username": username,
		"password": password,
	}

	jsonData, err = json.Marshal(signInData)
	assert.NoError(t, err)

	// Логируем запрос
	t.Logf("Sign in request: %s", string(jsonData))

	resp, err = http.Post(fmt.Sprintf("%s/user/sign_in", BaseURL), "application/json", bytes.NewBuffer(jsonData))
	assert.NoError(t, err)

	// Логируем ответ
	body, err = io.ReadAll(resp.Body)
	assert.NoError(t, err)
	t.Logf("Sign in response status: %d, body: %s", resp.StatusCode, string(body))

	assert.Equal(t, 200, resp.StatusCode)

	var signInResponse SignInResponse
	err = json.NewDecoder(bytes.NewReader(body)).Decode(&signInResponse)
	assert.NoError(t, err)
	assert.NotEmpty(t, signInResponse.Result.AccessToken)
	assert.NotEmpty(t, signInResponse.Result.RefreshToken)
}

func registerUser(t *testing.T, username, password, name string, role int) {
	signUpData := map[string]interface{}{
		"username": username,
		"password": password,
		"name":     name,
		"role":     role,
	}

	jsonData, err := json.Marshal(signUpData)
	assert.NoError(t, err)

	// Логируем запрос
	t.Logf("Sign up request: %s", string(jsonData))

	resp, err := http.Post(fmt.Sprintf("%s/user/sign_up", BaseURL), "application/json", bytes.NewBuffer(jsonData))
	assert.NoError(t, err)

	// Логируем ответ
	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	t.Logf("Sign up response status: %d, body: %s", resp.StatusCode, string(body))

	assert.Equal(t, 200, resp.StatusCode)

	var signUpResponse SignUpResponse
	err = json.NewDecoder(bytes.NewReader(body)).Decode(&signUpResponse)
	assert.NoError(t, err)
	assert.True(t, signUpResponse.Result.Success)

	// Ждем немного, чтобы убедиться, что данные сохранились
	time.Sleep(1 * time.Second)
}

func loginUser(t *testing.T, username, password string) string {
	signInData := map[string]interface{}{
		"username": username,
		"password": password,
	}

	jsonData, err := json.Marshal(signInData)
	assert.NoError(t, err)

	// Логируем запрос
	t.Logf("Login request: %s", string(jsonData))

	resp, err := http.Post(fmt.Sprintf("%s/user/sign_in", BaseURL), "application/json", bytes.NewBuffer(jsonData))
	assert.NoError(t, err)

	// Логируем ответ
	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	t.Logf("Login response status: %d, body: %s", resp.StatusCode, string(body))

	assert.Equal(t, 200, resp.StatusCode)

	var signInResponse SignInResponse
	err = json.NewDecoder(bytes.NewReader(body)).Decode(&signInResponse)
	assert.NoError(t, err)
	assert.NotEmpty(t, signInResponse.Result.AccessToken)
	assert.NotEmpty(t, signInResponse.Result.RefreshToken)

	return signInResponse.Result.AccessToken
}
