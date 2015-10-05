extern crate redis;

use super::Store;

#[allow(dead_code)]
pub struct RedisStore {
    cli: redis::Client,
}

impl RedisStore {
    #[allow(dead_code)]
    fn new(host: &str, port: u16, db: i64) -> RedisStore {
        let addr =  redis::ConnectionAddr::Tcp(host.to_string(), port);
        let cli = redis::Client::open(redis::ConnectionInfo{
            addr: Box::new(addr),
            db: db,
            passwd: None,
        }).unwrap();
        RedisStore{cli: cli}
    }
}

impl Store for RedisStore {
    fn get(&self, k: &str) -> String {
        match self.cli.get_connection(){
            Ok(con) => {
                match redis::cmd("GET").arg(k).query(&con) {
                    Ok(v) => v,
                    Err(_) => "".to_string(),
                }
            },
            Err(_) => "".to_string(),
        }
    }
    
    fn set(&self, k: &str, v: &str) {
        match self.cli.get_connection() {
            Ok(con) => {
                redis::cmd("SET").arg(k).arg(v).execute(&con);
            }
            Err(_) =>{},
        };        
    }
}


#[test]
fn test_redis_store() {
    let k = "test://hello";
    let v = "hello, ksana!";
    let s = RedisStore::new("127.0.0.1", 6379, 0);
    s.set(k, v);
    assert_eq!(v, s.get(k));
}
