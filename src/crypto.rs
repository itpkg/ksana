extern crate rand;
extern crate crypto;
extern crate rustc_serialize;

use std::cmp::PartialEq;
use std::borrow::BorrowMut;
use std::iter::Extend;
use self::crypto::{sha1, sha2, md5, digest, hmac};
use self::crypto::mac::{Mac};
use self::rand::Rng;
use self::rustc_serialize::base64::{ToBase64, FromBase64, STANDARD};
use self::rustc_serialize::hex::ToHex;

pub trait Encryptor {

    fn en_aes(&self, k: &str) -> String;
    fn de_aes(&self, k: &str) -> String;

    fn sum_hmac_md5(&self, k: &str) -> String;
    fn chk_hmac_md5(&self, s: &str, k: &str) -> bool;
          
    fn sum_sha1(&self) -> String;
    fn chk_sha1(&self, s: &str) -> bool;
    fn sum_sha256(&self) -> String;
    fn chk_sha256(&self, s: &str) -> bool;
    fn sum_sha512(&self) -> String;
    fn chk_sha512(&self, s: &str) -> bool;
    fn sum_md5(&self) -> String;
    fn chk_md5(&self, s: &str) -> bool;
    
    fn sum_ssha512(&self, l: usize) -> String;
    fn chk_ssha512(&self, s: &str, l: usize) -> bool;
    fn sum_ssha256(&self, l: usize) -> String;
    fn chk_ssha256(&self, s: &str, l: usize) -> bool;

    fn sum<T: digest::Digest>(&self, dig: &mut T) -> String;
    fn ssum<T: digest::Digest>(&self, dig: &mut T, l: usize) -> String;
    fn schk<T: digest::Digest>(&self, dig: &mut T, s: &str, l: usize) -> bool;

    fn sum_hmac<T: digest::Digest>(&self, dig:  T, k: &str) -> String;
  
}

impl Encryptor for String {

    fn en_aes(&self, k: &str) -> String {
        k.to_string()
    }
    fn de_aes(&self, k: &str) -> String {
        k.to_string()
    }
    
    fn sum_hmac_md5(&self,  k: &str) -> String {
        self.sum_hmac(md5::Md5::new(), k)            
    }
    fn chk_hmac_md5(&self, s:&str,  k: &str) -> bool {
        let d = s.to_string().sum_hmac_md5(k);
        self.eq(&d)        
    }
    
    fn sum_hmac<T: digest::Digest>(&self, dig:  T, k: &str) -> String {
        let mut hm = hmac::Hmac::new(dig, k.as_bytes());
        hm.input(&self.as_bytes());
        hm.result().code().to_hex()
    }

        
    fn sum_sha1(&self) -> String {        
        let mut dig = sha1::Sha1::new();        
        self.sum(&mut dig)        
    }
    fn chk_sha1(&self, s: &str) -> bool {
        let d = s.to_string().sum_sha1();
        self.eq(&d)
    }
    fn sum_sha256(&self) -> String {
        let mut dig = sha2::Sha256::new();
        self.sum(&mut dig)        
    }
    fn chk_sha256(&self, s: &str) -> bool {
        let d = s.to_string().sum_sha256();
        self.eq(&d)
    }
    fn sum_sha512(&self) -> String {
        let mut dig = sha2::Sha512::new();       
        self.sum(&mut dig)        
    }
    fn chk_sha512(&self, s: &str) -> bool {
        let d = s.to_string().sum_sha512();
        self.eq(&d)
    }
    fn sum_md5(&self) -> String {
        let mut dig = md5::Md5::new();        
        self.sum(&mut dig)        
    }
    fn chk_md5(&self, s: &str) -> bool {
        let d = s.to_string().sum_md5();
        self.eq(&d)
    }

    fn sum_ssha256(&self, l: usize) -> String {
        let mut dig = sha2::Sha256::new();
        self.ssum(&mut dig, l)        
    }
    fn chk_ssha256(&self, s: &str, l: usize) -> bool{
        let mut dig = sha2::Sha256::new();       
        self.schk(&mut dig, s, l)
    }
    fn sum_ssha512(&self, l: usize) -> String {
        let mut dig = sha2::Sha512::new();       
        self.ssum(&mut dig, l)        
    }
    fn chk_ssha512(&self, s: &str, l: usize) -> bool{
        let mut dig = sha2::Sha512::new();       
        self.schk(&mut dig, s, l)
    }
   
    fn sum<T: digest::Digest>(&self, dig: &mut T) -> String {
        dig.input_str(&self);
        dig.result_str()
    }
    fn ssum<T: digest::Digest>(&self, dig: &mut T, l: usize) -> String {
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
    fn schk<T: digest::Digest>(&self, dig: &mut T, s: &str, l: usize) -> bool {
        match self.from_base64() {
            Ok(buf) => {
                let (sha, salt) = buf.split_at(buf.len()-l);                
                let mut sha_v = vec![0u8; dig.block_size()/2];
                let mut sha_bs = sha_v.borrow_mut();
                
                dig.input_str(s);
                dig.input(salt);
                dig.result(sha_bs);
                
                sha == sha_bs
            },
            Err(_) => {
                false
            },
        }
    }
    
}

#[test]
fn test_encryptor() {
    let hello = "hello".to_string();
    let world = "world";
    let size = 6;
    let key = "123456";
    
    let s1 = hello.sum_sha1();
    assert!(s1.chk_sha1(&hello));
    assert!(!s1.chk_sha1(world));
    println!("sha1: {}", s1);
    
    let s256 = hello.sum_sha256();
    assert!(s256.chk_sha256(&hello));
    assert!(!s256.chk_sha256(world));
    println!("sha256: {}", s256);
    
    let s512 = hello.sum_sha512();
    assert!(s512.chk_sha512(&hello));
    assert!(!s512.chk_sha512(world));
    println!("sha512: {}", s512);
    
    let smd5 = hello.sum_md5();
    assert!(smd5.chk_md5(&hello));
    assert!(!smd5.chk_md5(world));
    println!("md5: {}", smd5);

    
    let ss256 = hello.sum_ssha256(size);
    assert!(ss256.chk_ssha256(&hello, size));
    assert!(!ss256.chk_ssha256(world, size));
    println!("run: doveadm pw -t {}{} -p {}", "{SSHA256}", ss256, hello);
    
    let ss512 = hello.sum_ssha512(size);    
    assert!(ss512.chk_ssha512(&hello, size));
    assert!(!ss512.chk_ssha512(world, size));
    println!("run: doveadm pw -t {}{} -p {}", "{SSHA512}", ss512, hello);

    let hmd5 = hello.sum_hmac_md5(key);    
    assert!(hmd5.chk_hmac_md5(&hello, key));
    assert!(!hmd5.chk_hmac_md5(world, key));
    println!("hmac md5: {}", hmd5);
    
}
