package entity

type Collection struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Count    int    `json:"count"`
	ParentID int    `json:"parentId"` // 0 = ルートコレクション
}

type CollectionCreateParams struct {
	Title    string
	ParentID int // 0 = ルート
}

type CollectionUpdateParams struct {
	Title *string // nilで変更なし
}
