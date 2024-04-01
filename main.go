package main

import (
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
}

// Операция вычитания векторов
func (v Vec3f) Subtract(other Vec3f) Vec3f {
	return Vec3f{v.X - other.X, v.Y - other.Y, v.Z - other.Z}
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
	len := math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
	return Vec3f{v.X / len, v.Y / len, v.Z / len}
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
func castRay(orig, dir Vec3f, sphere Sphere) Vec3f {
	hit, _ := sphere.RayIntersect(orig, dir)
	if !hit {
		return Vec3f{0.2, 0.7, 0.8} // background color
	}
	return Vec3f{0.4, 0.4, 0.3} // sphere color
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
func render(sphere Sphere) {
	const width, height = 1024, 768
	const fov = math.Pi / 3 // field of view
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	for j := 0; j < height; j++ {
		for i := 0; i < width; i++ {
			x := (2*(float64(i)+0.5)/float64(width) - 1) * math.Tan(fov/2) * float64(width) / float64(height)
			y := -(2*(float64(j)+0.5)/float64(height) - 1) * math.Tan(fov/2)
			dir := Vec3f{x, y, -1}.Normalize()
			col := castRay(Vec3f{0, 0, 0}, dir, sphere)
			img.Set(i, j, colorToRGBA(col))
		}
	}

	file, err := os.Create("out/out2.png")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	png.Encode(file, img)
}

func main() {
	// Задаем параметры сферы.
	sphere := Sphere{
		Center: Vec3f{X: 0, Y: 0, Z: -3}, // Центр сферы
		Radius: 1.0,                      // Радиус сферы
	}

	// Главная функция для генерации изображения
	render(sphere)
}
