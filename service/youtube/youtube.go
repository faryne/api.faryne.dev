package youtube

import (
	"context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
	"os"
)

type YoutubeInstance struct {
	service *youtube.Service
}

var (
	channelList = []string{"snippet", "statistics", "brandingSettings"}
	videoList   = []string{"liveStreamingDetails", "snippet"}
	searchList  = []string{"snippet"}
)

const (
	eventTypeUpcoming  = "upcoming"
	eventTypeLive      = "live"
	eventTypeCompleted = "completed"
)

func New(filename string) (*YoutubeInstance, error) {
	if _, err := os.Stat(filename); err != nil {
		return nil, err
	}
	youtube.NewService(context.TODO(), option.WithCredentials(&google.DefaultCredentials{}))
	client, err := youtube.NewService(context.Background(), option.WithCredentialsFile(filename))
	if err != nil {
		return nil, err
	}
	return &YoutubeInstance{
		service: client,
	}, nil
}

// GetChannelById 根據 ID 取得頻道資訊
func (s *YoutubeInstance) GetChannelById(id string, options ...string) (*youtube.ChannelListResponse, error) {
	if len(options) == 0 {
		return s.service.Channels.List(channelList).Id(id).Do()
	}
	return s.service.Channels.List(options).Id(id).Do()
}

// GetVideoById 根據 ID 取得影片資訊
func (s *YoutubeInstance) GetVideoById(videoId string, options ...string) (*youtube.VideoListResponse, error) {
	if len(options) == 0 {
		return s.service.Videos.List(videoList).Id(videoId).Do()
	}
	return s.service.Videos.List(options).Id().Do()
}

// 拉出各種狀態的直播影片
func (s *YoutubeInstance) getLives(eventType, channelId string, options ...string) (*youtube.SearchListResponse, error) {
	if len(options) == 0 {
		return s.service.Search.List(options).ChannelId(channelId).EventType(eventType).Type("video").Do()
	}
	return s.service.Search.List(searchList).ChannelId(channelId).EventType(eventType).Type("video").Do()
}

// GetUpcomingLive 拉出即將直播的影片
func (s *YoutubeInstance) GetUpcomingLive(channelId string, options ...string) (*youtube.SearchListResponse, error) {
	return s.getLives(eventTypeUpcoming, channelId, options...)
}

// GetCompletedLive 拉出已結束的直播
func (s *YoutubeInstance) GetCompletedLive(channelId string, options ...string) (*youtube.SearchListResponse, error) {
	return s.getLives(eventTypeCompleted, channelId, options...)
}

// GetNowLive 拉出正在直播的影片
func (s *YoutubeInstance) GetNowLive(channelId string, options ...string) (*youtube.SearchListResponse, error) {
	return s.getLives(eventTypeLive, channelId, options...)
}

// GetLiveMessages 取得聊天室即時訊息
func (s *YoutubeInstance) GetLiveMessages(livechatId string, callback func(*youtube.LiveChatMessageListResponse) error, options ...string) error {
	var resp *youtube.LiveChatMessageListResponse
	var err error
	if resp, err = s.service.LiveChatMessages.List(livechatId, options).Do(); err != nil {
		return err
	}
	return callback(resp)
}
