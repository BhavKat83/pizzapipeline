package main

import (
	"bytes"
	"context"
	b64 "encoding/base64"
	"errors"
	"fmt"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/jpoz/gomeme"
)

// InputData - Query Params for Lambda
type InputData struct {
	BottomText string `json:"bottomText"`
	InFile     string `json:"inFile"`
}

// Output - response body
type Output struct {
	StatusCode      int               `json:"statusCode"`
	Headers         map[string]string `json:"headers"`
	Body            string            `json:"body"`
	IsBase64Encoded bool              `json:"isBase64Encoded"`
}

// Version of gomeme
var Version = "1.1.0"

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// Download input base template for meme
func downloadFile(URL, fileName string) error {
	response, err := http.Get(URL)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return errors.New("Received non 200 response code")
	}
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	return nil
}

// HandleRequest - Lambda request handler
//func HandleRequest(inputData InputData) (Output, error) {
func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (Output, error) {

	config := gomeme.NewConfig()

	//config.BottomText = inputData.BottomText
	//inputImage := inputData.InFile
	inputImage := request.QueryStringParameters["inFile"]

	config.BottomText = request.QueryStringParameters["bottomText"]

	err := downloadFile(inputImage, "/tmp/meme.jpg")
	if err != nil {
		log.Fatal(err)
	}

	in, err := ioutil.ReadFile("/tmp/meme.jpg")
	if err != nil {
		fail("Could not open input", err)
	}

	meme := &gomeme.Meme{
		Config: config,
	}

	contentType := http.DetectContentType(in)
	buff := bytes.NewBuffer(in)

	switch contentType {
	case "image/gif":
		g, err := gif.DecodeAll(buff)
		if err != nil {
			fail("Failed to decode gif", err)
		}
		meme.Memeable = gomeme.GIF{g}
	case "image/jpeg":
		j, err := jpeg.Decode(buff)
		if err != nil {
			fail("Failed to decode jpeg", err)
		}
		meme.Memeable = gomeme.JPEG{j}
	case "image/png":
		p, err := png.Decode(buff)
		if err != nil {
			fail("Failed to decode png", err)
		}
		meme.Memeable = gomeme.PNG{p}
	default:
		fail(fmt.Sprintf("No idea what todo with a %s", contentType), nil)
	}

	out, err := os.Create("/tmp/output-meme.jpg")
	if err != nil {
		fail("Could not open output file", err)
	}

	check(meme.Write(out))

	buf := bytes.NewBuffer(nil)
	downloadedFile, _ := os.Open("/tmp/output-meme.jpg") // Error handling elided for brevity.
	io.Copy(buf, downloadedFile)                         // Error handling elided for brevity.
	downloadedFile.Close()

	o := Output{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type":        "application/jpg",
			"Content-Disposition": "attachment; filename=mymeme.jpg",
		},
		Body:            b64.StdEncoding.EncodeToString(buf.Bytes()),
		IsBase64Encoded: true,
	}

	return o, nil

}

// RandStringBytes - Generate random string for output file name
func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func fail(s string, e error) {
	fmt.Fprintf(os.Stderr, s)
	check(e)
	os.Exit(1)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	lambda.Start(HandleRequest)
}
