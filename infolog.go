package infolog

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/aws/aws-sdk-go/service/s3"
)

type LogConfig struct {
	Level         int
	Logfilepath   string
	Logfilehandle *os.File
	Logbucket     string
	Logbucketkey  string
	S3session     *s3.S3
}

var Loglevel int
var Logconf LogConfig

const (
	PANIC = 0
	ERROR = 1
	WARN  = 2
	INFO  = 3
	DEBUG = 4
)

func Log(level int, msg string, val ...interface{}) {
	if Loglevel >= level {
		fmt.Printf(msg, val...)
		fmt.Println(" ")
		logmesg := fmt.Sprintf(msg, val...)
		logmesgdata := logmesg + "\n"
		Logconf.Logfilehandle.WriteString(logmesgdata)
		if level == PANIC {
			fmt.Println("\nFatal error ...\n")
			Logconf.Logfilehandle.WriteString("\nFatal error ...\n")
			StoreLog(Logconf.Logbucket, Logconf.Logbucketkey)
			panic("Exiting")
		}
	}
}

func StoreLog(bucketname string, bucketkey string) (err error) {
	content, err := ioutil.ReadFile(Logconf.Logfilepath)
	fmt.Println("Storing logs to S3...")
	err = copyLogsToS3(Logconf.S3session, bucketname, bucketkey, content)
	if err != nil {
		fmt.Println("Error in log copy to S3 :: Error=%s", err)
		return err
	}
	return err
}

func copyLogsToS3(svc *s3.S3, bucket string, filename string, content []byte) (err error) {

	uploadResult, err := svc.PutObject(&s3.PutObjectInput{
		Body:   bytes.NewReader(content),
		Bucket: &bucket,
		Key:    &filename,
	})
	if err != nil {
		fmt.Println("Failed to upload data to %s/%s, %s, %s\n", bucket, filename, err, uploadResult)
		return err
	}
	fmt.Println("Successfully uploaded data to bucket %s  with key %s\n", bucket, filename)

	return
}
