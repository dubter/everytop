package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func YouTubeFetch() (YouTubeAPIResponse, error) {
	url := "https://yt-api.p.rapidapi.com/trending"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return YouTubeAPIResponse{}, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Add("x-rapidapi-key", "63f74429c5msh798c6b8a3f26e5ap13d290jsna70bc0b819a7")
	req.Header.Add("x-rapidapi-host", "yt-api.p.rapidapi.com")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return YouTubeAPIResponse{}, fmt.Errorf("error making request: %v", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return YouTubeAPIResponse{}, fmt.Errorf("error reading response body: %v", err)
	}

	var response YouTubeAPIResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return YouTubeAPIResponse{}, fmt.Errorf("failed to unmarshal YouTubeAPIResponse: %v", err)
	}

	return response, nil
}

func NewsFetch() (NewsResponse, error) {
	url := "https://news-api14.p.rapidapi.com/v2/trendings?topic=General&language=en"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("x-rapidapi-key", "63f74429c5msh798c6b8a3f26e5ap13d290jsna70bc0b819a7")
	req.Header.Add("x-rapidapi-host", "news-api14.p.rapidapi.com")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return NewsResponse{}, fmt.Errorf("error reading response body: %v", err)
	}

	var response NewsResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return NewsResponse{}, fmt.Errorf("failed to unmarshal YouTubeAPIResponse: %v", err)
	}

	return response, nil
}

func main() {
	bot, err := tgbotapi.NewBotAPI("7896611207:AAH6NNnNUYPlPAP99R3wJ-ELtj0R-TfmApk")
	if err != nil {
		log.Fatalf("Error creating bot: %v", err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil && update.Message.Text == "/start" {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "commands: \n/video\n/news")
			bot.Send(msg)
		}
		if update.Message != nil && update.Message.Text == "/video" {
			trendingData, err := YouTubeFetch()
			if err != nil {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Error fetching trending data: %v", err))
				bot.Send(msg)
				continue
			}

			if len(trendingData.Data) == 0 {
				return
			}

			idx := rand.Intn(len(trendingData.Data))

			finalMsg := fmt.Sprintf("https://youtube.com/watch?v=%s", trendingData.Data[idx].VideoId)

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, finalMsg)
			bot.Send(msg)
		}
		if update.Message != nil && update.Message.Text == "/news" {
			trendingData, err := NewsFetch()
			if err != nil {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Error fetching trending data: %v", err))
				bot.Send(msg)
				continue
			}

			if len(trendingData.Data) == 0 {
				return
			}

			idx := rand.Intn(len(trendingData.Data))

			finalMsg := fmt.Sprintf("%s\n%s", trendingData.Data[idx].Title, trendingData.Data[idx].Url)

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, finalMsg)
			bot.Send(msg)
		}
	}
}

type YouTubeAPIResponse struct {
	Data []YouTubeData `json:"data"`
	Msg  string        `json:"msg"`
}

type YouTubeData struct {
	Type             string `json:"type"`
	VideoId          string `json:"videoId,omitempty"`
	Title            string `json:"title"`
	ChannelTitle     string `json:"channelTitle,omitempty"`
	ChannelId        string `json:"channelId,omitempty"`
	ChannelHandle    string `json:"channelHandle,omitempty"`
	ChannelThumbnail []struct {
		Url    string `json:"url"`
		Width  int    `json:"width"`
		Height int    `json:"height"`
	} `json:"channelThumbnail,omitempty"`
	Description       string    `json:"description,omitempty"`
	ViewCount         string    `json:"viewCount,omitempty"`
	PublishedTimeText string    `json:"publishedTimeText,omitempty"`
	PublishDate       string    `json:"publishDate,omitempty"`
	PublishedAt       time.Time `json:"publishedAt,omitempty"`
	LengthText        string    `json:"lengthText,omitempty"`
	Thumbnail         []struct {
		Url    string `json:"url"`
		Width  int    `json:"width"`
		Height int    `json:"height"`
	} `json:"thumbnail,omitempty"`
	RichThumbnail []struct {
		Url    string `json:"url"`
		Width  int    `json:"width"`
		Height int    `json:"height"`
	} `json:"richThumbnail,omitempty"`
	Data []struct {
		Type          string `json:"type"`
		VideoId       string `json:"videoId"`
		Title         string `json:"title"`
		ViewCountText string `json:"viewCountText,omitempty"`
		Thumbnail     []struct {
			Url    string `json:"url"`
			Width  int    `json:"width"`
			Height int    `json:"height"`
		} `json:"thumbnail"`
		IsOriginalAspectRatio bool   `json:"isOriginalAspectRatio,omitempty"`
		Params                string `json:"params,omitempty"`
		PlayerParams          string `json:"playerParams,omitempty"`
		SequenceParams        string `json:"sequenceParams,omitempty"`
		ChannelTitle          string `json:"channelTitle,omitempty"`
		ChannelId             string `json:"channelId,omitempty"`
		ChannelThumbnail      []struct {
			Url    string `json:"url"`
			Width  int    `json:"width"`
			Height int    `json:"height"`
		} `json:"channelThumbnail,omitempty"`
		Description       string    `json:"description,omitempty"`
		ViewCount         string    `json:"viewCount,omitempty"`
		PublishedTimeText string    `json:"publishedTimeText,omitempty"`
		PublishDate       string    `json:"publishDate,omitempty"`
		PublishedAt       time.Time `json:"publishedAt,omitempty"`
		LengthText        string    `json:"lengthText,omitempty"`
		RichThumbnail     []struct {
			Url    string `json:"url"`
			Width  int    `json:"width"`
			Height int    `json:"height"`
		} `json:"richThumbnail,omitempty"`
		ChannelHandle string `json:"channelHandle,omitempty"`
	} `json:"data,omitempty"`
	Badges []string `json:"badges,omitempty"`
}

type NewsResponse struct {
	Success     bool `json:"success"`
	Size        int  `json:"size"`
	TotalHits   int  `json:"totalHits"`
	HitsPerPage int  `json:"hitsPerPage"`
	Page        int  `json:"page"`
	TotalPages  int  `json:"totalPages"`
	TimeMs      int  `json:"timeMs"`
	Data        []struct {
		Title         string    `json:"title"`
		Url           string    `json:"url"`
		Excerpt       string    `json:"excerpt"`
		Thumbnail     string    `json:"thumbnail"`
		Language      string    `json:"language"`
		Paywall       bool      `json:"paywall"`
		ContentLength int       `json:"contentLength"`
		Date          time.Time `json:"date"`
		Authors       []string  `json:"authors"`
		Keywords      []string  `json:"keywords"`
		Publisher     struct {
			Name    string `json:"name"`
			Url     string `json:"url"`
			Favicon string `json:"favicon"`
		} `json:"publisher"`
	} `json:"data"`
}
