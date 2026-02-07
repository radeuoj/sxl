use crate::token::Token;

#[derive(Debug, PartialEq)]
pub enum Expression {
    Ident {
        value: String,
    },
    Int {
        value: i32,
    },
    String {
        value: String,
    },
    Unary {
        op: Token,
        right: Box<Expression>,
    },
    Binary {
        op: Token,
        left: Box<Expression>,
        right: Box<Expression>,
    },
    Call {
        func: Box<Expression>,
        args: Vec<Expression>,
    },
}

impl std::fmt::Display for Expression {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        use Expression::*;

        let res = match self {
            Ident { value } => value.to_string(),
            Int { value } => value.to_string(),
            String { value } => format!("\"{value}\""),
            Unary { op, right } => format!("({op}{right})"),
            Binary { op, left, right } => format!("({left} {op} {right})"),
            Call { func, args } => format!("{func}({})", args
                .iter()
                .map(|arg| format!("{arg}"))
                .reduce(|acc, s| format!("{acc}, {s}"))
                .unwrap_or_default())
        };

        write!(f, "{res}")
    }
}

#[derive(Debug, PartialEq)]
pub enum Statement {
    Let { name: String, vtype: String, value: Option<Expression> },
    Return { value: Expression },
    Expression { value: Expression },
    Block { body: Vec<Statement> },
}

impl std::fmt::Display for Statement  {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        use Statement::*;

        let res = match self {
            Let { name, vtype, value } => {
                match value {
                    Some(value) => format!("let {}: {} = {value};",
                        name, vtype),
                    None => format!("let {}: {};",
                        name, vtype),
                }
            }
            Return { value } => format!("return {value};"),
            Expression { value } => format!("{value};"),
            Block { body } => format!("{{\n{}\n}}", body
                .iter()
                .map(|stmt| format!("{stmt}"))
                .reduce(|acc, stmt| format!("{acc}\n{stmt}"))
                .unwrap_or_default()),
        };

        write!(f, "{res}")
    }
}

pub struct Program {
    pub body: Vec<Statement>
}
