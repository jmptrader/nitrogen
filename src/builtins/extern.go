package builtins

import (
	"path/filepath"

	"github.com/nitrogen-lang/nitrogen/src/ast"
	"github.com/nitrogen-lang/nitrogen/src/eval"
	"github.com/nitrogen-lang/nitrogen/src/lexer"
	"github.com/nitrogen-lang/nitrogen/src/object"
	"github.com/nitrogen-lang/nitrogen/src/parser"
)

var included map[string]*ast.Program

func init() {
	eval.RegisterBuiltin("include", includeScript)
	eval.RegisterBuiltin("require", requireScript)
	eval.RegisterBuiltin("evalScript", evalScript)

	included = make(map[string]*ast.Program)
}

func includeScript(env *object.Environment, args ...object.Object) object.Object {
	return commonInclude(false, true, env, args...)
}

func requireScript(env *object.Environment, args ...object.Object) object.Object {
	return commonInclude(true, true, env, args...)
}

func evalScript(env *object.Environment, args ...object.Object) object.Object {
	cleanEnv := object.NewEnvironment()

	envvar, _ := env.Get("_ARGV")
	cleanEnv.CreateConst("_ARGV", envvar.Dup())

	envvar, _ = env.Get("_ENV")
	cleanEnv.CreateConst("_ENV", envvar.Dup())

	return commonInclude(false, false, cleanEnv, args...)
}

func commonInclude(require bool, save bool, env *object.Environment, args ...object.Object) object.Object {
	funcName := "include"
	if require {
		funcName = "require"
	}

	if ac := CheckMinArgs(funcName, 1, args...); ac != nil {
		return ac
	}

	filepathArg, ok := args[0].(*object.String)
	if !ok {
		return object.NewException("%s expected a string, got %s", funcName, args[0].Type().String())
	}

	once := false
	if len(args) > 1 {
		includeOnce, ok := args[1].(*object.Boolean)
		if !ok {
			return object.NewException("%s expected a boolean for second argument, got %s", funcName, args[1].Type().String())
		}
		once = includeOnce.Value
	}

	includedFile := filepath.Clean(filepath.Join(filepath.Dir(eval.GetCurrentScriptPath()), filepathArg.Value))

	program, exists := included[includedFile]
	if exists {
		if once || program == nil {
			return object.NullConst
		}
		return eval.Eval(program, object.NewEnclosedEnv(env))
	}

	l, err := lexer.NewFile(includedFile)
	if err != nil {
		if require {
			return object.NewException("including %s failed %s", includedFile, err.Error())
		}
		return object.NewError("including %s failed %s", includedFile, err.Error())
	}

	p := parser.New(l)
	program = p.ParseProgram()
	if len(p.Errors()) != 0 {
		if require {
			return object.NewException("including %s failed %s", includedFile, p.Errors()[0])
		}
		return object.NewError("including %s failed %s", includedFile, p.Errors()[0])
	}

	if save {
		if once {
			// Create the key, but don't save the parsed script since we don't need it anymore.
			included[includedFile] = nil
		} else {
			included[includedFile] = program
		}
	}
	return eval.Eval(program, object.NewEnclosedEnv(env))
}
