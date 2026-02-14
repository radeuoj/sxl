use std::io::Write;
use anyhow::{Result, anyhow, bail};
use crate::{ast::{BlockStmt, ExprKind, Expression, FuncDecl, Program, Statement, Symbol, ValueType}, environment::Environment, lexer::Lexer, token::Token};

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

    fn parse_expression(&mut self, bpow: BindingPower, env: &Environment) -> anyhow::Result<Expression> {
        let mut left = match self.next_token()? {
            Token::Ident(name) => self.parse_ident(name, env)?,
            Token::Int(lit) => self.parse_int(&lit)?,
            Token::String(lit) => self.parse_string(lit)?,
            op @ (Token::Minus | Token::Bang) => self.parse_unary_expression(op, env)?,
            token => bail!("invalid prefix operator {}", token),
        };

        while self.peek_token != Token::Semicolon
                && bpow < self.get_peek_binding_power() {
            left = match self.peek_token {
                Token::Equal | Token::NotEqual | Token::Lt | Token::Lte
                | Token::Gt | Token::Gte | Token::Plus | Token::Minus
                | Token::Asterisk | Token::Slash
                | Token::Assign => self.parse_binary_expression(left, env)?,
                Token::LParen => self.parse_call_expression(left, env)?,
                _ => return Ok(left),
            }
        }

        Ok(left)
    }

    fn parse_ident(&self, name: String, env: &Environment) -> anyhow::Result<Expression> {
        let vtype = env
            .get_vtype_of(&name)
            .ok_or_else(|| anyhow!("{} not found in current scope", name))?
            .clone();

        Ok(Expression {
            kind: ExprKind::Ident { value: name },
            vtype,
        })
    }

    fn parse_int(&self, lit: &str) -> anyhow::Result<Expression> {
        Ok(Expression { 
            kind: ExprKind::Int { value: lit.parse()? }, 
            vtype: ValueType::i32(),
        })
    }

    fn parse_string(&self, lit: String) -> Result<Expression> {
        Ok(Expression { 
            kind: ExprKind::String { value: lit }, 
            vtype: ValueType::str(),
        })
    }

    fn parse_unary_expression(&mut self, op: Token, env: &Environment) -> anyhow::Result<Expression> {
        let right = self.parse_expression(BindingPower::Unary, env)?;

        Ok(match &right.vtype {
            vtype @ ValueType::Type(vtypes) 
                if vtypes == "i32" => Expression {
                    kind: ExprKind::Unary {
                        op,
                        right: right.into(),
                    },
                    vtype: ValueType::i32(),
                },
            vtype => bail!("{} is not supported for {:?}", op, vtype),
        })
    }

    fn parse_binary_expression(&mut self, left: Expression, env: &Environment) -> anyhow::Result<Expression> {
        let op = self.next_token()?;
        let bpow = Parser::get_binding_power(&op);
        let right = self.parse_expression(bpow, env)?;

        Ok(match (&left.vtype, &right.vtype) {
            (ValueType::Type(left_vtypes), ValueType::Type(right_vtypes))
                if left_vtypes == "i32" && right_vtypes == "i32" => Expression {
                    kind: ExprKind::Binary {
                        op,
                        left: left.into(),
                        right: right.into(),
                    },
                    vtype: ValueType::i32(),
                },
            vtype => bail!("{} is not supported for {:?}", op, vtype),
        })
    }

    fn parse_call_expression(&mut self, left: Expression, env: &Environment) -> anyhow::Result<Expression> {
        let args = self.parse_call_arguments(env)?;

        Ok(match &left.vtype {
            ValueType::Func(decl) => {
                if args.len() != decl.params.len() {
                    bail!("{:?} expected {} arguments but got {}",
                        left, decl.params.len(), args.len());
                }

                for (arg, param) in args.iter().zip(decl.params.iter()) {
                    if arg.vtype != param.vtype {
                        bail!("parameter {} expected something of type {:?} but got {:?}",
                            param.name, param.vtype, arg);
                    }
                }

                Expression {
                    vtype: *decl.vtype.clone(),
                    kind: ExprKind::Call {
                        func: left.into(),
                        args,
                    },
                }
            }
            _ => bail!("{:?} is not a function", left),
        })
    }

    fn parse_call_arguments(&mut self, env: &Environment) -> anyhow::Result<Vec<Expression>> {
        self.next_token()?; // (
        let mut args = vec![];

        if self.peek_token == Token::RParen {
            self.next_token()?;
            return Ok(args);
        }

        loop {
            args.push(self.parse_expression(BindingPower::Lowest, env)?);
            if self.peek_token != Token::Comma { break; }
            self.next_token()?;
        }

        self.expect_peek(&Token::RParen)?;

        Ok(args)
    }

    fn parse_statement(&mut self, env: &mut Environment) -> anyhow::Result<Statement> {
        let res = match self.peek_token {
            Token::Let => {
                self.next_token()?; // let
                let name = self.expect_ident()?;

                if env.get_vtype_of(&name).is_some() {
                    bail!("{} already exists", name);
                }

                self.expect_peek(&Token::Colon)?;
                let vtype = self.expect_ident()?;

                if !env.does_vtype_exist(&vtype) {
                    bail!("{} not found in current scope", vtype);
                }

                let value = match self.peek_token {
                    Token::Assign => {
                        self.next_token()?; // =
                        let value = self.parse_expression(BindingPower::Lowest, env)?;

                        if ValueType::Type(vtype.clone()) != value.vtype {
                            bail!("{:?} is not of type {}", value, vtype);
                        }

                        Some(value)
                    }
                    _ => None,
                };

                env.push_symbol(Symbol { 
                    name: name.to_owned(), 
                    vtype: ValueType::Type(vtype.to_owned()) 
                })?;

                Statement::Let { name, vtype, value }
            }
            Token::Return => {
                self.next_token()?; // return
                let value = self.parse_expression(BindingPower::Lowest, env)?;
                Statement::Return { value }
            }
            Token::If => {
                self.next_token()?; // if
                let cond = self.parse_expression(BindingPower::Lowest, env)?;
                let then = self.parse_block_statement(env)?;

                let else_then = if self.peek_token == Token::Else {
                    self.next_token()?;
                    Some(self.parse_block_statement(env)?)
                } else {
                    None
                };

                return Ok(Statement::If { cond, then, else_then });
            }
            Token::LBrace => return Ok(self.parse_block_statement(env)?.into()),
            Token::Fn => {
                self.next_token()?; // fn
                let name = self.expect_ident()?;

                if env.get_vtype_of(&name).is_some() {
                    bail!("{} already exists", name);
                }

                let params = self.parse_func_params(env)?;

                self.expect_peek(&Token::Arrow)?;
                let vtype = self.expect_ident()?;

                if !env.does_vtype_exist(&vtype) {
                    bail!("{} not found in current scope", vtype);
                }

                let decl = FuncDecl {
                    name: name.clone(),
                    vtype: ValueType::Type(vtype).into(),
                    params: params.clone(),
                };

                env.push_symbol(Symbol { 
                    name, 
                    vtype: ValueType::Func(decl.clone()),
                })?;

                let mut env = Environment::from_parent(env);

                for param in &params {
                    env.push_symbol(param.clone())?;
                }

                let body = self.parse_block_statement(&env)?;

                return Ok(Statement::Func {
                    decl,
                    body,
                });
            }
            _ => Statement::Expression {
                value: self.parse_expression(BindingPower::Lowest, env)?
            },
        };

        self.expect_peek(&Token::Semicolon)?;

        Ok(res)
    }

    fn parse_block_statement(&mut self, env: &Environment) -> anyhow::Result<BlockStmt> {
        self.expect_peek(&Token::LBrace)?;
        let mut body = vec![];
        let mut errs = vec![];
        let mut env = Environment::from_parent(env);

        while ![Token::Eof, Token::RBrace].contains(&self.peek_token) {
            match self.parse_statement(&mut env) {
                Ok(stmt) => body.push(stmt),
                Err(err) => errs.push(err),
            }
        }

        if let Err(err) = self.expect_peek(&Token::RBrace) {
            errs.push(err);
        }

        if errs.is_empty() {
            Ok(BlockStmt(body))
        } else {
            bail!(errs.iter().map(|err| format!("{err}"))
                .reduce(|acc, err| format!("{acc}\n{err}")).unwrap_or_default());
        }
    }

    fn parse_func_params(&mut self, env: &Environment) -> anyhow::Result<Vec<Symbol>> {
        self.expect_peek(&Token::LParen)?;
        let mut params = vec![];

        if self.peek_token == Token::RParen {
            self.next_token()?;
            return Ok(params);
        }

        loop {
            params.push(self.parse_func_param(env)?);
            if self.peek_token != Token::Comma { break; }
            self.next_token()?;
        }

        self.expect_peek(&Token::RParen)?;

        Ok(params)
    }

    fn parse_func_param(&mut self, env: &Environment) -> anyhow::Result<Symbol> {
        let name = self.expect_ident()?;

        self.expect_peek(&Token::Colon)?;
        let vtype = self.expect_ident()?;

        if !env.does_vtype_exist(&vtype) {
            bail!("{} not found in current scope", vtype);
        }

        Ok(Symbol { name, vtype: ValueType::Type(vtype) })
    }

    pub fn parse_program(&mut self) -> anyhow::Result<Program> {
        let mut body = vec![];
        let mut errs = vec![];
        let mut env = Environment::new();

        env.push_vtype(ValueType::Type("i32".to_owned())).unwrap();
        env.push_vtype(ValueType::Type("str".to_owned())).unwrap();
        env.push_symbol(Symbol { name: "printf".to_owned(), vtype: ValueType::Func(FuncDecl { 
            name: "printf".to_owned(), 
            vtype: ValueType::Type("void".to_owned()).into(), 
            params: vec![Symbol {
                name: "str".to_owned(),
                vtype: ValueType::str(),
            }],
        }) }).unwrap();

        while self.peek_token != Token::Eof {
            match self.parse_statement(&mut env) {
                Ok(stmt) => body.push(stmt),
                Err(err) => errs.push(err),
            }
        }

        if !errs.is_empty() { bail!(errs.iter().map(|err| format!("{err}"))
            .reduce(|acc, err| format!("{acc}\n{err}")).unwrap_or_default()) }
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

        let mut env = Environment::new();

        match parser.parse_statement(&mut env) {
            Ok(expr) => println!("{:?}", expr),
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
}
