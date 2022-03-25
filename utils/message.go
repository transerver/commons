package utils

import (
	"github.com/transerver/commons/resp"
)

type Message struct {
	Success   bool
	MessageId string
	Code      resp.Code
	Data      interface{}
	Template  map[string]interface{}
}

type MessageOp interface{ apply(*Message) }
type MessageOpFunc func(*Message)
type success bool
type msg string
type msgCode resp.Code
type data struct{ v interface{} }
type template struct {
	messageId string
	data      map[string]interface{}
}

func (t template) apply(msg *Message)      { msg.MessageId = t.messageId; msg.Template = t.data }
func (d data) apply(msg *Message)          { msg.Data = d.v }
func (m msg) apply(msg *Message)           { msg.MessageId = string(m) }
func (f MessageOpFunc) apply(msg *Message) { f(msg) }
func (s success) apply(msg *Message) {
	msg.Success = bool(s)
	msg.Code = resp.CodeSuccess
	msg.MessageId = msg.Code.String()
}
func (c msgCode) apply(msg *Message) {
	msg.Code = resp.Code(c)
	msg.MessageId = msg.Code.String()
	if msg.Code == resp.CodeSuccess {
		msg.Success = true
	}
}

func WithSuccess() MessageOp               { return success(true) }
func WithCode(code resp.Code) MessageOp    { return msgCode(code) }
func WithMessage(message string) MessageOp { return msg(message) }
func WithData(v interface{}) MessageOp     { return data{v} }
func WithTemplate(messageId string, v map[string]interface{}) MessageOp {
	return template{messageId, v}
}

func Success() Message {
	return Message{Success: true}
}

func NewMessage(opts ...MessageOp) Message {
	msg := &Message{Success: false, Code: resp.CodeBaseErr, MessageId: resp.CodeBaseErr.String()}
	for _, opt := range opts {
		opt.apply(msg)
	}
	return *msg
}

func NewFailMessage(code resp.Code, msg string) Message {
	return Message{Success: false, Code: code, MessageId: msg}
}

func NewFailMessageWithTemplate(code resp.Code, msg string, template map[string]interface{}) Message {
	m := &Message{Success: false, Code: code}
	WithTemplate(msg, template).apply(m)
	return *m
}

func (m *Message) ToResponse() *resp.ResponseEntity {
	if len(m.Template) > 0 {
		return resp.NewEntity(m.Code, resp.WithTemplate(m.MessageId, m.Template))
	}
	return resp.NewEntity(m.Code, resp.WithMessage(m.MessageId))
}
