package rabbitmq

import (
	"context"
	"strings"
	"testing"

	agentstructs "github.com/MythicMeta/MythicContainer/agent_structs"
)

func newAgentRPCRequest(payloadType string) agentstructs.PTTaskAgentRPCMessage {
	return agentstructs.PTTaskAgentRPCMessage{
		TaskData: &agentstructs.PTTaskMessageAllData{
			PayloadType:        payloadType,
			CommandPayloadType: "command-augment",
			Task: agentstructs.PTTaskMessageTaskData{
				AgentTaskID: "agent-task-uuid",
			},
			Callback: agentstructs.PTTaskMessageCallbackData{
				ID: 42,
			},
		},
		Name:      "lookup",
		Arguments: map[string]any{"value": "test"},
	}
}

func TestPTTaskAgentRPCRequestRouteRegistered(t *testing.T) {
	for _, directMethod := range agentstructs.AllPayloadData.Get("").GetDirectMethods() {
		if directMethod.RabbitmqRoutingKey == PT_TASK_AGENT_RPC {
			return
		}
	}
	t.Fatal("agent RPC direct request route was not registered")
}

func TestInvokePTTaskAgentRPCRoutesToCallbackPayloadAndNormalizesIDs(t *testing.T) {
	payloadName := "agent-rpc-invoke-test"
	payloadData := agentstructs.AllPayloadData.Get(payloadName)
	payloadData.AddPayloadDefinition(agentstructs.PayloadType{Name: payloadName})
	payloadData.AddAgentRPCFunction(func(
		_ context.Context,
		task *agentstructs.PTTaskMessageAllData,
		name string,
		arguments any,
	) agentstructs.PTTaskAgentRPCMessageResponse {
		if task.PayloadType != payloadName {
			t.Fatalf("expected callback payload type, got %q", task.PayloadType)
		}
		if task.CommandPayloadType != "command-augment" {
			t.Fatalf("expected command payload metadata to remain available, got %q", task.CommandPayloadType)
		}
		if name != "lookup" {
			t.Fatalf("unexpected name: %q", name)
		}
		return agentstructs.PTTaskAgentRPCMessageResponse{
			CallbackID:  999,
			AgentTaskID: "developer-id",
			Status:      "custom-status",
			Output:      arguments,
		}
	})

	request := newAgentRPCRequest(payloadName)
	response := invokePTTaskAgentRPC(context.Background(), request)
	response = normalizePTTaskAgentRPCResponse(request, response)
	if response.CallbackID != 42 || response.AgentTaskID != "agent-task-uuid" {
		t.Fatalf("runtime must normalize response IDs from task context: %#v", response)
	}
	if response.Status != "custom-status" {
		t.Fatalf("status should pass through unchanged: %#v", response)
	}
	if agentstructs.AllPayloadData.Get(payloadName).GetRoutingKey(PT_TASK_AGENT_RPC) !=
		payloadName+"_pt_task_agent_rpc" {
		t.Fatalf("unexpected request routing key")
	}
}

func TestInvokePTTaskAgentRPCMissingFunctionAndPanicBecomeErrors(t *testing.T) {
	missingPayloadName := "agent-rpc-missing-test"
	agentstructs.AllPayloadData.Get(missingPayloadName).AddPayloadDefinition(
		agentstructs.PayloadType{Name: missingPayloadName},
	)
	missingRequest := newAgentRPCRequest(missingPayloadName)
	missingResponse := normalizePTTaskAgentRPCResponse(
		missingRequest,
		invokePTTaskAgentRPC(context.Background(), missingRequest),
	)
	if missingResponse.Status != "error" ||
		!strings.Contains(missingResponse.Output.(string), "not registered") {
		t.Fatalf("unexpected missing function response: %#v", missingResponse)
	}

	panicPayloadName := "agent-rpc-panic-test"
	panicPayloadData := agentstructs.AllPayloadData.Get(panicPayloadName)
	panicPayloadData.AddPayloadDefinition(agentstructs.PayloadType{Name: panicPayloadName})
	panicPayloadData.AddAgentRPCFunction(func(
		context.Context,
		*agentstructs.PTTaskMessageAllData,
		string,
		any,
	) agentstructs.PTTaskAgentRPCMessageResponse {
		panic("boom")
	})
	panicRequest := newAgentRPCRequest(panicPayloadName)
	panicResponse := normalizePTTaskAgentRPCResponse(
		panicRequest,
		invokePTTaskAgentRPC(context.Background(), panicRequest),
	)
	if panicResponse.Status != "error" ||
		!strings.Contains(panicResponse.Output.(string), "panicked: boom") {
		t.Fatalf("unexpected panic response: %#v", panicResponse)
	}
	if panicResponse.CallbackID != 42 || panicResponse.AgentTaskID != "agent-task-uuid" {
		t.Fatalf("panic response must remain correlated: %#v", panicResponse)
	}
}

func TestEnsurePTTaskAgentRPCResponseSerializable(t *testing.T) {
	response := ensurePTTaskAgentRPCResponseSerializable(agentstructs.PTTaskAgentRPCMessageResponse{
		CallbackID:  42,
		AgentTaskID: "agent-task-uuid",
		Status:      "success",
		Output:      make(chan int),
	})
	if response.Status != "error" {
		t.Fatalf("expected serialization failure status, got %#v", response)
	}
	output, ok := response.Output.(string)
	if !ok || !strings.Contains(output, "Failed to serialize agent RPC output") {
		t.Fatalf("unexpected serialization failure output: %#v", response.Output)
	}
}
