use crate::token::*;

pub struct Lexer {
    input: Vec<u8>,
    pos: usize,
}

impl Lexer {
    pub fn new(input: Vec<u8>) -> Self {
        Self { input, pos: 0 }
    }

    fn read_byte(&mut self) -> u8 {
        let byte = self.peek_byte();
        self.pos += 1;
        byte
    }

    fn peek_byte(&self) -> u8 {
        match self.input.get(self.pos) {
            Some(byte) => byte.to_owned(),
            None => 0,
        }
    }

    fn unread_byte(&mut self) {
        self.pos -= 1;
    }

    pub fn next_token(&mut self) -> Token {
        self.skip_whitespace();
        use Token::*;

        match self.read_byte() {
            b'=' => {
                if self.peek_byte() == b'=' {
                    self.read_byte();
                    Equal
                } else {
                    Assign
                }
            }
            b'<' => {
                if self.peek_byte() == b'=' {
                    self.read_byte();
                    Lte
                } else {
                    Lt
                }
            }
            b'>' => {
                if self.peek_byte() == b'=' {
                    self.read_byte();
                    Gte
                } else {
                    Gt
                }
            }
            b'+' => Plus,
            b'-' => Minus,
            b'*' => Asterisk,
            b'/' => {
                if self.peek_byte() == b'/' {
                    self.read_byte();
                    self.skip_comment();
                    self.next_token()
                } else {
                    Slash
                }
            }
            b'!' => {
                if self.peek_byte() == b'=' {
                    self.read_byte();
                    NotEqual
                } else {
                    Bang
                }
            }
            b',' => Comma,
            b':' => Colon,
            b';' => Semicolon,
            b'(' => LParen,
            b')' => RParen,
            b'{' => LBrace,
            b'}' => RBrace,
            b'\0' => Eof,
            ch => {
                if ch.is_ascii_digit() {
                    self.unread_byte();
                    Int(self.read_int().to_vec())
                } else if Self::is_ident_byte(ch) {
                    self.unread_byte();
                    let ident = self.read_ident();
                    Token::from_symbol(ident)
                } else {
                    Illegal
                }
            }
        }
    }

    fn skip_whitespace(&mut self) {
        while " \t\n\r".contains(self.peek_byte() as char) {
            self.read_byte();
        }
    }

    fn skip_comment(&mut self) {
        while !"\0\n".contains(self.peek_byte() as char) {
            self.read_byte();
        }
    }

    fn is_ident_byte(ch: u8) -> bool {
        ch.is_ascii_alphabetic() || ch.is_ascii_digit() || ch == b'_'
    }

    fn read_ident(&mut self) -> &[u8] {
        let start = self.pos;

        while Self::is_ident_byte(self.peek_byte()) {
            self.read_byte();
        }

        &self.input[start..self.pos]
    }

    fn read_int(&mut self) -> &[u8] {
        let start = self.pos;

        while self.peek_byte().is_ascii_digit() {
            self.read_byte();
        }

        &self.input[start..self.pos]
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_lexer() {
        let input = b"let a: int32 = b + c * (12 - _radu); // this is a comment
            return a + b;";

        use Token::*;
        let ok = vec![
            Let,
            Ident(b"a".to_vec()),
            Colon,
            Ident(b"int32".to_vec()),
            Assign,
            Ident(b"b".to_vec()),
            Plus,
            Ident(b"c".to_vec()),
            Asterisk,
            LParen,
            Int(b"12".to_vec()),
            Minus,
            Ident(b"_radu".to_vec()),
            RParen,
            Semicolon,
            Return,
            Ident(b"a".to_vec()),
            Plus,
            Ident(b"b".to_vec()),
            Semicolon,
        ];

        let mut lexer = Lexer::new(input.to_vec());
        let mut out = vec![];

        let mut tok = lexer.next_token();
        while tok != Eof {
            out.push(tok);
            tok = lexer.next_token();
        }

        assert_eq!(ok, out);
    }
}
