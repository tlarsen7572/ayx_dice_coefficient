package main_test

import (
	"encoding/xml"
	"github.com/tlarsen7572/goalteryx/sdk"
	"math"
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
	score = dc.CalculateDiceCoefficient(``, `Thomas Larsen`)
	if score != 0.0 {
		t.Fatalf(`expected 0 but got %v`, score)
	}
	score = dc.CalculateDiceCoefficient(`Thomas Larsen`, ``)
	if score != 0.0 {
		t.Fatalf(`expected 0 but got %v`, score)
	}
	score = dc.CalculateDiceCoefficient(`Thomas Larsen`, `a`)
	if score != 0.0 {
		t.Fatalf(`expected 0 but got %v`, score)
	}
	score = dc.CalculateDiceCoefficient(`Thomas Larsen`, `Thomas Larson`)
	if math.Abs(score-0.83333) > 0.0001 {
		t.Fatalf(`expected 0.833333333333333 but got %v`, score)
	}
	score = dc.CalculateDiceCoefficient(`Hello World`, `How are you`)
	if score != 0 {
		t.Fatalf(`expected 0 but got %v`, score)
	}
	score = dc.CalculateDiceCoefficient(`night`, `nacht`)
	if math.Abs(score-0.25) > 0.0001 {
		t.Fatalf(`expected 0.25 but got %v`, score)
	}
	score = dc.CalculateDiceCoefficient(`AA`, `AAAA`)
	if math.Abs(score-0.5) > 0.0001 {
		t.Fatalf(`expected 0.5 but got %v`, score)
	}
	score = dc.CalculateDiceCoefficient(`AAAA`, `AAAAAA`)
	if math.Abs(score-0.75) > 0.0001 {
		t.Fatalf(`expected 0.75 but got %v`, score)
	}
	score = dc.CalculateDiceCoefficient(`12121212`, `12345678`)
	if math.Abs(score-0.142857) > 0.0001 {
		t.Fatalf(`expected 0.142857 but got %v`, score)
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

	scores, ok := output.Data[`Match Score`]
	if !ok {
		t.Fatalf(`expected a Match Score field but it did not exist`)
	}

	if len(scores) != 6 {
		t.Fatalf(`expected 6 rows but got %v`, len(scores))
	}

	if scores[0] != 0.25 {
		t.Fatalf(`expected first row score to be 0.25 but got %v`, scores[0])
	}

	if scores[3] != 0.0 {
		t.Fatalf(`expected fourth row score to be 0 but got %v`, scores[3])
	}
}
