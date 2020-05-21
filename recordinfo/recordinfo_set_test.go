package recordinfo_test

import (
	"goalteryx/recordinfo"
	"strings"
	"testing"
	"time"
)

func TestSetValuesAndGenerateRecord(t *testing.T) {
	recordInfo := generateTestRecordInfo()
	setRecordInfoTestData(recordInfo)

	record, err := recordInfo.GenerateRecord()
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}

	byteVal, isNull, err := recordInfo.GetIntValueFrom(`ByteField`, record)
	checkExpectedGetValueFrom(t, byteVal, 1, isNull, false, err, nil)

	boolVal, isNull, err := recordInfo.GetBoolValueFrom(`BoolField`, record)
	checkExpectedGetValueFrom(t, boolVal, true, isNull, false, err, nil)

	int16Val, isNull, err := recordInfo.GetIntValueFrom(`Int16Field`, record)
	checkExpectedGetValueFrom(t, int16Val, 2, isNull, false, err, nil)

	int32Val, isNull, err := recordInfo.GetIntValueFrom(`Int32Field`, record)
	checkExpectedGetValueFrom(t, int32Val, 3, isNull, false, err, nil)

	int64Val, isNull, err := recordInfo.GetIntValueFrom(`Int64Field`, record)
	checkExpectedGetValueFrom(t, int64Val, 4, isNull, false, err, nil)

	fixedDecimalVal, isNull, err := recordInfo.GetFloatValueFrom(`FixedDecimalField`, record)
	checkExpectedGetValueFrom(t, fixedDecimalVal, 123.45, isNull, false, err, nil)

	floatVal, isNull, err := recordInfo.GetFloatValueFrom(`FloatField`, record)
	checkExpectedGetValueFrom(t, floatVal, 654.321, isNull, false, err, nil)

	doubleVal, isNull, err := recordInfo.GetFloatValueFrom(`DoubleField`, record)
	checkExpectedGetValueFrom(t, doubleVal, 909.33, isNull, false, err, nil)

	stringVal, isNull, err := recordInfo.GetStringValueFrom(`StringField`, record)
	checkExpectedGetValueFrom(t, stringVal, `ABCDEFG`, isNull, false, err, nil)

	wstringVal, isNull, err := recordInfo.GetStringValueFrom(`WStringField`, record)
	checkExpectedGetValueFrom(t, wstringVal, `CXVY`, isNull, false, err, nil)

	dateVal, isNull, err := recordInfo.GetDateValueFrom(`DateField`, record)
	expectedDate := time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC)
	checkExpectedGetValueFrom(t, dateVal, expectedDate, isNull, false, err, nil)

	dateTimeVal, isNull, err := recordInfo.GetDateValueFrom(`DateTimeField`, record)
	expectedDate = time.Date(2021, 3, 4, 5, 6, 7, 0, time.UTC)
	checkExpectedGetValueFrom(t, dateTimeVal, expectedDate, isNull, false, err, nil)
}

func TestCachedRecords(t *testing.T) {
	recordInfo := generateTestRecordInfo()
	setRecordInfoTestData(recordInfo)

	record1, _ := recordInfo.GenerateRecord()
	record2, _ := recordInfo.GenerateRecord()
	if record1 != record2 {
		t.Fatalf(`record1 and record2 are 2 different pointers`)
	}
}

func TestSetLongVarDataFieldsAndGenerateRecord(t *testing.T) {
	recordInfo := recordinfo.New()
	recordInfo.AddByteField(`ByteField`, ``)
	recordInfo.AddV_WStringField(`V_WStringField`, ``, 250)
	recordInfo.AddV_StringField(`V_StringField`, ``, 250)

	_ = recordInfo.SetIntField(`ByteField`, 1)
	_ = recordInfo.SetStringField(`V_StringField`, strings.Repeat(`B`, 200))
	_ = recordInfo.SetStringField(`V_WStringField`, strings.Repeat(`A`, 100))

	record, err := recordInfo.GenerateRecord()
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	value, isNull, err := recordInfo.GetStringValueFrom(`V_StringField`, record)
	checkExpectedGetValueFrom(t, value, strings.Repeat(`B`, 200), isNull, false, err, nil)

	value, isNull, err = recordInfo.GetStringValueFrom(`V_WStringField`, record)
	checkExpectedGetValueFrom(t, value, strings.Repeat(`A`, 100), isNull, false, err, nil)
}

func TestSetShortVarDataFieldsAndGenerateRecord(t *testing.T) {
	recordInfo := recordinfo.New()
	recordInfo.AddByteField(`ByteField`, ``)
	recordInfo.AddV_WStringField(`V_WStringField`, ``, 250)
	recordInfo.AddV_StringField(`V_StringField`, ``, 250)

	_ = recordInfo.SetIntField(`ByteField`, 1)
	_ = recordInfo.SetStringField(`V_StringField`, strings.Repeat(`B`, 100))
	_ = recordInfo.SetStringField(`V_WStringField`, strings.Repeat(`A`, 50))

	record, err := recordInfo.GenerateRecord()
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	value, isNull, err := recordInfo.GetStringValueFrom(`V_StringField`, record)
	checkExpectedGetValueFrom(t, value, strings.Repeat(`B`, 100), isNull, false, err, nil)

	value, isNull, err = recordInfo.GetStringValueFrom(`V_WStringField`, record)
	checkExpectedGetValueFrom(t, value, strings.Repeat(`A`, 50), isNull, false, err, nil)
}

func TestSetTinyVarDataFieldsAndGenerateRecord(t *testing.T) {
	recordInfo := recordinfo.New()
	recordInfo.AddByteField(`ByteField`, ``)
	recordInfo.AddV_WStringField(`V_WStringField`, ``, 250)
	recordInfo.AddV_StringField(`V_StringField`, ``, 250)

	_ = recordInfo.SetIntField(`ByteField`, 1)
	_ = recordInfo.SetStringField(`V_StringField`, `B`)
	_ = recordInfo.SetStringField(`V_WStringField`, `A`)

	record, err := recordInfo.GenerateRecord()
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	value, isNull, err := recordInfo.GetStringValueFrom(`V_StringField`, record)
	checkExpectedGetValueFrom(t, value, `B`, isNull, false, err, nil)

	value, isNull, err = recordInfo.GetStringValueFrom(`V_WStringField`, record)
	checkExpectedGetValueFrom(t, value, `A`, isNull, false, err, nil)
}

func TestSetEmptyVarDataFieldsAndGenerateRecord(t *testing.T) {
	recordInfo := recordinfo.New()
	recordInfo.AddByteField(`ByteField`, ``)
	recordInfo.AddV_WStringField(`V_WStringField`, ``, 250)
	recordInfo.AddV_StringField(`V_StringField`, ``, 250)

	_ = recordInfo.SetIntField(`ByteField`, 1)
	_ = recordInfo.SetStringField(`V_StringField`, ``)
	_ = recordInfo.SetStringField(`V_WStringField`, ``)

	record, err := recordInfo.GenerateRecord()
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	value, isNull, err := recordInfo.GetStringValueFrom(`V_StringField`, record)
	checkExpectedGetValueFrom(t, value, ``, isNull, false, err, nil)

	value, isNull, err = recordInfo.GetStringValueFrom(`V_WStringField`, record)
	checkExpectedGetValueFrom(t, value, ``, isNull, false, err, nil)
}

func TestSetNullVarDataFieldsAndGenerateRecord(t *testing.T) {
	recordInfo := recordinfo.New()
	recordInfo.AddByteField(`ByteField`, ``)
	recordInfo.AddV_WStringField(`V_WStringField`, ``, 250)
	recordInfo.AddV_StringField(`V_StringField`, ``, 250)

	_ = recordInfo.SetIntField(`ByteField`, 1)
	_ = recordInfo.SetFieldNull(`V_StringField`)
	_ = recordInfo.SetFieldNull(`V_WStringField`)

	record, err := recordInfo.GenerateRecord()
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	value, isNull, err := recordInfo.GetStringValueFrom(`V_StringField`, record)
	checkExpectedGetValueFrom(t, value, ``, isNull, true, err, nil)

	value, isNull, err = recordInfo.GetStringValueFrom(`V_WStringField`, record)
	checkExpectedGetValueFrom(t, value, ``, isNull, true, err, nil)
}

func generateTestRecordInfo() recordinfo.RecordInfo {
	recordInfo := recordinfo.New()
	recordInfo.AddByteField(`ByteField`, ``)
	recordInfo.AddBoolField(`BoolField`, ``)
	recordInfo.AddInt16Field(`Int16Field`, ``)
	recordInfo.AddInt32Field(`Int32Field`, ``)
	recordInfo.AddInt64Field(`Int64Field`, ``)
	recordInfo.AddFixedDecimalField(`FixedDecimalField`, ``, 19, 2)
	recordInfo.AddFloatField(`FloatField`, ``)
	recordInfo.AddDoubleField(`DoubleField`, ``)
	recordInfo.AddStringField(`StringField`, ``, 10)
	recordInfo.AddWStringField(`WStringField`, ``, 5)
	recordInfo.AddDateField(`DateField`, ``)
	recordInfo.AddDateTimeField(`DateTimeField`, ``)
	return recordInfo
}

func setRecordInfoTestData(recordInfo recordinfo.RecordInfo) {
	_ = recordInfo.SetIntField(`ByteField`, 1)
	_ = recordInfo.SetBoolField(`BoolField`, true)
	_ = recordInfo.SetIntField(`Int16Field`, 2)
	_ = recordInfo.SetIntField(`Int32Field`, 3)
	_ = recordInfo.SetIntField(`Int64Field`, 4)
	_ = recordInfo.SetFloatField(`FixedDecimalField`, 123.45)
	_ = recordInfo.SetFloatField(`FloatField`, 654.321)
	_ = recordInfo.SetFloatField(`DoubleField`, 909.33)
	_ = recordInfo.SetStringField(`StringField`, `ABCDEFG`)
	_ = recordInfo.SetStringField(`WStringField`, `CXVY`)
	_ = recordInfo.SetDateField(`DateField`, time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC))
	_ = recordInfo.SetDateField(`DateTimeField`, time.Date(2021, 3, 4, 5, 6, 7, 0, time.UTC))
}
