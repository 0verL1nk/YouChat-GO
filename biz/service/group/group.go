package group

import (
	"core/biz/dal/query"
	"errors"

	"gorm.io/gorm"
)

var (
	ErrGroupNotFound = errors.New("group not found")
)

func GetGroupUserIDs(groupID uint64, userID uint64) (ids []uint64, err error) {
	group, err := query.Q.Group.Where(query.Group.ID.Eq(uint(groupID))).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrGroupNotFound
		}
		return nil, err
	}
	for _, member := range group.Member {
		if member.ID == uint(userID) {
			continue
		}
		ids = append(ids, uint64(member.ID))
	}
	return
}
