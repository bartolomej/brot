package main

import (
	"errors"
	"fmt"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"math/cmplx"
	"math/rand"
	"os"
)

type SetParams struct {
	IntervalX [2]float32
	IntervalY [2]float32
	Step      float32
	Iter      int
	hue       Hue
}

type Hue struct {
	Start  float64
	Factor float64
}

type Scene struct {
	Name   string
	Params SetParams
}

type Config struct {
	Scenes []Scene
}

const outDir = "out"

func main() {
	app := &cli.App{
		Name:  "brot",
		Usage: "Render mandelbrot set images.",
		Action: func(c *cli.Context) error {
			initOutFolder()
			err, scene := setupConfig(c.Args().Get(0))
			if err != nil {
				log.Fatal(err)
			}
			log.Print("Rendering scene: " + c.Args().Get(0))
			drawComplexSet(scene.Params)
			return nil
		},
	}

	if len(os.Args) == 1 {
		log.Print("No args provided. Rendering test scene.")
		initOutFolder()
		params := SetParams{
			IntervalX: [2]float32{-2.1, 0.7},
			IntervalY: [2]float32{-1.2, 1.2},
			Step:      0.005,
			Iter:      20,
			hue: Hue{
				Start:  100,
				Factor: 20,
			},
		}
		drawComplexSet(params)
	} else {
		err := app.Run(os.Args)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func setupConfig(name string) (error, Scene) {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}
	var config Config
	err2 := viper.Unmarshal(&config)
	if err2 != nil {
		log.Fatal("Unable to unmarshal config")
	}
	for _, s := range config.Scenes {
		if s.Name == name {
			return nil, s
		}
	}
	return errors.New("No scene found"), Scene{}
}

func initOutFolder() {
	if _, err := os.Stat(outDir); os.IsNotExist(err) {
		err := os.Mkdir(outDir, 0775)
		if err != nil {
			log.Fatal("Could not create output directory")
		} else {
			log.Print("Created output directory")
		}
	}
}

func drawComplexSet(p SetParams) {
	spanX := p.IntervalX[1] - p.IntervalX[0]
	spanY := p.IntervalY[1] - p.IntervalY[0]
	setImage := image.NewRGBA(image.Rect(0, 0, int(spanX/p.Step), int(spanY/p.Step)))

	i, j := 0, 0
	for y := p.IntervalY[0]; y < p.IntervalY[1]; y += p.Step {
		for x := p.IntervalX[0]; x < p.IntervalX[1]; x += p.Step {
			n := mandelbrot(complex128(complex(x, y)), p.Iter)
			r, g, b := colorful.Hsl((n*p.hue.Factor)+p.hue.Start, 0.8, 0.6).RGB255()
			if math.IsNaN(n) {
				// color with black where points don't diverge
				r, g, b = 0, 0, 0
			}
			setImage.Set(j, i, color.RGBA{
				R: r,
				G: g,
				B: b,
				A: 255,
			})
			j++
		}
		i++
		j = 0
	}

	file, err := os.Create(fmt.Sprintf("%s/%d.png", outDir, rand.Int()))
	if err != nil {
		log.Fatal("Could not create image file")
	}

	if err := png.Encode(file, setImage); err != nil {
		log.Fatal("Could not encode png image")
	} else {
		log.Print("Image saved")
	}
}

func mandelbrot(c complex128, maxIter int) float64 {
	z, n, r, d := complex(0, 0), 0, 2.0, 2.0
	for cmplx.Abs(z) < r && n < maxIter {
		z = cmplx.Pow(z, complex(d, 0)) + c
		n++
	}
	// SMOOTH ITERATION COUNT TAKEN FROM:
	// http://linas.org/art-gallery/escape/smooth.html
	return float64(n) + 1 - math.Log(math.Log(cmplx.Abs(z)))/math.Log(d)
}
