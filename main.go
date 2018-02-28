//package s3tpa

// Options:
// -ls
// dir [remote-directory] [local-directory]  (list remote files)
// get [remote-file]
// put
// lcd (all of these commands expect an interactive shel-like environment)
// mdelete  (delete remote file)
// mdir
// mget
// mput
//    Expand wild cards in the list of local files given as argu‐
//    ments and do a put for each file in the resulting list.  See
//    glob for details of filename expansion.  Resulting file names
//    will then be processed according to ntrans and nmap settings.
// rmdir (delete remote directory)
// status
// verbose
// hash
package main

import (
	"flag"
	"fmt"
	"log"
	//"reflect"
	"os"
	_"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/aws"
)

var svc *s3.S3

func Ls(*s3.S3) {
	//result, err := svc.ListBuckets(nil)
	result, err := svc.ListBuckets(nil)
	if err != nil {
		//exitError("Unable to list buckets, %v", err)
		log.Fatalf("Unable to list buckets, %v", err)
	}

	fmt.Println("Buckets:")

	for _, b := range result.Buckets {
		fmt.Printf("* %s created on %s\n",
			aws.StringValue(b.Name), aws.TimeValue(b.CreationDate))
	}

}

func service() *s3.S3 {

	os.Setenv("AWS_SDK_LOAD_CONFIG", "true")
	os.Setenv("AWS_PROFILE", "default")
	//fmt.Println("AWS_SDK_LOAD_CONFIG:", os.Getenv("AWS_SDK_LOAD_CONFIG"))
	//fmt.Println("AWS_PROFILE:", os.Getenv("AWS_PROFILE"))

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	//svc := s3.New(sess)
	svc = s3.New(sess)

	return svc
}

func main() {
	flagB := flag.Bool("b", false, "")
	flagBuckets := flag.Bool("buckets", false, "")

	flag.Parse()

	svc := service()

	if *flagB == true || *flagBuckets == true {
		Ls(svc)
	}

}