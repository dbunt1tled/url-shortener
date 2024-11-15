package urlshort

type URLShort struct {
	ID    int64  `json:"id" jsonapi:"primary,url"`
	URL   string `json:"url" jsonapi:"attr,url"`
	Alias string `json:"alias" jsonapi:"attr,alias"`
}
