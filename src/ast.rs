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
    Let { name: Vec<u8>, vtype: Vec<u8>, value: Option<Expression> },
    Return { value: Expression },
    Expression { value: Expression },
}

impl std::fmt::Display for Statement  {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        use Statement::*;

        let res = match self {
            Let { name, vtype, value } => {
                match value {
                    Some(value) => format!("let {}: {} = {value};",
                        String::from_utf8(name.to_vec()).unwrap(),
                        String::from_utf8(vtype.to_vec()).unwrap()),
                    None => format!("let {}: {};",
                        String::from_utf8(name.to_vec()).unwrap(),
                        String::from_utf8(vtype.to_vec()).unwrap()),
                }
            }
            Return { value } => format!("return {value};"),
            Expression { value } => format!("{value};"),
        };

        write!(f, "{res}")
    }
}
