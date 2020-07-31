# Brot

A tool to produce aesthetically pleasing images of mandelbrot set from the command line.

![Black and white set](./images/black&white.png)

## Setup

Using this tool requires a working Go environment. 
See the [installation instructions for Go](https://golang.org/doc/install) or [read the articles on DigitalOcean](https://www.digitalocean.com/community/tutorial_series/how-to-install-and-set-up-a-local-programming-environment-for-go).
For further information read [this article](https://www.digitalocean.com/community/tutorials/how-to-build-and-install-go-programs) on building and installing go programs.

Move to your projects folder and clone the repo:
```shell script
$ git clone https://github.com/bartolomej/brot && cd brot
```

Download all dependencies:
```shell script
$ go get
```

Run go program with:
```shell script
$ go run main.go
```

...or build executable binaries to `/bin` directory:
```shell script
$ sh build.sh
```

You can also install the program to use it as a cli tool from anywhere in the system by running:
```shell script
$ go install
```

Note that the above commands will only work if you correctly configured the go environment.

## Usage

Once you've installed this program with `go install` you can run commands with `brot`.
Without any arguments the program will render example image.

Create `config.json` file for custom parameters. Example syntax:
```json
{
  "scenes": [
    {
      "name": "Mandelbrot",
      "params": {
        "intervalX": [-2.1, 0.7],
        "intervalY": [-1.2, 1.2],
        "step": 0.01,
        "iter": 50,
        "hue": {
          "start": 0,
          "factor": 10
        }
      }
    }
  ]
}
```

Render a custom scene from config by providing `name` as a first parameter to `brot` command:
```shell script
$ brot <scene-name>
```

Image outputs are stored in `/out` folder.