package health

import (
	"faceit-backend-test/internal/router"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"net/http"
	"runtime"
)

const route = "/health"

// controller it handles the health check
// @tag.name HealthController
type controller struct {
	db *sqlx.DB

	serviceStatus bool
}

var _ router.Controller = (*controller)(nil)

type ControllerOpts func(controller2 *controller)

func NewController(opts ...ControllerOpts) *controller {
	c := &controller{}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func WithDb(db *sqlx.DB) ControllerOpts {
	return func(c *controller) {
		c.db = db
	}
}

func (c *controller) Register(r *gin.RouterGroup) {
	r.GET(route, c.healthCheck)
}

// healthCheck godoc
// @Summary checks the status of the service
// @tags HealthController
// @Produce json
// @Success 200 {object} Response
// @Router /v1/health [get]
func (c *controller) healthCheck(ctx *gin.Context) {
	ctx.Header("Cache-Control", "no-cache")
	c.serviceStatus = true
	responseCode := http.StatusOK

	dbConn := c.healthCheckDatabase()

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	ctx.JSON(responseCode, Response{
		Status:             c.serviceStatus,
		DatabaseConnection: dbConn,
		MemoryStats: MemoryStats{
			CurrentAllocated:                bToMb(m.Alloc),
			TotalAllocated:                  bToMb(m.TotalAlloc),
			TotalMemory:                     bToMb(m.Sys),
			CompletedGarbageCollectorCycles: m.NumGC,
		},
	})
}

func (c *controller) healthCheckDatabase() DatabaseConnection {
	dbConnection := true
	dbErrMessage := ""
	dbErr := c.db.Ping()
	if dbErr != nil {
		c.serviceStatus = false
		dbConnection = false
		dbErrMessage = dbErr.Error()
	}

	return DatabaseConnection{
		Connected: dbConnection,
		Error:     dbErrMessage,
	}
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
