package recordinfo

import (
	"encoding/xml"
	"fmt"
	"strconv"
)

type xmlMetaInfo struct {
	XMLName    string      `xml:"MetaInfo"`
	Connection string      `xml:"connection,attr"`
	Fields     []*xmlField `xml:"RecordInfo>Field"`
}

type xmlField struct {
	Name   string `xml:"name,attr"`
	Source string `xml:"source,attr"`
	Size   string `xml:"size,attr"`
	Scale  string `xml:"scale,attr"`
	Type   string `xml:"type,attr"`
}

func FromXml(recordInfoXml string) (RecordInfo, error) {
	var metaInfo xmlMetaInfo
	err := xml.Unmarshal([]byte(recordInfoXml), &metaInfo)
	if err != nil {
		return nil, fmt.Errorf(`error creating RecordInfo from xml: %v`, err.Error())
	}
	recordInfo := New()
	for index, field := range metaInfo.Fields {
		switch field.Type {
		case byteType:
			recordInfo.AddByteField(field.Name, field.Source)
		case boolType:
			recordInfo.AddBoolField(field.Name, field.Source)
		case int16Type:
			recordInfo.AddInt16Field(field.Name, field.Source)
		case int32Type:
			recordInfo.AddInt32Field(field.Name, field.Source)
		case int64Type:
			recordInfo.AddInt64Field(field.Name, field.Source)
		case fixedDecimalType:
			size, err := strconv.Atoi(field.Size)
			if err != nil {
				return nil, fmt.Errorf(`error converting field %v size to an int.  Provided size was %v`, index, field.Size)
			}
			scale, err := strconv.Atoi(field.Scale)
			if err != nil {
				return nil, fmt.Errorf(`error converting field %v scale to an int.  Provided size was %v`, index, field.Size)
			}
			recordInfo.AddFixedDecimalField(field.Name, field.Source, size, scale)
		case floatType:
			recordInfo.AddFloatField(field.Name, field.Source)
		case doubleType:
			recordInfo.AddDoubleField(field.Name, field.Source)
		case stringType:
			size, err := strconv.Atoi(field.Size)
			if err != nil {
				return nil, fmt.Errorf(`error converting field %v size to an int.  Provided size was %v`, index, field.Size)
			}
			recordInfo.AddStringField(field.Name, field.Source, size)
		case wStringType:
			size, err := strconv.Atoi(field.Size)
			if err != nil {
				return nil, fmt.Errorf(`error converting field %v size to an int.  Provided size was %v`, index, field.Size)
			}
			recordInfo.AddWStringField(field.Name, field.Source, size)
		case v_StringType:
			size, err := strconv.Atoi(field.Size)
			if err != nil {
				return nil, fmt.Errorf(`error converting field %v size to an int.  Provided size was %v`, index, field.Size)
			}
			recordInfo.AddV_StringField(field.Name, field.Source, size)
		case v_WStringType:
			size, err := strconv.Atoi(field.Size)
			if err != nil {
				return nil, fmt.Errorf(`error converting field %v size to an int.  Provided size was %v`, index, field.Size)
			}
			recordInfo.AddV_WStringField(field.Name, field.Source, size)
		case dateType:
			recordInfo.AddDateField(field.Name, field.Source)
		case dateTimeType:
			recordInfo.AddDateTimeField(field.Name, field.Source)
		default:
			continue
		}
	}
	return recordInfo, nil
}

func (info *recordInfo) ToXml(connection string) (string, error) {
	fields := make([]*xmlField, 0)
	for _, field := range info.fields {
		fields = append(fields, &xmlField{
			Name:   field.Name,
			Source: field.Source,
			Size:   strconv.Itoa(field.Size),
			Scale:  strconv.Itoa(field.Precision),
			Type:   fieldTypeMap[field.Type],
		})
	}
	recordInfo := xmlMetaInfo{XMLName: `MetaInfo`, Connection: connection, Fields: fields}
	metaInfo, err := xml.Marshal(recordInfo)
	if err != nil {
		return ``, fmt.Errorf(`error converting recordinfo to xml: %v`, err.Error())
	}
	return string(metaInfo), nil
}
