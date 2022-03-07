package workerPool

import (
	db "YoutubeWorkerPool/dbConnect"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/bitly/go-simplejson"
	_ "github.com/joho/godotenv/autoload"
)

var UtubeKey = os.Getenv("YOUTUBE_KEY")

func UtuberDetail(UID []db.UtuberInfo) {
	jobs := make(chan db.UtuberInfo, 10)
	result := make(chan string, 10)
	for i := 1; i <= 10; i++ {
		go detailWorker(jobs, result)
	}
	go func() {
		for _, value := range UID {
			jobs <- value
		}
		close(jobs)
	}()
	for i := 0; i < len(UID); i++ {
		<-result
	}
}

func detailWorker(jobs <-chan db.UtuberInfo, result chan string) {
	for channelInfo := range jobs {
		videoList := []string{}
		response, _ := http.Get("https://www.googleapis.com/youtube/v3/playlistItems?part=snippet&playlistId=" + channelInfo.Playlist_id + "&key=" + UtubeKey + "&maxResults=50")
		body, _ := io.ReadAll(response.Body)
		js, _ := simplejson.NewJson(body)
		nextPageToken, _ := js.GetPath("nextPageToken").String()
		items, _ := js.GetPath("items").Array()
		for i, _ := range items {
			videoId, _ := js.GetPath("items").GetIndex(i).GetPath("snippet", "resourceId", "videoId").String()
			videoList = append(videoList, videoId)
		}
		for nextPageToken != "" {
			getViedoLise(&nextPageToken, channelInfo.Playlist_id, &videoList)
		}
		insertVideoInformation(channelInfo, videoList)
		result <- "ok"
	}
}
func getViedoLise(nextPageToken *string, playlistId string, videoList *[]string) {
	response, _ := http.Get("https://www.googleapis.com/youtube/v3/playlistItems?part=snippet&playlistId=" + playlistId + "&key=" + UtubeKey + "&maxResults=50&pageToken=" + *nextPageToken)
	body, _ := io.ReadAll(response.Body)
	js, _ := simplejson.NewJson(body)
	*nextPageToken, _ = js.GetPath("nextPageToken").String()
	items, _ := js.GetPath("items").Array()
	for i, _ := range items {
		videoId, _ := js.GetPath("items").GetIndex(i).GetPath("snippet", "resourceId", "videoId").String()
		*videoList = append(*videoList, videoId)
	}
}
func insertVideoInformation(channelInfo db.UtuberInfo, videoList []string) {
	commentCount := 0
	likeCount := 0

	for _, value := range videoList {
		response, _ := http.Get("https://www.googleapis.com/youtube/v3/videos?id=" + value + "&key=" + UtubeKey + "&part=statistics")
		body, _ := io.ReadAll(response.Body)
		js, _ := simplejson.NewJson(body)
		comment, _ := js.GetPath("items").GetIndex(0).GetPath("statistics", "commentCount").String()
		like, _ := js.GetPath("items").GetIndex(0).GetPath("statistics", "likeCount").String()
		commentInt, _ := strconv.Atoi(comment)
		likeInt, _ := strconv.Atoi(like)
		commentCount += commentInt
		likeCount += likeInt

	}
	response, _ := http.Get("https://www.googleapis.com/youtube/v3/channels?id=" + channelInfo.Channel_id + "&key=" + UtubeKey + "&part=statistics")
	body, _ := io.ReadAll(response.Body)
	js, _ := simplejson.NewJson(body)
	viewCountString, _ := js.GetPath("items").GetIndex(0).GetPath("statistics", "viewCount").String()
	subscriberCountString, _ := js.GetPath("items").GetIndex(0).GetPath("statistics", "subscriberCount").String()
	videoCountString, _ := js.GetPath("items").GetIndex(0).GetPath("statistics", "videoCount").String()
	viewCount, _ := strconv.Atoi(viewCountString)
	subscriberCount, _ := strconv.Atoi(subscriberCountString)
	videoCount, _ := strconv.Atoi(videoCountString)
	db.SetUtuberLog(channelInfo.Channel_id, commentCount, likeCount, viewCount, subscriberCount, videoCount)

}
