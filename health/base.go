package health

import "context"

type HealthCheck interface {
	Check(context.Context) (bool, error)
}
