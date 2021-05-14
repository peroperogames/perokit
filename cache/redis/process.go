package redis

import (
	"github.com/peroperogames/perokit/core/mapping"
	timex "github.com/peroperogames/perokit/core/time"
	"github.com/peroperogames/perokit/log"

	"strings"

	red "github.com/go-redis/redis"
)

func process(proc func(red.Cmder) error) func(red.Cmder) error {
	return func(cmd red.Cmder) error {
		start := timex.Now()

		defer func() {
			duration := timex.Since(start)
			if duration > slowThreshold {
				var buf strings.Builder
				for i, arg := range cmd.Args() {
					if i > 0 {
						buf.WriteByte(' ')
					}
					buf.WriteString(mapping.Repr(arg))
				}
				log.Info("[REDIS] slowcall on executing: %s  \n", buf.String())
			}
		}()

		return proc(cmd)
	}
}
