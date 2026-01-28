mod ast;
mod lexer;
mod parser;
mod token;

use lexer::Lexer;
use parser::Parser;
use parser::BindingPower;
use std::io::Write;

fn main() -> anyhow::Result<()> {
    println!("Welcome to temporary REPL");

    loop {
        print!("> ");
        std::io::stdout().flush()?;

        let mut buffer = String::new();
        std::io::stdin().read_line(&mut buffer)?;

        let lexer = Lexer::new(buffer.into_bytes());
        let mut parser = Parser::new(lexer);
        println!("{}", parser.parse_expression(BindingPower::Lowest)?);
    }
}
