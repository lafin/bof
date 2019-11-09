package vk

// Entripoints for the vk.com
const (
	AuthURL     = "https://oauth.vk.com"
	APIURL      = "https://api.vk.com"
	APIVersion  = "5.103"
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

// DocPreview - struct of json object the DocPreview
type DocPreview struct {
	Photo struct {
		Sizes []struct {
			Height float32 `json:"height"`
			Src    string  `json:"src"`
			Type   string  `json:"type"`
			Width  float32 `json:"width"`
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
	for _, size := range r.Photo.Sizes {
		if size.Type == "s" {
			return size.Src
		}
	}
	return ""
}

// Post - struct of json object the Item
type Post struct {
	Attachments []struct {
		Album struct {
			Created     int    `json:"created"`
			Description string `json:"description"`
			ID          string `json:"id"`
			OwnerID     int    `json:"owner_id"`
			Size        int    `json:"size"`
			Thumb       struct {
				AccessKey string `json:"access_key"`
				AlbumID   int    `json:"album_id"`
				Date      int    `json:"date"`
				ID        int    `json:"id"`
				OwnerID   int    `json:"owner_id"`
				Sizes     []struct {
					Height int    `json:"height"`
					Type   string `json:"type"`
					URL    string `json:"url"`
					Width  int    `json:"width"`
				} `json:"sizes"`
				Text string `json:"text"`
			} `json:"thumb"`
			Title   string `json:"title"`
			Updated int    `json:"updated"`
		} `json:"album"`
		Link struct {
			Description string `json:"description"`
			IsFavorite  bool   `json:"is_favorite"`
			Target      string `json:"target"`
			Title       string `json:"title"`
			URL         string `json:"url"`
		} `json:"link"`
		Doc struct {
			AccessKey string     `json:"access_key"`
			Date      int        `json:"date"`
			Ext       string     `json:"ext"`
			ID        int        `json:"id"`
			OwnerID   int        `json:"owner_id"`
			Preview   DocPreview `json:"preview"`
			Size      int        `json:"size"`
			Title     string     `json:"title"`
			Type      int        `json:"type"`
			URL       string     `json:"url"`
		} `json:"doc"`
		Photo struct {
			AccessKey string  `json:"access_key"`
			AlbumID   int     `json:"album_id"`
			Date      int     `json:"date"`
			ID        int     `json:"id"`
			Lat       float64 `json:"lat"`
			Long      float64 `json:"long"`
			OwnerID   int     `json:"owner_id"`
			PostID    int     `json:"post_id"`
			Sizes     []struct {
				Height int    `json:"height"`
				Type   string `json:"type"`
				URL    string `json:"url"`
				Width  int    `json:"width"`
			} `json:"sizes"`
			Text   string `json:"text"`
			UserID int    `json:"user_id"`
		} `json:"photo"`
		Type  string `json:"type"`
		Video struct {
			AccessKey   string `json:"access_key"`
			CanAdd      int    `json:"can_add"`
			Comments    int    `json:"comments"`
			Date        int    `json:"date"`
			Description string `json:"description"`
			Duration    int    `json:"duration"`
			FirstFrame  []struct {
				Height int    `json:"height"`
				URL    string `json:"url"`
				Width  int    `json:"width"`
			} `json:"first_frame"`
			Height int `json:"height"`
			ID     int `json:"id"`
			Image  []struct {
				Height      int    `json:"height"`
				URL         string `json:"url"`
				Width       int    `json:"width"`
				WithPadding int    `json:"with_padding"`
			} `json:"image"`
			IsFavorite bool   `json:"is_favorite"`
			OwnerID    int    `json:"owner_id"`
			Title      string `json:"title"`
			TrackCode  string `json:"track_code"`
			Type       string `json:"type"`
			UserID     int    `json:"user_id"`
			Views      int    `json:"views"`
			Width      int    `json:"width"`
		} `json:"video"`
	} `json:"attachments"`
	Comments struct {
		CanPost       int  `json:"can_post"`
		Count         int  `json:"count"`
		GroupsCanPost bool `json:"groups_can_post"`
	} `json:"comments"`
	CopyHistory []struct {
		Attachments []struct {
			Photo struct {
				AccessKey string `json:"access_key"`
				AlbumID   int    `json:"album_id"`
				Date      int    `json:"date"`
				ID        int    `json:"id"`
				OwnerID   int    `json:"owner_id"`
				PostID    int    `json:"post_id"`
				Sizes     []struct {
					Height int    `json:"height"`
					Type   string `json:"type"`
					URL    string `json:"url"`
					Width  int    `json:"width"`
				} `json:"sizes"`
				Text   string `json:"text"`
				UserID int    `json:"user_id"`
			} `json:"photo"`
			Type string `json:"type"`
		} `json:"attachments"`
		Date       int `json:"date"`
		FromID     int `json:"from_id"`
		ID         int `json:"id"`
		OwnerID    int `json:"owner_id"`
		PostSource struct {
			Platform string `json:"platform"`
			Type     string `json:"type"`
		} `json:"post_source"`
		PostType string `json:"post_type"`
		Text     string `json:"text"`
	} `json:"copy_history"`
	Copyright struct {
		ID   int    `json:"id"`
		Link string `json:"link"`
		Name string `json:"name"`
		Type string `json:"type"`
	} `json:"copyright"`
	Date       int  `json:"date"`
	Edited     int  `json:"edited"`
	FromID     int  `json:"from_id"`
	ID         int  `json:"id"`
	IsFavorite bool `json:"is_favorite"`
	IsPinned   int  `json:"is_pinned"`
	Likes      struct {
		CanLike    int `json:"can_like"`
		CanPublish int `json:"can_publish"`
		Count      int `json:"count"`
		UserLikes  int `json:"user_likes"`
	} `json:"likes"`
	MarkedAsAds int `json:"marked_as_ads"`
	OwnerID     int `json:"owner_id"`
	PostSource  struct {
		Platform string `json:"platform"`
		Type     string `json:"type"`
	} `json:"post_source"`
	PostType string `json:"post_type"`
	Reposts  struct {
		Count        int `json:"count"`
		UserReposted int `json:"user_reposted"`
	} `json:"reposts"`
	SignerID int    `json:"signer_id"`
	Text     string `json:"text"`
	Views    struct {
		Count int `json:"count"`
	} `json:"views"`
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
