package mythicrpc

import "testing"

func TestPTTaskMessageTaskDataRevertKeywordsForTaskSearch(t *testing.T) {
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
	got := task.RevertKeywords("asktgs /ticket:abc123 /nowrap", "args")
	want := "asktgs /ticket:@cred:12.credential /nowrap"
	if got != want {
		t.Fatalf("RevertKeywords() = %q, want %q", got, want)
	}
}
