package Netcat

import ("sync"
"net")

type Client struct {
	conn net.Conn
	name string
	Id   int
}

var HistoryMessage []string
var Clients        []Client
var ClientsNames []string
var mutex sync.Mutex

var logo        = []string{
	"          _nnnn_",
	"         dGGGGMMb",
	"        @p~qp~~qMb",
	"        M|@||@) M|",
	"        @,----.JM|",
	"       JS^\\__/  qKL",
	"      dZP        qKRb",
	"     dZP          qKKb",
	"    fZP            SMMb",
	"    HZM            MMMM",
	"    FqM            MMMM",
	"  __| \".        |\\dS\"qML",
	"  |    `.       | `' \\Zq",
	" _)      \\.___.,|     .'",
	" \\____   )MMMMMP|   .'",
	"      `-'       `--'",
}