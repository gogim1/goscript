[
  letrec (
    s1 = "(put "
    s2 = "\""
    s3 = "EVAL\\n\")"
  ) {
    (eval (concat (concat s1 s2) s3))
  }
  (put (eval (eval (eval (quote (quote (quote "hello world\n")))))))
  (put (eval (eval (eval (quote (quote (quote "hello world
")))))))
]