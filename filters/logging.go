package filters

import (
	"log"
	"time"

	"github.com/beego/beego/v2/server/web/context"
)

// LoggingFilter logs each incoming HTTP request and its execution time.
// It helps in debugging and monitoring API requests.
func LoggingFilter(ctx *context.Context) {
	start := time.Now()

	log.Printf("[REQUEST] %s %s | IP: %s",
		ctx.Input.Method(),
		ctx.Input.URI(),
		ctx.Input.IP(),
	)

	go func() {
		time.Sleep(1 * time.Millisecond)
		log.Printf("[DONE] %s %s | Duration: ~%v",
			ctx.Input.Method(),
			ctx.Input.URI(),
			time.Since(start),
		)
	}()
}
