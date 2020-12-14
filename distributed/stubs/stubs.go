package stubs

import "uk.ac.bris.cs/gameoflife/util"

var LogicOne  = "LogicOp.LogicTwo"
var Ticker = "LogicOp.SendCells"
var Key = "LogicOp.SendKey"




type Response struct {
	World [][]uint8
CellCount int
	CompletedTurns int
Turns int
	FinalAlive []util.Cell
}

type Request struct {
World[][]uint8
	CellCount int
	CompletedTurns int
	Turns int
    FinalAlive []util.Cell

}

type Ticks struct {
	Turns int
	CellCount int

}

type ResTicks struct {
	Turns int
	CellCount int

}

type Keys struct{
	Val int
}

type ResKeys struct {
	Val int
}
