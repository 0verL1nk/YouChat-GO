package chat

import (
	"context"
	"strconv"
	"time"

	"core/biz/cerrors"
	"core/biz/chttp"
	"core/biz/dal/model"
	"core/biz/dal/query"
	"core/biz/utils"
	chat "core/hertz_gen/chat"

	"github.com/cloudwego/hertz/pkg/app"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type GetConversationService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewGetConversationService(Context context.Context, RequestContext *app.RequestContext) *GetConversationService {
	return &GetConversationService{RequestContext: RequestContext, Context: Context}
}

func (h *GetConversationService) Run(req *chat.GetConversationReq) (resp *chat.GetConversationResponse, err error) {
	//defer func() {
	// hlog.CtxInfof(h.Context, "req = %+v", req)
	// hlog.CtxInfof(h.Context, "resp = %+v", resp)
	//}()

	// 获取默认值
	page, pageSize := utils.GetDefaultPageParam(int(req.Page), int(req.PageSize))
	after := utils.ParseTimestamp(req.After, time.Time{})
	// 检查会话是否存在
	exists, err := query.Group.WithContext(h.Context).Where(query.Group.ID.Eq(uint(req.GroupID))).Count()
	if err != nil {
		return nil, cerrors.Wrap(err, "err check group exists")
	}
	if exists == 0 {
		return nil, cerrors.ErrGroupNotFound
	}
	// 校验用户是否存在于会话中
	token, err := utils.GetTokenFromMiddleware(h.RequestContext)
	if err != nil {
		return nil, cerrors.Wrap(err, "err get token from middleware")
	}
	isInConversation, err := IsInConversation(h.Context, uint(token.UserId), uint(req.GroupID))
	if err != nil {
		return nil, err
	}
	if !isInConversation {
		return nil, cerrors.ErrUserNotInGroup
	}
	// 获取会话数量
	q := query.ChatMessage.WithContext(h.Context).Where(query.ChatMessage.ToId.Eq(uint(req.GroupID)))
	msgNumTotal, err := q.Count()
	if err != nil {
		return nil, cerrors.Wrap(err, "err get chat message")
	}
	msgs, err := q.Offset((page - 1) * pageSize).
		Limit(pageSize).
		Order(query.ChatMessage.CreatedAt.Desc()).
		Where(query.ChatMessage.CreatedAt.Gt(after)).
		Find()
	if err != nil {
		return nil, cerrors.Wrap(err, "err get chat message")
	}
	// 构建响应
	resp, err = constructConversationResponse(msgs)
	if err != nil {
		return nil, cerrors.Wrap(err, "err construct conversation response")
	}
	resp.Page = int64(page)
	resp.PageSize = int64(pageSize)
	resp.Total = msgNumTotal
	return
}

func IsInConversation(ctx context.Context, userID uint, groupID uint) (bool, error) {
	// 检查用户是否在群组中
	num, err := query.GroupMember.WithContext(ctx).Where(query.GroupMember.UserID.Eq(userID), query.GroupMember.GroupID.Eq(groupID)).Count()
	if err != nil {
		return false, err
	}
	// 如果返回的数量大于0，说明用户在群组中
	return num > 0, nil
}

func constructConversationResponse(msgs []*model.ChatMessage) (*chat.GetConversationResponse, error) {
	resp := &chat.GetConversationResponse{
		Msgs: make([]*chat.ChatMsg, 0, len(msgs)),
	}
	for _, msg := range msgs {
		res := &chat.ChatMsg{}
		res.From = strconv.FormatUint(uint64(msg.FromId), 10)
		res.To = strconv.FormatUint(uint64(msg.ToId), 10)
		res.Type = chat.ChatMsg_Type(msg.MsgType)
		res.CreatedAt = timestamppb.New(msg.CreatedAt)
		switch msg.MsgType {
		case model.MsgTypeText:
			res.Content = &chat.ChatMsg_Text{
				Text: msg.Content,
			}
			// TODO 其他类型
		}
		res.Code = chttp.MESSAGE_SUCCESS
		resp.Msgs = append(resp.Msgs, res)
	}
	return resp, nil
}
