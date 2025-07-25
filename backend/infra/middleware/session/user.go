// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package session

import "context"

type User struct {
	AppID int32  `json:"app_id,omitempty"`
	ID    string `json:"id"`
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
}

// userKeyType 定义自定义类型作为context键，避免键冲突
type userKeyType struct{}

var userKey = userKeyType{}

// UserIDInCtx returns the user ID from the context.
// Notice: NewSessionMD must be used in your service, or else this function always returns false.
func UserIDInCtx(ctx context.Context) (string, bool) {
	user, ok := ctx.Value(userKey).(*User)
	if !ok || user == nil {
		return "", false
	}
	return user.ID, ok
}

// UserIDInCtxOrEmpty returns the user id from the context, empty string if not present.
func UserIDInCtxOrEmpty(ctx context.Context) string {
	id, _ := UserIDInCtx(ctx)
	return id
}

// AppIDInCtx returns the app id from the context.
// Notice: NewSessionMD must be used in your service, or else this function always returns false.
func AppIDInCtx(ctx context.Context) (int32, bool) {
	user, ok := ctx.Value(userKey).(*User)
	if !ok || user == nil {
		return 0, false
	}
	return user.AppID, ok
}

// AppIDInCtxOrEmpty returns the app id from the context, 0 if not present.
func AppIDInCtxOrEmpty(ctx context.Context) int32 {
	id, _ := AppIDInCtx(ctx)
	return id
}

// UserInCtx returns the user from the context.
func UserInCtx(ctx context.Context) (*User, bool) {
	user, ok := ctx.Value(userKey).(*User)
	if !ok || user == nil {
		return nil, false
	}
	return user, ok
}

func WithCtxUser(ctx context.Context, user *User) context.Context {
	return context.WithValue(ctx, userKey, user)
}
