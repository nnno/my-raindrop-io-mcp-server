package entity

import "time"

type BookmarkType string

const (
	BookmarkTypeLink     BookmarkType = "link"
	BookmarkTypeArticle  BookmarkType = "article"
	BookmarkTypeImage    BookmarkType = "image"
	BookmarkTypeVideo    BookmarkType = "video"
	BookmarkTypeDocument BookmarkType = "document"
	BookmarkTypeAudio    BookmarkType = "audio"
)

type Bookmark struct {
	ID           int          `json:"id"`
	Title        string       `json:"title"`
	Link         string       `json:"link"`
	Excerpt      string       `json:"excerpt,omitempty"`
	Note         string       `json:"note,omitempty"`
	Tags         []string     `json:"tags,omitempty"`
	Type         BookmarkType `json:"type"`
	Created      time.Time    `json:"created"`
	LastUpdate   time.Time    `json:"lastUpdate"`
	CollectionID int          `json:"collectionId"`
	Highlights   []Highlight  `json:"highlights,omitempty"`
}

type Highlight struct {
	Text string `json:"text"`
	Note string `json:"note,omitempty"`
}

type SearchParams struct {
	CollectionID int
	Search       string
	Sort         string
	Page         int
	Perpage      int
}

type CreateParams struct {
	Link         string
	Title        string
	Tags         []string
	CollectionID int
	Note         string
}

type UpdateParams struct {
	Title        *string
	Tags         *[]string // nil = 変更なし, &[]string{} = タグ全削除
	CollectionID *int
	Note         *string
}
