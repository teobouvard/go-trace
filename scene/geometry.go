package scene

import (
	"github.com/teobouvard/gotrace/space"
)

// Geometry ??
type Geometry interface {
	Hit(ray Ray, tMin float64, tMax float64) (hit bool, dist float64, pos space.Vec3, normal space.Vec3)
}
