package main

const QueLen = 10
type CData interface {}

type CQueue struct {
	queArr [QueLen]CData
	front int
	rear int
}
