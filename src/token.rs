pub enum Token {
    Illegal,
    Eof,

    Ident(String),
    Int(String),

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
    pub fn from_symbol(symbol: &str) -> Self {
        use Token::*;
        match symbol {
            "let" => Let,
            "fn" => Fn,
            "if" => If,
            "else" => Else,
            "return" => Return,
            "true" => True,
            "false" => False,
            _ => Ident(symbol.to_string()),
        }
    }
}
