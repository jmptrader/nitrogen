package vm

import (
	"fmt"

	"github.com/nitrogen-lang/nitrogen/src/compiler"
	"github.com/nitrogen-lang/nitrogen/src/object"
)

type blockType byte

const (
	loopBlockT blockType = iota
	tryBlockT
)

type block interface {
	blockType() blockType
}

type forLoopBlock struct {
	start, iter, end int
}

func (b *forLoopBlock) blockType() blockType { return loopBlockT }

type Frame struct {
	lastFrame  *Frame
	code       *compiler.CodeBlock
	stack      []object.Object
	sp         int
	blockStack []block
	bp         int
	Env        *object.Environment
	pc         int
}

func (f *Frame) pushStack(obj object.Object) {
	if f.sp == len(f.stack) {
		panic("Stack overflow")
	}
	f.stack[f.sp] = obj
	f.sp++
}

func (f *Frame) popStack() object.Object {
	if f.sp == 0 {
		panic("Stack exhausted")
	}
	f.sp--
	return f.stack[f.sp]
}

func (f *Frame) getFrontStack() object.Object {
	return f.stack[f.sp-1]
}

func (f *Frame) printStack() {
	for i := f.sp; i >= 0; i-- {
		fmt.Printf(" %d: %s\n", i, f.stack[i].Inspect())
	}
}

func (f *Frame) pushBlock(b block) {
	if f.bp == len(f.blockStack) {
		panic("Block stack overflow")
	}
	f.blockStack[f.bp] = b
	f.bp++
}

func (f *Frame) popBlock() block {
	if f.bp == 0 {
		panic("Block stack exhausted")
	}
	f.bp--
	return f.blockStack[f.bp]
}

func (f *Frame) popBlockUntil(bt blockType) block {
	for f.blockStack[f.bp-1].blockType() != bt {
		f.popBlock()
	}
	return f.blockStack[f.bp-1]
}

func (f *Frame) getCurrentBlock() block {
	return f.blockStack[f.bp-1]
}

type frameStack struct {
	head   *frameStackElement
	length int
}

type frameStackElement struct {
	val  *Frame
	prev *frameStackElement
}

func newFrameStack() *frameStack {
	return &frameStack{}
}

func (s *frameStack) Push(val *Frame) {
	s.head = &frameStackElement{
		val:  val,
		prev: s.head,
	}
	s.length++
}

func (s *frameStack) GetFront() *Frame {
	if s.head == nil {
		return nil
	}
	return s.head.val
}

func (s *frameStack) Pop() *Frame {
	if s.head == nil {
		return nil
	}
	r := s.head.val
	s.head = s.head.prev
	s.length--
	return r
}

func (s *frameStack) Len() int {
	return s.length
}