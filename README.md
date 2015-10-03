ksana(a web framwork for rust)
---

## Usage

    cargo new demo --bin
    cd demo
	cat >> Cargo.toml <<EOF
	
	[dependencies]
	ksana = { git = "https://github.com/itpkg/ksana.git" }

    EOF

## Editor

### vim
 * see [https://github.com/rust-lang/rust.vim](https://github.com/VundleVim/Vundle.vim)

### emacs

#### notes
 * open directory: C-x d
 * close buffer: C-x k
 * last buffer: C-x b
 * all buffers: C-x C-b
 * exit: C-x C-c
 * save: C-x C-s
 * Copy: M-w
 * Cut: C-w
 * Parse: C-y
 * Debug mode: emacs --debug-init

#### rust-mode

    wget https://raw.githubusercontent.com/itpkg/ksana/master/.emacs ~/.emacs
