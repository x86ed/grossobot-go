package main

const (
	//EQUALS 0
	EQUALS = iota
	//CONTAINS 1
	CONTAINS
	//DOESNTEQUAL 2
	DOESNTEQUAL
	//DOESNTCONTAIN 3
	DOESNTCONTAIN
	//REGEX 4
	REGEX
)

var quizJudge = []string{
	"329451587422519297",
	"202218126987755523", //lsd
	"313719596886523904",
}

var judgeChannel = "734466604393431170"
var role = "327840671811502081"
var badRoles = []string{
	"359852475181694976",
	"636246344855453696",
	"716536970309927033",
	"416424719487860736",
	"697875957964603403",
}
