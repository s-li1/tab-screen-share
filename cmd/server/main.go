package main

import (
	"bufio"
	"compress/zlib"
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	pb "github.com/s-li1/remarkable-screen-share/proto"
  "github.com/s-li1/remarkable-screen-share/internal/stream"
	"google.golang.org/grpc"
	"google.golang.org/grpc/encoding/gzip"
)

const port = ":2000"

func init() {
	err := gzip.SetLevel(zlib.BestSpeed)
	if err != nil {
		panic(err)
	}
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Printf("listening on port %s", port)
	defer lis.Close()

	if len(os.Args) < 2 {
		log.Fatalf("Usage: Provide process ID of xochitl %v", os.Args[0])
	}

	pid := os.Args[1]

	file, err := os.OpenFile(filepath.Join("/proc", pid, "mem"), os.O_RDONLY, os.ModeDevice)
	if err != nil {
		log.Fatal("cannot open file: ", err)
	}
	defer file.Close()

	addr, err := getFrameBufferPointer(pid)

	s := stream.NewServer(file, addr+8)
	s.Start()

	grpcServer := grpc.NewServer()

	pb.RegisterStreamServer(grpcServer, s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to server: %v", err)
	}
}

func getFrameBufferPointer(pid string) (int64, error) {
	file, err := os.OpenFile(filepath.Join("/proc", pid, "maps"), os.O_RDONLY, os.ModeDevice)
	if err != nil {
		log.Fatal("cannot open file: ", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	scanner.Split(bufio.ScanWords)
	var addr int64
	scanAddr := false
	for scanner.Scan() {
		if scanAddr {
			hex := strings.Split(scanner.Text(), "-")[0]
			addr, err = strconv.ParseInt("0x"+hex, 0, 64)
			break
		}

		if scanner.Text() == "/dev/fb0" {
			scanAddr = true
		}
	}

	return addr, err
}
