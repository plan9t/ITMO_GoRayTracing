package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
)

type Vec3f struct {
	X, Y, Z float64
}

type Sphere struct {
	Center Vec3f
	Radius float64
	Color  Vec3f
}

type Light struct {
	Position  Vec3f
	Intensity float64
}

func NewLight(position Vec3f, intensity float64) *Light {
	return &Light{Position: position, Intensity: intensity}
}

// Операция сложения векторов
func (v Vec3f) Add(other Vec3f) Vec3f {
	return Vec3f{v.X + other.X, v.Y + other.Y, v.Z + other.Z}
}

// Операция вычитания векторов
func (v Vec3f) Subtract(other Vec3f) Vec3f {
	return Vec3f{v.X - other.X, v.Y - other.Y, v.Z - other.Z}
}

// Операция умножения вектора на скаляр
func (v Vec3f) MulScalar(scalar float64) Vec3f {
	return Vec3f{v.X * scalar, v.Y * scalar, v.Z * scalar}
}

// Скалярное произведение векторов
func (v Vec3f) Dot(other Vec3f) float64 {
	return v.X*other.X + v.Y*other.Y + v.Z*other.Z
}

// Квадрат длины вектора
func (v Vec3f) Length2() float64 {
	return v.Dot(v)
}

// Нормализация вектора
func (v Vec3f) Normalize() Vec3f {
	sqrt := math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
	return Vec3f{v.X / sqrt, v.Y / sqrt, v.Z / sqrt}
}

// Пересечение луча со сферой
func (s *Sphere) RayIntersect(orig, dir Vec3f) (bool, float64) {
	L := s.Center.Subtract(orig)
	tca := L.Dot(dir)
	d2 := L.Length2() - tca*tca
	if d2 > s.Radius*s.Radius {
		return false, 0
	}
	thc := math.Sqrt(s.Radius*s.Radius - d2)
	t0 := tca - thc
	t1 := tca + thc
	if t0 < 0 {
		t0 = t1
	}
	if t0 < 0 {
		return false, 0
	}
	return true, t0
}

// castRay определяет цвет луча.
func castRay(orig, dir Vec3f, spheres []Sphere, lights []Light) Vec3f {
	closestDist := math.MaxFloat64
	var hitSphere *Sphere
	for i := range spheres {
		hit, dist := spheres[i].RayIntersect(orig, dir)
		if hit && dist < closestDist {
			closestDist = dist
			hitSphere = &spheres[i]
		}
	}

	if hitSphere == nil {
		return Vec3f{0.2, 0.7, 0.8} // background color
	}

	// Точка пересечения луча со сферой
	point := orig.Add(dir.MulScalar(closestDist))
	// Нормаль в точке пересечения
	N := point.Subtract(hitSphere.Center).Normalize()
	// Диффузная интенсивность света
	diffuseLightIntensity := 0.0

	for _, light := range lights {
		lightDir := light.Position.Subtract(point).Normalize()
		diffuseLightIntensity += light.Intensity * math.Max(0, lightDir.Dot(N))
	}

	// Возвращаем цвет сферы, умноженный на интенсивность света
	return hitSphere.Color.MulScalar(diffuseLightIntensity)
}

// colorToRGBA преобразует Vec3f в color.RGBA.
func colorToRGBA(c Vec3f) color.RGBA {
	return color.RGBA{
		R: uint8(math.Max(0, math.Min(255, c.X*255))),
		G: uint8(math.Max(0, math.Min(255, c.Y*255))),
		B: uint8(math.Max(0, math.Min(255, c.Z*255))),
		A: 255, // Полная непрозрачность
	}
}

// render - генерация изображения.
func render(spheres []Sphere, lights []Light) {
	const width, height = 1024, 768
	const fov = math.Pi / 3 // field of view
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	for j := 0; j < height; j++ {
		for i := 0; i < width; i++ {
			x := (2*(float64(i)+0.5)/float64(width) - 1) * math.Tan(fov/2) * float64(width) / float64(height)
			y := -(2*(float64(j)+0.5)/float64(height) - 1) * math.Tan(fov/2)
			dir := Vec3f{x, y, -1}.Normalize()
			col := castRay(Vec3f{0, 0, 0}, dir, spheres, lights)
			img.Set(i, j, colorToRGBA(col))
		}
	}

	file, err := os.Create("out/out9.png")
	if err != nil {
		panic(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Printf("Close error")
		}
	}(file)

	err = png.Encode(file, img)
	if err != nil {
		fmt.Printf("Encode error")
	}
}

func main() {
	lights := []Light{
		*NewLight(Vec3f{X: 1.0, Y: 2.0, Z: 3.0}, 1.4),
		*NewLight(Vec3f{X: 3.0, Y: -2.0, Z: -3.0}, 1.0),
	}

	// Инициализация сцены с несколькими сферами
	spheres := []Sphere{
		{Center: Vec3f{X: 0, Y: 0, Z: -3}, Radius: 0.8, Color: Vec3f{X: 0.4, Y: 0.4, Z: 0.3}},
		{Center: Vec3f{X: 2, Y: 0, Z: -4}, Radius: 0.5, Color: Vec3f{X: 0.7, Y: 0.3, Z: 0.5}},
		{Center: Vec3f{X: -2, Y: 0, Z: -5}, Radius: 1.2, Color: Vec3f{X: 0.3, Y: 0.6, Z: 0.7}},
	}

	render(spheres, lights)
}
