package dbConnect

//utuberList
type UtuberInfo struct {
	Channel_id  string
	Playlist_id string
}

func (v UtuberInfo) TableName() string {
	return "youtuber_info"
}
func UtuberList() (channelID []UtuberInfo) {
	db := DbConnect("KOL")
	db.Model(&UtuberInfo{}).Pluck("playlist_id,channel_id", &channelID)
	//close connect
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	return channelID

}
