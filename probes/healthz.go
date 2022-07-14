package probes

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var (
	// Status of the Service
	STATUS int = 500

	// Status Codes
	LIVENESS int = 500
)

func StartHealthz(_port string, _log *log.Logger) {

	_log.Printf("%s | %s\n", "INFO", "Start Healthz-Probe")

	corsconfig := cors.DefaultConfig()

	corsconfig.AllowAllOrigins = true

	gin.SetMode("release")

	gin.DisableConsoleColor()

	gin.DefaultWriter = _log.Writer()

	ginrouter := gin.New()

	ginrouter.Use(cors.New(corsconfig))
	ginrouter.Use(gin.LoggerWithFormatter(func(params gin.LogFormatterParams) string {
		return fmt.Sprintf(
			"[jr] %s - Status: %d | Method: %s | Path: %s |  %dms | response: %dBytes | RemoteAddr: %s\n",
			fmt.Sprintf(
				"%d/%02d/%02d %02d:%02d:%02d.%.6s",
				time.Now().Year(),
				time.Now().Month(),
				time.Now().Day(),
				time.Now().Hour(),
				time.Now().Minute(),
				time.Now().Second(),
				strconv.Itoa(time.Now().Nanosecond())),
			params.StatusCode,
			params.Method,
			params.Path,
			params.Latency,
			params.BodySize,
			params.Request.RemoteAddr)
	}))

	ginrouter.Use(gin.Recovery())

	httpserver := &http.Server{
		Addr:    ":" + _port,
		Handler: ginrouter,
	}

	ginrouter.Handle("GET", "/livez", func(ctx *gin.Context) {

		switch LIVENESS {
		case 200:
			ctx.String(200, "Liveness: %s", "OK")
		default:
			ctx.String(500, "Liveness: %s", "ERROR")
		}
	})

	if err := httpserver.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
		_log.Printf("%s | %s\n", "ERROR", err.Error())
	}

	_log.Printf("%s | %s\n", "WARNING", "HTTP-Controller stopped working.")
}
