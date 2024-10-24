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
        if (not (list)) then (void)
        else &head list
    }
    tail = lambda (list) {
        if (not (list)) then (void)
        else &tail list
    }
    show = lambda (list) {
        if (null list) then (void)
        else [
            (put (head list) "\n")
            (show (tail list))
        ]
    }
) {
    (show (cons 5 (cons 4 (cons 3 (cons 2 (cons 1 (empty)))))))
}