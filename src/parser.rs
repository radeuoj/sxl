use std::io::Write;
use anyhow::bail;
use crate::{ast::{Expression, Statement, Program}, lexer::Lexer, token::Token};

pub struct Parser {
    lexer: Lexer,
    peek_token: Token,
}

#[derive(PartialEq, PartialOrd)]
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
    pub fn new(mut lexer: Lexer) -> anyhow::Result<Self> {
        Ok(Self {
            peek_token: lexer.next_token()?,
            lexer,
        })
    }

    fn next_token(&mut self) -> anyhow::Result<Token> {
        Ok(std::mem::replace(&mut self.peek_token, self.lexer.next_token()?))
    }

    fn expect_peek(&mut self, token: &Token) -> anyhow::Result<()> {
        if self.peek_token == *token {
            self.next_token()?;
            Ok(())
        } else {
            anyhow::bail!("expected {} but got {}", token, self.peek_token)
        }
    }

    fn expect_ident(&mut self) -> anyhow::Result<String> {
        match self.next_token()? {
            Token::Ident(value) => Ok(value),
            token => anyhow::bail!("expected identifier but got {}", token)
        }
    }

    fn get_binding_power(token: &Token) -> BindingPower {
        use BindingPower::*;

        match token {
            Token::Equal | Token::NotEqual | Token::Lt | Token::Gt
            | Token::Lte | Token::Gte => Equals,
            Token::Plus | Token::Minus => Sum,
            Token::Asterisk | Token::Slash => Product,
            Token::Assign => Assign,
            Token::LParen => Call,
            _ => Lowest,
        }
    }

    fn get_peek_binding_power(&self) -> BindingPower {
        Self::get_binding_power(&self.peek_token)
    }

    fn parse_expression(&mut self, bpow: BindingPower) -> anyhow::Result<Expression> {
        let mut left = match self.next_token()? {
            Token::Ident(name) => Expression::Ident { value: name },
            Token::Int(lit) => self.parse_int(&lit)?,
            Token::String(lit) => Expression::String { value: lit },
            op @ (Token::Minus | Token::Bang) => self.parse_unary_expression(op)?,
            token => bail!("invalid prefix operator {}", token),
        };

        while self.peek_token != Token::Semicolon
                && bpow < self.get_peek_binding_power() {
            left = match self.peek_token {
                Token::Equal | Token::NotEqual | Token::Lt | Token::Lte
                | Token::Gt | Token::Gte | Token::Plus | Token::Minus
                | Token::Asterisk | Token::Slash
                | Token::Assign => self.parse_binary_expression(left)?,
                Token::LParen => self.parse_call_expression(left)?,
                _ => return Ok(left),
            }
        }

        Ok(left)
    }

    fn parse_int(&self, lit: &str) -> anyhow::Result<Expression> {
        Ok(Expression::Int {
            value: i32::from_str_radix(lit, 10)?
        })
    }

    fn parse_unary_expression(&mut self, op: Token) -> anyhow::Result<Expression> {
        Ok(Expression::Unary {
            op,
            right: Box::new(self.parse_expression(BindingPower::Unary)?)
        })
    }

    fn parse_binary_expression(&mut self, left: Expression) -> anyhow::Result<Expression> {
        let op = self.next_token()?;
        let bpow = Parser::get_binding_power(&op);
        let right = self.parse_expression(bpow)?;

        Ok(Expression::Binary {
            op,
            left: Box::new(left),
            right: Box::new(right),
        })
    }

    fn parse_call_expression(&mut self, left: Expression) -> anyhow::Result<Expression> {
        Ok(Expression::Call {
            func: left.into(),
            args: self.parse_call_arguments()?,
        })
    }

    fn parse_call_arguments(&mut self) -> anyhow::Result<Vec<Expression>> {
        self.next_token()?; // (
        let mut args = vec![];

        if self.peek_token == Token::RParen {
            self.next_token()?;
            return Ok(args);
        }

        loop {
            args.push(self.parse_expression(BindingPower::Lowest)?);
            if self.peek_token != Token::Comma { break; }
            self.next_token()?;
        }

        self.expect_peek(&Token::RParen)?;

        Ok(args)
    }

    fn parse_statement(&mut self) -> anyhow::Result<Statement> {
        let res = match self.peek_token {
            Token::Let => {
                self.next_token()?; // let
                let name = self.expect_ident()?;

                self.expect_peek(&Token::Colon)?;
                let vtype = self.expect_ident()?;

                let value = match self.peek_token {
                    Token::Assign => {
                        self.next_token()?; // =
                        Some(self.parse_expression(BindingPower::Lowest)?)
                    }
                    _ => None,
                };

                Statement::Let { name, vtype, value }
            }
            Token::Return => {
                self.next_token()?; // return
                let value = self.parse_expression(BindingPower::Lowest)?;
                Statement::Return { value }
            }
            Token::If => {
                self.next_token()?; // if
                let cond = self.parse_expression(BindingPower::Lowest)?;
                let then = self.parse_block_statement()?;

                let else_then = if self.peek_token == Token::Else {
                    self.next_token()?;
                    Some(self.parse_block_statement()?)
                } else {
                    None
                };

                return Ok(Statement::If {
                    cond,
                    then: then.into(),
                    else_then: else_then.map(Box::new),
                });
            }
            Token::LBrace => self.parse_block_statement()?,
            _ => Statement::Expression {
                value: self.parse_expression(BindingPower::Lowest)?
            },
        };

        self.expect_peek(&Token::Semicolon)?;

        Ok(res)
    }

    fn parse_block_statement(&mut self) -> anyhow::Result<Statement> {
        self.expect_peek(&Token::LBrace)?;
        let mut body = vec![];
        let mut errs = vec![];

        while ![Token::Eof, Token::RBrace].contains(&self.peek_token) {
            match self.parse_statement() {
                Ok(stmt) => body.push(stmt),
                Err(err) => errs.push(err),
            }
        }

        if let Err(err) = self.expect_peek(&Token::RBrace) {
            errs.push(err);
        }

        if errs.is_empty() {
            Ok(Statement::Block { body })
        } else {
            bail!("{:?}", errs)
        }
    }

    pub fn parse_program(&mut self) -> anyhow::Result<Program> {
        let mut body = vec![];
        let mut errs = vec![];

        while self.peek_token != Token::Eof {
            match self.parse_statement() {
                Ok(stmt) => body.push(stmt),
                Err(err) => errs.push(err),
            }
        }

        if !errs.is_empty() { anyhow::bail!(format!("{:?}", errs)) }
        Ok(Program { body })
    }
}

#[allow(unused)]
pub fn repl() {
    println!("Welcome to temporary REPL");

    loop {
        print!("> ");
        std::io::stdout().flush().unwrap();

        let mut buffer = String::new();
        std::io::stdin().read_line(&mut buffer).unwrap();

        let lexer = Lexer::new(buffer.into_bytes());
        let mut parser = Parser::new(lexer).unwrap();

        match parser.parse_statement() {
            Ok(expr) => println!("{expr}"),
            Err(err) => eprintln!("{err}"),
        };
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_expect_peek() -> anyhow::Result<()> {
        let input = b"let a: int32 = 10;";
        let lexer = Lexer::new(input.to_vec());
        let mut parser = Parser::new(lexer)?;

        use Token::*;
        parser.expect_peek(&Let)?;
        parser.expect_peek(&Ident("a".to_string()))?;
        parser.expect_peek(&Colon)?;
        parser.expect_peek(&Ident("int32".to_string()))?;
        parser.expect_peek(&Assign)?;
        parser.expect_peek(&Int("10".to_string()))?;
        parser.expect_peek(&Semicolon)?;
        parser.expect_peek(&Eof)?;

        Ok(())
    }

    #[test]
    fn test_unary_expressions() -> anyhow::Result<()> {
        let input = b"!!-!100";
        let lexer = Lexer::new(input.to_vec());
        let mut parser = Parser::new(lexer)?;

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
