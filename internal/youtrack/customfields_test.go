package youtrack

import (
	"encoding/json"
	"testing"
)

func makeRaw(t *testing.T, v any) json.RawMessage {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("makeRaw: %v", err)
	}
	return json.RawMessage(b)
}

func field(typ, name string, value any) map[string]any {
	return map[string]any{
		"$type": typ,
		"name":  name,
		"value": value,
	}
}

func fieldRaw(typ, name string, value json.RawMessage) map[string]any {
	return map[string]any{
		"$type": typ,
		"name":  name,
		"value": value,
	}
}

func encodeFields(t *testing.T, fields []map[string]any) []json.RawMessage {
	t.Helper()
	raw := make([]json.RawMessage, len(fields))
	for i, f := range fields {
		b, err := json.Marshal(f)
		if err != nil {
			t.Fatalf("encodeFields[%d]: %v", i, err)
		}
		raw[i] = json.RawMessage(b)
	}
	return raw
}

func TestDecodeCustomFields_SingleEnum(t *testing.T) {
	raw := encodeFields(t, []map[string]any{
		field("SingleEnumIssueCustomField", "Priority", map[string]any{"name": "Major"}),
	})
	got := DecodeCustomFields(raw)
	if len(got) != 1 || got[0].Name != "Priority" || got[0].Value != "Major" {
		t.Errorf("unexpected: %+v", got)
	}
}

func TestDecodeCustomFields_SingleEnum_Null(t *testing.T) {
	raw := encodeFields(t, []map[string]any{
		field("SingleEnumIssueCustomField", "Priority", nil),
	})
	got := DecodeCustomFields(raw)
	if len(got) != 1 || got[0].Value != "" {
		t.Errorf("unexpected: %+v", got)
	}
}

func TestDecodeCustomFields_MultiEnum(t *testing.T) {
	raw := encodeFields(t, []map[string]any{
		field("MultiEnumIssueCustomField", "Tags", []map[string]any{
			{"name": "backend"},
			{"name": "urgent"},
		}),
	})
	got := DecodeCustomFields(raw)
	if len(got) != 1 || got[0].Value != "backend, urgent" {
		t.Errorf("unexpected: %+v", got)
	}
}

func TestDecodeCustomFields_MultiEnum_Null(t *testing.T) {
	raw := encodeFields(t, []map[string]any{
		field("MultiEnumIssueCustomField", "Tags", nil),
	})
	got := DecodeCustomFields(raw)
	if len(got) != 1 || got[0].Value != "" {
		t.Errorf("unexpected: %+v", got)
	}
}

func TestDecodeCustomFields_MultiEnum_Empty(t *testing.T) {
	raw := encodeFields(t, []map[string]any{
		field("MultiEnumIssueCustomField", "Tags", []map[string]any{}),
	})
	got := DecodeCustomFields(raw)
	if len(got) != 1 || got[0].Value != "" {
		t.Errorf("unexpected: %+v", got)
	}
}

func TestDecodeCustomFields_State(t *testing.T) {
	raw := encodeFields(t, []map[string]any{
		field("StateIssueCustomField", "State", map[string]any{"name": "In Progress", "isResolved": false}),
	})
	got := DecodeCustomFields(raw)
	if len(got) != 1 || got[0].Value != "In Progress" {
		t.Errorf("unexpected: %+v", got)
	}
}

func TestDecodeCustomFields_State_Resolved(t *testing.T) {
	raw := encodeFields(t, []map[string]any{
		field("StateIssueCustomField", "State", map[string]any{"name": "Fixed", "isResolved": true}),
	})
	got := DecodeCustomFields(raw)
	if len(got) != 1 || got[0].Value != "Fixed" {
		t.Errorf("unexpected: %+v", got)
	}
}

func TestDecodeCustomFields_State_Null(t *testing.T) {
	raw := encodeFields(t, []map[string]any{
		field("StateIssueCustomField", "State", nil),
	})
	got := DecodeCustomFields(raw)
	if len(got) != 1 || got[0].Value != "" {
		t.Errorf("unexpected: %+v", got)
	}
}

func TestDecodeCustomFields_SingleUser(t *testing.T) {
	raw := encodeFields(t, []map[string]any{
		field("SingleUserIssueCustomField", "Assignee", map[string]any{
			"login":    "jdoe",
			"fullName": "John Doe",
		}),
	})
	got := DecodeCustomFields(raw)
	if len(got) != 1 || got[0].Value != "John Doe (jdoe)" {
		t.Errorf("unexpected: %+v", got)
	}
}

func TestDecodeCustomFields_SingleUser_Null(t *testing.T) {
	raw := encodeFields(t, []map[string]any{
		field("SingleUserIssueCustomField", "Assignee", nil),
	})
	got := DecodeCustomFields(raw)
	if len(got) != 1 || got[0].Value != "" {
		t.Errorf("unexpected: %+v", got)
	}
}

func TestDecodeCustomFields_MultiUser(t *testing.T) {
	raw := encodeFields(t, []map[string]any{
		field("MultiUserIssueCustomField", "Reviewers", []map[string]any{
			{"login": "alice", "fullName": "Alice Smith"},
			{"login": "bob", "fullName": "Bob Jones"},
		}),
	})
	got := DecodeCustomFields(raw)
	if len(got) != 1 || got[0].Value != "Alice Smith (alice), Bob Jones (bob)" {
		t.Errorf("unexpected: %+v", got)
	}
}

func TestDecodeCustomFields_MultiUser_Null(t *testing.T) {
	raw := encodeFields(t, []map[string]any{
		field("MultiUserIssueCustomField", "Reviewers", nil),
	})
	got := DecodeCustomFields(raw)
	if len(got) != 1 || got[0].Value != "" {
		t.Errorf("unexpected: %+v", got)
	}
}

func TestDecodeCustomFields_SingleVersion(t *testing.T) {
	raw := encodeFields(t, []map[string]any{
		field("SingleVersionIssueCustomField", "Fix version", map[string]any{"name": "v2.0"}),
	})
	got := DecodeCustomFields(raw)
	if len(got) != 1 || got[0].Value != "v2.0" {
		t.Errorf("unexpected: %+v", got)
	}
}

func TestDecodeCustomFields_SingleVersion_Null(t *testing.T) {
	raw := encodeFields(t, []map[string]any{
		field("SingleVersionIssueCustomField", "Fix version", nil),
	})
	got := DecodeCustomFields(raw)
	if len(got) != 1 || got[0].Value != "" {
		t.Errorf("unexpected: %+v", got)
	}
}

func TestDecodeCustomFields_MultiVersion(t *testing.T) {
	raw := encodeFields(t, []map[string]any{
		field("MultiVersionIssueCustomField", "Affected versions", []map[string]any{
			{"name": "v1.0"},
			{"name": "v1.1"},
		}),
	})
	got := DecodeCustomFields(raw)
	if len(got) != 1 || got[0].Value != "v1.0, v1.1" {
		t.Errorf("unexpected: %+v", got)
	}
}

func TestDecodeCustomFields_MultiVersion_Null(t *testing.T) {
	raw := encodeFields(t, []map[string]any{
		field("MultiVersionIssueCustomField", "Affected versions", nil),
	})
	got := DecodeCustomFields(raw)
	if len(got) != 1 || got[0].Value != "" {
		t.Errorf("unexpected: %+v", got)
	}
}

func TestDecodeCustomFields_Period(t *testing.T) {
	raw := encodeFields(t, []map[string]any{
		field("PeriodIssueCustomField", "Spent time", map[string]any{"presentation": "2w 3d"}),
	})
	got := DecodeCustomFields(raw)
	if len(got) != 1 || got[0].Value != "2w 3d" {
		t.Errorf("unexpected: %+v", got)
	}
}

func TestDecodeCustomFields_Period_Null(t *testing.T) {
	raw := encodeFields(t, []map[string]any{
		field("PeriodIssueCustomField", "Spent time", nil),
	})
	got := DecodeCustomFields(raw)
	if len(got) != 1 || got[0].Value != "" {
		t.Errorf("unexpected: %+v", got)
	}
}

func TestDecodeCustomFields_Text(t *testing.T) {
	raw := encodeFields(t, []map[string]any{
		field("TextIssueCustomField", "Notes", map[string]any{"text": "some markdown"}),
	})
	got := DecodeCustomFields(raw)
	if len(got) != 1 || got[0].Value != "some markdown" {
		t.Errorf("unexpected: %+v", got)
	}
}

func TestDecodeCustomFields_Text_Null(t *testing.T) {
	raw := encodeFields(t, []map[string]any{
		field("TextIssueCustomField", "Notes", nil),
	})
	got := DecodeCustomFields(raw)
	if len(got) != 1 || got[0].Value != "" {
		t.Errorf("unexpected: %+v", got)
	}
}

func TestDecodeCustomFields_Simple_String(t *testing.T) {
	raw := encodeFields(t, []map[string]any{
		field("SimpleIssueCustomField", "Ticket URL", "https://example.com"),
	})
	got := DecodeCustomFields(raw)
	if len(got) != 1 || got[0].Value != "https://example.com" {
		t.Errorf("unexpected: %+v", got)
	}
}

func TestDecodeCustomFields_Simple_Number(t *testing.T) {
	raw := encodeFields(t, []map[string]any{
		field("SimpleIssueCustomField", "Story Points", 8),
	})
	got := DecodeCustomFields(raw)
	if len(got) != 1 || got[0].Value != "8" {
		t.Errorf("unexpected: %+v", got)
	}
}

func TestDecodeCustomFields_Simple_Bool(t *testing.T) {
	raw := encodeFields(t, []map[string]any{
		field("SimpleIssueCustomField", "Is blocked", true),
	})
	got := DecodeCustomFields(raw)
	if len(got) != 1 || got[0].Value != "true" {
		t.Errorf("unexpected: %+v", got)
	}
}

func TestDecodeCustomFields_Simple_Null(t *testing.T) {
	raw := encodeFields(t, []map[string]any{
		field("SimpleIssueCustomField", "Story Points", nil),
	})
	got := DecodeCustomFields(raw)
	if len(got) != 1 || got[0].Value != "" {
		t.Errorf("unexpected: %+v", got)
	}
}

func TestDecodeCustomFields_Date(t *testing.T) {
	// 2024-03-15 00:00:00 UTC = 1710460800000 ms
	raw := encodeFields(t, []map[string]any{
		field("DateIssueCustomField", "Due date", int64(1710460800000)),
	})
	got := DecodeCustomFields(raw)
	if len(got) != 1 || got[0].Value != "2024-03-15" {
		t.Errorf("unexpected: %+v", got)
	}
}

func TestDecodeCustomFields_Date_Null(t *testing.T) {
	raw := encodeFields(t, []map[string]any{
		field("DateIssueCustomField", "Due date", nil),
	})
	got := DecodeCustomFields(raw)
	if len(got) != 1 || got[0].Value != "" {
		t.Errorf("unexpected: %+v", got)
	}
}

func TestDecodeCustomFields_Unknown_Type(t *testing.T) {
	raw := encodeFields(t, []map[string]any{
		field("SomeNewFutureField", "Mystery", map[string]any{"foo": "bar"}),
	})
	got := DecodeCustomFields(raw)
	if len(got) != 1 || got[0].Name != "Mystery" {
		t.Errorf("unexpected: %+v", got)
	}
	if got[0].Value == "" {
		t.Errorf("expected non-empty raw JSON for unknown type, got empty string")
	}
}

func TestDecodeCustomFields_MultipleFields(t *testing.T) {
	raw := encodeFields(t, []map[string]any{
		field("SingleEnumIssueCustomField", "Priority", map[string]any{"name": "Critical"}),
		field("StateIssueCustomField", "State", map[string]any{"name": "Open", "isResolved": false}),
		field("SingleUserIssueCustomField", "Assignee", nil),
	})
	got := DecodeCustomFields(raw)
	if len(got) != 3 {
		t.Fatalf("expected 3 fields, got %d", len(got))
	}
	if got[0].Name != "Priority" || got[0].Value != "Critical" {
		t.Errorf("field 0: %+v", got[0])
	}
	if got[1].Name != "State" || got[1].Value != "Open" {
		t.Errorf("field 1: %+v", got[1])
	}
	if got[2].Name != "Assignee" || got[2].Value != "" {
		t.Errorf("field 2: %+v", got[2])
	}
}

func TestDecodeCustomFields_InvalidJSON_Skipped(t *testing.T) {
	raw := []json.RawMessage{
		json.RawMessage(`not-valid-json`),
		makeRaw(t, field("SingleEnumIssueCustomField", "Priority", map[string]any{"name": "Low"})),
	}
	got := DecodeCustomFields(raw)
	if len(got) != 1 || got[0].Value != "Low" {
		t.Errorf("unexpected: %+v", got)
	}
}

func TestDecodeCustomFields_Empty(t *testing.T) {
	got := DecodeCustomFields([]json.RawMessage{})
	if len(got) != 0 {
		t.Errorf("expected empty slice, got %+v", got)
	}
}

func TestDecodeCustomFields_Nil(t *testing.T) {
	got := DecodeCustomFields(nil)
	if len(got) != 0 {
		t.Errorf("expected empty slice, got %+v", got)
	}
}

func TestDecodeCustomFields_UnknownType_NullValue(t *testing.T) {
	raw := encodeFields(t, []map[string]any{
		field("SomeNewFutureField", "Mystery", nil),
	})
	got := DecodeCustomFields(raw)
	if len(got) != 1 || got[0].Value != "null" {
		t.Errorf("unexpected: %+v", got)
	}
}

func TestDecodeCustomFields_MultiEnum_Single(t *testing.T) {
	raw := encodeFields(t, []map[string]any{
		field("MultiEnumIssueCustomField", "Tags", []map[string]any{
			{"name": "solo"},
		}),
	})
	got := DecodeCustomFields(raw)
	if len(got) != 1 || got[0].Value != "solo" {
		t.Errorf("unexpected: %+v", got)
	}
}

func TestDecodeCustomFields_MultiUser_Single(t *testing.T) {
	raw := encodeFields(t, []map[string]any{
		field("MultiUserIssueCustomField", "Reviewers", []map[string]any{
			{"login": "carol", "fullName": "Carol White"},
		}),
	})
	got := DecodeCustomFields(raw)
	if len(got) != 1 || got[0].Value != "Carol White (carol)" {
		t.Errorf("unexpected: %+v", got)
	}
}

func TestDecodeCustomFields_Simple_FloatNumber(t *testing.T) {
	raw := encodeFields(t, []map[string]any{
		field("SimpleIssueCustomField", "Weight", 3.14),
	})
	got := DecodeCustomFields(raw)
	if len(got) != 1 || got[0].Value == "" {
		t.Errorf("unexpected: %+v", got)
	}
}
