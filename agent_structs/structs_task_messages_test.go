package agentstructs

import "testing"

func TestPTTaskMessageTaskDataRevertKeywordsScalarReplacement(t *testing.T) {
	task := PTTaskMessageTaskData{
		KeywordResolution: []PTTaskKeywordResolution{
			{
				Raw:            "@cred:12.credential",
				ValueType:      "string",
				ExpandedValue:  "abc123",
				ParameterNames: []string{"args"},
			},
		},
	}
	got := task.RevertKeywords("asktgs /ticket:abc123 /nowrap")
	want := "asktgs /ticket:@cred:12.credential /nowrap"
	if got != want {
		t.Fatalf("RevertKeywords() = %q, want %q", got, want)
	}
}

func TestPTTaskMessageTaskDataRevertKeywordsLongestFirst(t *testing.T) {
	task := PTTaskMessageTaskData{
		KeywordResolution: []PTTaskKeywordResolution{
			{Raw: "@token:1.value", ValueType: "string", ExpandedValue: "abc"},
			{Raw: "@token:2.value", ValueType: "string", ExpandedValue: "abc123"},
		},
	}
	got := task.RevertKeywords("value=abc123 next=abc")
	want := "value=@token:2.value next=@token:1.value"
	if got != want {
		t.Fatalf("RevertKeywords() = %q, want %q", got, want)
	}
}

func TestPTTaskMessageTaskDataRevertKeywordsParameterNameFiltering(t *testing.T) {
	task := PTTaskMessageTaskData{
		KeywordResolution: []PTTaskKeywordResolution{
			{
				Raw:            "@cred:12.credential",
				ValueType:      "string",
				ExpandedValue:  "shared",
				ParameterNames: []string{"args"},
			},
			{
				Raw:            "@file:7.filename",
				ValueType:      "string",
				ExpandedValue:  "shared",
				ParameterNames: []string{"filename"},
			},
		},
	}
	got := task.RevertKeywords("shared", "filename")
	want := "@file:7.filename"
	if got != want {
		t.Fatalf("RevertKeywords() = %q, want %q", got, want)
	}
}

func TestPTTaskMessageTaskDataRevertKeywordsStructured(t *testing.T) {
	task := PTTaskMessageTaskData{
		KeywordResolution: []PTTaskKeywordResolution{
			{
				Raw:            "@cred:12",
				ValueType:      "structured",
				ParameterNames: []string{"cred"},
			},
		},
	}
	got := task.RevertKeywords(map[string]interface{}{"id": 12}, "cred")
	if got != "@cred:12" {
		t.Fatalf("RevertKeywords() = %q, want %q", got, "@cred:12")
	}
}

func TestPTTaskMessageTaskDataRevertKeywordsPassthrough(t *testing.T) {
	task := PTTaskMessageTaskData{}
	got := task.RevertKeywords("unchanged")
	if got != "unchanged" {
		t.Fatalf("RevertKeywords() = %q, want unchanged", got)
	}
}
