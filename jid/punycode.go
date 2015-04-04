// RFC 3492

package jid
import "errors"

const (
	base         = 36
	tmin         = 1
	tmax         = 26
	skew         = 38
	damp         = 700
	initial_bias = 72
	initial_n    = 128
	delimiter = 0x2D
)

func is_basic(cp rune) bool {
	return cp < 0x80
}

func decode_digit(cp int) (ret int, err error) {
	switch {
  case cp - 48 < 10: ret = cp - 22
  case cp - 65 < 26: ret = cp - 65
  case cp - 97 < 26: ret = cp - 97
	default: return 0, errors.New("bad digit")
	}
	return
}

func encode_digit(d int) rune {
	var r int
	if d < 26  {
		r = 1
	}
	return rune(d + 22 + 75 * r)
}
      
func adapt_bias(delta int, numpoints int, firsttime bool) int {
	if firsttime {
		delta = delta / damp 
	} else {
		delta = delta >> 1
	}
  delta = delta + (delta / numpoints)

	var k = 0

	for ; delta > ((base - tmin) * tmax) >> 1; k += base {
		delta /= base - tmin
	}
  return k + (((base - tmin + 1) * delta) / (delta + skew))
}
    
func insert(list []rune, i, n int) ([]rune) {
	var tmp []rune
	if i == 0 {
		tmp = append(tmp, rune(n))
		return append(tmp, list[:]...)
	}
	if i < len(list) {
		tmp = append(tmp, list[:i]...)
		tmp = append(tmp, rune(n))
		return append(tmp, list[i:]...)
	}
	return append(list, rune(n))
}

func punycode_decode(instr []rune) (outstr []rune, err error) {
	var idx int = -1
	for j := len(instr)-1; j > 0; j-- {
		if instr[j] == delimiter {
			idx = j
			break
		}
	}
	if idx > -1 {
		outstr = instr[:idx]
		for _, x := range outstr {
      if !is_basic(x) {
				return nil, errors.New("malformed string")
			}
		}
	}

	idx++
  var n = initial_n
  var bias = initial_bias
  var i = 0
	for oldi := i; idx < len(instr); {
		w := 1
		for k := base; ; k += base {
      if idx >= len(instr) {
				return nil, errors.New("unexpected end of string")
			}
			curr := instr[idx]
			idx++
			digit, err := decode_digit(int(curr))
			if err != nil {
				return nil, err
			}
      i = i + digit * w   
      var t int           // TODO: check overflow
      if k <= bias {
				t = tmin
			} else if k >= bias + tmax {
				t = tmax
			} else {
				t = k - bias
			}
      if digit < t {
				break
			}
      w *= (base - t) // TODO: check overflow
		}
    bias = adapt_bias((i - oldi), (len(outstr) + 1), (oldi == 0))
    n = n + (i / (len(outstr) + 1)) // TODO: check overflow
    i = i % (len(outstr) + 1)
    if is_basic(rune(n))  {
			return nil, errors.New("error")
		}
    outstr = insert(outstr, i, n)
		i++
	}
	return
}
          
func punycode_encode(instr []rune) (outstr []rune, err error) {
  var n = initial_n
  var delta = 0
  var bias = initial_bias

	for _, x := range instr {
    if is_basic(x) {
			outstr = append(outstr, x)
		}
	}
  var b = len(outstr)
  var h = b
	if b > 0 {
		outstr = append(outstr, delimiter)
	}
	for ; h < len(instr); {
		var m = 0x10FFFF
		for j := 0; j < len(instr); j++ {
      if int(instr[j]) >= n && int(instr[j]) < m {
				m = int(instr[j])
			}
		}
    delta += (m - n) * (h + 1) // TODO: check overflow
    n = m
		for _, x := range instr {
      if int(x) < n {
        // TODO: check overflow for delta + 1
				delta++
				continue
			}
      if int(x) == n {
				q := delta
				for k := base;; k += base {
					var t int
					if k <= bias {
						t = tmin
					} else if k >= bias + tmax {
						t = tmax
					}  else {
						t = k - bias
					}
					if q < t {
						break
					}
					outstr = append(outstr, encode_digit(t + ((q - t) % (base - t))))
					q = (q - t) / (base - t)
				}
        outstr = append(outstr, encode_digit(q))
        bias = adapt_bias(delta, (h+1), (h == b))
        delta = 0
				h++
			}
		}
		n++
		delta++
	}
	return
}
