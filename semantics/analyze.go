package semantics

import (
	"fmt"

	"niklaskorz.de/nklang/ast"
)

func AnalyzeLookups(p *ast.Program) error {
	globalScope := &definitionScope{definitions: make(definitionSet)}

	for _, n := range p.Statements {
		if err := analyzeStatement(globalScope, n); err != nil {
			return err
		}
	}

	return nil
}

func analyzeStatement(scope *definitionScope, n ast.Statement) error {
	switch s := n.(type) {
	case *ast.IfStatement:
		if s.Condition != nil {
			if err := analyzeExpression(scope, s.Condition); err != nil {
				return err
			}
		}
		ds := scope.newScope()
		for _, n := range s.Statements {
			if err := analyzeStatement(ds, n); err != nil {
				return err
			}
		}
		if s.ElseBranch != nil {
			if err := analyzeStatement(scope, s.ElseBranch); err != nil {
				return err
			}
		}
	case *ast.WhileStatement:
		if err := analyzeExpression(scope, s.Condition); err != nil {
			return err
		}
		ds := scope.newScope()
		for _, n := range s.Statements {
			if err := analyzeStatement(ds, n); err != nil {
				return err
			}
		}
	case *ast.DeclarationStatement:
		if err := analyzeExpression(scope, s.Value); err != nil {
			return err
		}
		if scope.definitions.has(s.Identifier) {
			return fmt.Errorf("Redeclaration of %s in same scope", s.Identifier)
		}
		scope.declare(s.Identifier)
	case *ast.AssignmentStatement:
		if err := analyzeExpression(scope, s.Value); err != nil {
			return err
		}
		scopeIndex := scope.lookup(s.Identifier, 0)
		if scopeIndex == -1 {
			return fmt.Errorf("%s must be declared before assignment", s.Identifier)
		}
		s.ScopeIndex = scopeIndex
	case *ast.ReturnStatement:
		if err := analyzeExpression(scope, s.Expression); err != nil {
			return err
		}
	case *ast.ExpressionStatement:
		if err := analyzeExpression(scope, s.Expression); err != nil {
			return err
		}
	}

	return nil
}

func analyzeExpression(scope *definitionScope, n ast.Expression) error {
	switch e := n.(type) {
	case *ast.IfExpression:
		if e.Condition != nil {
			if err := analyzeExpression(scope, e.Condition); err != nil {
				return err
			}
		}
		if err := analyzeExpression(scope, e.Value); err != nil {
			return err
		}
		if e.ElseBranch != nil {
			if err := analyzeExpression(scope, e.ElseBranch); err != nil {
				return err
			}
		}
	case *ast.BinaryOperationExpression:
		if err := analyzeExpression(scope, e.A); err != nil {
			return err
		}
		if err := analyzeExpression(scope, e.B); err != nil {
			return err
		}
	case *ast.LookupExpression:
		scopeIndex := scope.lookup(e.Identifier, 0)
		if scopeIndex == -1 {
			return fmt.Errorf("%s must be declared before usage", e.Identifier)
		}
		e.ScopeIndex = scopeIndex
	case *ast.CallExpression:
		if err := analyzeExpression(scope, e.Callee); err != nil {
			return err
		}
		for _, p := range e.Parameters {
			if err := analyzeExpression(scope, p); err != nil {
				return err
			}
		}
	case *ast.Function:
		ds := scope.newScope()
		for _, p := range e.Parameters {
			ds.declare(p)
		}
		for _, s := range e.Statements {
			if err := analyzeStatement(ds, s); err != nil {
				return err
			}
		}
	}

	return nil
}