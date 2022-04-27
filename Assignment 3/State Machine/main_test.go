package main

import (
	"testing"
)
import "github.com/stretchr/testify/assert"

func setup() {
	for k := range StateTable {
		delete(StateTable, k)
	}
	for k := range TransitionTable {
		delete(TransitionTable, k)
	}
}

func TestPathNotFound(t *testing.T) {
	setup()
	path = "this_file_does_not_exist.machine"
	assert.Errorf(t, parse(), "Opening a nonexistent file should give an error!")
}

func TestParseSuccess(t *testing.T) {
	setup()
	path = "test_machines/test_machine_success_parsable.machine"
	assert.Nil(t, parse(), "valid state machine not parsed correctly")

	expectedStateTable := make(map[string]state)
	init1(expectedStateTable)

	expectedTransitionTable := make(map[key]value)
	init2(expectedTransitionTable, expectedStateTable)

	assert.Equal(t, expectedStateTable, StateTable, "state table not as expected")
	assert.Equal(t, expectedTransitionTable, TransitionTable, "transition table not as expected")

}

func TestParseFail_SourceStateNotExists(t *testing.T) {
	setup()
	path = "test_machines/test_machine_fail_transition_invalid_source.machine"
	assert.Error(t, parse(), "expected an Error when parsing an invalid file")
}

func TestParseFail_DestinationStateNotExists(t *testing.T) {
	setup()
	path = "test_machines/test_machine_fail_transition_invalid_destination.machine"
	assert.Error(t, parse(), "expected an Error when parsing an invalid file")
}

func TestParseFail_TwoStateDefinitions(t *testing.T) {
	setup()
	path = "test_machines/test_machine_fail_identical_state_definition.machine"
	assert.Error(t, parse(), "expected an Error when parsing a file with identical states")
}

func TestParseFail_TwoTransitionDefinitions(t *testing.T) {
	setup()
	path = "test_machines/test_machine_fail_transition_invalid_different_destination.machine"
	assert.Error(t, parse(), "expected an Error when parsing a file with ambiguous transition definitions")
}

func TestValidateSuccess(t *testing.T) {
	setup()
	path = "test_machines/test_machine_success.machine"
	err := parse()
	if err != nil {
		t.Fatal("Parsing failed")
	}
	s, err := isValid(StateTable, TransitionTable)
	expected := state{stateName: "Park", stateText: "\n  The transmission is in \"park\".\n  " +
		"(Drive) Put the transmission into \"drive\"\n  (Leave) Leave the car (quit)\n",
		stateType: startState}

	assert.Equal(t, expected, s, "start states are not the same")
	assert.Nil(t, err, "the validity check yielded an error, none was expected")
}

func TestValidateFail_TwoStartStates(t *testing.T) {
	setup()
	path = "test_machines/test_machine_success_parsable.machine"
	err := parse()
	if err != nil {
		t.Fatal("Parsing failed")
	}
	_, err = isValid(StateTable, TransitionTable)
	assert.NotNil(t, err, "expected error when validating a state machine with two start states")
}

func TestValidateFail_NoEndState(t *testing.T) {
	setup()
	path = "test_machines/test_machine_fail_no_end_state.machine"
	err := parse()
	if err != nil {
		t.Fatal("Parsing failed")
	}
	_, err = isValid(StateTable, TransitionTable)
	assert.NotNil(t, err, "expected error when validating a state machine with no end states")
}

func TestValidateFail_NoRunPossible(t *testing.T) {
	setup()
	path = "test_machines/test_machine_fail_no_run_possible.machine"
	err := parse()
	if err != nil {
		t.Fatal("Parsing failed")
	}
	_, err = isValid(StateTable, TransitionTable)
	assert.NotNil(t, err, "expected error when validating a state machine with no run from the start"+
		"state to any end states")
}

func init1(expectedStateTable map[string]state) {
	expectedStateTable["Start"] = state{
		stateName: "Start",
		stateText: "foo",
		stateType: startState,
	}
	expectedStateTable["Start2"] = state{
		stateName: "Start2",
		stateText: "foo",
		stateType: startState,
	}
	expectedStateTable["state1"] = state{
		stateName: "state1",
		stateText: "foo",
		stateType: normalState,
	}
	expectedStateTable["state2"] = state{
		stateName: "state2",
		stateText: " fo\nfo\nfo\n",
		stateType: normalState,
	}
	expectedStateTable["End"] = state{
		stateName: "End",
		stateText: "x",
		stateType: endState,
	}
	expectedStateTable["End2"] = state{
		stateName: "End2",
		stateText: "\n            fo\n            ",
		stateType: endState,
	}
}
func init2(expectedTransitionTable map[key]value, table map[string]state) {
	expectedTransitionTable[key{state: table["Start"], action: "foo"}] = value{state: table["state1"], description: "blabla"}
	expectedTransitionTable[key{state: table["Start"], action: "foofoo"}] = value{state: table["state2"], description: "blabla"}
	expectedTransitionTable[key{state: table["state1"], action: "fo"}] = value{state: table["End"], description: "huhu"}
	expectedTransitionTable[key{state: table["state2"], action: "gogo"}] = value{state: table["End2"], description: "huhu"}
}
