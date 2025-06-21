package chat

import (
	"context"
	"time"

	"core/biz/cerrors"
	"core/biz/chttp"
	"core/biz/dal/model"
	"core/biz/dal/query"
	"core/biz/utils"
	chat "core/hertz_gen/chat"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/google/uuid"
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
	page, pageSize := utils.GetDefaultPageParam(req.Page, req.PageSize)
	after := utils.ParseUnixTime(req.After, time.Unix(0, 0))
	// parse UUID
	groupID, err := uuid.Parse(req.GroupID)
	if err != nil {
		return nil, cerrors.Wrap(err, "err parse groupID")
	}

	// 检查会话是否存在
	exists, err := query.Group.WithContext(h.Context).Where(query.Group.ID.Eq(groupID)).Count()
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
	isInConversation, err := IsInConversation(h.Context, token.UserId, groupID)
	if err != nil {
		return nil, err
	}
	if !isInConversation {
		return nil, cerrors.ErrUserNotInGroup
	}
	// 获取会话数量
	q := query.ChatMessage.WithContext(h.Context).Where(query.ChatMessage.ToId.Eq(groupID))
	msgs, count, err := q.
		Order(query.ChatMessage.CreatedAt.Asc()).
		Where(query.ChatMessage.CreatedAt.Gt(after)).
		FindByPage((page-1)*pageSize, pageSize)
	if err != nil {
		hlog.Error(cerrors.Wrap(err, "err get chat message"))
		return nil, cerrors.ErrFetchDataFromDatabase
	}
	hlog.Debugf("get conversation msgs: %v, page: %d, pageSize: %d, total: %d, after:%v", msgs, page, pageSize, count, after)
	// 构建响应
	resp, err = constructConversationResponse(msgs)
	if err != nil {
		return nil, cerrors.Wrap(err, "err construct conversation response")
	}
	resp.Page = int64(page)
	resp.PageSize = int64(pageSize)
	resp.Total = count
	return
}

func IsInConversation(ctx context.Context, userID uuid.UUID, groupID uuid.UUID) (bool, error) {
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
		res := &chat.ChatMsg{
			Id:        msg.ID.String(),
			From:      msg.FromId.String(),
			To:        msg.ToId.String(),
			Type:      chat.ChatMsg_Type(msg.MsgType),
			CreatedAt: msg.CreatedAt.Unix(),
		}
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
