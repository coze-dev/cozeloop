package errno

import (
	"github.com/coze-dev/coze-loop/backend/pkg/errorx"
)

const (
	baseRespExtraAffectStableKey = "biz_err_affect_stability"
	affectStableValue            = "1"
)

func isStableStatusErrorStable(err errorx.StatusError) bool {
	val, ok := err.Extra()[baseRespExtraAffectStableKey]
	if !ok {
		return true
	}
	return val != affectStableValue
}

func ParseStatusError(err error) (code int32, stable, ok bool) {
	if err == nil {
		return 0, true, false
	}
	if serr, ok := errorx.FromStatusError(err); ok {
		return serr.Code(), isStableStatusErrorStable(serr), true
	}
	return 0, false, false
}
