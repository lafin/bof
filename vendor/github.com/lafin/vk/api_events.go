package vk

// Entripoints for the vk.com
const (
	AuthURL     = "https://oauth.vk.com"
	APIURL      = "https://api.vk.com"
	APIVersion  = "5.50"
	Permissions = "wall,groups,friends,photos,video,docs"
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

type DocPreview struct {
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
}

// GetSmallPreview - return preview with type "s" for gif's
func (r *DocPreview) GetSmallPreview() string {
	sizes := r.Photo.Sizes
	for i := 0; i < len(sizes); i++ {
		size := sizes[i]
		if size.Type == "s" {
			return size.Src
		}
	}
	return ""
}

// Post - struct of json object the Item
type Post struct {
	Attachments []struct {
		Video struct {
			AccessKey   string `json:"access_key"`
			CanAdd      int    `json:"can_add"`
			Comments    int    `json:"comments"`
			Date        int    `json:"date"`
			Description string `json:"description"`
			Duration    int    `json:"duration"`
			ID          int    `json:"id"`
			OwnerID     int    `json:"owner_id"`
			Photo130    string `json:"photo_130"`
			Photo320    string `json:"photo_320"`
			Photo640    string `json:"photo_640"`
			Platform    string `json:"platform"`
			Title       string `json:"title"`
			Views       int    `json:"views"`
		} `json:"video"`
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
			Preview   DocPreview
			Size      int    `json:"size"`
			Title     string `json:"title"`
			Type      int    `json:"type"`
			URL       string `json:"url"`
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
	Error struct {
		ErrorCode     int    `json:"error_code"`
		ErrorMsg      string `json:"error_msg"`
		RequestParams []struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		} `json:"request_params"`
	} `json:"error"`
	Response struct {
		PostID int `json:"post_id"`
	} `json:"response"`
}

// ResponseUsersOfGroup - struct of response list users of group
type ResponseUsersOfGroup struct {
	Response struct {
		Count int                   `json:"count"`
		Items []ResponseUserOfGroup `json:"items"`
	} `json:"response"`
}

// ResponseUserOfGroup - struct of response list user of group
type ResponseUserOfGroup struct {
	FirstName   string `json:"first_name"`
	ID          int    `json:"id"`
	LastName    string `json:"last_name"`
	Deactivated string `json:"deactivated"`
	LastSeen    struct {
		Platform int `json:"platform"`
		Time     int `json:"time"`
	} `json:"last_seen"`
}

// IsBanned - returned true if the user was banned
func (r *ResponseUserOfGroup) IsBanned() bool {
	return r.Deactivated == "banned"
}

// IsDeleted - returned true if the user was deleted
func (r *ResponseUserOfGroup) IsDeleted() bool {
	return r.Deactivated == "deleted"
}

// ResponseRemoveUser - struct of response status of removing the user
type ResponseRemoveUser struct {
	Response int `json:"response"`
}
