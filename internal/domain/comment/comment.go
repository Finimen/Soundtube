package comment

import (
	"errors"
)

type Comment struct {
	id              int
	content         string
	isResponse      bool
	responeTargetID int
}

func (c *Comment) ID() int               { return c.id }
func (c *Comment) Content() string       { return c.content }
func (c *Comment) IsResponse() bool      { return c.isResponse }
func (c *Comment) ResponseTargetID() int { return c.responeTargetID }

func NewComment(content string, isResponse bool, responseTargetID int) (*Comment, error) {
	if content == "" {
		return nil, errors.New("content is requered")
	}
	return &Comment{
		content:         content,
		isResponse:      isResponse,
		responeTargetID: responseTargetID,
	}, nil
}

func RestoreCommentFromStorage(id int, content string, isResponse bool, responseTargetID int, likes, dislikes int) *Comment {
	return &Comment{
		id:              id,
		content:         content,
		isResponse:      isResponse,
		responeTargetID: responseTargetID,
	}
}
