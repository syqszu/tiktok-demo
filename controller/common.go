package controller

type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

type Video struct {
	Id            int64  `gorm:"primaryKey" json:"id,omitempty"`
	AuthorID      int64  `gorm:"foreignKey:Id" json:"-"` // 使用外键AuthorID关联User表中的Id字段
	Author        User   `gorm:"foreignKey:AuthorID; constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"author"`
	PlayUrl       string `gorm:"size:255" json:"play_url,omitempty"`
	CoverUrl      string `gorm:"size:255" json:"cover_url,omitempty"`
	FavoriteCount int64  `json:"favorite_count,omitempty"`
	CommentCount  int64  `json:"comment_count,omitempty"`
	IsFavorite    bool   `gorm:"-" json:"is_favorite"` // 在返回时根据用户是否收藏该视频进行赋值，不保存在数据库中
	UploadTime    int64  `json:"-"`
	Title         string `json:"title"` // 视频标题
}

type Comment struct {
	Id         int64  `json:"id,omitempty"`
	UserID     int64  `gorm:"foreignkey:UserRefer"` // 该字段将作为外键与User表中的Id字段关联
	User       User   `json:"user"`
	VideoID    int64  `json:"video_id,omitempty"` // 添加的字段
	Content    string `json:"content,omitempty"`
	CreateDate string `json:"create_date,omitempty"`
}

type User struct {
	Id              int64   `gorm:"primary_key" json:"id,omitempty"`
	Name            string  `json:"name,omitempty"`
	FollowCount     int64   `json:"follow_count,omitempty"`
	FollowerCount   int64   `json:"follower_count,omitempty"`
	Avatar          string  `json:"avatar,omitempty"`
	BackgroundImage string  `json:"background_image,omitempty"`
	Signature       string  `json:"signature,omitempty"`
	TotalFavorited  int64   `json:"total_favorited,omitempty"`
	WorkCount       int64   `json:"work_count,omitempty"`
	FavoriteCount   int64   `json:"favorite_count,omitempty"`
	Token           string  `json:"token,omitempty"`
	FavoritedVideos []Video `gorm:"many2many:video_favorites;" json:"-"`
}

type Message struct {
	Id         int64  `json:"id,omitempty"`
	Content    string `json:"content,omitempty"`
	CreateTime string `json:"create_time,omitempty"`
}

type MessageSendEvent struct {
	UserId     int64  `json:"user_id,omitempty"`
	ToUserId   int64  `json:"to_user_id,omitempty"`
	MsgContent string `json:"msg_content,omitempty"`
}

type MessagePushEvent struct {
	FromUserId int64  `json:"user_id,omitempty"`
	MsgContent string `json:"msg_content,omitempty"`
}
