package validation

import (
	"testing"
)

type mainSubject struct {
	StringField string
	IntField    int

	NestedStructField nestedStruct

	NestedStructFieldPtr *nestedStruct

	SliceOfStringsField []string

	SliceOfStructsField []nestedStruct
}

// we don't care about testing rules here, they are tested in rule_test.go
// we just use simple rules to prove validation works
func (mainSubject) ValidationRules() RuleSet {
	return RuleSet{
		"StringField":          []Rule{NotZeroValue{}},
		"IntField":             []Rule{NotZeroValue{}},
		"NestedStructField":    []Rule{NestedValid{}},
		"NestedStructFieldPtr": []Rule{NestedValid{}},
		"SliceOfStringsField":  []Rule{NotZeroValue{}},
		"SliceOfStructsField":  []Rule{NotZeroValue{}},
	}
}

type nestedStruct struct {
	StringField string
	IntField    int

	DeeplyNestedStructField deeplyNestedStruct

	SliceOfStringsField []string

	SliceOfDeeplyNestedStructsField []deeplyNestedStruct
}

func (nestedStruct) ValidationRules() RuleSet {
	return RuleSet{
		"StringField":                     []Rule{NotZeroValue{}},
		"IntField":                        []Rule{NotZeroValue{}},
		"DeeplyNestedStructField":         []Rule{NestedValid{}},
		"SliceOfStringsField":             []Rule{NotZeroValue{}},
		"SliceOfDeeplyNestedStructsField": []Rule{NotZeroValue{}},
	}
}

type deeplyNestedStruct struct {
	StringField string
	IntField    int
}

func (deeplyNestedStruct) ValidationRules() RuleSet {
	return RuleSet{
		"StringField": []Rule{NotZeroValue{}},
		"IntField":    []Rule{NotZeroValue{}},
	}
}

// this build a struct with all fields containing zero values and validates it
func TestNothingIsValid(t *testing.T) {
	failures, err := Validate(mainSubject{
		StringField:          "",
		IntField:             0,
		NestedStructField:    nestedStruct{},
		NestedStructFieldPtr: &nestedStruct{},
		SliceOfStringsField:  nil,
		SliceOfStructsField: []nestedStruct{{
			StringField:                     "",
			IntField:                        0,
			DeeplyNestedStructField:         deeplyNestedStruct{},
			SliceOfStringsField:             nil,
			SliceOfDeeplyNestedStructsField: nil,
		}},
	})
	if err != nil {
		t.Error("did not expect an error")
	}

	noExpectedFailures := 21
	if len(failures) != noExpectedFailures {
		t.Errorf("want: %d failures; got: %d failures, %v", noExpectedFailures, len(failures), failures)
	}

	tester := func(fieldName string, noFailsWant int) {
		fs, ok := failures[fieldName]
		if !ok && 0 != noFailsWant {
			t.Errorf("expected failures for '%s'", fieldName)
		}
		if len(fs) != noFailsWant {
			t.Errorf("number of failures for '%s' want: %d; got: %d", fieldName, noFailsWant, len(fs))
		}
	}

	tester("StringField", 1)
	tester("IntField", 1)
	tester("NestedStructField", 0)
	tester("NestedStructField.StringField", 1)
	tester("NestedStructField.IntField", 1)
	tester("NestedStructField.DeeplyNestedStructField", 0)
	tester("NestedStructField.DeeplyNestedStructField.StringField", 1)
	tester("NestedStructField.DeeplyNestedStructField.IntField", 1)
	tester("NestedStructField.SliceOfStringsField", 1)
	tester("NestedStructField.SliceOfDeeplyNestedStructsField", 1)
	tester("NestedStructFieldPtr", 0)
	tester("NestedStructFieldPtr.StringField", 1)
	tester("NestedStructFieldPtr.IntField", 1)
	tester("NestedStructFieldPtr.DeeplyNestedStructField", 0)
	tester("NestedStructFieldPtr.DeeplyNestedStructField.StringField", 1)
	tester("NestedStructFieldPtr.DeeplyNestedStructField.IntField", 1)
	tester("NestedStructFieldPtr.SliceOfStringsField", 1)
	tester("NestedStructFieldPtr.SliceOfDeeplyNestedStructsField", 1)
	tester("SliceOfStringsField", 1)
	tester("SliceOfStructsField", 0)
	tester("SliceOfStructsField.[0].StringField", 1)
	tester("SliceOfStructsField.[0].IntField", 1)
	tester("SliceOfStructsField.[0].DeeplyNestedStructField", 0)
	tester("SliceOfStructsField.[0].SliceOfStringsField", 1)
	tester("SliceOfStructsField.[0].SliceOfDeeplyNestedStructsField", 1)
}

// this build a struct with all fields filled in with non zero values and validates it
func TestEverythingIsValid(t *testing.T) {
	failures, err := Validate(mainSubject{
		StringField: "valid_string",
		IntField:    10,
		NestedStructField: nestedStruct{
			StringField: "valid_string",
			IntField:    100,
			DeeplyNestedStructField: deeplyNestedStruct{
				StringField: "valid",
				IntField:    1000,
			},
			SliceOfStringsField:             []string{"valid", "slice"},
			SliceOfDeeplyNestedStructsField: make([]deeplyNestedStruct, 0),
		},
		NestedStructFieldPtr: &nestedStruct{
			StringField: "valid_string",
			IntField:    100,
			DeeplyNestedStructField: deeplyNestedStruct{
				StringField: "valid",
				IntField:    1000,
			},
			SliceOfStringsField:             []string{"valid", "slice"},
			SliceOfDeeplyNestedStructsField: make([]deeplyNestedStruct, 0),
		},
		SliceOfStringsField: []string{"val", "id", "_string"},
		SliceOfStructsField: []nestedStruct{{
			StringField: "nested_valid_string",
			IntField:    100,
			DeeplyNestedStructField: deeplyNestedStruct{
				StringField: "valid",
				IntField:    1000,
			},
			SliceOfStringsField:             []string{"nested_", "valid", "_string"},
			SliceOfDeeplyNestedStructsField: make([]deeplyNestedStruct, 0),
		}},
	})
	if err != nil {
		t.Error("did not expect an error")
	}

	noExpectedFailures := 0
	if len(failures) != noExpectedFailures {
		t.Errorf("want: %d failures; got: %d failures, %v", noExpectedFailures, len(failures), failures)
	}

	tester := func(fieldName string, noFailsWant int) {
		fs, ok := failures[fieldName]
		if !ok && 0 != noFailsWant {
			t.Errorf("expected failures for '%s'", fieldName)
		}
		if len(fs) != noFailsWant {
			t.Errorf("number of failures for '%s' want: %d; got: %d", fieldName, noFailsWant, len(fs))
		}
	}

	tester("StringField", 0)
	tester("IntField", 0)
	tester("NestedStructField", 0)
	tester("NestedStructField.StringField", 0)
	tester("NestedStructField.IntField", 0)
	tester("NestedStructField.DeeplyNestedStructField", 0)
	tester("NestedStructField.DeeplyNestedStructField.StringField", 0)
	tester("NestedStructField.DeeplyNestedStructField.IntField", 0)
	tester("NestedStructField.SliceOfStringsField", 0)
	tester("NestedStructField.SliceOfDeeplyNestedStructsField", 0)
	tester("SliceOfStringsField", 0)
	tester("SliceOfStructsField", 0)
	tester("SliceOfStructsField.[0].StringField", 0)
	tester("SliceOfStructsField.[0].IntField", 0)
	tester("SliceOfStructsField.[0].DeeplyNestedStructField", 0)
	tester("SliceOfStructsField.[0].SliceOfStringsField", 0)
	tester("SliceOfStructsField.[0].SliceOfDeeplyNestedStructsField", 0)
}
