use std::process::Command;
use crate::{compiler::Compiler, lexer::Lexer, parser::Parser};

mod ast;
mod compiler;
mod lexer;
mod parser;
mod token;

enum Mode {
    CompileAndRun,
    LexerRepl,
    ParserRepl,
}

impl Mode {
    fn from_args(args: std::env::Args) -> Self {
        match args.skip(1).next()
            .unwrap_or("--compile-and-run".to_string()).as_str() {
            "--compile-and-run" => Mode::CompileAndRun,
            "--lexer-repl" => Mode::LexerRepl,
            "--parser-repl" => Mode::ParserRepl,
            _ => unreachable!(),
        }
    }

    fn run(self) -> anyhow::Result<()> {
        Ok(match self {
            Mode::CompileAndRun => {
                let input = std::fs::read("test.sxl")?;
                let lexer = Lexer::new(input);
                let mut parser = Parser::new(lexer)?;
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
            }
            Mode::LexerRepl => lexer::repl(),
            Mode::ParserRepl => parser::repl(),
        })
    }
}

fn main() -> anyhow::Result<()> {
    let mode = Mode::from_args(std::env::args());
    mode.run()
}
