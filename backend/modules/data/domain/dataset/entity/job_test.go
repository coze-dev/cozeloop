// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package entity

import "testing"

// 测试 JobStatus 的 String 方法
func TestJobStatus_String(t *testing.T) {
	testCases := []struct {
		status   JobStatus
		expected string
	}{
		{JobStatus_Undefined, "Undefined"},
		{JobStatus_Pending, "Pending"},
		{JobStatus_Running, "Running"},
		{JobStatus_Completed, "Completed"},
		{JobStatus_Failed, "Failed"},
		{JobStatus_Cancelled, "Cancelled"},
		{JobStatus(999), "<UNSET>"},
	}

	for _, tc := range testCases {
		result := tc.status.String()
		if result != tc.expected {
			t.Errorf("JobStatus.String() = %s; want %s", result, tc.expected)
		}
	}
}

// 测试 FileFormat 的 String 方法
func TestFileFormat_String(t *testing.T) {
	testCases := []struct {
		format   FileFormat
		expected string
	}{
		{FileFormat_JSONL, "JSONL"},
		{FileFormat_Parquet, "Parquet"},
		{FileFormat_CSV, "CSV"},
		{FileFormat_ZIP, "ZIP"},
		{FileFormat(999), "<UNSET>"},
	}

	for _, tc := range testCases {
		result := tc.format.String()
		if result != tc.expected {
			t.Errorf("FileFormat.String() = %s; want %s", result, tc.expected)
		}
	}
}

// 测试 IsJobTerminal 函数
func TestIsJobTerminal(t *testing.T) {
	testCases := []struct {
		status   JobStatus
		expected bool
	}{
		{JobStatus_Completed, true},
		{JobStatus_Failed, true},
		{JobStatus_Cancelled, true},
		{JobStatus_Undefined, false},
		{JobStatus_Pending, false},
		{JobStatus_Running, false},
	}

	for _, tc := range testCases {
		result := IsJobTerminal(tc.status)
		if result != tc.expected {
			t.Errorf("IsJobTerminal(%d) = %v; want %v", tc.status, result, tc.expected)
		}
	}
}

// 测试 JobType 的 String 方法
func TestJobType_String(t *testing.T) {
	testCases := []struct {
		jobType  JobType
		expected string
	}{
		{JobType_ImportFromFile, "ImportFromFile"},
		{JobType_ExportToFile, "ExportToFile"},
		{JobType_ExportToDataset, "ExportToDataset"},
		{JobType(999), "<UNSET>"},
	}

	for _, tc := range testCases {
		result := tc.jobType.String()
		if result != tc.expected {
			t.Errorf("JobType.String() = %s; want %s", result, tc.expected)
		}
	}
}

// 测试 IOJob 的 IsSetStartedAt 方法
func TestIOJob_IsSetStartedAt(t *testing.T) {
	var startedAt int64 = 1234567890
	testCases := []struct {
		job      *IOJob
		expected bool
	}{
		{&IOJob{StartedAt: &startedAt}, true},
		{&IOJob{StartedAt: nil}, false},
	}

	for _, tc := range testCases {
		result := tc.job.IsSetStartedAt()
		if result != tc.expected {
			t.Errorf("IOJob.IsSetStartedAt() = %v; want %v", result, tc.expected)
		}
	}
}

// 测试 IOJob 的 IsSetEndedAt 方法
func TestIOJob_IsSetEndedAt(t *testing.T) {
	var endedAt int64 = 1234567890
	testCases := []struct {
		job      *IOJob
		expected bool
	}{
		{&IOJob{EndedAt: &endedAt}, true},
		{&IOJob{EndedAt: nil}, false},
	}

	for _, tc := range testCases {
		result := tc.job.IsSetEndedAt()
		if result != tc.expected {
			t.Errorf("IOJob.IsSetEndedAt() = %v; want %v", result, tc.expected)
		}
	}
}

// 测试 JobTypeFromString 函数
func TestJobTypeFromString(t *testing.T) {
	testCases := []struct {
		input       string
		expected    JobType
		expectedErr bool
	}{
		{"ImportFromFile", JobType_ImportFromFile, false},
		{"ExportToFile", JobType_ExportToFile, false},
		{"ExportToDataset", JobType_ExportToDataset, false},
		{"InvalidType", JobType(0), true},
	}

	for _, tc := range testCases {
		result, err := JobTypeFromString(tc.input)
		if (err != nil) != tc.expectedErr {
			t.Errorf("JobTypeFromString(%s) error = %v; wantErr %v", tc.input, err, tc.expectedErr)
		}
		if result != tc.expected {
			t.Errorf("JobTypeFromString(%s) = %v; want %v", tc.input, result, tc.expected)
		}
	}
}

// 测试 JobStatusFromString 函数
func TestJobStatusFromString(t *testing.T) {
	testCases := []struct {
		input       string
		expected    JobStatus
		expectedErr bool
	}{
		{"Undefined", JobStatus_Undefined, false},
		{"Pending", JobStatus_Pending, false},
		{"Running", JobStatus_Running, false},
		{"Completed", JobStatus_Completed, false},
		{"Failed", JobStatus_Failed, false},
		{"Cancelled", JobStatus_Cancelled, false},
		{"InvalidStatus", JobStatus(0), true},
	}

	for _, tc := range testCases {
		result, err := JobStatusFromString(tc.input)
		if (err != nil) != tc.expectedErr {
			t.Errorf("JobStatusFromString(%s) error = %v; wantErr %v", tc.input, err, tc.expectedErr)
		}
		if result != tc.expected {
			t.Errorf("JobStatusFromString(%s) = %v; want %v", tc.input, result, tc.expected)
		}
	}
}

// 测试 IOJob 的 GetID 方法
func TestIOJob_GetID(t *testing.T) {
	testCases := []struct {
		job      *IOJob
		expected int64
	}{
		{&IOJob{ID: 123}, 123},
		{nil, 0},
	}

	for _, tc := range testCases {
		result := tc.job.GetID()
		if result != tc.expected {
			t.Errorf("IOJob.GetID() = %d; want %d", result, tc.expected)
		}
	}
}

// 测试 IOJob 的 SetID 方法
func TestIOJob_SetID(t *testing.T) {
	job := &IOJob{}
	newID := int64(456)
	job.SetID(newID)
	if job.ID != newID {
		t.Errorf("IOJob.SetID(%d) failed; got ID %d", newID, job.ID)
	}
}
