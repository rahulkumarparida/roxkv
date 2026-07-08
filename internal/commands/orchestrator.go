package commands

type functionOccuerence struct {
	FuncName string `json:"funcName"`
	Params   string `json:"params"`
	Count    int    `json:"count"`
}

var Output string

var listOfTools []functionOccuerence = []functionOccuerence{
	{
		"GetSnapshot",
		"1 parameters: fromPast by default 5 mins depending",
		0,
	},
	{
		"DisplaySnapshot",
		"no parametrs",
		0,
	},
	{
		"GetKeys",
		"no parameters",
		0,
	},
}

func Iterator() {

}
