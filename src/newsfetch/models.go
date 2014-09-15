package main

import "time"

const (
	TypeHTML  = "html"
	TypeText  = "text"
	TypeImage = "image"
	TypeVideo = "video"
	TypeAudio = "audio"
	TypeRSS   = "rss"
	TypeXML   = "xml"
	TypeAtom  = "atom"
	TypeJSON  = "json"
	TypePPT   = "ptt"
	TypeLink  = "link"
	TypeError = "error"
)

type Options struct {
	MaxWidth     int
	MaxHeight    int
	Width        int
	Words        int
	Chars        int
	WMode        bool
	AllowScripts bool
	NoStyle      bool
	Autoplay     bool
	VideoSrc     bool
	Frame        bool
	Secure       bool
}

type Author struct {
	ID   []byte `json:"id"`
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Entity struct {
	ID    []byte `json:"id"`
	Count int    `json:"count"`
	Name  string `json:"name"`
}

type Image struct {
	ID      []byte  `json:"id"`
	Caption string  `json:"caption"`
	URL     string  `json:"url"`
	Width   int     `json:"width"`
	Height  int     `json:"height"`
	Entropy float32 `json:"entropy"`
	Size    int     `json:"size"`
}

type Keyword struct {
	ID    []byte `json:"id"`
	Score int    `json:"score"`
	Name  string `json:"name"`
}

type Related struct {
	ID              []byte  `json:"id"`
	Description     string  `json:"description"`
	Title           string  `json:"title"`
	URL             string  `json:"url"`
	ThumbnailWidth  int     `json:"thumbnail_width"`
	Score           float32 `json:"score"`
	ThumbnailHeight int     `json:"thumbnail_height"`
	ThumbnailURL    string  `json:"thumbnail_url"`
}

type Response struct {
	ID              []byte    `json:"id"`
	OriginalURL     string    `json:"original_url"`
	URL             string    `json:"url"`
	Type            string    `json:"type"`
	ProviderName    string    `json:"provider_name"`
	ProviderURL     string    `json:"provider_url"`
	ProviderDisplay string    `json:"provider_display"`
	FaviconURL      string    `json:"favicon_url"`
	Title           string    `json:"title"`
	Description     string    `json:"description"`
	Date            time.Time `json:"date"`
	Authors         []Author  `json:"authors"`
	Published       int64     `json:"published,omitempty"`
	Lead            string    `json:"lead"`
	Content         string    `json:"content"`
	Keywords        []Keyword `json:"keywords"`
	Entities        []Entity  `json:"entities"`
	RelatedArticles []Related `json:"related,omitempty"`
	Images          []Image   `json:"images"`
}
