package main

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

type Response struct {
	OriginalUrl     string    `json:"original_url"`
	Url             string    `json:"url"`
	Type            string    `json:"type"`
	Title           string    `json:"title"`
	Description     string    `json:"description"`
	Safe            bool      `json:"safe"`
	ProviderUrl     string    `json:"provider_url"`
	ProviderName    string    `json:"provider_name"`
	ProviderDisplay string    `json:"provider_display"`
	FaviconUrl      string    `json:"favicon_url"`
	Language        string    `json:"language"`
	Published       int64     `json:"published"`
	Offset          string    `json:"offset"`
	Lead            string    `json:"lead"`
	Keywords        []Keyword `json:"keywords"`
	Entities        []Entity  `json:"entities"`
	Related         []Related `json:"related"`
	Images          []Image   `json:"images"`
}

type Entity struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

type Image struct {
	Url     string  `json:"url"`
	Height  int     `json:"height"`
	Width   int     `json:"width"`
	Colors  []Color `json:"colors"`
	Entropy float32 `json:"entropy"`
	Size    int     `json:"size"`
}

type Color struct {
	Color  []int   `json:"color"`
	Weight float64 `json:"weight"`
}

type Keyword struct {
	Name  string `json:"name"`
	Score int    `json:"score"`
}

type Related struct {
	Url             string `json:"url"`
	Title           string `json:"title"`
	Description     string `json:"description"`
	ThumbnailWidth  int    `json:"thumbnail_width"`
	ThumbnailHeight int    `json:"thumbnail_height"`
	ThumbnailUrl    string `json:"thumbnail_url"`
}
