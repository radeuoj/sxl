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

    fn next_token() -> Token {
        todo!()
    }
}
