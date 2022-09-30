package binchunk

import (
	"encoding/binary"
	"fmt"
	"math"
)

type reader struct {
	data []byte
}

func (self *reader) readByte() byte {
	b := self.data[0]
	fmt.Printf("%x\n", b)
	self.data = self.data[1:]
	return b
}
func (self *reader) readBytes(n uint) []byte {
	bytes := self.data[:n]
	self.data = self.data[n:]
	fmt.Printf("%x\n", bytes)
	return bytes
}
func (self *reader) readUint32() uint32 {
	i := binary.LittleEndian.Uint32(self.data)

	// trace
	bytes := self.data[:4]
	fmt.Printf("%x\n", bytes)

	self.data = self.data[4:]

	return i
}
func (self *reader) readUint64() uint64 {
	i := binary.LittleEndian.Uint64(self.data)
	self.data = self.data[8:]
	return i
}
func (self *reader) readLuaNumber() float64 {
	return math.Float64frombits(self.readUint64())
}

func (self *reader) readString() string {
	size := uint(self.readUint32())
	if size == 0 {
		return ""
	}
	bytes := self.readBytes(size)
	return string(bytes) // todo
}

func (self *reader) checkHeader() {
	if string(self.readBytes(4)) != LUA_SINGATURE {
		panic("not a precompiled chunk!")
	}
	if self.readByte() != LUAC_VERSION {
		panic("version mismatch!")
	}

	if self.readByte() != LUAC_FORMAT {
		panic("endianness mismatch!")
	}
	if self.readByte() != LUAC_DATA {
		panic("corrupted!")
	}
	if self.readByte() != CINT_SIZE {
		panic("int size mismatch!")
	}
	if self.readByte() != CSIZET_SIZE {
		panic("size_t size mismatch!")
	}
	if self.readByte() != INSTRUCTION_SIZE {
		panic("instruction size mismatch!")
	}
	if self.readByte() != LUA_NUMBER_SIZE {
		panic("lua_Number size mismatch!")
	}
	if self.readByte() != LUAC_INTFlag {
		panic("int flag mismatch!")
	}
}
func (self *reader) readProto(parentSource string) *Prototype {
	source := self.readString()
	if source == "" {
		source = parentSource
	}
	return &Prototype{
		Source:          source,
		LineDefined:     self.readUint32(),
		LastLineDefined: self.readUint32(),
		UpvalueNum:      self.readByte(),
		NumParams:       self.readByte(),
		IsVararg:        self.readByte(),
		MaxStackSize:    self.readByte(),
		Code:            self.readCode(),
		Constants:       self.readConstants(),
		Protos:          self.readProtos(source),
		LineInfo:        self.readLineInfo(),
		LocVars:         self.readLocVars(),
		EndLine:         self.readUint32(),
	}
}

func (self *reader) readCode() []uint32 {
	code := make([]uint32, self.readUint32())
	for i := range code {
		code[i] = self.readUint32()
	}
	return code
}
func (self *reader) readConstants() []interface{} {
	constants := make([]interface{}, self.readUint32())
	for i := range constants {
		constants[i] = self.readConstant()
	}
	return constants
}

func (self *reader) readConstant() interface{} {
	switch self.readByte() {
	case TAG_NIL:
		return nil
	case TAG_BOOLEAN:
		return self.readByte() != 0
	case TAG_NUMBER:
		return self.readLuaNumber()
	case TAG_SHORT_STR, TAG_LONG_STR:
		return self.readString()
	default:
		panic("corrupted!") // todo
	}
}
func (self *reader) readProtos(parentSource string) []*Prototype {
	protos := make([]*Prototype, self.readUint32())
	for i := range protos {
		protos[i] = self.readProto(parentSource)
	}
	return protos
}
func (self *reader) readLineInfo() []uint32 {
	lineInfo := make([]uint32, self.readUint32())
	for i := range lineInfo {
		lineInfo[i] = self.readUint32()
	}
	return lineInfo
}
func (self *reader) readLocVars() []LocVar {
	locVars := make([]LocVar, self.readUint32())
	for i := range locVars {
		locVars[i] = LocVar{
			VarName: self.readString(),
			StartPC: self.readUint32(),
			EndPC:   self.readUint32(),
		}
	}
	return locVars
}
