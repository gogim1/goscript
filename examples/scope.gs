letrec (
  f = letrec (lex = 1 Dyn = 101) {
    lambda () { [(put lex "\n") (put Dyn "\n")] }
  }
) {
  letrec (lex = 3 Dyn = 303) { (f) }
}
