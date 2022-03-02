package dbConnect

import "time"

//utuberList
type UtuberLog struct {
	Id              string
	ViewCount       int `gorm:"column:ViewCount"`
	VideoCount      int `gorm:"column:VideoCount"`
	SubscriberCount int `gorm:"column:SubscriberCount"`
	LikeCount       int `gorm:"column:LikeCount"`
	CommentCount    int `gorm:"column:CommentCount"`
	Date            string
}

func (v UtuberLog) TableName() string {
	return "youtuber_trend_log"
}
func SetUtuberLog(channelID string, commentCount int, likeCount int, viewCount int, subscriberCount int, videoCount int) {
	db := DbConnect("KOL")
	date := time.Now().Format("2006-01-02")
	logInfo := UtuberLog{Id: channelID, ViewCount: viewCount, VideoCount: videoCount, SubscriberCount: subscriberCount, LikeCount: likeCount, CommentCount: commentCount, Date: date}
	db.Create(&logInfo)
	//close connect
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

}
