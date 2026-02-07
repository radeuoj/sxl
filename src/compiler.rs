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
typedef const char* ccp;

int main() {{
    int res;
    {}

    printf("%d\n", res);
}}
            "#, program.body.iter().fold(String::new(),
                |acc, stmt| format!("{acc}\n{}",
                    &self.compile_statement(stmt, 1))
            )
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
                Expression { value } => format!("{};",
                    self.compile_expression(&value)),
                Block { body } => format!("{{\n{}\n{indent_str}}}", body
                    .iter()
                    .map(|stmt| self.compile_statement(stmt, indent + 1))
                    .reduce(|acc, stmt| format!("{acc}\n{stmt}"))
                    .unwrap_or_default()),
                _ => todo!(),
            }
        )
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
            _ => todo!(),
            // Call { func, args } => format!("{func}({})", args
            //     .iter()
            //     .map(|arg| format!("{arg}"))
            //     .reduce(|acc, s| format!("{acc}, {s}"))
            //     .unwrap_or_default())
        }
    }
}
