// +build linux,cgo darwin,cgo

package imports

import (
	"path/filepath"
	"plugin"

	"github.com/nitrogen-lang/nitrogen/src/eval"
	"github.com/nitrogen-lang/nitrogen/src/moduleutils"
	"github.com/nitrogen-lang/nitrogen/src/object"
)

func init() {
	eval.RegisterBuiltin("module", importModule)
	eval.RegisterBuiltin("modulesSupported", moduleSupport)
}

func moduleSupport(i object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	return object.NativeBoolToBooleanObj(true)
}

func importModule(i object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	if ac := moduleutils.CheckMinArgs("module", 1, args...); ac != nil {
		return ac
	}

	filepathArg, ok := args[0].(*object.String)
	if !ok {
		return object.NewException("module expected a string, got %s", args[0].Type().String())
	}

	required := false
	if len(args) > 1 {
		requiredArg, ok := args[1].(*object.Boolean)
		if !ok {
			return object.NewException("module expected a boolean for second argument, got %s", args[1].Type().String())
		}
		required = requiredArg.Value
	}

	// Return already registered, named module
	if module := eval.GetModule(filepathArg.Value); module != nil {
		return module
	}

	includedPath := filepath.Clean(filepath.Join(filepath.Dir(i.GetCurrentScriptPath()), filepathArg.Value))
	if !fileExists(includedPath) {
		if required {
			return object.NewException("Module %s not found", filepathArg.Value)
		}
		return object.NewError("Module %s not found", filepathArg.Value)
	}

	p, err := plugin.Open(includedPath)
	if err != nil {
		if required {
			return object.NewException("%s", err)
		}
		return object.NewError("%s", err)
	}

	// Check module name
	moduleNameSym, err := p.Lookup("ModuleName")
	if err != nil {
		// The module didn't declare a name
		return object.NewException("Invalid module %s", filepathArg.Value)
	}

	if module := eval.GetModule(*(moduleNameSym.(*string))); module != nil {
		return module
	}
	return object.NullConst
}
