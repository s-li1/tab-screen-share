package main

import (
	"bytes"
	"compress/zlib"
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"io/ioutil"
	"log"
	"net/http"
	"sync"

	_ "github.com/joho/godotenv/autoload"
	"github.com/mattn/go-mjpeg"
	"github.com/s-li1/remarkable-screen-share/internal/compression"
	"github.com/s-li1/remarkable-screen-share/internal/notion"
	pb "github.com/s-li1/remarkable-screen-share/proto"
	"google.golang.org/grpc"
	ggzip "google.golang.org/grpc/encoding/gzip"
)

const (
	SAddr = "10.0.0.75:2000"
	BAddr = ":8080"
)

func init() {
	err := ggzip.SetLevel(zlib.BestSpeed)
	if err != nil {
		panic(err)
	}
}

func main() {
	log.Printf("dialing on %v", SAddr)
	conn, err := grpc.Dial(SAddr, grpc.WithInsecure(), grpc.WithMaxMsgSize(8*1024*1024), grpc.WithDefaultCallOptions(grpc.UseCompressor("gzip")))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()

  imageBufPool := sync.Pool {
    New: func() interface{} {
      return &image.Gray{}
    },
  }

	mjpegStream := mjpeg.NewStream()
	go GetImage(mjpegStream, conn, &imageBufPool)

	response, err := notion.UploadFile()
	if err != nil {
		log.Fatalf("client: error making http request: %s\n", err)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(responseData))

	mux := http.NewServeMux()

	mux.HandleFunc("/video", compression.CreateGzipHandler(mjpegStream))
  mux.HandleFunc("/download", func(w http.ResponseWriter, r *http.Request) {
 // if err := download(&imageBufPool); err != nil {
 //   log.Fatalln("Download Failed")
 // }

    log.Println("Download Succeeded")
  })

	log.Printf("listening on %v, registered /video", BAddr)
	err = http.ListenAndServe(BAddr, mux)
	if err != nil {
		log.Fatal(err)
	}
}

func uploadFile(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "Uploading File")
}

func GetImage(mjpegStream *mjpeg.Stream, conn *grpc.ClientConn, imageBufPool *sync.Pool) {
	c := pb.NewStreamClient(conn)

	var img image.Gray
	for {
		response, err := c.GetImage(context.Background(), &pb.Input{})
		if err != nil {
			log.Fatalf("Error when calling GetImage: %s", err)
		}

		img.Rect = image.Rect(0, 0, int(response.Width), int(response.Height))
		img.Stride = int(response.Width)
		img.Pix = response.ImageData

		if err := updateMJPEGStream(mjpegStream, &img); err != nil {
			log.Fatal(err)
		}

    imageBufPool.Put(img)
	}
}

func updateMJPEGStream(mjpegStream *mjpeg.Stream, img *image.Gray) error {
	var b bytes.Buffer
	if err := jpeg.Encode(&b, img, nil); err != nil {
		return err
	}

	return mjpegStream.Update(b.Bytes())
}

