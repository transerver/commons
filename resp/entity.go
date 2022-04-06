package resp

import "github.com/Charliego93/go-i18n"

type ResponseEntity struct {
	Code    Code        `json:"code" xml:"code"`
	Message string      `json:"message" xml:"message"`
	Payload interface{} `json:"payload" xml:"payload"`
}

type EntityOp interface{ apply(*ResponseEntity) }
type EntityOpFunc func(*ResponseEntity)
type message string
type payload struct{ v interface{} }
type config i18n.LocalizeConfig
type noI18n string

func (n noI18n) apply(entity *ResponseEntity)       { entity.Message = string(n) }
func (p payload) apply(entity *ResponseEntity)      { entity.Payload = p.v }
func (m message) apply(entity *ResponseEntity)      { entity.Message = i18n.MustTr(string(m)) }
func (i Code) apply(entity *ResponseEntity)         { entity.Code = i }
func (f EntityOpFunc) apply(entity *ResponseEntity) { f(entity) }
func (c *config) apply(entity *ResponseEntity) {
	entity.Message = i18n.MustTr((*i18n.LocalizeConfig)(c))
}

func WithPayload(v interface{}) EntityOp                 { return payload{v} }
func WithMessage(msg string) EntityOp                    { return message(msg) }
func WithNonI18nMsg(msg string) EntityOp                 { return noI18n(msg) }
func WithTemplateConfig(c *i18n.LocalizeConfig) EntityOp { return (*config)(c) }
func WithTemplate(messageId string, template map[string]interface{}) EntityOp {
	return (*config)(&i18n.LocalizeConfig{MessageID: messageId, TemplateData: template})
}

func NewEntity(opts ...EntityOp) *ResponseEntity {
	re := &ResponseEntity{}
	for _, opt := range opts {
		opt.apply(re)
	}
	return re
}

func Msg(code Code, message string) *ResponseEntity {
	return NewEntity(code, WithMessage(message))
}

func FailMsg(msg string) *ResponseEntity {
	return Msg(CodeBaseErr, msg)
}

func FailParams(msg string) *ResponseEntity {
	return Msg(CodeParamErr, msg)
}

func Failf(messageId string, template map[string]interface{}) *ResponseEntity {
	return NewEntity(CodeBaseErr, WithTemplate(messageId, template))
}

func SuccessMsg() *ResponseEntity {
	return Success(nil)
}

func Success(payload interface{}) *ResponseEntity {
	entity := NewEntity(CodeSuccess, WithMessage("操作成功"))
	if payload != nil {
		entity.Payload = payload
	}
	return entity
}
