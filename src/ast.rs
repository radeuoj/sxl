use crate::token::Token;

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
