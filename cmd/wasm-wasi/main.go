package main

import (
	"fmt"
	"unsafe"
	"sync"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/parser"
)

var (
	nodeMap   = make(map[uint32]*ast.Node)
	nodeMapMu sync.RWMutex
	nextNodeID uint32 = 1
)

func storeNode(n *ast.Node) uint32 {
	if n == nil {
		return 0
	}
	nodeMapMu.Lock()
	defer nodeMapMu.Unlock()
	id := nextNodeID
	nextNodeID++
	nodeMap[id] = n
	return id
}

func getNode(id uint32) *ast.Node {
	nodeMapMu.RLock()
	defer nodeMapMu.RUnlock()
	return nodeMap[id]
}

var memoryPool [][]byte

//go:wasmexport wasm_malloc
func wasm_malloc(size int32) int32 {
	b := make([]byte, size)
	memoryPool = append(memoryPool, b)
	ptr := &b[0]
	return int32(uintptr(unsafe.Pointer(ptr)))
}

//go:wasmexport wasm_free
func wasm_free(ptr int32, size int32) {
}

//go:wasmexport ParseSource
func ParseSource(filenamePtr int32, filenameLen int32, sourcePtr int32, sourceLen int32) (rootId int32) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("ParseSource panic: %v\n", r)
			rootId = 0
		}
	}()

	
	

	
	filename := string(unsafe.Slice((*byte)(unsafe.Pointer(uintptr(filenamePtr))), filenameLen))
	sourceText := string(unsafe.Slice((*byte)(unsafe.Pointer(uintptr(sourcePtr))), sourceLen))
	

	opts := ast.SourceFileParseOptions{
		FileName: filename,
	}
	
	scriptKind := core.ScriptKindTS

	sourceFile := parser.ParseSourceFile(opts, sourceText, scriptKind)
	if sourceFile == nil {
		return 0
	}
	return int32(storeNode(&sourceFile.Node))
}

//go:wasmexport GetNodeChildren
func GetNodeChildren(nodeId int32, outArrayPtr int32) (count int32) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("GetNodeChildren panic: %v\n", r)
			count = -1
		}
	}()

	n := getNode(uint32(nodeId))
	if n == nil {
		return 0
	}

	var children []*ast.Node
	for child := range n.IterChildren() {
		children = append(children, child)
	}

	count = int32(len(children))
	if count == 0 {
		return 0
	}

	size := count * 4
	b := make([]byte, size)
	ptr := &b[0]
	
	outPtrObj := (*int32)(unsafe.Pointer(uintptr(outArrayPtr)))
	*outPtrObj = int32(uintptr(unsafe.Pointer(ptr)))
	
	slice := unsafe.Slice((*int32)(unsafe.Pointer(ptr)), count)
	for i, child := range children {
		slice[i] = int32(storeNode(child))
	}

	return count
}

//go:wasmexport GetNodeKind
func GetNodeKind(nodeId int32) (kind int32) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("GetNodeKind panic: %v\n", r)
			kind = -1
		}
	}()

	n := getNode(uint32(nodeId))
	if n == nil {
		return 0 
	}

	return int32(n.Kind)
}

//go:wasmexport GetNodeText
func GetNodeText(nodeId int32, outPtr int32) (textLen int32) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("GetNodeText panic: %v\n", r)
			textLen = -1
		}
	}()

	n := getNode(uint32(nodeId))
	if n == nil {
		return 0
	}

	text := fmt.Sprintf("NodeKind:%v", n.Kind)

	length := int32(len(text))
	if length == 0 {
		return 0
	}

	b := make([]byte, length)
	copy(b, []byte(text))
	
	outPtrObj := (*int32)(unsafe.Pointer(uintptr(outPtr)))
	*outPtrObj = int32(uintptr(unsafe.Pointer(&b[0])))
	
	return length
}

//go:wasmexport CheckDiagnostics
func CheckDiagnostics() (errId int32) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("CheckDiagnostics panic: %v\n", r)
			errId = 0
		}
	}()
	return 0
}

func main() {}
