package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rekognition"
)

type RequestBody struct {
	Image string
}

var svc *rekognition.Rekognition

func main() {
	http.Handle("/", http.FileServer(http.Dir("./public")))
	http.HandleFunc("/check", checkHandler)
	http.HandleFunc("/upload", UploadHandler)
	http.ListenAndServe(":3000", nil)
}

func init() {
	sess := session.New(&aws.Config{
		Region: aws.String("us-east-1"),
	})
	svc = rekognition.New(sess)
}
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Allowed POST method only", http.StatusMethodNotAllowed)
		return
	}
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fileSrc, _, err := r.FormFile("upload")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer fileSrc.Close()
	bytes, err := ioutil.ReadAll(fileSrc)
	if err != nil {
		fmt.Println("err", err.Error())
	}
	input := &rekognition.DetectModerationLabelsInput{
		Image: &rekognition.Image{
			Bytes: bytes,
		},
	}
	result, err := svc.DetectModerationLabels(input)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	// Shoot the response back to the front-end
	output, err := json.Marshal(result)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(output)
}

func checkHandler(w http.ResponseWriter, r *http.Request) {

	if r.Body == nil {
		http.Error(w, "request body expected", 400)
		return
	}
	// referencing the struct we created earlier
	var parsed RequestBody

	if err := json.NewDecoder(r.Body).Decode(&parsed); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	// Decode the string sent from the front-end
	decodedImage, err := base64.StdEncoding.DecodeString(parsed.Image)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	// Send the request to Rekognition
	input := &rekognition.DetectLabelsInput{
		Image: &rekognition.Image{
			Bytes: decodedImage,
		},
	}
	result, err := svc.DetectLabels(input)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	// Shoot the response back to the front-end
	output, err := json.Marshal(result)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(output)
}
