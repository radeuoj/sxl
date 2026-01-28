use anyhow::bail;

use crate::{ast::Expression, lexer::Lexer, token::Token};

pub struct Parser {
    lexer: Lexer,
    peek_token: Token,
}

pub enum BindingPower {
    Lowest,
    Assign,
    Equals,
    Sum,
    Product,
    Unary,
    Call,
}

impl Parser {
    pub fn new(mut lexer: Lexer) -> Self {
        Self {
            peek_token: lexer.next_token(),
            lexer,
        }
    }

    fn next_token(&mut self) -> Token {
        std::mem::replace(&mut self.peek_token, self.lexer.next_token())
    }

    fn expect_peek(&mut self, token: &Token) -> anyhow::Result<()> {
        if self.peek_token == *token {
            self.next_token();
            Ok(())
        } else {
            anyhow::bail!("expected {} but got {}", token, self.peek_token)
        }
    }

    fn get_binding_power(token: &Token) -> Option<BindingPower> {
        use BindingPower::*;

        match token {
            Token::Equal | Token::NotEqual | Token::Lt | Token::Gt
            | Token::Lte | Token::Gte => Some(Equals),
            Token::Plus | Token::Minus => Some(Sum),
            Token::Asterisk | Token::Slash => Some(Product),
            Token::Assign => Some(Assign),
            Token::LParen => Some(Call),
            _ => None,
        }
    }

    pub fn parse_expression(&mut self, bpow: BindingPower) -> anyhow::Result<Expression> {
        let left = match self.next_token() {
            Token::Ident(name) => Expression::Ident { value: name },
            Token::Int(lit) => self.parse_int(&lit)?,
            op @ (Token::Minus | Token::Bang) => self.parse_unary_expression(op)?,
            token => bail!("invalid prefix operator {}", token),
        };

        Ok(left)
    }

    fn parse_int(&self, lit: &[u8]) -> anyhow::Result<Expression> {
        Ok(Expression::Int {
            value: i32::from_str_radix(str::from_utf8(lit)?, 10)?
        })
    }

    fn parse_unary_expression(&mut self, op: Token) -> anyhow::Result<Expression> {
        Ok(Expression::Unary {
            op,
            right: Box::new(self.parse_expression(BindingPower::Unary)?)
        })
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_expect_peek() -> anyhow::Result<()> {
        let input = b"let a: int32 = 10;";
        let lexer = Lexer::new(input.to_vec());
        let mut parser = Parser::new(lexer);

        use Token::*;
        parser.expect_peek(&Let)?;
        parser.expect_peek(&Ident(b"a".to_vec()))?;
        parser.expect_peek(&Colon)?;
        parser.expect_peek(&Ident(b"int32".to_vec()))?;
        parser.expect_peek(&Assign)?;
        parser.expect_peek(&Int(b"10".to_vec()))?;
        parser.expect_peek(&Semicolon)?;
        parser.expect_peek(&Eof)?;

        Ok(())
    }

    #[test]
    fn test_unary_expressions() -> anyhow::Result<()> {
        let input = b"!!-!100";
        let lexer = Lexer::new(input.to_vec());
        let mut parser = Parser::new(lexer);

        assert_eq!(parser.parse_expression(BindingPower::Product)?, Expression::Unary {
            op: Token::Bang,
            right: Box::new(Expression::Unary {
                op: Token::Bang,
                right: Box::new(Expression::Unary {
                    op: Token::Minus,
                    right: Box::new(Expression::Unary {
                        op: Token::Bang,
                        right: Box::new(Expression::Int { value: 100 })
                    })
                })
            })
        });

        Ok(())
    }
}
