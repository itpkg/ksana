extern crate rand;
extern crate crypto;
extern crate rustc_serialize;

use std::borrow::BorrowMut;
use std::iter::Extend;
use self::crypto::{sha1, sha2, md5, digest};
use self::rand::Rng;
use self::rustc_serialize::base64::{ToBase64, STANDARD};

pub trait Encryptor {
    
    fn sha1(&self) -> String;
    fn sha256(&self) -> String;
    fn sha512(&self) -> String;
    fn md5(&self) -> String;
    
    fn ssha512(&self, l: usize) -> String;   
    fn ssha256(&self, l: usize) -> String;

    fn encrypt<T: digest::Digest>(&self, dig: &mut T) -> String;
    fn encrypt_with_slat<T: digest::Digest>(&self, dig: &mut T, l: usize) -> String;
}

impl Encryptor for String {
        
    fn sha1(&self) -> String {        
        let mut dig = sha1::Sha1::new();        
        self.encrypt(&mut dig)        
    }
    fn sha256(&self) -> String {
        let mut dig = sha2::Sha256::new();
        self.encrypt(&mut dig)        
    }
    fn sha512(&self) -> String {
        let mut dig = sha2::Sha512::new();       
        self.encrypt(&mut dig)        
    }
    fn md5(&self) -> String {
        let mut dig = md5::Md5::new();        
        self.encrypt(&mut dig)        
    }

    fn ssha256(&self, l: usize) -> String {
        let mut dig = sha2::Sha256::new();
        self.encrypt_with_slat(&mut dig, l)        
    }
    fn ssha512(&self, l: usize) -> String {
        let mut dig = sha2::Sha512::new();       
        self.encrypt_with_slat(&mut dig, l)        
    }
   
    fn encrypt<T: digest::Digest>(&self, dig: &mut T) -> String {
        dig.input_str(&self);
        dig.result_str()
    }
    fn encrypt_with_slat<T: digest::Digest>(&self, dig: &mut T, l: usize) -> String {
        let salt: String = rand::thread_rng().gen_ascii_chars().take(l).collect();
        let mut buf = "".to_string();
        buf.push_str(&self);
        buf.push_str(&salt);
        
        let mut sha_v = vec![0u8; dig.block_size()/2];
        let mut sha_bs = sha_v.borrow_mut();
        dig.input_str(&buf);
        dig.result(sha_bs);

        let mut res = sha_bs.to_vec();
        res.extend(salt.as_bytes().into_iter());
        res.to_base64(STANDARD)
    }

}

#[test]
fn test_encryptor() {
    let hello = "hello".to_string();
    let size = 6;
    
    let s1 = hello.sha1();
    assert!(!s1.is_empty());
    println!("sha1: {}", s1);
    
    let s256 = hello.sha256();
    assert!(!s256.is_empty());
    println!("sha256: {}", s256);
    
    let s512 = hello.sha512();
    assert!(!s512.is_empty());
    println!("sha512: {}", s512);
    
    let md5 = hello.md5();
    assert!(!md5.is_empty());
    println!("md5: {}", md5);

    
    let ss256 = hello.ssha256(size);
    assert!(!ss256.is_empty());
    println!("run: doveadm pw -t {}{} -p {}", "{SSHA256}", ss256, hello);
    
    let ss512 = hello.ssha512(size);
    assert!(!ss512.is_empty());
    println!("run: doveadm pw -t {}{} -p {}", "{SSHA512}", ss512, hello);
}
