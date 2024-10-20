package runtime

var (
	trueValue     = NewNumber(1, 1)
	falseValue    = NewNumber(0, 1)
	voidValue     = NewVoid()
	emptyStrValue = NewString("")
)

var numValue = map[int]Value{
	-1: NewNumber(-1, 1),
	-2: NewNumber(-2, 1),
	-3: NewNumber(-3, 1),
	-4: NewNumber(-4, 1),
	-5: NewNumber(-5, 1),
	0:  NewNumber(0, 1),
	1:  NewNumber(1, 1),
	2:  NewNumber(2, 1),
	3:  NewNumber(3, 1),
	4:  NewNumber(4, 1),
	5:  NewNumber(5, 1),
	6:  NewNumber(6, 1),
	7:  NewNumber(7, 1),
	8:  NewNumber(8, 1),
	9:  NewNumber(9, 1),
	10: NewNumber(10, 1),
	11: NewNumber(11, 1),
	12: NewNumber(12, 1),
	13: NewNumber(13, 1),
	14: NewNumber(14, 1),
	15: NewNumber(15, 1),
	16: NewNumber(16, 1),
	17: NewNumber(17, 1),
	18: NewNumber(18, 1),
	19: NewNumber(19, 1),
}
