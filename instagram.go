package main

import (
	"net/http"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"strconv"
)

type InstagramSearchResults struct {
	Data []struct {
		Attribution interface{} `json:"attribution,omitempty"`
		Caption     struct {
			CreatedTime string `json:"created_time,omitempty"`
			From        struct {
				FullName       string `json:"full_name,omitempty"`
				ID             string `json:"id,omitempty"`
				ProfilePicture string `json:"profile_picture,omitempty"`
				Username       string `json:"username,omitempty"`
			} `json:"from"`
			ID   string `json:"id,omitempty"`
			Text string `json:"text,omitempty"`
		} `json:"caption"`
		Comments struct {
			Count int `json:"count,omitempty"`
			Data  []struct {
				CreatedTime string `json:"created_time,omitempty"`
				From        struct {
					FullName       string `json:"full_name,omitempty"`
					ID             string `json:"id,omitempty"`
					ProfilePicture string `json:"profile_picture,omitempty"`
					Username       string `json:"username,omitempty"`
				} `json:"from"`
				ID   string `json:"id,omitempty"`
				Text string `json:"text,omitempty"`
			} `json:"data"`
		} `json:"comments"`
		CreatedTime string `json:"created_time,omitempty"`
		Filter      string `json:"filter,omitempty"`
		ID          string `json:"id,omitempty"`
		Images      struct {
			LowResolution struct {
				Height int    `json:"height,omitempty"`
				URL    string `json:"url,omitempty"`
				Width  int    `json:"width,omitempty"`
			} `json:"low_resolution"`
			StandardResolution struct {
				Height int    `json:"height,omitempty"`
				URL    string `json:"url,omitempty"`
				Width  int    `json:"width,omitempty"`
			} `json:"standard_resolution,omitempty"`
			Thumbnail struct {
				Height int    `json:"height,omitempty"`
				URL    string `json:"url,omitempty"`
				Width  int    `json:"width,omitempty"`
			} `json:"thumbnail,omitempty"`
		} `json:"images,omitempty"`
		Likes struct {
			Count int `json:"count,omitempty"`
			Data  []struct {
				FullName       string `json:"full_name,omitempty"`
				ID             string `json:"id,omitempty"`
				ProfilePicture string `json:"profile_picture,omitempty"`
				Username       string `json:"username,omitempty"`
			} `json:"data"`
		} `json:"likes"`
		Link     string      `json:"link,omitempty"`
		Location interface{} `json:"location,omitempty"`
		Tags     []string    `json:"tags,omitempty"`
		Type     string      `json:"type,omitempty"`
		User     struct {
			FullName       string `json:"full_name,omitempty"`
			ID             string `json:"id,omitempty"`
			ProfilePicture string `json:"profile_picture,omitempty"`
			Username       string `json:"username,omitempty"`
		} `json:"user"`
		UsersInPhoto []interface{} `json:"users_in_photo,omitempty"`
	} `json:"data"`
	Meta struct {
		Code int `json:"code,omitempty"`
	} `json:"meta,omitempty"`
	Pagination struct {
		MinTagID           string `json:"min_tag_id,omitempty"`
		NextMaxTagID       string `json:"next_max_tag_id,omitempty"`
		NextURL            string `json:"next_url,omitempty"`
	} `json:"pagination"`
}

type InstagramPhoto struct {
	ID 			string
	URL 		string
	Username 	string
	Caption     string
	CreatedTime int64
}

func PhotosWithHashTag(hashTag string, maxTagId string, lastCreatedTime int64, clientID string)([]*InstagramPhoto, string, int64) {
	var url = fmt.Sprintf("https://api.instagram.com/v1/tags/%s/media/recent?client_id=%s&min_tag_id=%s&count=200", hashTag, clientID, maxTagId)
	log.Debug(fmt.Sprintf("Instagram API URL is %s", url))
	resp, err := http.Get(url)
	if err != nil {
		log.Error(fmt.Sprintf("Could not get photos from Instagram API: %v", err))
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(fmt.Sprintf("Could not read body of response from Instagram: %v", err))
	}
	var data InstagramSearchResults
	json.Unmarshal(body, &data)

	instagramPhotos := []*InstagramPhoto{}

	var nextMaxTagId  = data.Pagination.NextMaxTagID
	if (len(nextMaxTagId) == 0) {
		nextMaxTagId = maxTagId
	}
	log.Debug(fmt.Sprintf("nextMaxTagId before = %s", nextMaxTagId))
	log.Debug(fmt.Sprintf("maxTagId = %s", maxTagId))

	nextLastCreatedTime := int64(0)

	for _, entry := range data.Data {
		createdTime, _ := strconv.ParseInt(entry.CreatedTime, 10, 64)

		if (createdTime >= lastCreatedTime) {
			var instagramPhoto = &InstagramPhoto{
				entry.ID,
				entry.Images.StandardResolution.URL,
				entry.User.Username,
				entry.Caption.Text,
				createdTime,
			}
			instagramPhotos = append(instagramPhotos, instagramPhoto)
			if (createdTime > nextLastCreatedTime) {
				nextLastCreatedTime = createdTime
			}
		}
	}

	return instagramPhotos, nextMaxTagId, nextLastCreatedTime
}