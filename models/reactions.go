package models

import (
	"github.com/jinzhu/gorm"
)

type Reaction interface {
	GetMessageID() int
	GetUserID() int
}

type LikeReaction struct {
	gorm.Model
	MessageID int
	UserID    int
}

func (r *LikeReaction) GetMessageID() int {
	return r.MessageID
}

func (r *LikeReaction) GetUserID() int {
	return r.UserID
}

type LolReaction struct {
	gorm.Model
	MessageID int
	UserID    int
}

func (r *LolReaction) GetMessageID() int {
	return r.MessageID
}

func (r *LolReaction) GetUserID() int {
	return r.UserID
}

type FacepalmReaction struct {
	gorm.Model
	MessageID int
	UserID    int
}

func (r *FacepalmReaction) GetMessageID() int {
	return r.MessageID
}

func (r *FacepalmReaction) GetUserID() int {
	return r.UserID
}
