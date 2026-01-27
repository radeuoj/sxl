#[derive(Debug, Clone, PartialEq)]
pub enum Token {
    Illegal,
    Eof,

    Ident(Vec<u8>),
    Int(Vec<u8>),

    Assign,
    Plus,
    Minus,
    Asterisk,
    Slash,
    Bang,
    Equal,
    NotEqual,
    Lt,
    Gt,
    Lte,
    Gte,

    Comma,
    Colon,
    Semicolon,

    LParen,
    RParen,
    LBrace,
    RBrace,

    Let,
    Fn,
    If,
    Else,
    Return,
    True,
    False,
}

impl Token {
    pub fn from_symbol(symbol: &[u8]) -> Self {
        use Token::*;
        match symbol {
            b"let" => Let,
            b"fn" => Fn,
            b"if" => If,
            b"else" => Else,
            b"return" => Return,
            b"true" => True,
            b"false" => False,
            _ => Ident(symbol.to_vec()),
        }
    }
}
