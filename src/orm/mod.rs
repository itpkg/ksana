
pub trait Db {
}

pub trait Model {
    fn table_name() -> String;
    fn find<T: Model> (id: usize) -> T;
    fn delete(id: usize);
    fn update(id: usize);
}

pub struct Query {
    script: String,
}

impl Query {
    fn new() -> Query {
        Query{script: "".to_string()}
    }
}

