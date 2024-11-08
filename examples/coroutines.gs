letrec (
  getcc = lambda () {
    (callcc lambda (k) { (k k) })
  }
  task = lambda (yield) {
    [
    letrec (c = (getcc)) {
      if (iscont c) then (yield c) # jump to main
      else (void)
    }
    (put "task 1\n")
    letrec (c = (getcc)) {
      if (iscont c) then (yield c)
      else (void)
    }
    (put "task 2\n")
    letrec (c = (getcc)) {
      if (iscont c) then (yield c)
      else (void)
    }
    (put "task 3\n")
    ]
  }
  ) {
  letrec (
    c = (callcc task)
  ) {
    if (iscont c) then [
    (put "main\n")
    (c (void)) # jump to task
    ]
    else (void)
  }
}