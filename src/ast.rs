use std::ops::Deref;

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

#[derive(Debug, PartialEq)]
pub struct FuncParam {
    pub name: String,
    pub vtype: String,
}

#[derive(Debug, PartialEq)]
pub struct BlockStmt(pub Vec<Statement>);

impl Deref for BlockStmt {
    type Target = Vec<Statement>;

    fn deref(&self) -> &Self::Target {
        &self.0
    }
}

impl From<BlockStmt> for Statement {
    fn from(value: BlockStmt) -> Self {
        Statement::Block { body: value }
    }
}

#[derive(Debug, PartialEq)]
pub enum Statement {
    Let { name: String, vtype: String, value: Option<Expression> },
    Return { value: Expression },
    If { cond: Expression, then: BlockStmt, else_then: Option<BlockStmt> },
    Expression { value: Expression },
    Block { body: BlockStmt },
    Func { decl: FuncDecl, body: BlockStmt },
}

#[derive(Debug, PartialEq)]
pub struct FuncDecl {
    pub name: String,
    pub vtype: String,
    pub params: Vec<FuncParam>
}

pub struct Program {
    pub body: Vec<Statement>
}
