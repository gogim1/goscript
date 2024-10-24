letrec (
  empty = lambda () {
    lambda () { 0 }
  }
  cons = lambda (head tail) {
    lambda () { 1 } 
  }
  null = lambda (list) {
      if (not (list)) then 1
      else 0
  }
  head = lambda (list) {
      if (null list) then (void)
      else &head list
  }
  tail = lambda (list) {
      if (not (list)) then (void)
      else &tail list
  }
  last = lambda (list) {[
    if (null list) then (void) 
    else if (null (tail list)) then (head list) else (last (tail list))
  ]}
  init = lambda (list) {
    if (null list) then (void) 
    else if (null (tail list)) then (empty) else (cons (head list) (init (tail list)))
  }
  get = lambda (idx list) {
    if (or (null list) (lt idx 0)) then (void)
    else if (eq idx 0) then (head list) 
    else (get (sub idx 1) (tail list))
  }
  set = lambda (idx value list) {
    if (or (null list) (lt idx 0)) then (void) 
    else if (eq idx 0) then (cons value (tail list)) 
    else (cons (head list) (set (sub idx 1) value (tail list)))
  }
  sum = lambda (list) {
    if (null list) then 0
    else (add (head list) (sum (tail list)))
  }
  maximun = lambda (list) {
    if (null list) then (void)
    else if (null (tail list)) then (head list)
    else letrec (maxtail = (maximun (tail list))) {
      if (gt (head list) maxtail) then (head list)
      else maxtail
    } 
  }
  length = lambda (list) {
    if (null list) then 0 
    else (add (length (tail list)) 1)
  }
  replicate = lambda (size value) {
    if (eq size 0) then (empty) else (cons value (newlist (sub size 1)))
  }
  show = lambda (list) {
      if (null list) then (void)
      else [
          (put (head list) ", ")
          (show (tail list))
      ]
  }
  newlist = lambda (size) {(replicate size 0)}
) {[
  (reg "empty" empty)
  (reg "cons" cons)
  (reg "null" null)
  (reg "head" head)
  (reg "tail" tail)
  (reg "last" last)
  (reg "init" init)
  (reg "get" get)
  (reg "set" set)
  (reg "sum" sum)
  (reg "maximun" maximun)
  (reg "length" length)
  (reg "replicate" replicate)
  (reg "show" show)
  (reg "newlist" newlist)
]}