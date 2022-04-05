package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
)

func main() {
	// --- pipeline --
	// a go routine reads the json
	// the other downloads the images
	// the last writes the back to the data folder

	ch1 := ReadJson(os.Args[1])
	ch2 := GetImagesFromWeb(ch1)
	done := SaveImages(ch2)
	<-done
}

func ReadJson(filename string) <-chan string {
	out := make(chan string, 10)
	fmt.Println("Read json")

	go func() {

		var bytes []byte
		var images []string

		fd, err := os.Open(filename)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("File " + filename + " opened")
		bytes, err = ioutil.ReadAll(fd)
		if err != nil {
			fmt.Println(err)
			return
		}
		err = json.Unmarshal(bytes, &images)
		if err != nil {
			fmt.Println(err)
			return
		}
		for _, i := range images {
			fmt.Println(i)
			out <- i
		}
		close(out)
	}()

	return out
}

func GetImagesFromWeb(urls <-chan string) <-chan FileImage {
	out := make(chan FileImage, 10)
	fmt.Println("Get images")

	go func() {
		for image := range urls {
			response, err := http.Get(image)
			if err != nil {
				fmt.Println(err)
			}
			defer response.Body.Close()

			//io.Copy(os.Stdout, response.Body)
			out <- FileImage{response.Body, image}
		}
		close(out)
	}()
	return out
}

type FileImage struct {
	http_response io.Reader
	url           string
}

func SaveImages(input <-chan FileImage) <-chan struct{} {
	done := make(chan struct{})

	go func() {
		defer close(done)
		fmt.Println("Save images")
		for image := range input {
			file, err := os.Create("data/" + path.Base(image.url))
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			defer file.Close()
			io.Copy(file, image.http_response)
		}
	}()

	return done
}
