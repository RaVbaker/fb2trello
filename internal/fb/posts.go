package fb

type Result struct {
	Data []Post `json:"data"`
}

type shares struct {
	Link string `json:"link"`
}

type SharedPostLinks struct {
	Data []shares `json:"data"`
}

type media struct {
	Src    string `json:"src"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

type attachments struct {
	Title       string `json:"title"`
	Description string `json:"description"` // optional
	Url         string `json:"url"`
	Media       media  `json:"media"` // optional
	Type        string `json:"type"`
}

type PostAttachments struct {
	Data []attachments `json:"data"`
}

type messageTags struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Post struct {
	Id          string          `json:"id"`
	Picture     string          `json:"full_picture"`
	Link        string          `json:"permalink_url"`
	StatusType  string          `json:"status_type"`
	CreatedTime string          `json:"created_time"`
	Message     string          `json:"message"`
	SharedPosts SharedPostLinks `json:"sharedposts"`
	Attachments PostAttachments `json:"attachments"`
	MessageTags []messageTags   `json:"message_tags"`
}
