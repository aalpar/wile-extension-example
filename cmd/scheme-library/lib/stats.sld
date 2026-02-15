;; A statistics library written in pure Scheme.
;;
;; Demonstrates that Scheme libraries can use standard operations (+, -, *, /)
;; available from the engine's registry. Go-registered functions live in the
;; engine's environment and are composed with library exports at the call site.
(define-library (stats)
  (export mean variance describe)
  (begin
    ;; Helper: sum elements of a list
    (define (sum lst)
      (if (null? lst)
          0
          (+ (car lst) (sum (cdr lst)))))

    ;; Helper: length of a list
    (define (len lst)
      (if (null? lst)
          0
          (+ 1 (len (cdr lst)))))

    ;; Mean: sum / count
    (define (mean lst)
      (if (null? lst)
          0
          (/ (sum lst) (len lst))))

    ;; Variance: E[(x - mu)^2]
    (define (variance lst)
      (if (null? lst)
          0
          (let ((mu (mean lst))
                (n (len lst)))
            (/ (sum (map (lambda (x)
                           (let ((d (- x mu)))
                             (* d d)))
                         lst))
               n))))

    ;; Describe: returns a list of (name . value) pairs summarizing the data
    (define (describe lst)
      (list (cons 'count (len lst))
            (cons 'mean (mean lst))
            (cons 'variance (variance lst))))))
