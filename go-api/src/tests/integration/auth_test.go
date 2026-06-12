package integration_test

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/jordanhuaman/go-api/src/tests/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type envelope struct {
	Status  string          `json:"status"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

func TestRegister_Success(t *testing.T) {
	db := setupDB(t)
	app := testutil.NewTestApp(db, jwtSec)

	body := `{"email":"new@test.com","username":"newuser","password":"pass123"}`
	req, _ := http.NewRequest("POST", "/api/auth/register", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var env envelope
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&env))
	assert.Equal(t, "success", env.Status)

	var data map[string]string
	require.NoError(t, json.Unmarshal(env.Data, &data))
	assert.NotEmpty(t, data["id"])
	assert.Equal(t, "new@test.com", data["email"])
}

func TestRegister_DuplicateEmail(t *testing.T) {
	db := setupDB(t)
	app := testutil.NewTestApp(db, jwtSec)

	testutil.SeedUser(t, db, "dup@test.com", "pass123")

	body := `{"email":"dup@test.com","username":"another","password":"pass456"}`
	req, _ := http.NewRequest("POST", "/api/auth/register", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusConflict, resp.StatusCode)

	var env envelope
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&env))
	assert.Equal(t, "error", env.Status)
	assert.Contains(t, env.Message, "Email already in use")
}

func TestLogin_Success(t *testing.T) {
	db := setupDB(t)
	app := testutil.NewTestApp(db, jwtSec)

	testutil.SeedUser(t, db, "login@test.com", "mypassword")

	body := `{"email":"login@test.com","password":"mypassword"}`
	req, _ := http.NewRequest("POST", "/api/auth/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var env envelope
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&env))
	assert.Equal(t, "success", env.Status)

	var data map[string]string
	require.NoError(t, json.Unmarshal(env.Data, &data))
	assert.NotEmpty(t, data["access_token"])
	assert.NotEmpty(t, data["refresh_token"])
}

func TestLogin_InvalidPassword(t *testing.T) {
	db := setupDB(t)
	app := testutil.NewTestApp(db, jwtSec)

	testutil.SeedUser(t, db, "secure@test.com", "correctpass")

	body := `{"email":"secure@test.com","password":"wrongpass"}`
	req, _ := http.NewRequest("POST", "/api/auth/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestRefreshToken_Success(t *testing.T) {
	db := setupDB(t)
	app := testutil.NewTestApp(db, jwtSec)

	user := testutil.SeedUser(t, db, "refresh@test.com", "pass123")
	tok := testutil.SeedRefreshToken(t, db, user.ID.String(), longTTL)

	body := `{"refresh_token":"` + tok.Token + `"}`
	req, _ := http.NewRequest("POST", "/api/auth/refresh-token", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var env envelope
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&env))
	assert.Equal(t, "success", env.Status)

	var data map[string]string
	require.NoError(t, json.Unmarshal(env.Data, &data))
	assert.NotEmpty(t, data["token"])
	assert.NotEmpty(t, data["refresh_token"])
	assert.NotEqual(t, tok.Token, data["refresh_token"], "refresh token should be rotated")
}

func TestRefreshToken_Revoked(t *testing.T) {
	db := setupDB(t)
	app := testutil.NewTestApp(db, jwtSec)

	user := testutil.SeedUser(t, db, "revoked@test.com", "pass123")
	tok := testutil.SeedRefreshToken(t, db, user.ID.String(), longTTL)
	db.Model(&tok).Update("revoked", true)

	body := `{"refresh_token":"` + tok.Token + `"}`
	req, _ := http.NewRequest("POST", "/api/auth/refresh-token", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestLogout(t *testing.T) {
	db := setupDB(t)
	app := testutil.NewTestApp(db, jwtSec)

	user := testutil.SeedUser(t, db, "logout@test.com", "pass123")
	token := testutil.GenerateToken(jwtSec, user.ID.String())

	req, _ := http.NewRequest("POST", "/api/auth/logout", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var env envelope
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&env))
	assert.Equal(t, "success", env.Status)
	assert.Contains(t, env.Message, "Success logout")
}
