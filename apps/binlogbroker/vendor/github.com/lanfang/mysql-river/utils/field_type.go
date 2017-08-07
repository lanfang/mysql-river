package utils

import (
//"reflect"
)

const (
	// bool
	TypeBooleanField = 1 << iota

	// string
	TypeCharField

	// string
	TypeTextField

	// time.Time
	TypeDateField
	// time.Time
	TypeDateTimeField

	// int8
	TypeBitField
	// int16
	TypeSmallIntegerField
	// int32
	TypeIntegerField
	// int64
	TypeBigIntegerField
	// uint8
	TypePositiveBitField
	// uint16
	TypePositiveSmallIntegerField
	// uint32
	TypePositiveIntegerField
	// uint64
	TypePositiveBigIntegerField

	// float64
	TypeFloatField
	// float64
	TypeDecimalField

	RelForeignKey
	RelOneToOne
	RelManyToMany
	RelReverseOne
	RelReverseMany
)

const (
	IsIntegerField        = ^-TypePositiveBigIntegerField >> 4 << 5
	IsPostiveIntegerField = ^-TypePositiveBigIntegerField >> 8 << 9
	IsRelField            = ^-RelReverseMany >> 14 << 15
	IsFieldType           = ^-RelReverseMany<<1 + 1
)

type Fielder interface {
	String() string
	FieldType() int
	SetRaw(interface{}) error
	RawValue() interface{}
}
