use std::collections::HashMap;
use anyhow::{Result, bail};

use crate::ast::{Symbol, ValueType};

#[derive(Debug)]
pub struct Environment<'a> {
    parent: Option<&'a Environment<'a>>,
    symbols: HashMap<String, Symbol>,
    vtypes: HashMap<String, ()>,
}

impl<'a> Environment<'a> {
    pub fn new() -> Self {
        Self { 
            parent: None,
            symbols: HashMap::new(),
            vtypes: HashMap::new(), 
        }
    }

    pub fn from_parent(parent: &'a Environment) -> Self {
        let mut env: Environment<'a> = Self::new();
        env.parent = Some(parent);
        env
    }

    pub fn get_vtype_of(&self, name: &str) -> Option<&ValueType> {
        self.symbols.get(name)       
            .map(|symbol| &symbol.vtype)
            .or_else(|| self.parent?.get_vtype_of(name))
    }

    pub fn does_vtype_exist(&self, vtype: &str) -> bool {
        self.vtypes.contains_key(vtype) ||
        self.parent.map_or(false, 
            |parent| parent.does_vtype_exist(vtype))
    }

    pub fn push_symbol(&mut self, symbol: Symbol) -> Result<()> {
        if self.symbols.contains_key(&symbol.name) {
            bail!("{} already exists", symbol.name);
        }

        self.symbols.insert(symbol.name.to_owned(), symbol);
        Ok(())
    }

    pub fn push_vtype(&mut self, vtype: ValueType) -> Result<()> {
        match &vtype {
            ValueType::Type(name) => {
                if self.vtypes.contains_key(name) {
                    bail!("{} already exists", name);
                }

                self.vtypes.insert(name.to_owned(), ());
            }
            ValueType::Func(decl) => {
                if self.vtypes.contains_key(&decl.name) {
                    bail!("{} already exists", decl.name);
                }

                self.vtypes.insert(decl.name.to_owned(), ());
            }
        };

        Ok(())
    }
}