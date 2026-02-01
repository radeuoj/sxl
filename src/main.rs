use std::process::Command;

use crate::{compiler::Compiler, lexer::Lexer, parser::Parser};

mod ast;
mod compiler;
mod lexer;
mod parser;
mod token;

fn main() -> anyhow::Result<()> {
    let input = std::fs::read("test.sxl")?;
    let lexer = Lexer::new(input);
    let mut parser = Parser::new(lexer);
    let program = parser.parse_program()?;
    let compiler = Compiler::new();
    let output = compiler.compile_program(program);
    std::fs::write("test.c", output)?;

    Command::new("clang")
        .args(["-o", "test.exe", "test.c"])
        .spawn()?
        .wait()?;

    let path = std::env::current_dir()?;
    Command::new(path.join("test.exe"))
        .spawn()?
        .wait()?;

    Ok(())
}
