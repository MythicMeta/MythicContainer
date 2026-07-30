package agentstructs

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"
)

func TestPTTaskAgentRPCWireContract(t *testing.T) {
	request := PTTaskAgentRPCMessage{
		TaskData: &PTTaskMessageAllData{
			PayloadType:        "apollo",
			CommandPayloadType: "command-augment",
			Task: PTTaskMessageTaskData{
				AgentTaskID: "agent-task-uuid",
			},
			Callback: PTTaskMessageCallbackData{
				ID: 42,
			},
		},
		Name: "lookup",
		Arguments: map[string]any{
			"nested": []any{"one", float64(2), true, nil},
		},
	}
	encoded, err := json.Marshal(request)
	if err != nil {
		t.Fatalf("failed to marshal agent RPC request: %v", err)
	}
	decoded := PTTaskAgentRPCMessage{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		t.Fatalf("failed to unmarshal agent RPC request: %v", err)
	}
	if decoded.TaskData == nil || decoded.TaskData.PayloadType != "apollo" {
		t.Fatalf("unexpected task data: %#v", decoded.TaskData)
	}
	if decoded.TaskData.CommandPayloadType != "command-augment" {
		t.Fatalf("expected command payload type to survive, got %q", decoded.TaskData.CommandPayloadType)
	}
	if decoded.Name != "lookup" {
		t.Fatalf("unexpected function name: %q", decoded.Name)
	}
	if !reflect.DeepEqual(decoded.Arguments, request.Arguments) {
		t.Fatalf("unexpected arguments: %#v", decoded.Arguments)
	}

	response := PTTaskAgentRPCMessageResponse{
		CallbackID:  42,
		AgentTaskID: "agent-task-uuid",
		Status:      "custom",
		Output:      map[string]any{"value": nil},
	}
	encoded, err = json.Marshal(response)
	if err != nil {
		t.Fatalf("failed to marshal agent RPC response: %v", err)
	}
	wireResponse := map[string]any{}
	if err = json.Unmarshal(encoded, &wireResponse); err != nil {
		t.Fatalf("failed to unmarshal agent RPC response: %v", err)
	}
	if wireResponse["callback_id"] != float64(42) ||
		wireResponse["agent_task_id"] != "agent-task-uuid" ||
		wireResponse["status"] != "custom" {
		t.Fatalf("unexpected response wire shape: %#v", wireResponse)
	}
	if _, ok := wireResponse["output"]; !ok {
		t.Fatalf("response output must always be present: %#v", wireResponse)
	}
}

func TestAgentRPCFunctionRegistration(t *testing.T) {
	payloadName := "agent-rpc-registration-test"
	payloadData := AllPayloadData.Get(payloadName)
	payloadData.AddPayloadDefinition(PayloadType{Name: payloadName})
	expectedOutput := map[string]any{"ok": true}
	payloadData.AddAgentRPCFunction(func(
		_ context.Context,
		task *PTTaskMessageAllData,
		name string,
		arguments any,
	) PTTaskAgentRPCMessageResponse {
		if task.PayloadType != payloadName || name != "lookup" || arguments != "value" {
			t.Fatalf("unexpected hook arguments: %#v %q %#v", task, name, arguments)
		}
		return PTTaskAgentRPCMessageResponse{
			Status: "success",
			Output: expectedOutput,
		}
	})

	agentRPCFunction := payloadData.GetAgentRPCFunction()
	if agentRPCFunction == nil {
		t.Fatal("expected registered agent RPC function")
	}
	response := agentRPCFunction(context.Background(), &PTTaskMessageAllData{
		PayloadType: payloadName,
	}, "lookup", "value")
	if response.Status != "success" || !reflect.DeepEqual(response.Output, expectedOutput) {
		t.Fatalf("unexpected hook response: %#v", response)
	}
}
