// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package entity

type Options struct {
	// Temperature is the temperature for the model, which controls the randomness of the model.
	Temperature *float32
	// MaxTokens is the max number of tokens, if reached the max tokens, the model will stop generating, and mostly return an finish reason of "length".
	MaxTokens *int
	// Model is the model name.
	Model *string
	// TopP is the top p for the model, which controls the diversity of the model.
	TopP *float32
	// Stop is the stop words for the model, which controls the stopping condition of the model.
	Stop []string
	// Tools is a list of tools the model may call.
	Tools []*ToolInfo
	// ToolChoice controls which tool is called by the model.
	ToolChoice *ToolChoice
}

type Option struct {
	apply func(opts *Options)
}

func ApplyOptions(base *Options, opts ...Option) *Options {
	if base == nil {
		base = &Options{}
	}
	for _, opt := range opts {
		opt.apply(base)
	}
	return base
}

func WithTemperature(t float32) Option {
	return Option{
		apply: func(opts *Options) {
			opts.Temperature = &t
		},
	}
}

func WithMaxTokens(m int) Option {
	return Option{
		apply: func(opts *Options) {
			opts.MaxTokens = &m
		},
	}
}

func WithModel(m string) Option {
	return Option{
		apply: func(opts *Options) {
			opts.Model = &m
		},
	}
}

func WithTopP(t float32) Option {
	return Option{
		apply: func(opts *Options) {
			opts.TopP = &t
		},
	}
}

func WithStop(s []string) Option {
	return Option{
		apply: func(opts *Options) {
			opts.Stop = s
		},
	}
}

func WithTools(t []*ToolInfo) Option {
	return Option{
		apply: func(opts *Options) {
			opts.Tools = t
		},
	}
}

func WithToolChoice(t *ToolChoice) Option {
	return Option{
		apply: func(opts *Options) {
			opts.ToolChoice = t
		},
	}
}
