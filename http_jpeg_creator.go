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

func readDimensionParams(req *http.Request) (width, height uint32, err error) {
	// width
	if width, err = readUInt32Param(req, "width"); err != nil {
		return 0, 0, err
	}

	// height
	if height, err = readUInt32Param(req, "height"); err != nil {
		return 0, 0, err
	}

	return width, height, nil
}

func readRgbParams(req *http.Request) (width, height, red, green, blue uint32, err error){
	if width, height, err = readDimensionParams(req); err != nil {
		return 0, 0, 0, 0, 0, err
	}

	// red
	if red, err = readUInt32Param(req, "red"); err != nil {
		return 0, 0, 0, 0, 0, err
	}

	// green
	if green, err = readUInt32Param(req, "green"); err != nil {
		return 0, 0, 0, 0, 0, err
	}

	// blue
	if red, err = readUInt32Param(req, "blue"); err != nil {
		return 0, 0, 0, 0, 0, err
	}

	return width, height, red, green, blue, nil
}

func readUInt32Param(req *http.Request, paramName string) (uint32, error){
	str := req.URL.Query().Get(paramName)
	log.Println(paramName + ": " + str)
	i, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint32(i), nil
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
