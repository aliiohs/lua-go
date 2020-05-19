package luogo

import (
	"encoding/binary"
	"math"
)

type reader struct {
	data []byte
}

func (r *reader) CheckHeader() {
	if string(r.readBytes(4)) != LUA_SIGNATURE {
		panic("not a precompiled chunk!")
	} else if r.readByte() != LUAC_VERSION {
		panic("version mismatch!")
	} else if r.readByte() != LUAC_FORMAT {
		panic("format mismatch!")
	} else if string(r.readBytes(6)) != LUAC_DATA {
		panic("corrupted!")
	} else if r.readByte() != CINT_SIZE {
		panic("int size mismatch!")
	} else if r.readByte() != CSIZET_SIZE {
		panic("size_t size mismatch!")
	} else if r.readByte() != INSTRUCTION_SIZE {
		panic("instruction size mismatch!")
	} else if r.readByte() != LUA_INTEGER_SIZE {
		panic("lua integer mismatch!")
	} else if r.readByte() != LUA_NUMBER_SIZE {
		panic("lua number size mismatch!")
	} else if r.readLuaInteger() != LUAC_INT {
		panic("endianness mismatch!")
	} else if r.readLuaNumber() != LUAC_NUM {
		panic("float format mismatch!")
	}
}

func (r *reader) readByte() byte {
	b := r.data[0]
	r.data = r.data[1:]
	return b
}
func (r *reader) readBytes(size uint) []byte {
	b := r.data[:size]
	r.data = r.data[size:]
	return b
}
func (r *reader) readUint32() uint32 {
	i32 := binary.LittleEndian.Uint32(r.data)
	r.data = r.data[4:]
	return i32
}
func (r *reader) readUint64() uint64 {
	i64 := binary.LittleEndian.Uint64(r.data)
	r.data = r.data[8:]
	return i64
}

func (r *reader) readLuaInteger() int64 {
	return int64(r.readUint64())
}

func (r *reader) readLuaNumber() float64 {
	return math.Float64frombits(r.readUint64())
}

func (r *reader) readString() string {
	size := uint(r.readByte())
	if size == 0 {
		return ""
	}
	if size == 0xFF {
		size = uint(r.readUint64())
	}
	bytes := r.readBytes(size - 1)
	return string(bytes)
}

func (r *reader) readProto(parentSource string) *ProtoType {
	source := r.readString()
	if source == "" {
		source = parentSource
	}
	return &ProtoType{
		Source:          source,
		LineDefined:     r.readUint32(),
		LastLineDefined: r.readUint32(),
		NumParams:       r.readByte(),
		IsVararg:        r.readByte(),
		MaxStackSize:    r.readByte(),
		Code:            r.readCode(),
		Constants:       r.readConstants(),
		UpValues:        r.readUpvalue(),
		Protos:          r.readProtos(source),
		LineInfo:        r.readLineInfo(),
		LocVars:         r.readLocVars(),
		UpValueNames:    r.readUpValueNames(),
	}
}

func (r *reader) readCode() []uint32 {
	code := make([]uint32, r.readUint32())
	for i := range code {
		code[i] = r.readUint32()
	}
	return code
}

func (r *reader) readConstants() []interface{} {
	constant := make([]interface{}, r.readUint32())
	for i := range constant {
		constant[i] = r.readConstant()
	}
	return constant
}

func (r *reader) readConstant() interface{} {
	switch r.readByte() {
	case TAG_NIL:
		return nil
	case TAG_BOOLEAN:
		return r.readByte() != 0
	case TAG_INTEGER:
		return r.readLuaInteger()
	case TAG_NUMBER:
		return r.readLuaNumber()
	case TAG_SHORT_STRING:
		return r.readString()
	case TAG_LONG_STRING:
		return r.readString()
	default:
		panic("corrected!")
	}
}

func (r *reader) readUpvalue() []UpValue {
	upValues := make([]UpValue, r.readUint32())
	for i := range upValues {
		upValues[i] = UpValue{
			Instack: r.readByte(),
			Idx:     r.readByte(),
		}
	}
	return upValues
}

func (r *reader) readProtos(parentSource string) []*ProtoType {
	protos := make([]*ProtoType, r.readUint32())
	for i := range protos {
		protos[i] = r.readProto(parentSource)
	}
	return protos
}

func (r *reader) readLineInfo() []uint32 {
	lineInfos := make([]uint32, r.readUint32())
	for i := range lineInfos {
		lineInfos[i] = r.readUint32()
	}
	return lineInfos
}

func (r *reader) readLocVars() []LocVar {
	locvars := make([]LocVar, r.readUint32())
	for i := range locvars {
		locvars[i] = LocVar{
			VarName: r.readString(),
			StartPC: r.readUint32(),
			EndPC:   r.readUint32(),
		}
	}
	return locvars

}

func (r *reader) readUpValueNames() []string {
	upValueNames := make([]string, r.readUint32())
	for i := range upValueNames {
		upValueNames[i] = r.readString()
	}
	return upValueNames
}
