package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"log"
	"strconv"
	"image"
	"image/color"
	"image/jpeg"
	"bytes"
	"math/rand"
	"time"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/create/rgb/rnd", createRgbRndHandler).Methods("GET")
	r.HandleFunc("/create/rgb", createRgbHandler).Methods("GET")
	http.Handle("/", r)
	http.ListenAndServe(":8000", nil)
}

func createRgbRndHandler(rw http.ResponseWriter, req *http.Request) {
	imgResponseHandler(rw, req, func(req *http.Request) (image.Image, error) {
		width, height, err := readDimensionParams(req)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		img := createRndImage(width, height)
		return img, nil
	})
}

func createRgbHandler(rw http.ResponseWriter, req *http.Request) {
	imgResponseHandler(rw, req, func(req *http.Request) (image.Image, error) {
		w, h, r, g, b, err := readRgbParams(req)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		img := createImageOfColor(w, h, r, g, b)
		return img, nil
	})
}

func imgResponseHandler(rw http.ResponseWriter, req *http.Request, fn func(req *http.Request) (image.Image, error)) {
	img, err := fn(req)
	if err != nil {
		log.Println("could not create image")
	}

	buf := new(bytes.Buffer)
	if err := jpeg.Encode(buf, img, &jpeg.Options{Quality: 100}); err != nil {
		log.Println("unable to encode image")
		return
	}

	rw.Header().Set("Content-Type", "image/jpeg")
	rw.Header().Set("Content-Length", strconv.Itoa(buf.Len()))

	if _, err := rw.Write(buf.Bytes()); err != nil {
		log.Println("unable to write image")
		return
	}
}

func readDimensionParams(req *http.Request) (width, height int, err error) {
	// width
	wStr := req.URL.Query().Get("width")
	log.Println("width: " + wStr)
	w, err := strconv.ParseInt(wStr, 10, 64)
	if err != nil {
		return 0, 0, err
	}
	width = int(w)

	// height
	hStr := req.URL.Query().Get("height")
	log.Println("height: " + hStr)
	h, err := strconv.ParseInt(hStr, 10, 64)
	if err != nil {
		return 0, 0, err
	}
	height = int(h)

	return width, height, nil
}

func readRgbParams(req *http.Request) (width, height int, red, green, blue uint32, err error){
	width, height, err = readDimensionParams(req)
	if err != nil {
		return 0, 0, 0, 0, 0, err
	}

	// red
	rStr := req.URL.Query().Get("red")
	log.Println("red: " + rStr)
	r, err := strconv.ParseInt(rStr, 10, 64)
	if err != nil {
		return 0, 0, 0, 0, 0, err
	}
	red = uint32(r)

	// green
	gStr := req.URL.Query().Get("green")
	log.Println("green: " + gStr)
	g, err := strconv.ParseInt(gStr, 10, 64)
	if err != nil {
		return 0, 0, 0, 0, 0, err
	}
	green = uint32(g)

	// blue
	bStr := req.URL.Query().Get("blue")
	log.Println("blue: " + bStr)
	b, err := strconv.ParseInt(bStr, 10, 64)
	if err != nil {
		return 0, 0, 0, 0, 0, err
	}
	blue = uint32(b)

	return width, height, red, green, blue, nil
}

func createRndImage(width, height int) (image.Image) {
	rgba := image.NewRGBA(image.Rect(0, 0, width, height))

	rand.Seed(time.Now().UnixNano())

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			c := color.RGBA{R: uint8(rand.Uint32()), G: uint8(rand.Uint32()), B: uint8(rand.Uint32()), A: uint8(255)}
			rgba.Set(x, y, c)
		}
	}

	img := rgba
	return img
}

func createImageOfColor(width, height int, red, green, blue uint32) (image.Image) {
	c := color.RGBA{R: uint8(red), G: uint8(green), B: uint8(blue), A: uint8(255)}
	rgba := image.NewRGBA(image.Rect(0, 0, width, height))

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			rgba.Set(x, y, c)
		}
	}

	img := rgba
	return img
}
