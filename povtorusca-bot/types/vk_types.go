package types

import "encoding/json"

// VKCallback represents the main callback structure from VK
type VKCallback struct {
	Type    string          `json:"type"`
	EventID string          `json:"event_id"`
	V       string          `json:"v"`
	Object  json.RawMessage `json:"object"`
	GroupID int             `json:"group_id"`
	Secret  string          `json:"secret"`
}

// WallPost represents a wall post object from VK
type WallPost struct {
	ID             int             `json:"id"`
	FromID         int             `json:"from_id"`
	OwnerID        int             `json:"owner_id"`
	Date           int             `json:"date"`
	PostType       string          `json:"post_type"`
	Text           string          `json:"text"`
	Attachments    []Attachment    `json:"attachments"`
	PostAuthorData *PostAuthorData `json:"post_author_data"`
}

// PostAuthorData represents author information
type PostAuthorData struct {
	Author    int `json:"author"`
	Publisher int `json:"publisher"`
}

// Attachment represents an attachment in a VK post
type Attachment struct {
	Type  string `json:"type"`
	Photo *Photo `json:"photo,omitempty"`
}

// Photo represents a photo attachment
type Photo struct {
	ID        int        `json:"id"`
	OwnerID   int        `json:"owner_id"`
	OrigPhoto *PhotoSize `json:"orig_photo,omitempty"`
}

// PhotoSize represents photo dimensions and URL
type PhotoSize struct {
	Type   string `json:"type"`
	URL    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}
