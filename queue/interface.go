package queue

type Msg interface {
	Id() string
	Body() interface{}
}

type TaskQueue interface {
	GetMsgExclusively(topic string, size int) ([]Msg, error)
	AckMsg(topic string, msg Msg, done bool) error
	WriteMsgs(topic string, objs []interface{}) error
}
