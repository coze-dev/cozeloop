package errno

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/coze-dev/coze-loop/backend/pkg/lang/conv"
)

func TestError(t *testing.T) {
	err := NewMQRetryErr("test")
	assert.True(t, NeedMQRetry(err))

	err = errors.New("cause")
	err = WrapMQRetryErr(err)
	assert.True(t, NeedMQRetry(err))

	err = errors.New("cause")
	assert.False(t, NeedMQRetry(err))
}

func TestConv(t *testing.T) {
	t.Log(conv.UnsafeBytesToString(nil))
}

func TestPersistent(t *testing.T) {
	err := NewTargetResultErr("test")
	ok, _ := ParseTurnOtherErr(err)
	assert.False(t, ok)

	err = DeserializeErr(conv.UnsafeStringToBytes(SerializeErr(err)))
	ei, ok := ParseErrImpl(err)
	assert.True(t, ok)
	assert.Equal(t, ei.Code, targetResultErrCode)
	assert.Equal(t, "test", ei.Msg)

	err = NewTurnOtherErr("other msg", errors.New("err content"))
	ok, msg := ParseTurnOtherErr(err)
	assert.True(t, ok)
	assert.Equal(t, "other msg", msg)

	err = DeserializeErr(conv.UnsafeStringToBytes(SerializeErr(err)))
	ok, msg = ParseTurnOtherErr(err)
	assert.True(t, ok)
	assert.Equal(t, "other msg", msg)
}

func TestErrImpl_CauseMsg(t *testing.T) {
	err := &ErrImpl{
		Code:  turnOtherErrCode,
		Msg:   "msg",
		Cause: errors.New("cause"),
	}
	assert.Equal(t, "cause", err.CauseMsg())
	t.Log(error(err).Error())

	recursionErr := err.SetErrMsg("update msg").SetCause(err)
	t.Log(recursionErr.Error())

	err1 := &ErrImpl{
		Code:  turnOtherErrCode,
		Msg:   "msg",
		Cause: errors.New("cause"),
	}
	cerr := CloneErr(err1)
	err1 = err1.SetErrMsg("update msg").SetCause(cerr)
	t.Log(err1.Error())

	err2 := &ErrImpl{
		Code: turnOtherErrCode,
		Msg:  "msg",
	}
	err2.Cause = err2
	t.Log(err2.Error())
}
