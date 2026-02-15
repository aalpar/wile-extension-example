;; A statistics library written in pure Scheme.
;;
;; Demonstrates that Scheme libraries can use standard operations (+, -, *, /)
;; available from the engine's registry. Go-registered functions live in the
;; engine's environment and are composed with library exports at the call site.
;;
;; Edge case: empty lists return 0 for mean/variance (not an error).
(define-library (stats)
  (export mean variance describe)
  (begin
    ;; Helper: sum elements of a list (tail-recursive)
    (define (sum lst)
      (let loop ((lst lst) (acc 0))
        (if (null? lst)
            acc
            (loop (cdr lst) (+ acc (car lst))))))

    ;; Helper: length of a list (tail-recursive)
    (define (len lst)
      (let loop ((lst lst) (acc 0))
        (if (null? lst)
            acc
            (loop (cdr lst) (+ acc 1)))))

    ;; Mean: sum / count. Returns 0 for empty lists.
    (define (mean lst)
      (if (null? lst)
          0
          (/ (sum lst) (len lst))))

    ;; Variance: E[(x - mu)^2]. Returns 0 for empty lists.
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
