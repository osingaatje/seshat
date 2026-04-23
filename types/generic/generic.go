package generic

import (
	"math"
)

type UTMLXY struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// location on some grid or whatever
type Vector2D struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

func (v Vector2D) New(utmlPos UTMLXY) Vector2D {
	return Vector2D{
		X: utmlPos.X,
		Y: utmlPos.Y,
	}
}
func (v Vector2D) NewInt(x int, y int) Vector2D {
	return Vector2D{
		X: float64(x),
		Y: float64(y),
	}
}
func (v Vector2D) Add(vO Vector2D) Vector2D {
	return Vector2D{
		X: v.X + vO.X,
		Y: v.Y + vO.Y,
	}
}
func (v Vector2D) Sub(vO Vector2D) Vector2D {
	return Vector2D{
		X: v.X - vO.X,
		Y: v.Y - vO.Y,
	}
}
func (v Vector2D) Div(factor float64) Vector2D {
	return Vector2D{
		X: v.X / factor,
		Y: v.Y / factor,
	}
}
func (v Vector2D) Mul(factor float64) Vector2D {
	return Vector2D{
		X: v.X * factor,
		Y: v.Y * factor,
	}
}

func (v Vector2D) MulComponents(X float64, Y float64) Vector2D {
	return Vector2D{
		X: v.X * X,
		Y: v.Y * Y,
	}
}

// euclidean distance. d(p,q) = sqrt((p1 - q1)^2 + (p2 - q2)^2)
func (v Vector2D) Dist(vO Vector2D) float64 {
	return math.Sqrt(math.Pow(v.X-vO.X, 2) + math.Pow(v.Y-vO.Y, 2))
}
