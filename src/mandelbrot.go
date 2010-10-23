package main

import (
  "cmath"
  "image"
  "image/png"
  "os"
)

type Settings struct {
  z0 complex128
  n complex128
  repeat int

  viewCenter complex128
  viewWidth float64
  viewHeight float64

  imageWidth int
  imageHeight int
}

var settings = &Settings{
  z0 : complex128(cmplx(0, 0)),
  n : complex128(cmplx(2.0, 0)),
  repeat : 30,
  //viewCenter : complex128(cmplx(-0.75, 0)),
  viewCenter : complex128(cmplx(0, 0)),
  viewWidth : 3.5,
  viewHeight : 3.5,
  imageWidth : 300,
  imageHeight : 300,
}

func cmplx128(r float64, i float64) complex128 {
  return complex128(cmplx(r, i))
}

func checkMandelbrot(c complex128, res chan int) {
  z := settings.z0
  for i := 0; i < settings.repeat; i++ {
    z = cmath.Pow(z, settings.n) + c
    if 2.0 < cmath.Abs(cmath.Pow(z, 2)) {
      res <- i + 1
      return
    }
  }
  res <- -1
}

func calculate(results [][]chan int) {
  var comp complex128
  cr := real(settings.viewCenter) - settings.viewWidth / 2
  ci := imag(settings.viewCenter) - settings.viewHeight / 2
  dr := settings.viewWidth / float64(settings.imageWidth)
  di := settings.viewHeight / float64(settings.imageHeight)

  for y := 0; y < settings.imageHeight; y++ {
    results[y] = make([]chan int, settings.imageWidth)
    for x := 0; x < settings.imageWidth; x++ {
      results[y][x] = make(chan int)
      comp = cmplx128(cr + dr * float64(x), ci + di * float64(y))
      go checkMandelbrot(comp, results[y][x])
    }
  }
}

func generateImage(results [][]chan int) *image.RGBA {
  rgba := image.NewRGBA(settings.imageWidth, settings.imageHeight)
  for y := 0; y < settings.imageHeight; y++ {
    for x := 0; x < settings.imageWidth; x++ {
      res := <-results[y][x]
      if res < 0 {
        rgba.Set(x, y, image.RGBAColor{0, 0, 0, 255})
      } else {
        depth := uint8(25 * res)
        rgba.Set(x, y, image.RGBAColor{0, depth, depth, 255})
      }
    }
  }
  return rgba
}

func main() {
  var results [][]chan int
  results = make([][]chan int, settings.imageHeight)
  calculate(results)
  image := generateImage(results)
  png.Encode(os.Stdout, image)
}
