// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0
import { z } from 'zod';

const inputArgumentSchema = z.array(
  z.object({
    key: z.string(),
    value: z
      .union([
        z.string(),
        z.object({
          content: z.string(),
        }),
      ])
      .optional(),
    source: z.string(),
  }),
);

export const promptInputSchema = z.object({
  arguments: z.union([inputArgumentSchema, z.null()]).optional(),
  templates: z.array(
    z.object({
      role: z.string(),
      content: z.string().optional(),
      reasoning_content: z.string().optional(),
      parts: z
        .array(
          z.object({
            type: z.string(),
            text: z.string().optional(),
            image_url: z
              .object({
                name: z.string().optional(),
                url: z.string(),
                detail: z.string().optional(),
              })
              .optional(),
            file_url: z
              .object({
                name: z.string().optional(),
                url: z.string(),
                detail: z.string().optional(),
                suffix: z.string().optional(),
              })
              .optional(),
          }),
        )
        .optional(),
      name: z.string().optional(),
      tool_calls: z
        .array(
          z.object({
            id: z.string().optional(),
            type: z.string(),
            function: z.object({
              name: z.string(),
              arguments: z.string().optional(),
            }),
          }),
        )
        .optional(),
    }),
  ),
});

const userPromptOutputSchema = z.object({
  prompts: z.union([
    z
      .array(
        z.object({
          role: z.string(),
          content: z.string().optional(),
        }),
      )
      .optional(),
    z.null(),
  ]),
});

const servicePromptOutputSchema = z.union([
  z.array(
    z.object({
      role: z.string(),
      content: z.string().optional(),
    }),
  ),
  z.null(),
]);

export const promptOutputSchema = z.union([
  userPromptOutputSchema,
  servicePromptOutputSchema,
]);

export type PromptInputSchema = z.infer<typeof promptInputSchema>;
export type PromptOutputSchema = z.infer<typeof promptOutputSchema>;
