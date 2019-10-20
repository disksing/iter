package iter

func AllOf(first, last Iter, p func(Iter)bool) bool { 
	return FindIfNot(first, last, p) == last
}

func AnyOf(first, last Iter, p func(Iter)bool) bool {
	return FindIf(first, last, p) != last
}

func NoneOf(first, last Iter, p func(Iter)bool) bool { 
	return FindIf(first, last, p) == last
}

func ForEach(first, last Iter, f func(Iter)) func(Iter) {
	for ; first != last; first=first.Next() {
		f(first)
	}
	return f
}

func ForEachN(first Iter, n int, f func(Iter)) func(Iter) {
	for i := 0; i < n; i++ {
		f(first)
		first = first.Next()
	}
	return f
}

func Count(first, last Iter) int {
	var ret int
	for ; first != last; first = first.Next() {
		ret++
	}
	return ret
}

func CountIf(first, last Iter, p func(Iter)bool) int {
	var ret int
	for ; first != last; first=first.Next() {
		if p(first) {
			ret++
		}
	}
	return ret
}

func MisMatch(first1, last1, first2, last2 Iter) (Iter, Iter) {
	return MismatchIf(first1, last1, first2, last2, eq2)
}

func MismatchIf(first1, last1, first2, last2 Iter, eq func(Iter, Iter) bool) (Iter, Iter) {
	for first1 != last1 && first2 != last2 && eq(first1, first2) {
		first1, first2 = first1.Next(), first2.Next()
	}
	return first1, first2
}

func Find(first, last Iter, v interface{}) Iter {
	return FindIf(first, last, eqv(v))
}

func FindIf(first, last Iter, p func(Iter)bool) Iter {
	for ; first != last; first = first.Next() {
		if p(first) {
			return first
		}
	}
	return last
}

func FindIfNot(first, last Iter, p func(Iter)bool) Iter {
	for ; first != last; first = first.Next() {
		if !p(first) {
			return first
		}
	}
	return last
}

func FindEnd(first, last, sFirst, sLast) Iter {
	return FindEndIf(first, last, sFirst, sLast, eq2)
}

func FindEndIf(first, last, sFirst, sLast Iter, eq func(Iter,Iter)bool) Iter {
	if sFirst == sLast {
		return last
	}
	result := last
	for {
		if newResult := SearchIf(first, last, sFirst, sLast, eq); newResult == last {
			break
		} else {
			result = newResult
			first = result.Next()
		}
	}
	return result
}

func Search(first, last, sFirst, sLast Iter) Iter {
	return SearchIf(first, last, sFirst, sLast, eq2)
}

func SearchIf(first, last, sFirst, sLast Iter, eq func(Iter,Iter)bool) Iter {
	for {
		it := first
		for sIt := sFirst;;sIt,it=sIt.Next(),it.Next() {
			if sIt == sLast {
				return first
			}
			if it == last {
				return last
			}
			if !eq(it,sIt) {
				break
			}
		}
		first = first.Next()
	}
}

func SearchN(first, last Iter, count int, v interface{}) Iter {
	return SearchNIf(first, last, count,eqv(v))
}

func SearchNIf(first, last Iter, count int, p func(Iter)bool) Iter {
	if count <= 0 {
		return first
	}
	for ; first != last; first = first.Next() {
		if !p(first) {
			continue
		}
		candidate := first
		var curCount int
		for {
			curCount++
			if curCount >= count {
				return candidate
			}
			if first = first.Next(); first == last {
				return last
			}
			if !p(first) {
				break
			}
		}
	}
	return last
}

func eqv(v interface{}) func(Iter)bool {
	return func(it Iter)bool {return it.Get() == v}
}

func eq2(it1, it2 Iter) bool {
	return it1.Get() == it2.Get()
}