package model

type WSMessage struct {
	ID      string      `json:"id"`
	From    string      `json:"from"`
	To      string      `json:"to"`
	Type    MessageType `json:"type"`
	Code    int         `json:"code"`
	Content string      `json:"content"`
}

type ChatMsgImage struct {
	ImageName  string `protobuf:"bytes,1,opt,name=image_name,json=imageName,proto3" json:"image_name" form:"image_name" query:"image_name"` // 对应 TS 的 ImageName
	Url        string `protobuf:"bytes,2,opt,name=url,proto3" json:"url" form:"url" query:"url"`
	Compressed string `protobuf:"bytes,3,opt,name=compressed,proto3" json:"compressed" form:"compressed" query:"compressed"` // base64 缩略图
}

type ChatMsgAudio struct {
	AudioName string `protobuf:"bytes,1,opt,name=audio_name,json=audioName,proto3" json:"audio_name" form:"audio_name" query:"audio_name"`
	Url       string `protobuf:"bytes,2,opt,name=url,proto3" json:"url" form:"url" query:"url"`
	Length    string `protobuf:"bytes,3,opt,name=length,proto3" json:"length" form:"length" query:"length"` // "day:hour:minute:second"
}

type ChatMsgVideo struct {
	VideoName string `protobuf:"bytes,1,opt,name=video_name,json=videoName,proto3" json:"video_name" form:"video_name" query:"video_name"`
	Url       string `protobuf:"bytes,2,opt,name=url,proto3" json:"url" form:"url" query:"url"`
	Cover     string `protobuf:"bytes,3,opt,name=cover,proto3" json:"cover" form:"cover" query:"cover"` // base64 缩略图
	Length    string `protobuf:"bytes,4,opt,name=length,proto3" json:"length" form:"length" query:"length"`
}
