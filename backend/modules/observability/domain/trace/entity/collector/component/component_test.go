// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package component

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/coze-dev/cozeloop/backend/modules/observability/domain/trace/entity/collector/component/mocks"
)

func TestID_String(t *testing.T) {
	tests := []struct {
		name string
		id   ID
		want string
	}{
		{
			name: "empty name",
			id:   ID{typeVal: "testType", nameVal: ""},
			want: "testType",
		},
		{
			name: "with name",
			id:   ID{typeVal: "testType", nameVal: "testName"},
			want: "testType/testName",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.id.String())
		})
	}
}

func TestID_Type(t *testing.T) {
	id := ID{typeVal: "testType"}
	assert.Equal(t, Type("testType"), id.Type())
}

func TestID_UnmarshalText(t *testing.T) {
	tests := []struct {
		name    string
		text    string
		want    ID
		err     bool
		errText string
	}{
		{
			name: "valid with name",
			text: "testType/testName",
			want: ID{typeVal: "testType", nameVal: "testName"},
		},
		{
			name: "valid without name",
			text: "testType",
			want: ID{typeVal: "testType"},
		},
		{
			name:    "empty",
			text:    "",
			err:     true,
			errText: "id must not be empty",
		},
		{
			name:    "empty type",
			text:    "/testName",
			err:     true,
			errText: "the part before / should not be empty",
		},
		{
			name:    "empty name",
			text:    "testType/",
			err:     true,
			errText: "the part after / should not be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var id ID
			err := id.UnmarshalText([]byte(tt.text))
			if tt.err {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errText)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, id)
			}
		})
	}
}

func TestComponentLifecycle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockComp := mocks.NewMockComponent(ctrl)

	ctx := context.Background()

	mockComp.EXPECT().Start(ctx).Return(nil)
	mockComp.EXPECT().Shutdown(ctx).Return(nil)

	assert.NoError(t, mockComp.Start(ctx))
	assert.NoError(t, mockComp.Shutdown(ctx))
}

func TestFactoryMap(t *testing.T) {
	factoryMap := make(map[Type]Factory)

	factory := &mockFactory{typ: "testType"}
	factoryMap[factory.Type()] = factory

	assert.Equal(t, factory, factoryMap["testType"])
}

type mockFactory struct {
	typ Type
}

func (m *mockFactory) Type() Type {
	return m.typ
}

func (m *mockFactory) CreateDefaultConfig() Config {
	return nil
}

func TestCreateDefaultConfigFunc(t *testing.T) {
	called := false
	f := CreateDefaultConfigFunc(func() Config {
		called = true
		return nil
	})

	f.CreateDefaultConfig()
	assert.True(t, called)
}
