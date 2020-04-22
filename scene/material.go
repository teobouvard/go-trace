package scene

import (
	"github.com/teobouvard/gotrace/space"
)

// Material define the way actors interact with a ray
type Material interface {
	Scatter(ray Ray, record HitRecord) (scatters bool, attenuation space.Vec3, scattered Ray)
}
