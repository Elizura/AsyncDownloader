package models

import (
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

func (downloadRequest *DownloadRequest) GetPeice(start_byte int, end_byte int) error {
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

	err = downloadRequest.WriteToAFile(start_byte, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func (downloadRequest *DownloadRequest) WriteToAFile(start_byte int, respBody io.ReadCloser) error {
	fileName := fmt.Sprintf("%s_%d", downloadRequest.FileName, start_byte)
	file, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	_, err = io.Copy(file, respBody)
	if err != nil {
		panic(err)
	}
	println(fmt.Sprintf("Wrote chunk %v to file", start_byte))

	return nil
}
