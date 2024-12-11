package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	proto "thumbnail/internal/proto/proto"

	"google.golang.org/grpc"
)

func main() {
	async := flag.Bool("async", false, "Execute requests asynchronously")
	flag.Parse()

	if len(flag.Args()) == 0 {
		log.Fatalf("Usage: %s [--async] <url1,url2,...>", os.Args[0])
	}

	urls := strings.Split(flag.Arg(0), ",")

	conn, err := grpc.NewClient("localhost:50051")
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	client := proto.NewThumbnailServiceClient(conn)

	if *async {
		getThumbnailsAsync(client, urls)
	} else {
		getThumbnailsSync(client, urls)
	}
}

func getThumbnailsSync(client proto.ThumbnailServiceClient, urls []string) {
	fmt.Println("Executing in synchronous mode...")

	req := &proto.ThumbnailsRequest{
		Urls: urls,
	}

	resp, err := client.GetThumbnails(context.Background(), req)
	if err != nil {
		log.Printf("Failed to get thumbnails: %v", err)
		return
	}

	for i, thumbnail := range resp.Results {
		fmt.Printf("Thumbnail for %s received: %d bytes\n", urls[i], len(thumbnail.ImageData))
	}
}

func getThumbnailsAsync(client proto.ThumbnailServiceClient, urls []string) {
	fmt.Println("Executing in asynchronous mode...")
	var wg sync.WaitGroup

	for _, url := range urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			req := &proto.ThumbnailRequest{Url: url}
			resp, err := client.GetThumbnail(context.Background(), req)
			if err != nil {
				log.Printf("Failed to get thumbnail for URL %s: %v", url, err)
				return
			}
			fmt.Printf("Thumbnail for %s received: %d bytes\n", url, len(resp.ImageData))
		}(url)
	}
	wg.Wait()
}
