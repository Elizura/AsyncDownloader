package manager

import (
	"fmt"
	httpClient "main/http_client"
	"main/models"
	"net/url"
	"path"
	"sync"
)

func Download(url url.URL) {
	// init http client
	urlString := url.String()
	httpClient := httpClient.CreateClient()
	request, err := httpClient.CreateNewRequest("HEAD", urlString, nil, nil)
	if err != nil {
		panic(err)
	}
	// make HEAD request to get the file size
	response, err := httpClient.DoRequest(request)
	if err != nil {
		panic(err)
	}
	// getContentLength  in bytes
	contentLength := response.ContentLength
	fileName := path.Base(url.Path)
	chuncks := 5
	chunckSize := contentLength / int64(chuncks)

	// create a download request

	downLoadRequest := models.DownloadRequest{
		Url:        urlString,
		FileName:   fileName,
		Chuncks:    chuncks,
		TotalSize:  int(contentLength),
		ChunkSize:  int(chunckSize),
		HttpClient: httpClient,
	}

	// get FileName, Chuncks, TotalSize, ChunkSize

	chunckArray := downLoadRequest.SplitIntoChuncks()

	var wg sync.WaitGroup

	for _, chunk := range chunckArray {
		wg.Add(1)
		go func(start int, end int) {
			defer wg.Done()
			err := downLoadRequest.GetPeice(start, end)
			if err != nil {
				panic(err)
			}
		}(chunk[0], chunk[1])
	}
	wg.Wait()
	fmt.Println(chunckArray)
}