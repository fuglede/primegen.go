// The primegen package generates prime numbers in order.  It uses the Sieve of
// Atkin instead of the traditional Sieve of Eratosthenes.
//
// The primegen Go package is a port of D. J. Bernstein's original implementation
// in C (http://cr.yp.to/primegen.html).
package primegen

import (
	"log"
)

const (
	_PRIMEGEN_WORDS = 8192 // Assumes 32K L1 data cache
	_B32            = _PRIMEGEN_WORDS
	_B              = _PRIMEGEN_WORDS * 32
)

type Primegen struct {
	buf  [16][]uint32
	p    [512]uint64 /* p[num-1] ... p[0], in that order */
	num  int
	pos  int /* next entry to use in buf; _PRIMEGEN_WORDS to restart */
	base uint64
	_L   uint64
}

// New returns a new primegen.Primegen prime number generator.
func New() *Primegen {
	pg := new(Primegen)
	for i := 0; i < len(pg.buf); i++ {
		pg.buf[i] = make([]uint32, _PRIMEGEN_WORDS)
	}
	pg.Reset()
	return pg
}

// Reset resets the generator at the first prime number.
func (pg *Primegen) Reset() {
	pg._L = 1
	pg.base = 60

	pg.pos = _PRIMEGEN_WORDS

	pg.p[0] = 59
	pg.p[1] = 53
	pg.p[2] = 47
	pg.p[3] = 43
	pg.p[4] = 41
	pg.p[5] = 37
	pg.p[6] = 31
	pg.p[7] = 29
	pg.p[8] = 23
	pg.p[9] = 19
	pg.p[10] = 17
	pg.p[11] = 13
	pg.p[12] = 11
	pg.p[13] = 7
	pg.p[14] = 5
	pg.p[15] = 3
	pg.p[16] = 2

	pg.num = 17
}

// Peek returns the next prime number, without advancing the generator.
func (pg *Primegen) Peek() uint64 {
	for pg.num == 0 {
		pg.fill()
	}

	return pg.p[pg.num-1]
}

// Next returns the next prime number and advances the generator.
func (pg *Primegen) Next() uint64 {
	p := pg.Peek()
	pg.num--
	return p
}

// SkipTo advances the generator to generate prime numbers >= x.
func (pg *Primegen) SkipTo(x uint64) {
	for {
		for pg.num != 0 {
			if pg.p[pg.num-1] >= x {
				return
			}
			pg.num--
		}

		pos := pg.pos
		for (pos < _B32) && (pg.base+1920 < x) {
			pg.base += 1920
			pos++
		}
		pg.pos = pos
		if pos == _B32 {
			for pg.base+_B*60 < x {
				pg._L += _B
				pg.base += _B * 60
			}
		}

		pg.fill()
	}
}

const ndigits = 32

// CountClasses finds the number of consecutive prime numbers between low and high,
// whose mod 3 classes are (1, 1), (1, 2), (2, 1), and (2, 2).
func CountClasses(low, high uint64) {
	sieve := New()
	sieve.SkipTo(low)
	x := sieve.Next()
	cnew := x % 3
	var cold, c11, c12, c21, c22 uint64
	for x < high {
		cold = cnew
		cnew = x % 3
		if cold == 1 && cnew == 1 {
			c11 += 1
		} else if cold == 1 && cnew == 2 {
			c12 += 1
		} else if cold == 2 && cnew == 1 {
			c21 += 1
		} else if cold == 2 && cnew == 2 {
			c22 += 1
		}

		x = sieve.Next()
	}
	log.Println("In class (1, 1):")
	log.Println(c11)
	log.Println("In class (1, 2):")
	log.Println(c12)
	log.Println("In class (2, 1):")
	log.Println(c21)
	log.Println("In class (2, 2):")
	log.Println(c22-1)
}
