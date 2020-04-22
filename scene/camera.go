package scene

import "github.com/teobouvard/gotrace/space"

// Camera represents the camera of the scene
type Camera struct {
	origin     space.Vec3
	horizontal space.Vec3
	vertical   space.Vec3
	swCorner   space.Vec3
}

// NewCamera creates a camera
func NewCamera() Camera {
	return Camera{
		origin:     space.NewVec3(0, 0, 0),
		horizontal: space.NewVec3(4, 0, 0),
		vertical:   space.NewVec3(0, 2, 0),
		swCorner:   space.NewVec3(-2, -1, -1),
	}
}

// LookAt returns the Ray when the camera looks at (u, v)
func (c Camera) LookAt(u float64, v float64) Ray {
	hOffset := space.Scale(c.horizontal, u)
	vOffset := space.Scale(c.vertical, v)
	return NewRay(c.origin, space.Add(c.swCorner, hOffset, vOffset))
}
