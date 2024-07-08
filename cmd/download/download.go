package download

import (
	"fmt"
	"github.com/imMohika/gohangyourself/net"
	"github.com/pterm/pterm"
	"io"
	"os"
)

type WriteCounter struct {
	downloaded uint64
	progress   *pterm.ProgressbarPrinter
}

func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.downloaded += uint64(n)
	wc.progress.Add(n)
	return n, nil
}

func FromURL(url string, fileName string) error {
	resp := net.Request(url, "failed to download "+fileName)

	out, err := os.Create(fileName + ".tmp")
	if err != nil {
		return err
	}

	total := resp.ContentLength

	progress, _ := pterm.DefaultProgressbar.WithTotal(int(total)).WithTitle("Downloading...").Start()

	counter := &WriteCounter{
		progress: progress,
	}
	if _, err = io.Copy(out, io.TeeReader(resp.Body, counter)); err != nil {
		if err := out.Close(); err != nil {
			return fmt.Errorf("failed to close file %s: %w", fileName+".tmp", err)
		}
		return err
	}

	fmt.Print("\n")

	if err := out.Close(); err != nil {
		return fmt.Errorf("failed to close file %s: %w", fileName+".tmp", err)
	}

	// progress.Stop never returns an error but yeah!
	if _, err := progress.Stop(); err != nil {
		return fmt.Errorf("failed to stop progressbar: %w", err)
	}

	if err = os.Rename(fileName+".tmp", fileName); err != nil {
		return err
	}
	return nil
}
