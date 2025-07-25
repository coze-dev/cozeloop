// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0
import { z } from 'zod';

export const modelInputSchema = z.object({
  tools: z
    .array(
      z.object({
        type: z.string(),
        function: z.object({
          name: z.string(),
          description: z.string().optional(),
          parameters: z.object({
            required: z.array(z.string()).optional(),
            properties: z
              .record(
                z.object({
                  description: z.string().optional(),
                  type: z.string().optional(),
                }),
              )
              .optional(),
          }),
        }),
      }),
    )
    .optional(),
  messages: z.array(
    z.object({
      role: z.string(),
      content: z.string().optional(),
      reasoning_content: z.string().optional(),
      tool_calls: z
        .array(
          z.object({
            type: z.string(),
            function: z.object({
              name: z.string(),
              arguments: z.string().optional(),
            }),
          }),
        )
        .optional(),
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
    }),
  ),
});

export const modelOutputSchema = z.object({
  choices: z.array(
    z.object({
      index: z.number().optional(),
      message: z.object({
        role: z.string(),
        content: z.string().optional(),
        reasoning_content: z.string().optional(),
        tool_calls: z
          .array(
            z.object({
              type: z.string(),
              function: z.object({
                name: z.string(),
                arguments: z.string().optional(),
              }),
            }),
          )
          .optional(),
        parts: z
          .array(
            z.object({
              type: z.string(),
              text: z.string().optional(),
              file_url: z
                .object({
                  name: z.string().optional(),
                  url: z.string(),
                  detail: z.string().optional(),
                  suffix: z.string().optional(),
                })
                .optional(),
              image_url: z
                .object({
                  name: z.string().optional(),
                  url: z.string().optional(),
                  detail: z.string().optional(),
                })
                .optional(),
            }),
          )
          .optional(),
      }),
    }),
  ),
});

export type ModelInputSchema = z.infer<typeof modelInputSchema>;
export type ModelOutputSchema = z.infer<typeof modelOutputSchema>;
