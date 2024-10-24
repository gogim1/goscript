letrec (
  try = lambda (try catch finally) {
    letrec (
      cc = (callcc lambda (k) { (k k) })
    ) {
      if (iscont cc) then [(try cc) (finally)] 
      else [(catch cc) (finally)]
    }
  }
  body = lambda (v) {(
    try
    lambda (throw) {[
      (put "enter `try` block\n") 
      if (eq v 0) then (throw "message") else (void)
      (put "exit `try` block\n") 
    ]} 
    lambda (exception) {[ 
      (put "enter `catch` block\n" )
      (put "Exception: " exception "\n" )
      (put "exit `catch` block\n" )
    ]} 
    lambda () {[
      (put "enter `finally` block\n") 
      (put "exit `finally` block\n") 
    ]}
  )}
) {[
    (body 0)
    (body 1)
  ]
}