use crate::token::Token;

#[derive(Debug, PartialEq)]
pub enum Expression {
    Ident {
        value: Vec<u8>,
    },
    Int {
        value: i32,
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
            Ident { value } => String::from_utf8(value.to_vec()).unwrap(),
            Int { value } => format!("{value}"),
            Unary { op, right } => format!("({op}{right})"),
            Binary { op, left, right } => format!("({left} {op} {right}"),
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
    Let { name: Vec<u8>, value: Expression },
    Return { value: Expression },
    Expression { value: Expression },
}
