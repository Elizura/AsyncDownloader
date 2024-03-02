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

	err = downloadRequest.WriteToAFile(idx, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func (downloadRequest *DownloadRequest) WriteToAFile(idx int, respBody io.ReadCloser) error {
	fileName := fmt.Sprintf("%s_%d", downloadRequest.FileName, idx)
	file, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	_, err = io.Copy(file, respBody)
	if err != nil {
		panic(err)
	}

	defer file.Close()
	downloadRequest.MergeFiles()
	println(fmt.Sprintf("Wrote chunk %v to file", idx))

	return nil
}

func (downloadRequest *DownloadRequest) MergeFiles() error {
	downloadedFile, err := os.Create(downloadRequest.FileName)
	if err != nil {
		return err
	}
	defer downloadedFile.Close()
	for idx := 0; idx < downloadRequest.Chuncks; idx++ {
		chunkFile, err := os.Open(fmt.Sprintf("%s_%d", downloadRequest.FileName, idx))
		if err != nil {
			return err
		}
		defer chunkFile.Close()
		_, err = io.Copy(downloadedFile, chunkFile)
		if err != nil {
			return err
		}
		println(fmt.Sprintf("Merged chunk %v to file", idx))
	}
	return nil
}
