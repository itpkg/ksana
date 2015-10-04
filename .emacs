(set-default-font "WenQuanYi Zen Hei 12")

(custom-set-variables
   '(initial-frame-alist (quote ((fullscreen . maximized)))))


;(setf inhibit-splash-screen t)
;(switch-to-buffer (get-buffer-create "emtpy"))
;(delete-other-windows)

(tool-bar-mode -1)
(menu-bar-mode -1)
(setq inhibit-startup-message t)

(setq-default cursor-type 'bar)
(blink-cursor-mode -1)
(column-number-mode t)

(defalias 'yes-or-no-p 'y-or-n-p)

;(setq make-backup-files nil)

;git clone git://jblevins.org/git/markdown-mode.git ~/.emacs.d/markdown-mode
(add-to-list 'load-path "~/.emacs.d/markdown-mode")
(autoload 'markdown-mode "markdown-mode"
   "Major mode for editing Markdown files" t)
(add-to-list 'auto-mode-alist '("\\.md\\'" . markdown-mode))


;git clone https://github.com/sellout/emacs-color-theme-solarized.git ~/.emacs.d/emacs-color-theme-solarized
(add-to-list 'custom-theme-load-path "~/.emacs.d/emacs-color-theme-solarized")
(load-theme 'solarized t)

;git clone https://github.com/rust-lang/rust-mode.git ~/.emacs.d/rust-mode
(add-to-list 'load-path "~/.emacs.d/rust-mode/")
(autoload 'rust-mode "rust-mode" nil t)
(add-to-list 'auto-mode-alist '("\\.rs\\'" . rust-mode))

