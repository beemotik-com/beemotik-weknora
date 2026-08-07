package router

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestFrontendStaticDoesNotInterceptResourceGrant(t *testing.T) {
	gin.SetMode(gin.TestMode)
	webDir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(webDir, "index.html"), []byte("spa"), 0o600))
	t.Setenv("WEKNORA_WEB_DIR", webDir)

	r := gin.New()
	serveFrontendStatic(r)
	r.GET("/r/:token", func(c *gin.Context) { c.String(http.StatusOK, "resource") })

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/r/token", nil)
	r.ServeHTTP(recorder, request)

	require.Equal(t, http.StatusOK, recorder.Code)
	require.Equal(t, "resource", recorder.Body.String())
}

func TestServeFrontendStaticStandardEditionWithWebDir(t *testing.T) {
	gin.SetMode(gin.TestMode)
	webDir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(webDir, "index.html"), []byte("<html>WeKnora SPA</html>"), 0o600))
	t.Setenv("WEKNORA_WEB_DIR", webDir)

	r := gin.New()
	serveFrontendStatic(r)

	// Test GET /
	recRoot := httptest.NewRecorder()
	reqRoot := httptest.NewRequest(http.MethodGet, "/", nil)
	r.ServeHTTP(recRoot, reqRoot)

	require.Equal(t, http.StatusOK, recRoot.Code)
	require.Contains(t, recRoot.Body.String(), "<html>WeKnora SPA</html>")

	// Test GET /dashboard (SPA fallback)
	recSPA := httptest.NewRecorder()
	reqSPA := httptest.NewRequest(http.MethodGet, "/dashboard", nil)
	r.ServeHTTP(recSPA, reqSPA)

	require.Equal(t, http.StatusOK, recSPA.Code)
	require.Contains(t, recSPA.Body.String(), "<html>WeKnora SPA</html>")
}

func TestServeFrontendStaticStandardEditionWithoutWebDir(t *testing.T) {
	gin.SetMode(gin.TestMode)
	nonExistentDir := filepath.Join(t.TempDir(), "nonexistent")
	t.Setenv("WEKNORA_WEB_DIR", nonExistentDir)

	r := gin.New()
	serveFrontendStatic(r)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	r.ServeHTTP(recorder, request)

	// When webDir/index.html doesn't exist, static middleware is skipped, resulting in 404
	require.Equal(t, http.StatusNotFound, recorder.Code)
}
