package youtrack

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type CustomField struct {
	Name  string
	Value string
}

type rawCustomField struct {
	Type  string          `json:"$type"`
	Name  string          `json:"name"`
	Value json.RawMessage `json:"value"`
}

type enumValue struct {
	Name string `json:"name"`
}

type stateValue struct {
	Name       string `json:"name"`
	IsResolved bool   `json:"isResolved"`
}

type userValue struct {
	Login    string `json:"login"`
	FullName string `json:"fullName"`
}

type periodValue struct {
	Presentation string `json:"presentation"`
}

type textValue struct {
	Text string `json:"text"`
}

func isNull(raw json.RawMessage) bool {
	return len(raw) == 0 || string(raw) == "null"
}

func decodeSingleEnum(raw json.RawMessage) string {
	if isNull(raw) {
		return ""
	}
	var v enumValue
	if err := json.Unmarshal(raw, &v); err != nil {
		return string(raw)
	}
	return v.Name
}

func decodeMultiEnum(raw json.RawMessage) string {
	if isNull(raw) {
		return ""
	}
	var vs []enumValue
	if err := json.Unmarshal(raw, &vs); err != nil {
		return string(raw)
	}
	names := make([]string, len(vs))
	for i, v := range vs {
		names[i] = v.Name
	}
	return strings.Join(names, ", ")
}

func decodeState(raw json.RawMessage) string {
	if isNull(raw) {
		return ""
	}
	var v stateValue
	if err := json.Unmarshal(raw, &v); err != nil {
		return string(raw)
	}
	return v.Name
}

func formatUser(u userValue) string {
	return fmt.Sprintf("%s (%s)", u.FullName, u.Login)
}

func decodeSingleUser(raw json.RawMessage) string {
	if isNull(raw) {
		return ""
	}
	var v userValue
	if err := json.Unmarshal(raw, &v); err != nil {
		return string(raw)
	}
	return formatUser(v)
}

func decodeMultiUser(raw json.RawMessage) string {
	if isNull(raw) {
		return ""
	}
	var vs []userValue
	if err := json.Unmarshal(raw, &vs); err != nil {
		return string(raw)
	}
	parts := make([]string, len(vs))
	for i, v := range vs {
		parts[i] = formatUser(v)
	}
	return strings.Join(parts, ", ")
}

func decodeMultiVersion(raw json.RawMessage) string {
	if isNull(raw) {
		return ""
	}
	var vs []enumValue
	if err := json.Unmarshal(raw, &vs); err != nil {
		return string(raw)
	}
	names := make([]string, len(vs))
	for i, v := range vs {
		names[i] = v.Name
	}
	return strings.Join(names, ", ")
}

func decodePeriod(raw json.RawMessage) string {
	if isNull(raw) {
		return ""
	}
	var v periodValue
	if err := json.Unmarshal(raw, &v); err != nil {
		return string(raw)
	}
	return v.Presentation
}

func decodeText(raw json.RawMessage) string {
	if isNull(raw) {
		return ""
	}
	var v textValue
	if err := json.Unmarshal(raw, &v); err != nil {
		return string(raw)
	}
	return v.Text
}

func decodeSimple(raw json.RawMessage) string {
	if isNull(raw) {
		return ""
	}
	s := string(raw)
	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		var str string
		if err := json.Unmarshal(raw, &str); err == nil {
			return str
		}
	}
	return s
}

func decodeDate(raw json.RawMessage) string {
	if isNull(raw) {
		return ""
	}
	var ms int64
	if err := json.Unmarshal(raw, &ms); err != nil {
		return string(raw)
	}
	return time.UnixMilli(ms).UTC().Format("2006-01-02")
}

func DecodeCustomFields(raw []json.RawMessage) []CustomField {
	fields := make([]CustomField, 0, len(raw))
	for _, elem := range raw {
		var r rawCustomField
		if err := json.Unmarshal(elem, &r); err != nil {
			continue
		}
		var value string
		switch r.Type {
		case "SingleEnumIssueCustomField":
			value = decodeSingleEnum(r.Value)
		case "MultiEnumIssueCustomField":
			value = decodeMultiEnum(r.Value)
		case "StateIssueCustomField":
			value = decodeState(r.Value)
		case "SingleUserIssueCustomField":
			value = decodeSingleUser(r.Value)
		case "MultiUserIssueCustomField":
			value = decodeMultiUser(r.Value)
		case "SingleVersionIssueCustomField", "SingleBuildIssueCustomField", "SingleOwnedIssueCustomField", "SingleGroupIssueCustomField":
			value = decodeSingleEnum(r.Value)
		case "MultiVersionIssueCustomField":
			value = decodeMultiVersion(r.Value)
		case "MultiBuildIssueCustomField", "MultiOwnedIssueCustomField", "MultiGroupIssueCustomField":
			value = decodeMultiEnum(r.Value)
		case "PeriodIssueCustomField":
			value = decodePeriod(r.Value)
		case "TextIssueCustomField":
			value = decodeText(r.Value)
		case "SimpleIssueCustomField":
			value = decodeSimple(r.Value)
		case "DateIssueCustomField":
			value = decodeDate(r.Value)
		default:
			value = string(r.Value)
		}
		fields = append(fields, CustomField{Name: r.Name, Value: value})
	}
	return fields
}
