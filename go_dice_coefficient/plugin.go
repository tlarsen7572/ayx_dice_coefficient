package main

import (
	"encoding/xml"
	"github.com/tlarsen7572/goalteryx/sdk"
)

type Configuration struct {
	Text1       string
	Text2       string
	OutputField string
}

type DiceCoefficientPlugin struct {
	configuration Configuration
	provider      sdk.Provider
	hasError      bool
	text1         sdk.IncomingStringField
	text2         sdk.IncomingStringField
	outgoingInfo  *sdk.OutgoingRecordInfo
	scoreField    sdk.OutgoingFloatField
	output        sdk.OutputAnchor
}

func (p *DiceCoefficientPlugin) Init(provider sdk.Provider) {
	p.provider = provider
	err := xml.Unmarshal([]byte(provider.ToolConfig()), &p.configuration)
	if err != nil {
		p.sendError(err)
		return
	}
	p.output = provider.GetOutputAnchor(`Output`)
}

func (p *DiceCoefficientPlugin) OnInputConnectionOpened(connection sdk.InputConnection) {
	var err error
	incomingInfo := connection.Metadata()
	p.text1, err = incomingInfo.GetStringField(p.configuration.Text1)
	if err != nil {
		p.sendError(err)
		return
	}
	p.text2, err = incomingInfo.GetStringField(p.configuration.Text2)
	if err != nil {
		p.sendError(err)
		return
	}
	editor := incomingInfo.Clone()
	outgoingField := editor.AddDoubleField(p.configuration.OutputField, `Dice Coefficient (Go)`)
	p.outgoingInfo = editor.GenerateOutgoingRecordInfo()
	p.scoreField, _ = p.outgoingInfo.FloatFields[outgoingField]
	p.output.Open(p.outgoingInfo)
}

func (p *DiceCoefficientPlugin) OnRecordPacket(connection sdk.InputConnection) {
	packet := connection.Read()
	for packet.Next() {
		p.outgoingInfo.CopyFrom(packet.Record())
		text1, isNull1 := p.text1.GetValue(packet.Record())
		text2, isNull2 := p.text2.GetValue(packet.Record())
		if isNull1 || isNull2 {
			p.scoreField.SetFloat(0)
			p.output.Write()
			continue
		}
		score := CalculateDiceCoefficient(text1, text2)
		p.scoreField.SetFloat(score)
		p.output.Write()
	}
	p.output.UpdateProgress(connection.Progress())
}

func (p *DiceCoefficientPlugin) OnComplete() {

}

func (p *DiceCoefficientPlugin) sendError(err error) {
	p.hasError = true
	p.provider.Io().Error(err.Error())
}
