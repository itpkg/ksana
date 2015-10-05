extern crate rand;
extern crate rustc_serialize;
extern crate uuid as _uuid;
extern crate crypto;

use std::borrow::BorrowMut;
use std::iter::Extend;
use self::rand::Rng;
use self::rustc_serialize::json;
use self::rustc_serialize::base64::{ToBase64, FromBase64, STANDARD, FromBase64Error};
use self::rustc_serialize::{Encodable, Decodable};
use self::crypto::digest::Digest;
use self::crypto::sha1::Sha1;
use self::crypto::sha2::Sha256;
use self::crypto::sha2::Sha512;

fn ssha(s: &str, l: usize, dig: &mut Digest) -> String {
    
    let salt = random_string(l);
    let mut pwd = s.to_string();
    pwd.push_str(&salt);

    dig.input_str(&pwd);

    let mut val = vec![0u8; dig.block_size()/2];
    let mut ssh = val.borrow_mut();
    dig.result(ssh);

    let mut res = ssh.to_vec();
    res.extend(salt.as_bytes().into_iter());

    return to_base64(res.borrow_mut());
}

pub fn ssha512(s: &str, l: usize) -> String {
    let mut dig = Sha512::new();
    ssha(s, l, &mut dig)
}

pub fn ssha256(s: &str, l: usize) -> String {
    let mut dig = Sha256::new();
    ssha(s, l, &mut dig)
}

pub fn sha1(s: &str) -> String {
    let mut sha = Sha1::new();
    sha.input_str(s);
    sha.result_str()
}

pub fn sha256(s: &str) -> String {
    let mut sha = Sha256::new();
    sha.input_str(s);
    sha.result_str()
}

pub fn sha512(s: &str) -> String {
    let mut sha = Sha512::new();
    sha.input_str(s);
    sha.result_str()
}

#[test]
fn test_sha() {    
    let hello = "hello";
    let  s512 = sha512(hello);
    assert!(!s512.is_empty());
    let  s256 = sha256(hello);
    assert!(!s256.is_empty());
    let  s1 = sha1(hello);
    assert!(!s1.is_empty());
    println!("sha512 => {}\nsha256 => {}\nsha1 => {}", s512, s256, s1);

    let ss512 = ssha512(hello, 6);
    assert!(!ss512.is_empty());
    println!("doveadm pw -t {}{} -p {}", "{SSHA512}", ss512, hello);
    let ss256 = ssha256(hello, 6);
    assert!(!ss256.is_empty());
    println!("doveadm pw -t {}{} -p {}","{SSHA256}", ss256, hello);
        
}

pub fn to_json<T: Encodable>(o: &T) -> String {
    json::encode(o).unwrap()
}

pub fn from_json<T: Decodable>(s: &str) -> T {
    json::decode(s).unwrap()
}

#[test]
fn test_json(){
    let val1 :f32 = 1.1;
    let js = to_json(&val1);
    assert!(!js.is_empty());
    println!("To json: {}", js);
    let val2 :f32 = from_json(&js);
    assert_eq!(val1, val2);
    println!("From json: {} == {}", val1, val2);

}

pub fn to_base64(o: &[u8]) -> String {
    o.to_base64(STANDARD)
}

pub fn from_base64(s: &String) -> Result<Vec<u8>, FromBase64Error> {
    s.from_base64()
}

#[test]
fn test_base64(){
    let val1  = "hello";
    let s = to_base64(val1.to_string().as_bytes());
    println!("To Base64: {}", s);
    assert!(!s.is_empty());
    match from_base64(&s) {
        Ok(_) => {1},
        Err(_) => { assert!(false); 0}
    };

}

pub fn random_string(len: usize) -> String {
    rand::thread_rng().gen_ascii_chars().take(len).collect()
}

pub fn random_bytes(dest: &mut [u8]) {   
   rand::thread_rng().fill_bytes(dest);   
}

#[test]
fn test_random() {
    let size :usize = 16;
    let str = random_string(size);
    assert!(!str.is_empty());
    assert_eq!(str.len(), size);
    println!("Random string: {}", str);

    let mut bys = [0u8; 32];
    random_bytes(&mut bys);    
    println!("Random bytes: {:?}", &bys[..]);
}

pub fn uuid() -> String {
    _uuid::Uuid::new_v4().to_hyphenated_string()
}

#[test]
fn test_uuid() {
    let uid = uuid();
    assert!(!uid.is_empty());
    assert_eq!(uid.len(), 36);
    println!("UUID: {}", uid);
}





