package group

import (
	"core/biz/dal/query"
	"errors"

	"github.com/google/uuid"
)

var (
	ErrGroupNotFound = errors.New("group not found")
)

func GetGroupUserIDs(groupID uuid.UUID, userID uuid.UUID) (ids []uuid.UUID, err error) {
	groupMems, err := query.Q.GroupMember.Where(query.GroupMember.GroupID.Eq(groupID)).Find()
	if err != nil {
		return nil, err
	}
	if len(groupMems) == 0 {
		return nil, ErrGroupNotFound
	}
	for _, mem := range groupMems {
		if mem.UserID == userID {
			continue
		}
		ids = append(ids, mem.UserID)
	}
	return
}
