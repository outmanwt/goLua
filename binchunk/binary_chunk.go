package binchunk

type binaryChunk struct {
	header                  // 头部
	sizeUpvalues byte       // 主函数upvalue数量
	mainFunc     *Prototype // 主函数原型
}

type header struct {
	signature       [4]byte // .Lua
	version         byte    // 51 版本号
	luacData        byte    // 00 官方保留
	endianness      byte    // 01,字节序，1表示Little Endian，0表示Big Endian
	intSize         byte    // 04
	sizeSize        byte    // 04
	instructionSize byte    // 04
	numberSize      byte    // 08
	intFlag         byte    // 0=float 1=int
}

const (
	LUA_SINGATURE    = "\x1bLua"
	LUAC_VERSION     = 0x51
	LUAC_FORMAT      = 0
	LUAC_DATA        = 1
	CINT_SIZE        = 4
	CSIZET_SIZE      = 4
	INSTRUCTION_SIZE = 4
	LUA_NUMBER_SIZE  = 8
	LUAC_INTFlag     = 0
)
const (
	TAG_NIL       = 0x00
	TAG_BOOLEAN   = 0x01
	TAG_NUMBER    = 0x03
	TAG_INTEGER   = 0x13
	TAG_SHORT_STR = 0x04
	TAG_LONG_STR  = 0x14
)

// function prototype
type Prototype struct {
	Source          string // debug
	LineDefined     uint32
	LastLineDefined uint32
	UpvalueNum      byte
	NumParams       byte
	IsVararg        byte
	MaxStackSize    byte
	// codeSize        uint32
	Code []uint32
	// constantSize    unit32
	Constants []interface{}
	// NumSubFun uint32
	Protos []*Prototype
	// LenInstruction	uint32
	LineInfo []uint32 // debug
	// NumLoc	uint32
	// LenLoc	uint32
	LocVars []LocVar // debug
	EndLine uint32
}

type LocVar struct {
	VarName string
	StartPC uint32
	EndPC   uint32
}

func Undump(data []byte) *Prototype {
	reader := &reader{data}
	reader.checkHeader()
	return reader.readProto("")
}
