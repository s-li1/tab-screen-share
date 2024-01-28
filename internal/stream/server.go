package stream

import (
	"context"
	"io"
	"sync"
	"time"
  pb "github.com/s-li1/remarkable-screen-share/proto"
)

type Server struct {
	pb.UnimplementedStreamServer
	imagePool sync.Pool
	r         io.ReaderAt
	fbPtr     int64
	runnable  chan struct{}
}

func (s *Server) Start() {
	go func(c chan struct{}) {
		for {
			c <- struct{}{}
			time.Sleep(200 * time.Millisecond)
		}
	}(s.runnable)
}

func NewServer(r io.ReaderAt, addr int64) *Server {
	return &Server{
		imagePool: sync.Pool{
			New: func() interface{} {
				return &pb.Image{
					Width:     ScreenWidth,
					Height:    ScreenHeight,
					ImageData: make([]byte, ScreenHeight*ScreenWidth*2),
				}
			},
		},
		r:        r,
		fbPtr:    addr,
		runnable: make(chan struct{}),
	}
}

func (s *Server) GetImage(ctx context.Context, in *pb.Input) (*pb.Image, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-s.runnable:
		img := s.imagePool.Get().(*pb.Image)
		defer s.imagePool.Put(img)
		_, err := s.r.ReadAt(img.ImageData, s.fbPtr)
		if err != nil {
			return nil, err
		}
		copiedData := make([]uint8, len(img.ImageData))
		copy(copiedData, img.ImageData)
		flipImage(copiedData, img.ImageData, int(img.Width), int(img.Height))
		return img, nil
	}
}

func flipImage(src, dst []uint8, w, h int) {
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			srcIndex := (y*w + x) * 2 // every second byte is useful
			dstIndex := (h-y-1)*w + x // unflip position
			dst[dstIndex] = src[srcIndex] * 17
		}
	}
}
