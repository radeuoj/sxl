use crate::{ast::{Expression, Program, Statement}};

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

int main() {{
{}
}}
            "#, program.body.iter()
                .map(|stmt| self.compile_statement(stmt, 1))
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
                If { cond, then, else_then } => format!("if ({}) {{\n{}\n{indent_str}}} {}",
                    self.compile_expression(cond),
                    self.compile_statements(then.unwrap_block().unwrap(), indent + 1),
                    match else_then {
                        Some(else_then) => format!("else {{\n{}\n{indent_str}}}",
                            self.compile_statements(else_then.unwrap_block().unwrap(), indent + 1)),
                        None => "".to_string(),
                    }),
                Expression { value } => format!("{};",
                    self.compile_expression(&value)),
                Block { body } => format!("{{\n{}\n{indent_str}}}",
                    self.compile_statements(body, indent + 1)),
            }
        )
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
                self.compile_expression(&right)),
            Binary { op, left, right } => format!("{} {op} {}",
                self.compile_expression(&left),
                self.compile_expression(&right)),
            Call { func, args } => format!("{}({})",
                self.compile_expression(func),
                args.iter()
                    .map(|arg| self.compile_expression(arg))
                    .reduce(|acc, s| format!("{acc}, {s}"))
                    .unwrap_or_default()),
        }
    }
}
