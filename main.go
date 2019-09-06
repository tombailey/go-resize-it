package main

import (
	"errors"
	"fmt"
	"github.com/noelyahan/impexp"
	"github.com/noelyahan/mergi"
	"image"
	"image/png"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

func main() {
	http.HandleFunc("/", resizeImageHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

type ResizeParameters struct {
	imageUrl string
	width uint
	height uint
}

func resizeImageHandler(response http.ResponseWriter, request *http.Request) {
	resizeParameters, err := parseQueryParameters(request.URL.Query())
	if err != nil {
		fmt.Println("Failed to resize image.", err)
		response.WriteHeader(400)
	} else {
		image, err := resize(resizeParameters)
		if err != nil {
			fmt.Println("Failed to resize image.", err)
			response.WriteHeader(500)
		} else {
			response.WriteHeader(200)
			err = png.Encode(response, image)
			if err != nil {
				fmt.Println("Failed to write response.", err)
				response.WriteHeader(500)
			}
		}
	}
}

func parseQueryParameters(values url.Values) (ResizeParameters, error) {
	var resizeParameters ResizeParameters

	imageUrl := values.Get("imageUrl")
	if len(imageUrl) == 0 {
		return resizeParameters, errors.New("imageUrl should be provided")
	}

	width, err := strconv.ParseUint(values.Get("width"), 10, 64)
	if err != nil {
		return resizeParameters, errors.New("width should be a uint")
	}

	height, err := strconv.ParseUint(values.Get("height"), 10, 64)
	if err != nil {
		return resizeParameters, errors.New("height should be a uint")
	}

	resizeParameters = ResizeParameters{
		imageUrl: imageUrl,
		width:    uint(width),
		height:   uint(height),
	}

	return resizeParameters, nil
}

func resize(resizeParameters ResizeParameters) (image.Image, error) {
	image, err := mergi.Import(impexp.NewURLImporter(resizeParameters.imageUrl))
	if err != nil {
		return nil, err
	} else {
		return mergi.Resize(image, resizeParameters.width, resizeParameters.height)
	}
}
