letrec (
  f = lambda (x) {
    if (lt x 0) then (void)
    else [(put x "\n") (g (sub x 1))]
  }
  g = lambda (x) {
    if (lt x 0) then (void)
    else [(put x "\n") (f (sub x 1))]
  }
) {
  (f 5)
}