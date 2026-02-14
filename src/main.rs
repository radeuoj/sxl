use std::process::Command;
use anyhow::{Context, bail};

use crate::{compiler::Compiler, lexer::Lexer, parser::Parser};

mod ast;
mod compiler;
mod lexer;
mod parser;
mod token;
mod environment;

enum Mode {
    Compile { file: String },
    CompileAndRun { file: String },
    LexerRepl,
    ParserRepl,
}

impl Mode {
    fn from_args(args: std::env::Args) -> anyhow::Result<Self> {
        let mut file = None;
        let mut run_mode = false;

        for arg in args.skip(1) {
            match arg.as_str() {
                "--compile" => (),
                "--run" => run_mode = true,
                "--lexer-repl" => return Ok(Mode::LexerRepl),
                "--parser-repl" => return Ok(Mode::ParserRepl),
                arg => file = Some(arg.to_string()),
            };
        }

        let Some(file) = file else {
            bail!("No file attached!")
        };

        Ok(if run_mode {
            Mode::CompileAndRun { file }
        } else {
            Mode::Compile { file }
        })
    }

    fn compile_file(&self, file: &str) -> anyhow::Result<()> {
        let input = std::fs::read(file)?;
        let lexer = Lexer::new(input);
        let mut parser = Parser::new(lexer)?;

        match parser.parse_program() {
            Ok(program) => {
                let compiler = Compiler::new();
                let output = compiler.compile_program(program);
                std::fs::write(format!("{file}.c"), output)?;
                Ok(())
            }
            Err(err) => {
                eprintln!("{}", err);
                bail!("Compilation failed");
            },
        }
    }

    fn run(self) -> anyhow::Result<()> {
        match self {
            Mode::Compile { ref file } => self.compile_file(file)?,
            Mode::CompileAndRun { ref file } => {
                self.compile_file(file)?;
                let exe = if file.ends_with(".sxl") {
                    &file[..file.len()-4]
                } else {
                    file
                };

                Command::new("clang")
                    .args(["-o", exe, &format!("{file}.c")])
                    .spawn().with_context(|| "Clang not found")?
                    .wait()?;

                let path = std::env::current_dir()?;
                Command::new(path.join(exe))
                    .spawn()?
                    .wait()?;
            }
            Mode::LexerRepl => lexer::repl(),
            Mode::ParserRepl => parser::repl(),
        }

        Ok(())
    }
}

fn main() -> anyhow::Result<()> {
    let mode = Mode::from_args(std::env::args());
    mode?.run()
}
