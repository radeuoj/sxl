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

int main() {{
    int res;
    {}

    printf("%d\n", res);
}}
            "#, program.body.iter().fold(String::new(),
                |acc, stmt| format!("{acc}\n    {}",
                    &self.compile_statement(stmt))
            )
        )
    }

    fn compile_statement(&self, stmt: &Statement) -> String {
        use Statement::*;

        match stmt {
            Let { name, vtype, value } => {
                match value {
                    Some(value) => format!("{} {} = {value};",
                        str::from_utf8(vtype).unwrap(),
                        str::from_utf8(name).unwrap()),
                    None => format!("{} {};",
                        str::from_utf8(vtype).unwrap(),
                        str::from_utf8(name).unwrap()),
                }
            }
            Expression { value } => format!("{};",
                self.compile_expression(&value)),
            _ => todo!(),
        }
    }

    fn compile_expression(&self, expr: &Expression) -> String {
        use Expression::*;

        match expr {
            Ident { value } => String::from_utf8(value.to_vec()).unwrap(),
            Int { value } => value.to_string(),
            Unary { op, right } => format!("{op}{}",
                self.compile_expression(&right)),
            Binary { op, left, right } => format!("{} {op} {}",
                self.compile_expression(&left),
                self.compile_expression(&right)),
            Call { func, args } => format!("{func}({})", args
                .iter()
                .map(|arg| format!("{arg}"))
                .reduce(|acc, s| format!("{acc}, {s}"))
                .unwrap_or_default())
        }
    }
}
