package errno

const (
	mqRetryErrCode = 1

	targetResultErrCode    = 11
	evaluatorResultErrCode = 12
	turnOtherErrCode       = 13

	ServiceInternalErrMsg = "Server internal error"
)

func NeedMQRetry(err error) bool {
	ei, ok := ParseErrImpl(err)
	if ok && ei.Code == mqRetryErrCode {
		return true
	}
	return false
}

func NewMQRetryErr(msg string) error {
	return &ErrImpl{
		Code: mqRetryErrCode,
		Msg:  msg,
	}
}

func WrapMQRetryErr(err error) error {
	return &ErrImpl{
		Code:  mqRetryErrCode,
		Cause: err,
	}
}

func WrapTargetResultErr(err error) error {
	return &ErrImpl{
		Code:  targetResultErrCode,
		Cause: err,
	}
}

func NewTargetResultErr(msg string) error {
	return &ErrImpl{
		Code: targetResultErrCode,
		Msg:  msg,
	}
}

func WrapEvaluatorResultErr(err error) error {
	return &ErrImpl{
		Code:  evaluatorResultErrCode,
		Cause: err,
	}
}

func NewEvaluatorResultErr(msg string) error {
	return &ErrImpl{
		Code: evaluatorResultErrCode,
		Msg:  msg,
	}
}

func WrapTurnOtherErr(err error) error {
	return &ErrImpl{
		Code:  turnOtherErrCode,
		Cause: err,
	}
}

func NewTurnOtherErr(msg string, err error) error {
	return &ErrImpl{
		Code:  turnOtherErrCode,
		Msg:   msg,
		Cause: err,
	}
}

func ParseTurnOtherErr(err error) (bool, string) {
	ei, ok := ParseErrImpl(err)
	if ok && ei.Code == turnOtherErrCode {
		return true, ei.ErrMsg()
	}
	return false, ""
}
