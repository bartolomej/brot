package main

import (
	"errors"
	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
	"image"
	"image/color"
	"image/png"
	"log"
	"math/cmplx"
	"os"
)

type SetParams struct {
	IntervalX [2]float32
	IntervalY [2]float32
	Step      float32
	Iter      int
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
		Name:  "mandelbrot",
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

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
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
			iterations := mandelbrot(complex128(complex(x, y)), p.Iter)
			//color := colorful.Hsl(iterations, 0, (50 - iterations) * 2)
			c := color.RGBA{
				R: 255 - uint8(iterations*5),
				G: 255 - uint8(iterations*5),
				B: 255 - uint8(iterations*5),
				A: 255,
			}
			setImage.Set(j, i, c)
			j++
		}
		i++
		j = 0
	}

	file, err := os.Create(outDir + "/black&white.png")
	if err != nil {
		log.Fatal("Could not create image file")
	}

	if err := png.Encode(file, setImage); err != nil {
		log.Fatal("Could not encode png image")
	} else {
		log.Print("Image saved")
	}
}

func mandelbrot(p complex128, maxIter int) float64 {
	z, iter := p, 0
	for cmplx.Abs(z) < 2 && iter < maxIter {
		z = cmplx.Pow(z, 2) + p
		iter++
	}
	return float64(iter)
}
