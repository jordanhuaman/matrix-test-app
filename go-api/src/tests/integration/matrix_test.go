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

func TestCreateQR_Success(t *testing.T) {
	db := setupDB(t)
	app := testutil.NewTestApp(db, jwtSec)

	user := testutil.SeedUser(t, db, "matrix@test.com", "pass123")
	token := testutil.GenerateToken(jwtSec, user.ID.String())

	body := `{"data":[[1,2],[3,4],[5,6]]}`
	req, _ := http.NewRequest("POST", "/api/matrix/qr", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var env envelope
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&env))
	assert.Equal(t, "success", env.Status)

	var data map[string]interface{}
	require.NoError(t, json.Unmarshal(env.Data, &data))
	assert.NotEmpty(t, data["id"])
	assert.Equal(t, "completed", data["status"])
	require.NotNil(t, data["qrResult"])
	require.NotNil(t, data["statistics"])
}

func TestCreateQR_EmptyMatrix(t *testing.T) {
	db := setupDB(t)
	app := testutil.NewTestApp(db, jwtSec)

	user := testutil.SeedUser(t, db, "empty@test.com", "pass123")
	token := testutil.GenerateToken(jwtSec, user.ID.String())

	body := `{"data":[]}`
	req, _ := http.NewRequest("POST", "/api/matrix/qr", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestCreateQR_NonRectangular(t *testing.T) {
	db := setupDB(t)
	app := testutil.NewTestApp(db, jwtSec)

	user := testutil.SeedUser(t, db, "jagged@test.com", "pass123")
	token := testutil.GenerateToken(jwtSec, user.ID.String())

	body := `{"data":[[1,2],[3,4,5]]}`
	req, _ := http.NewRequest("POST", "/api/matrix/qr", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	var env envelope
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&env))
	assert.Equal(t, "error", env.Status)
}

func TestCreateQR_Unauthorized(t *testing.T) {
	db := setupDB(t)
	app := testutil.NewTestApp(db, jwtSec)

	body := `{"data":[[1,2]]}`
	req, _ := http.NewRequest("POST", "/api/matrix/qr", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestListResults(t *testing.T) {
	db := setupDB(t)
	app := testutil.NewTestApp(db, jwtSec)

	user := testutil.SeedUser(t, db, "list@test.com", "pass123")
	token := testutil.GenerateToken(jwtSec, user.ID.String())
	testutil.SeedMatrixResult(t, db, user.ID)

	req, _ := http.NewRequest("GET", "/api/matrix/qr?page=1&limit=20", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var env envelope
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&env))

	var data map[string]interface{}
	require.NoError(t, json.Unmarshal(env.Data, &data))
	assert.NotNil(t, data["results"])
	assert.NotNil(t, data["total"])
}

func TestGetResultByID(t *testing.T) {
	db := setupDB(t)
	app := testutil.NewTestApp(db, jwtSec)

	user := testutil.SeedUser(t, db, "getbyid@test.com", "pass123")
	token := testutil.GenerateToken(jwtSec, user.ID.String())
	result := testutil.SeedMatrixResult(t, db, user.ID)

	req, _ := http.NewRequest("GET", "/api/matrix/qr/"+result.ID.String(), nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var env envelope
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&env))
	assert.Equal(t, "success", env.Status)

	var data map[string]interface{}
	require.NoError(t, json.Unmarshal(env.Data, &data))
	assert.Equal(t, result.ID.String(), data["id"])
	assert.NotNil(t, data["qrResult"])
	assert.NotNil(t, data["statistics"])
}

func TestGetResultByID_NotFound(t *testing.T) {
	db := setupDB(t)
	app := testutil.NewTestApp(db, jwtSec)

	user := testutil.SeedUser(t, db, "notfound@test.com", "pass123")
	token := testutil.GenerateToken(jwtSec, user.ID.String())

	req, _ := http.NewRequest("GET", "/api/matrix/qr/00000000-0000-0000-0000-000000000000", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}
