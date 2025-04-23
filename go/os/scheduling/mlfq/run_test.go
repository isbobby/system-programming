package main

import (
	"testing"
	"time"
)

func TestMLFQDemotion(t *testing.T) {
	inputs := []*Job{
		NewJob(1, 0, []JobInput{CPUInstruction, CPUInstruction, CPUInstruction, CPUInstruction, CPUInstruction, CPUInstruction}),
	}

	MLFQConfig := MLFQConfig{
		QueueConfigs: []QueueConfig{
			{Priority: 2, TimeAllotment: 1},
			{Priority: 1, TimeAllotment: 2},
			{Priority: 0, TimeAllotment: 3},
		},
		ResetInterval: 5,
		QueueSize:     100,
	}

	out := RunSystem(&MLFQConfig, inputs, time.Duration(2*time.Second), true)

	expectedSystemBehaviour := []AuditLog{
		{Action: EXEC, JobID: "1"},
		{Action: EXPIRE, JobID: "1"},
		{Action: EXEC, JobID: "1"},
		{Action: EXEC, JobID: "1"},
		{Action: EXPIRE, JobID: "1"},
		{Action: EXEC, JobID: "1"},
		{Action: EXEC, JobID: "1"},
		{Action: EXEC, JobID: "1"},
		{Action: COMPLETE, JobID: "1"},
	}

	mismatched := -1
	for i := range expectedSystemBehaviour {
		if expectedSystemBehaviour[i].Action != out[i].Action || expectedSystemBehaviour[i].JobID != out[i].JobID {
			mismatched = i
			break
		}
	}

	if mismatched != -1 {
		t.Errorf("Mismatch detected at index %v", mismatched)
	}
}

func TestMLFQBasicCaseNoIO(t *testing.T) {
	inputs := []*Job{
		NewJob(1, 0, []JobInput{CPUInstruction, CPUInstruction, CPUInstruction, CPUInstruction}),
		NewJob(2, 0, []JobInput{CPUInstruction, CPUInstruction, CPUInstruction, CPUInstruction}),
	}

	MLFQConfig := MLFQConfig{
		QueueConfigs: []QueueConfig{
			{Priority: 4, TimeAllotment: 2},
			{Priority: 3, TimeAllotment: 3},
			{Priority: 2, TimeAllotment: 4},
			{Priority: 1, TimeAllotment: 5},
			{Priority: 0, TimeAllotment: 6},
		},
		ResetInterval: 5,
		QueueSize:     100,
	}

	// TODO: manually configured system run time and clock delay
	out := RunSystem(&MLFQConfig, inputs, time.Duration(2*time.Second), true)

	expectedSystemBehaviour := []AuditLog{
		{Action: EXEC, JobID: "1"},
		{Action: EXEC, JobID: "1"},
		{Action: EXPIRE, JobID: "1"},
		{Action: EXEC, JobID: "2"},
		{Action: EXEC, JobID: "2"},
		{Action: EXPIRE, JobID: "2"},
		{Action: EXEC, JobID: "1"},
		{Action: EXEC, JobID: "1"},
		{Action: COMPLETE, JobID: "1"},
		{Action: EXEC, JobID: "2"},
		{Action: EXEC, JobID: "2"},
		{Action: COMPLETE, JobID: "2"},
	}

	mismatched := -1
	for i := range expectedSystemBehaviour {
		if expectedSystemBehaviour[i].Action != out[i].Action || expectedSystemBehaviour[i].JobID != out[i].JobID {
			mismatched = i
			break
		}
	}

	if mismatched != -1 {
		t.Errorf("Mismatch detected at index %v, expected:%v, got:%v", mismatched, expectedSystemBehaviour[mismatched], out[mismatched])
	}
}

func TestMLFQOnlyIO(t *testing.T) {
	inputs := []*Job{
		NewJob(1, 0, []JobInput{IOInstruction, IOInstruction}),
	}

	MLFQConfig := MLFQConfig{
		QueueConfigs: []QueueConfig{
			{Priority: 4, TimeAllotment: 2},
			{Priority: 3, TimeAllotment: 3},
			{Priority: 2, TimeAllotment: 4},
			{Priority: 1, TimeAllotment: 5},
			{Priority: 0, TimeAllotment: 6},
		},
		ResetInterval: 5,
		QueueSize:     100,
	}

	out := RunSystem(&MLFQConfig, inputs, time.Duration(2*time.Second), true)

	expectedSystemBehaviour := []AuditLog{
		{Action: SWAP, JobID: "1"},
		{Action: SWAP, JobID: "1"},
	}

	if len(expectedSystemBehaviour) != len(out) {
		t.Errorf("Mismatch in number of CPU actions, expected:%v, got:%v", len(expectedSystemBehaviour), len(out))
		return
	}

	mismatched := -1
	for i := range expectedSystemBehaviour {
		if expectedSystemBehaviour[i].Action != out[i].Action || expectedSystemBehaviour[i].JobID != out[i].JobID {
			mismatched = i
			break
		}
	}

	if mismatched != -1 {
		t.Errorf("Mismatch detected at index %v, expected:%v, got:%v", mismatched, expectedSystemBehaviour[mismatched], out[mismatched])
	}
}

func TestMLFQMultipleIOJobs(t *testing.T) {
	inputs := []*Job{
		NewJob(1, 0, []JobInput{IOInstruction}),
		NewJob(2, 0, []JobInput{IOInstruction}),
		NewJob(3, 0, []JobInput{CPUInstruction, CPUInstruction}),
	}

	MLFQConfig := MLFQConfig{
		QueueConfigs: []QueueConfig{
			{Priority: 4, TimeAllotment: 2},
			{Priority: 3, TimeAllotment: 3},
			{Priority: 2, TimeAllotment: 4},
			{Priority: 1, TimeAllotment: 5},
			{Priority: 0, TimeAllotment: 6},
		},
		ResetInterval: 5,
		QueueSize:     100,
	}

	out := RunSystem(&MLFQConfig, inputs, time.Duration(2*time.Second), true)

	expectedSystemBehaviour := []AuditLog{
		{Action: SWAP, JobID: "1"},
		{Action: SWAP, JobID: "2"},
		{Action: SWAP, JobID: "3"},
	}

	if len(expectedSystemBehaviour) != len(out) {
		t.Errorf("Mismatch in number of CPU actions, expected:%v, got:%v", len(expectedSystemBehaviour), len(out))
		return
	}

	mismatched := -1
	for i := range expectedSystemBehaviour {
		if expectedSystemBehaviour[i].Action != out[i].Action || expectedSystemBehaviour[i].JobID != out[i].JobID {
			mismatched = i
			break
		}
	}

	if mismatched != -1 {
		t.Errorf("Mismatch detected at index %v, expected:%v, got:%v", mismatched, expectedSystemBehaviour[mismatched], out[mismatched])
	}
}
