package utilities

import (
	"bit-jobs-api/shared/response"
	"net/http"
	"os"

	// Additional imports needed for examples below
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
)

func UploadImage(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		response.SendJSONMessage(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	key := os.Getenv("SPACES_KEY")
	secret := os.Getenv("SPACES_SECRET")

	// Parse multipart form, 10 << 20 specifies a maximum
	// upload of 10 MB files.
	r.ParseMultipartForm(10 << 20)

	file, handler, err := r.FormFile("file")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		response.SendJSONMessage(w, http.StatusInternalServerError, response.FriendlyError)
		return
	}
	defer file.Close()

	s3Config := &aws.Config{
		Credentials: credentials.NewStaticCredentials(key, secret, ""),
		Endpoint:    aws.String(os.Getenv("SPACES_ENDPOINT")),
		Region:      aws.String("us-east-1"),
	}

	newSession, err := session.NewSession(s3Config)
	if err != nil {
		fmt.Println(err.Error())
		response.SendJSONMessage(w, http.StatusInternalServerError, response.FriendlyError)
	}

	s3Client := s3.New(newSession)

	keyname := fmt.Sprintf("%s-%s", uuid.NewString(), handler.Filename)
	object := s3.PutObjectInput{
		Bucket: aws.String(os.Getenv("SPACES_BUCKET")),
		Key:    aws.String(keyname),
		Body:   file,
		ACL:    aws.String("public-read"),
	}
	_, err = s3Client.PutObject(&object)

	if err != nil {
		fmt.Println(err.Error())
		response.SendJSONMessage(w, http.StatusInternalServerError, response.FriendlyError)
	}

	response.SendJSON(w, map[string]string{"url": fmt.Sprintf("https://%s.%s/%s", os.Getenv("SPACES_BUCKET"), os.Getenv("SPACES_ENDPOINT"), keyname)})
}
