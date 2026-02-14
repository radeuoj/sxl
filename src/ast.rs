use std::ops::Deref;
use crate::token::Token;

#[derive(Debug, Clone, PartialEq)]
pub enum ValueType {
    Type(String),
    Func(FuncDecl)
}

impl ValueType {
    pub fn i32() -> Self {
        Self::Type("i32".to_owned())
    }

    pub fn str() -> Self {
        Self::Type("str".to_owned())
    }
}

#[derive(Debug, Clone, PartialEq)]
pub struct Symbol {
    pub name: String,
    pub vtype: ValueType,
}

#[derive(Debug, PartialEq)]
pub enum ExprKind {
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
pub struct Expression {
    pub kind: ExprKind,
    pub vtype: ValueType,
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

#[derive(Debug, Clone, PartialEq)]
pub struct FuncDecl {
    pub name: String,
    pub vtype: Box<ValueType>,
    pub params: Vec<Symbol>
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

pub struct Program {
    pub body: Vec<Statement>
}
