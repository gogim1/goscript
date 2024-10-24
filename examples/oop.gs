letrec (
  makeBase = lambda () {
    lambda (call) {
      letrec (
        F = lambda () { 1 }
        G = lambda () { 2 }
      ) {
        (call)
      }
    }
  }
  makeDerived = lambda () {
    letrec (base = (makeBase)) { # reuse the base class (inheritance)
      lambda (call) {
        (base
          lambda () {
            letrec (F = lambda () { 100 } # override F
            ) {
              (call) 
            }
          }
        )
      }
    }
  }
) {
  letrec (
    base = (makeBase)
    derived = (makeDerived)
  ) {
    [
      (put (base lambda () { (F) }) "\n")
      (put (base lambda () { (G) }) "\n")
      (put (derived lambda () { (F) }) "\n")
      (put (derived lambda () { (G) }) "\n")
    ]
  }
}