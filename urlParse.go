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
// -bo -- ordinal position of the bucket you want to operate on in the list returned by --buckets
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
	_ "strings"
	"net/url"
)

var svc *s3.S3

func service() *s3.S3 {

	os.Setenv("AWS_SDK_LOAD_CONFIG", "true")
	os.Setenv("AWS_PROFILE", "default")
	//fmt.Println("AWS_SDK_LOAD_CONFIG:", os.Getenv("AWS_SDK_LOAD_CONFIG"))
	//fmt.Println("AWS_PROFILE:", os.Getenv("AWS_PROFILE"))

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	svc = s3.New(sess)

	return svc
}

func BucketList(*s3.S3) {
	//result, err := svc.ListBuckets(nil)
	result, err := svc.ListBuckets(nil)
	if err != nil {
		//exitError("Unable to list buckets, %v", err)
		log.Fatalf("Unable to list buckets, %v", err)
	}

	fmt.Println("Buckets:")

	for _, b := range result.Buckets {
		//fmt.Printf("* %s created on %s\n",
		//	aws.StringValue(b.Name), aws.TimeValue(b.CreationDate))
		fmt.Println(aws.StringValue(b.Name))
	}

}

func Ls(svc *s3.S3, bucketName string) {
	// TODO:
	// -- validate that bucket exists and is available, or throw error
	// -- operate on proper nouns, or ordinal position from BucketList
	// --if bucketName has a / separated path return list of items starting from the final element

	//_ = svc

	//fmt.Println(bucketName)

	s, err := url.Parse(bucketName)
	if err != nil {
		log.Fatal(err)
	}
	b := s.Host
	p := s.Path
	fmt.Println("bucket:", b)
	fmt.Println("path:", p)

	if len(p) == 0 {
		fmt.Println("No Path, list bucket contents")
		//sub := strings.Split(bucketName, "/")
		//host, path := sub[0], sub[1]

		//if len(sub) == 1 {
		//fmt.Println(sub)

		params := &s3.ListObjectsInput{
			//Bucket: aws.String("dlts-s3-stan"),
			Bucket: aws.String(bucketName),
		}

		resp, _ := svc.ListObjects(params)
		for _, key := range resp.Contents {
			fmt.Println(*key.Key)
		}

		//}
	}
}

func main() {
	flagB := flag.Bool("b", false, "")
	_ = flagB
	flagBuckets := flag.Bool("buckets", false, "")
	flagLs := flag.Bool("ls", false, "")

	flag.Parse()

	svc := service()

	if *flagBuckets == true {
		BucketList(svc)
	}
	if *flagLs == true {
		var bucketName string

		if len(flag.Args()) == 1 {
			for index, val := range flag.Args() {
				fmt.Println(index, ":", val)
				bucketName = val
			}
			//fmt.Println("os.Args[1]:", os.Args[1])
			////fmt.Println("flag.Args():", flag.Args())
			Ls(svc, bucketName)
		}
	}

}