package main_test

import (
	"encoding/xml"
	"github.com/tlarsen7572/goalteryx/sdk"
	"testing"
)
import dc "go_dice_coefficient"

func TestConfigUnmarshal(t *testing.T) {
	configStr := `<Configuration>
  <Text1>First</Text1>
  <Text2>Second</Text2>
  <OutputField>Match Score</OutputField>
</Configuration>`
	config := dc.Configuration{}
	err := xml.Unmarshal([]byte(configStr), &config)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	if config.Text1 != `First` {
		t.Fatalf(`expected 'First' but got '%v'`, config.Text1)
	}
	if config.Text2 != `Second` {
		t.Fatalf(`expected 'Second' but got '%v'`, config.Text1)
	}
	if config.OutputField != `Match Score` {
		t.Fatalf(`expected 'Match Score' but got '%v'`, config.Text1)
	}
}

func TestScoreCalc(t *testing.T) {
	score := dc.CalculateDiceCoefficient(`Thomas Larsen`, `Thomas Larsen`)
	if score != 1.0 {
		t.Fatalf(`expected 1 but got %v`, score)
	}
}

func TestEndToEnd(t *testing.T) {
	config := `<Configuration>
  <Text1>Text1</Text1>
  <Text2>Text2</Text2>
  <OutputField>Match Score</OutputField>
</Configuration>`
	plugin := &dc.DiceCoefficientPlugin{}
	runner := sdk.RegisterToolTest(plugin, 1, config)
	runner.ConnectInput(`Input`, `testdata.txt`)
	output := runner.CaptureOutgoingAnchor(`Output`)
	runner.SimulateLifecycle()

	_, ok := output.Data[`Match Score`]
	if !ok {
		t.Fatalf(`expected a Match Score field but it did not exist`)
	}
}
