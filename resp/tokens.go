package resp

type Token int

const (
	TokenEOF Token = iota
	TokenSet
	TokenGet
	TokenEcho
	TokenPing
	TokenDel
	TokenExists
	TokenEx
	TokenSub
	TokenPub
	TokenUnSub
	TokenArg
)

type Symbol rune

const (
	SymbolString     Symbol = '+'
	SymbolError      Symbol = '-'
	SymbolInt        Symbol = ':'
	SymbolBulkString Symbol = '$'
	SymbolArray      Symbol = '*'
	SymbolCR         Symbol = '\r'
	SymbolLF         Symbol = '\n'
)

const (
	CmdSet    = "SET"
	CmdGet    = "GET"
	CmdPing   = "PING"
	CmdEcho   = "ECHO"
	CmdDel    = "DEL"
	CmdExists = "EXISTS"
	CmdSub    = "SUBSCRIBE"
	CmdPub    = "PUBLISH"
	CmdUnSub  = "UNSUBSCRIBE"
)
