package main

/*
#include "implementation.h"
*/
import "C"
import (
	"encoding/xml"
	"fmt"
	"goalteryx/api"
	"goalteryx/output_connection"
	"goalteryx/recordinfo"
	"unsafe"
)

func main() {}

//export AlteryxGoPlugin
func AlteryxGoPlugin(toolId C.int, xmlProperties unsafe.Pointer, engineInterface unsafe.Pointer, pluginInterface unsafe.Pointer) C.long {
	myPlugin := &MyNewPlugin{
		Output1: output_connection.New(int(toolId), `Output1`),
		Blah:    output_connection.New(int(toolId), `Blah`),
	}
	return C.long(api.ConfigurePlugin(myPlugin, int(toolId), xmlProperties, engineInterface, pluginInterface))
}

type MyNewPlugin struct {
	ToolId  int
	Field   string
	Output1 output_connection.OutputConnection
	Blah    output_connection.OutputConnection
}

type ConfigXml struct {
	Field string `xml:"Field"`
}

func (plugin *MyNewPlugin) Init(toolId int, config string) bool {
	plugin.ToolId = toolId
	var c ConfigXml
	err := xml.Unmarshal([]byte(config), &c)
	if err != nil {
		api.OutputMessage(toolId, api.Error, err.Error())
		return false
	}
	plugin.Field = c.Field
	return true
}

func (plugin *MyNewPlugin) PushAllRecords(recordLimit int) bool {
	return true
}

func (plugin *MyNewPlugin) Close(hasErrors bool) {

}

func (plugin *MyNewPlugin) AddIncomingConnection(connectionType string, connectionName string) api.IncomingInterface {
	return &MyPluginIncomingInterface{Parent: plugin}
}

func (plugin *MyNewPlugin) AddOutgoingConnection(connectionName string, connectionInterface *api.ConnectionInterfaceStruct) bool {
	if connectionName == `Output1` {
		plugin.Output1.Add(connectionInterface)
	} else {
		plugin.Blah.Add(connectionInterface)
	}
	return true
}

type MyPluginIncomingInterface struct {
	Parent   *MyNewPlugin
	inInfo   recordinfo.RecordInfo
	blahInfo recordinfo.RecordInfo
}

func (ii *MyPluginIncomingInterface) Init(recordInfoIn string) bool {
	var err error
	ii.inInfo, err = recordinfo.FromXml(recordInfoIn)
	if err != nil {
		api.OutputMessage(ii.Parent.ToolId, api.Error, err.Error())
		return false
	}
	ii.blahInfo = recordinfo.New()
	ii.blahInfo.AddByteField(`hello`, `goalteryx`)

	err = ii.Parent.Output1.Init(ii.inInfo)
	if err != nil {
		api.OutputMessage(ii.Parent.ToolId, api.Error, err.Error())
		return false
	}
	err = ii.Parent.Blah.Init(ii.blahInfo)
	if err != nil {
		api.OutputMessage(ii.Parent.ToolId, api.Error, err.Error())
		return false
	}
	return true
}

func (ii *MyPluginIncomingInterface) PushRecord(record unsafe.Pointer) bool {
	var value interface{}
	var isNull bool
	var err error
	value, isNull, err = ii.inInfo.GetInterfaceValueFrom(ii.Parent.Field, record)
	if err != nil {
		api.OutputMessage(ii.Parent.ToolId, api.Error, err.Error())
		return false
	}
	if isNull {
		api.OutputMessage(ii.Parent.ToolId, api.TransientInfo, fmt.Sprintf(`[%v] is null`, ii.Parent.Field))
	} else {
		api.OutputMessage(ii.Parent.ToolId, api.TransientInfo, fmt.Sprintf(`[%v] is %v`, ii.Parent.Field, value))
	}

	for index := 0; index < ii.inInfo.NumFields(); index++ {
		field, err := ii.inInfo.GetFieldByIndex(index)
		if err != nil {
			api.OutputMessage(ii.Parent.ToolId, api.Error, err.Error())
			return false
		}
		value, isNull, err := ii.inInfo.GetInterfaceValueFrom(field.Name, record)
		if err != nil {
			api.OutputMessage(ii.Parent.ToolId, api.Error, err.Error())
			return false
		}
		if isNull {
			_ = ii.inInfo.SetFieldNull(field.Name)
			continue
		}
		err = ii.inInfo.SetFromInterface(field.Name, value)
		if err != nil {
			api.OutputMessage(ii.Parent.ToolId, api.Error, err.Error())
			return false
		}
	}

	outputRecord, err := ii.inInfo.GenerateRecord()
	if err != nil {
		api.OutputMessage(ii.Parent.ToolId, api.Error, err.Error())
		return false
	}
	ii.Parent.Output1.PushRecord(outputRecord)
	byteVal, isNull, err := ii.inInfo.GetByteValueFrom(`ByteField`, record)
	if err != nil {
		api.OutputMessage(ii.Parent.ToolId, api.Error, err.Error())
		return false
	}
	err = ii.blahInfo.SetByteField(`hello`, byteVal)
	if err != nil {
		api.OutputMessage(ii.Parent.ToolId, api.Error, err.Error())
		return false
	}
	blahRecord, err := ii.blahInfo.GenerateRecord()
	if err != nil {
		api.OutputMessage(ii.Parent.ToolId, api.Error, err.Error())
		return false
	}
	ii.Parent.Blah.PushRecord(blahRecord)
	return true
}

func (ii *MyPluginIncomingInterface) UpdateProgress(percent float64) {

}

func (ii *MyPluginIncomingInterface) Close() {
	ii.Parent.Output1.Close()
	ii.Parent.Blah.Close()
}

func (ii *MyPluginIncomingInterface) Free() {

}
