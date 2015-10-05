mod redis;

pub trait Store {
    fn get(&self, k: &str) -> String;
    fn set(&self, k: &str, v: &str);
}
