use crate::ast::{BlockStmt, Expression, FuncDecl, Program, Statement, Symbol, ValueType};

pub struct Compiler {

}

impl Compiler {
    pub fn new() -> Self {
        Self {}
    }

    pub fn compile_program(&self, program: Program) -> String {
        format!(r#"// compiled from SXL
#include <stdio.h>
typedef const char* str;

{}
            "#, program.body.iter()
                .map(|stmt| self.compile_statement(stmt, 0))
                .reduce(|acc, stmt| format!("{acc}\n{stmt}"))
                .unwrap_or_default()
        )
    }

    fn compile_statement(&self, stmt: &Statement, indent: i32) -> String {
        use Statement::*;

        let indent_str = "    ".repeat(indent as usize);
        format!("{}{}",
            indent_str,
            match stmt {
                Let { name, vtype, value } => {
                    match value {
                        Some(value) => format!("{} {} = {};",
                            vtype, name, self.compile_expression(value)),
                        None => format!("{} {};",
                            vtype, name),
                    }
                }
                Return { value } => format!("return {};",
                    self.compile_expression(value)),
                If { cond, then, else_then } => format!("if ({}) {} {}",
                    self.compile_expression(cond),
                    self.compile_block_statement(then, indent),
                    match else_then {
                        Some(else_then) => format!("else {{\n{}\n{indent_str}}}",
                            self.compile_statements(else_then, indent + 1)),
                        None => "".to_string(),
                    }),
                Expression { value } => format!("{};",
                    self.compile_expression(value)),
                Block { body } => self.compile_block_statement(body, indent),
                Func { decl, body } => format!("{} {}",
                    self.compile_func_decl(decl),
                    self.compile_block_statement(body, indent)),
            }
        )
    }

    fn compile_block_statement(&self, block: &BlockStmt, indent: i32) -> String {
        format!("{{\n{}\n{}}}",
            self.compile_statements(block, indent + 1),
            "    ".repeat(indent as usize))
    }

    fn compile_statements(&self, stmts: &[Statement], indent: i32) -> String {
        stmts.iter()
            .map(|stmt| self.compile_statement(stmt, indent))
            .reduce(|acc, stmt| format!("{acc}\n{stmt}"))
            .unwrap_or_default()
    }

    fn compile_expression(&self, expr: &Expression) -> String {
        use Expression::*;

        match expr {
            Ident { value } => value.to_string(),
            Int { value } => value.to_string(),
            String { value } => format!("\"{value}\""),
            Unary { op, right } => format!("{op}{}",
                self.compile_expression(right)),
            Binary { op, left, right } => format!("{} {op} {}",
                self.compile_expression(left),
                self.compile_expression(right)),
            Call { func, args } => format!("{}({})",
                self.compile_expression(func),
                args.iter()
                    .map(|arg| self.compile_expression(arg))
                    .reduce(|acc, s| format!("{acc}, {s}"))
                    .unwrap_or_default()),
        }
    }

    fn compile_func_decl(&self, decl: &FuncDecl) -> String {
        format!("{} {}({})", 
            decl.vtype, decl.name,
            decl.params.iter()
                .map(|param| self.compile_symbol(param))
                .reduce(|acc, s| format!("{acc}, {s}"))
                .unwrap_or_default())
    }

    fn compile_symbol(&self, symbol: &Symbol) -> String {
        match &symbol.vtype {
            ValueType::Type(vtype) => format!("{} {}",
                vtype, symbol.name),
            ValueType::Func(_) => todo!(),
        }
    }
}
