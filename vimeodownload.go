package vimeodownload

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/raitonoberu/vimego"
)

// Download downloads the video from the given URL.
func Download(ctx context.Context, id, referer, out string) error {
	if out == "" {
		out = fmt.Sprintf("%s.mp4", id)
	}
	u := fmt.Sprintf("https://vimeo.com/%s", id)
	video, err := vimego.NewVideo(u)
	if err != nil {
		return fmt.Errorf("couldn't create video %s: %w", u, err)
	}
	if referer != "" {
		video.Header["Referer"] = []string{referer}
	}
	formats, err := video.Formats()
	if err != nil {
		return fmt.Errorf("couldn't get video formats: %w", err)
	}

	if err := downloadFile(out, formats.Progressive[1].URL); err != nil {
		return fmt.Errorf("couldn't download video %s: %w", u, err)
	}
	return nil
}

func downloadFile(filepath string, url string) (err error) {
	log.Println("📹 start downloading", url)
	defer log.Println("📹 finished downloading", url)

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Writer the body to file
	if _, err = io.Copy(out, resp.Body); err != nil {
		return err
	}
	return nil
}
