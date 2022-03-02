package main

import (
	db "YoutubeWorkerPool/dbConnect"
	worker "YoutubeWorkerPool/workerPool"
)

func main() {

	utuberList := db.UtuberList()
	worker.UtuberDetail(utuberList)
}
