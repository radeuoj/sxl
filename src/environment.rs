use crate::ast::FuncDecl;

#[derive(Debug)]
pub enum ValueType {
    Type(String),
    Func(FuncDecl)
}

#[derive(Debug)]
pub struct Symbol {
    ident: String,
    vtype: ValueType,
}