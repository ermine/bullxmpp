package jid

import (
	"errors"
)

func get_bst(arr [][]rune, x rune, s int, e int) (ret rune) {
  if s > e || x < arr[s][0] || x > arr[e][1] {
		ret = -1
	} else {
		mid := s + (e-s)/2
		if x < arr[mid][0] {
			ret = get_bst(arr, x, s, mid-1)
		} else if x > arr[mid][1] {
			ret = get_bst(arr, x, mid+1, e)
		} else {
			ret = arr[mid][2]
		}
	}
	return
}

func get_casemap(x rune) rune {
  if x > 0x40 && x < 0x5b {
    return casemap[0][2]
	}
  return get_bst(casemap, x, 1, casemap_max_idx)
}

func get_b1_b2(x rune) (ret []rune) {
  v := get_casemap(x)
	switch {
	case v == -1: 
		ret = make([]rune, 1)
		ret[0] = x
  case v == 0: 
		ret = nil
  case v & b_mc == b_mc: 
		ret = fmap[v >> b_shift]
	default: 
		ret = make([]rune, 1)
		ret[0] = x + (v >> b_shift)
	}
	return
}

func get_b1(x rune) rune {
	v := get_casemap(x)
  if v == 0 {
		return -1
  }
  return x
}          

// http://unicode.org/reports/tr15/
const (
	hangulSBase = 0xAC00
	hangulLBase = 0x1100
	hangulVBase = 0x1161
	hangulTBase = 0x11A7        
	hangulLCount = 19
	hangulVCount = 21
	hangulTCount = 28
	hangulNCount = hangulVCount * hangulTCount // 588
	hangulSCount = hangulLCount * hangulNCount // 11172
)

type decom_with_cclass struct {
	rune
	int
}

func get_decomp(x rune) (ret []decomp_data) {
  if x > 0x009F && x < 0x2FA1E {
		s, e := 0, dmap_max_idx
		for {
			if s > e {
				ret = make([]decomp_data, 1)
				ret[0].r = x
				ret[0].cc = 0
				return
			}
			mid := s + (e - s) / 2
			var ddata data = dmap[mid]
			switch {
			case x < ddata.first:
				e = mid-1
				continue
			case x > rune(int(ddata.first) + len(ddata.arr) - 1):
				s = mid+1;
				continue
			}
			data := ddata.arr[x - ddata.first]
			i := 0
			ret = make([]decomp_data, len(data))
			for d := range data {
				ret[i].r = rune(d >> 8)
				ret[i].cc = d & 0xFF
				i++
			}
			return
		}
	}
	ret = make([]decomp_data, 1)
	ret[0].r = x
	ret[0].cc = 0
	return
}

func compatibility_decompose(x rune) (ret []decomp_data) {
  sindex := x - hangulSBase
  if sindex < 0 || sindex >= hangulSCount {
    ret = get_decomp(x)
  } else {
    l := hangulLBase + sindex / hangulNCount
    v := hangulVBase + (sindex % hangulNCount) / hangulTCount
    t := hangulTBase + sindex % hangulTCount
		if t != hangulTBase {
			ret = make([]decomp_data, 2)
			ret[2].r = t
			ret[2].cc = 0
		} else {
			ret = make([]decomp_data, 3)
		}
		ret[0].r = l
		ret[0].cc = 0
		ret[1].r = v
		ret[1].cc = 0
	}
	return
}

func canonical_order(d []decomp_data) (ret []decomp_data) {
	if d == nil {
		return
	}
	var prev decomp_data = d[0]
	ret = make([]decomp_data, len(d))
	j := 0

	for i := 1; i < len(d); {
		if d[i].cc == 0 || prev.cc <= d[i].cc {
			ret[j] = prev
			prev = d[i]
			i++
			j++
		} else if j == 0 {
			ret[j] = d[i]
			j++
			i++
		} else {
			j--
			prev = ret[j]
			var tmp []decomp_data
			tmp = append(tmp, d[0:j-1]...)
			tmp = append(tmp, prev)
			tmp = append(tmp, d[i:]...)
			d = tmp
		}
	}
	ret[j] = prev
	return
}

func decompose(rs []rune) []decomp_data {
	var decomps []decomp_data

	for _, x := range rs {
    var ds []decomp_data = compatibility_decompose(x)
		if ds != nil {
			decomps = append(decomps, ds[0:]...)
		}
	}
  return canonical_order(decomps)
}

func compose_hangul(ch1 rune, ch2 rune) rune {
  // check if two current characters are L and V
  lindex := int(ch1) - hangulLBase
  vindex := int(ch2) - hangulVBase
  if (lindex >= 0 && lindex < hangulLCount) &&
    (vindex >= 0 && vindex < hangulVBase) {
    // make syllable of form LV
    return rune(hangulSBase + (lindex * hangulVCount + vindex) * hangulTCount)
  } else {
    // 2. check to see if two current characters are LV and T
    sindex := ch1 - hangulSBase
    tindex := ch2 - hangulTBase
    if (sindex >= 0 && sindex < hangulSCount &&
      (sindex % hangulTCount) == 0) &&
      (tindex > 0 && tindex < hangulTCount) {
      // make syllable of form LVT
      return rune(ch1 + tindex)
    }
	}
  // if neither case was true
  return -1
}

func get_comp_branch(branch [][2]rune, ch1 rune, s int, e int) (ret rune) {
	if s > e {
		ret = -1
	} else {
		mid := s + (e-s) / 2
		var amid [2]rune = branch[mid]
		switch {
		case amid[0] == ch1: ret = amid[1]
		case amid[0] > ch1: ret = get_comp_branch(branch, ch1, s, mid-1)
		default: ret = get_comp_branch(branch, ch1, mid+1, e)
		}
	}
	return
}

func get_comp(ch1 rune, ch2 rune, s int, e int) (ret rune) {
	if s > e {
		ret = -1
	} else {
		mid := s + (e - s) / 2
		var amid comp_data = comp_map[mid]
		if amid.ch2 == ch2 {
			return get_comp_branch(amid.arr, ch1, 0, len(amid.arr)-1)
		} else if ch2 < amid.ch2 {
			ret = get_comp(ch1, ch2, s, mid-1)
		} else {
			ret = get_comp(ch1, ch2, mid+1, e)
		}
	}
	return
}

func composeTwo(ch1 rune, ch2 rune) (newch rune) {
  newch = compose_hangul(ch1, ch2)
  if newch == -1 && (ch2 > 767 && ch2 < 12443) {
    newch = get_comp(ch1, ch2, 0, comps_max_idx)
	}
	return
}

func compose(rs []decomp_data) (ret []rune) {
	if rs == nil {
		return
	}
	var prev decomp_data = rs[0]
	var comps []rune
	for _, curr := range rs[1:] {
		var newch rune
		if prev.cc == 0 || curr.cc > prev.cc {
			newch = composeTwo(prev.r, curr.r)
		} else {
			newch = -1
		}
		if newch > -1 {
			prev.r = newch
    } else if curr.cc == 0 {
			ret = append(ret, prev.r)
			ret = append(ret, comps[0:]...)
			comps = nil
			prev.r = curr.r
			prev.cc = 0
		} else {
			comps = append(comps, curr.r)
			prev.cc = curr.cc
		}
	}
	ret = append(ret, prev.r)
	ret = append(ret, comps[0:]...)
	return
}

func nfkc(rs []rune) (ret []rune) {
	return compose(decompose(rs))
}

func nodeprep(rs []rune) ([]rune) {
	var res []rune
	for _, x := range rs {
    cs := get_b1_b2(x)
		if cs != nil {
			res = append(res, cs[0:]...)
		}
	}
  return nfkc(res)
}        

func resourceprep(rs []rune) ([]rune) {
	var res []rune
	for _, x := range rs {
		c := get_b1(x)
		if c != -1 {
			res = append(res, x)
		}
	}
  return nfkc(res)
}
  
// Check prohibited symbols and bidi
func check_prohibits(p int, rs []rune) (err error) {
	if rs == nil {
		return
	}
	var dir int
  var v rune = get_bst(prohibits, rs[0], 0, prohibits_max_idx)
  if v == -1 {
		dir = b_l
	} else if int(v) & p == p {
    return errors.New("prohibited symbol")
	} else if int(v) & b_randal == b_randal {
		dir = b_randal
	} else {
		dir = b_l
	}
	last_dir := dir
	for _, x := range rs[1:] {
    v = get_bst(prohibits, x, 0, prohibits_max_idx)
    if v == -1 {
			last_dir = b_l
		} else if int(v) & p == p {
			return errors.New("prohibited symbol")
		} else if int(v) & dir != dir && 
			(v & b_randal == b_randal || v & b_l == b_l) {
			return errors.New("invalid bidi")
    }
	}
	if last_dir != dir {
		err = errors.New("invalid bidi")
	}
	return
}

func StrongNodeprep(rs []rune) (normalized []rune, err error) {
  normalized = nodeprep(rs)
  err = check_prohibits(b_nodeprep_prohibit, normalized)
  return
}

func StrongResourceprep(rs []rune) (normalized []rune, err error)  {
  normalized = resourceprep(rs)
  err = check_prohibits(b_resourceprep_prohibit, normalized)
  return
}
    
// for domains
func nameprep (rs []rune) ([]rune) {
	var res []rune
	for _, x := range rs {
    cs := get_b1_b2(x)
		if cs != nil {
			res = append(res, cs[0:]...)
		}
	}
  return nfkc(res)
}

func StrongNameprep(rs []rune) (normalized []rune, err error) {
	normalized = nameprep(rs)
  err = check_prohibits(b_nameprep_prohibit, normalized)
  return
}










