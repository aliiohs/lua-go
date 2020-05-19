package luogo

// 源文件常量
const (
	LUA_SIGNATURE    = "\x1bLUA"
	LUAC_VERSION     = 0x53
	LUAC_FORMAT      = 0
	LUAC_DATA        = "\x19\x93\r\n\x1a\n"
	CINT_SIZE        = 4
	CSIZET_SIZE      = 8
	INSTRUCTION_SIZE = 4
	LUA_INTEGER_SIZE = 8
	LUA_NUMBER_SIZE  = 8
	LUAC_INT         = 0x5678 //检测二进制chunk大小端，intel为小端
	LUAC_NUM         = 370.5  //检测浮点数格式是否正确，默认IEEE 754
)

// 常量
const (
	TAG_NIL          = 0x00
	TAG_BOOLEAN      = 0x01
	TAG_NUMBER       = 0x03
	TAG_INTEGER      = 0x13
	TAG_SHORT_STRING = 0x04
	TAG_LONG_STRING  = 0x14
)

type UpValue struct {
	Instack byte
	Idx     byte
}

type LocVar struct {
	VarName string
	StartPC uint32
	EndPC   uint32
}

type ProtoType struct {
	Source          string        //源文件名
	LineDefined     uint32        //普通函数大于0，主函数=0
	LastLineDefined uint32        //普通函数大于0，主函数=0
	NumParams       byte          //固定参数个数
	IsVararg        byte          //是否有变长参数，主函数是Vararg
	MaxStackSize    byte          //寄存器数量
	Code            []uint32      //指令表：每个指令占4字节
	Constants       []interface{} // 常量表
	UpValues        []UpValue
	Protos          []*ProtoType
	LineInfo        []uint32
	LocVars         []LocVar
	UpValueNames    []string
}

type BinChunk struct {
	header
	sizeUPValues byte
	mainFunc     *ProtoType
}

type header struct {
	signature       [4]byte
	Version         byte
	format          byte
	luacData        [6]byte
	cintSize        byte
	sizetSize       byte
	instructionSize byte
	luaIntegerSize  byte
	luaNumberSize   byte
	luaInt          int64
	luaNum          float64
}

func Undump(data []byte) *ProtoType {

	r := reader{data}
	r.CheckHeader()
	r.readByte()
	return r.readProto("")
}
