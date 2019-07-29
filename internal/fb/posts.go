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

type attachments struct {
	Title string `json:"title"`
	Description string `json:"description"`
	Url string `json:"url"`
} 

type PostAttachments struct {
	Data []attachments `json:"data"`	
} 

type Post struct {
	Picture string `json:"picture"`
	Link string `json:"permalink_url"`
	CreatedTime string `json:"created_time"`
	Message string `json:"message"`
	SharedPosts SharedPostLinks `json:"shared_posts"`
	Attachments PostAttachments `json:"attachments"`
}