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
	"net/url"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	_ "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
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
	// -- if bucketName has a / separated path return list of items starting from the final element

	//_ = svc

	fmt.Println("bucketName from Ls():", bucketName)

	// s is for String - maybe not the best here, u might be more appropriate for url
	fmt.Println("buckerName:", bucketName)
	s, err := url.Parse(bucketName)
	if err != nil {
		log.Fatal(err)
	}
	b := s.Host
	p := s.Path
	fmt.Println("bucket:", b)
	fmt.Println("path:", p)

	if len(p) != 0 {
		fmt.Println("Path exists:", p)
		// i is for input
		i := &s3.ListObjectsV2Input{
			Bucket: aws.String(strings.Trim(b, "/")),
			Prefix: aws.String("test/"),
			// For now lets limit query to 2 items
			//MaxKeys: aws.Int64(2),
		}

		resp, _ := svc.ListObjectsV2(i)
		for _, key := range resp.Contents {
			fmt.Println(*key.Key)
		}

	} else {
		fmt.Println("No Path, list bucket contents")

		// i is for input
		i := &s3.ListObjectsV2Input{
			Bucket: aws.String(b),
			// For now lets limit query to 2 items
			MaxKeys: aws.Int64(2),
		}

		resp, _ := svc.ListObjectsV2(i)
		for _, key := range resp.Contents {
			fmt.Println(*key.Key)
		}
	}
}

// Lls local ls
func Lls(p string) {
	fmt.Println("print p", p)
	fi, err := os.Stat(p)
	if err != nil {
		fmt.Println(err)
		return
	}
	switch mode := fi.Mode(); {
	case mode.IsDir():
		// do directory stuff
		fmt.Println("It's a directory")
		fmt.Println("Do the file walk and list all keys")
	case mode.IsRegular():
		s := fi.Size()
		if fi.Size() > 10240 {
			// Multipart uploads can be performed on objects from 5 MB to 5 TB
			// each part has to be >= 5MB
			fmt.Println(s, " is greater than 10240")
			fmt.Println(s, " send file as multi part upload")
		} else {
			fmt.Println(fi.Name(), "is < 10240 bytes, send as single file")
		}
		//fmt.Println("Size", fi.Size(), "Bytes")
		//fmt.Println("Mode", fi.Mode())
	}

}

func main() {
	flagB := flag.Bool("b", false, "")
	_ = flagB
	flagBuckets := flag.Bool("buckets", false, "")
	flagLs := flag.Bool("ls", false, "")
	flagLls := flag.Bool("lls", false, "-lls list local files")

	flag.Parse()

	svc := service()

	if *flagBuckets == true {
		BucketList(svc)
	}
	if *flagLs == true {
		var bucketName string
		//fmt.Println(*flagLs)

		if len(flag.Args()) == 1 {
			for index, val := range flag.Args() {
				fmt.Println(index, ":", val)
				bucketName = val
			}
			//fmt.Println("os.Args[1]:", os.Args[1])
			//fmt.Println("flag.Args():", flag.Args())

			//fmt.Println("bucketName:", bucketName)

			Ls(svc, bucketName)
		}
	}
	if *flagLls == true {
		var path string
		if len(flag.Args()) == 1 {
			for _, val := range flag.Args() {
				//fmt.Println(index, ":", val)
				path = val
			}
		} else if len(flag.Args()) == 0 {
			path = fmt.Sprint("print all files names")
		}
		Lls(path)
	}

}
