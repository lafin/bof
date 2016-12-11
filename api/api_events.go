package api

// Entripoints for the vk.com
const (
	AuthURL    = "https://oauth.vk.com"
	APIURL     = "https://api.vk.com"
	APIVersion = "5.50"
)

// Groups - struct of json object the Groups
type Groups struct {
	Response []Group
}

// Group - struct of json object the Group
type Group struct {
	AdminLevel int `json:"admin_level"`
	ID         int `json:"id"`
	IsAdmin    int `json:"is_admin"`
	IsClosed   int `json:"is_closed"`
	IsMember   int `json:"is_member"`
	Links      []struct {
		ID       int    `json:"id"`
		Name     string `json:"name"`
		Photo100 string `json:"photo_100"`
		Photo50  string `json:"photo_50"`
		URL      string `json:"url"`
	} `json:"links"`
	Name       string `json:"name"`
	Photo100   string `json:"photo_100"`
	Photo200   string `json:"photo_200"`
	Photo50    string `json:"photo_50"`
	ScreenName string `json:"screen_name"`
	Type       string `json:"type"`
}

// Post - struct of json object the Item
type Post struct {
	Attachments []struct {
		Photo struct {
			AccessKey string `json:"access_key"`
			AlbumID   int    `json:"album_id"`
			Date      int    `json:"date"`
			Height    int    `json:"height"`
			ID        int    `json:"id"`
			OwnerID   int    `json:"owner_id"`
			Photo130  string `json:"photo_130"`
			Photo604  string `json:"photo_604"`
			Photo75   string `json:"photo_75"`
			Text      string `json:"text"`
			UserID    int    `json:"user_id"`
			Width     int    `json:"width"`
		} `json:"photo"`
		Doc struct {
			AccessKey string `json:"access_key"`
			Date      int    `json:"date"`
			Ext       string `json:"ext"`
			ID        int    `json:"id"`
			OwnerID   int    `json:"owner_id"`
			Preview   struct {
				Photo struct {
					Sizes []struct {
						Height int    `json:"height"`
						Src    string `json:"src"`
						Type   string `json:"type"`
						Width  int    `json:"width"`
					} `json:"sizes"`
				} `json:"photo"`
				Video struct {
					FileSize int    `json:"file_size"`
					Height   int    `json:"height"`
					Src      string `json:"src"`
					Width    int    `json:"width"`
				} `json:"video"`
			} `json:"preview"`
			Size  int    `json:"size"`
			Title string `json:"title"`
			Type  int    `json:"type"`
			URL   string `json:"url"`
		} `json:"doc"`
		Type string `json:"type"`
	} `json:"attachments"`
	Comments struct {
		CanPost int `json:"can_post"`
		Count   int `json:"count"`
	} `json:"comments"`
	Date     int `json:"date"`
	FromID   int `json:"from_id"`
	ID       int `json:"id"`
	IsPinned int `json:"is_pinned"`
	Likes    struct {
		CanLike    int `json:"can_like"`
		CanPublish int `json:"can_publish"`
		Count      int `json:"count"`
		UserLikes  int `json:"user_likes"`
	} `json:"likes"`
	OwnerID    int `json:"owner_id"`
	PostSource struct {
		Type string `json:"type"`
	} `json:"post_source"`
	PostType string `json:"post_type"`
	Reposts  struct {
		Count        int `json:"count"`
		UserReposted int `json:"user_reposted"`
	} `json:"reposts"`
	Text string `json:"text"`
}

// Posts - struct of json object the Posts
type Posts struct {
	Response struct {
		Count int    `json:"count"`
		Items []Post `json:"items"`
	} `json:"response"`
}

// ResponseRepost - struct of response after repost of post
type ResponseRepost struct {
	Response struct {
		LikesCount   int `json:"likes_count"`
		PostID       int `json:"post_id"`
		RepostsCount int `json:"reposts_count"`
		Success      int `json:"success"`
	} `json:"response"`
}

// ResponsePost - struct of response after post of post
type ResponsePost struct {
	Response struct {
		PostID int `json:"post_id"`
	} `json:"response"`
}
