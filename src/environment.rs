use std::collections::HashMap;
use crate::ast::{Symbol, ValueType};

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
        self.symbols.get(name).map(|symbol| &symbol.vtype)
    }

    pub fn does_vtype_exist(&self, vtype: &str) -> bool {
        self.symbols.contains_key(vtype)
    }
}