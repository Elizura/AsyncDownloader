package models

import (
	"bytes"
	"fmt"
	"io"
	httpclient "main/http_client"
	"os"
)

type DownloadRequest struct {
	Url        string
	FileName   string
	Chuncks    int
	TotalSize  int
	ChunkSize  int
	HttpClient *httpclient.HTTPClient
	Pieces     [5][]byte
}

func (dr *DownloadRequest) SplitIntoChuncks() [][2]int {
	chunkRanges := make([][2]int, dr.Chuncks)
	start := 0
	end := dr.ChunkSize
	for i := 0; i < dr.Chuncks; i++ {
		chunkRanges[i] = [2]int{start, end}
		start = end + 1
		end = end + dr.ChunkSize
	}

	return chunkRanges
}

func (downloadRequest *DownloadRequest) GetPeice(idx, start_byte int, end_byte int) error {
	header := map[string]string{
		"Range": "bytes=" + fmt.Sprintf("%d", start_byte) + "-" + fmt.Sprintf("%d", end_byte),
	}
	request, err := downloadRequest.HttpClient.CreateNewRequest("GET", downloadRequest.Url, header, nil)
	if err != nil {
		return err
	}
	resp, err := downloadRequest.HttpClient.DoRequest(request)
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, resp.Body)
	if err != nil {
		return err
	}
	println(fmt.Sprintf("Getting chunk %v", idx))
	downloadRequest.Pieces[idx] = buf.Bytes()

	return nil
}

func (downloadRequest *DownloadRequest) MergeFiles() error {
	downloadedFile, err := os.Create(downloadRequest.FileName)
	if err != nil {
		return err
	}
	defer downloadedFile.Close()
	for idx := 0; idx < downloadRequest.Chuncks; idx++ {
		_, err := downloadedFile.Write(downloadRequest.Pieces[idx])
		if err != nil {
			return err
		}
		println(fmt.Sprintf("Merged chunk %v to file", idx))
	}
	return nil
}
