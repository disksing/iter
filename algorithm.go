package iter

// func Equal(first1, last1, first2, last2 ForwardIter) bool {
// 	return EqualIf(first1, last1, first2, last2, eqv2)
// }

// func EqualIf(first1, last1, first2, last2 ForwardIter, eq func(Iter, Iter) bool) bool {
// 	for ; ne(first1, last1); first1, first2 = first1.Next(), first2.Next() {
// 		if !eq(first1, first2) {
// 			return false
// 		}
// 	}
// 	return true
// }

// func IsPermutation(first1, last1, first2, last2 ForwardIter) bool {
// 	return IsPermutationIf(first1, last1, first2, last2, eqv2)
// }

// func IsPermutationIf(first1, last1, first2, last2 ForwardIter, eq func(Iter, Iter) bool) bool {
// 	for ; ne(first1, last1); first1, first2 = first1.Next(), first2.Next() {
// 		if !eq(first1, first2) {
// 			break
// 		}
// 	}
// 	if first1 == last1 {
// 		return true
// 	}

// 	l1 := Distance(first1, last1)
// 	if l1 == 1 {
// 		return false
// 	}
// 	last2 = AdvanceN(last2, l1).(ForwardIter)
// 	for i := first1; i != last1; i = i.Next() {
// 		match := first1
// 		for ; match != i; match = match.Next() {
// 			if eq(match, i) {
// 				break
// 			}
// 		}
// 		if match == i {
// 			var c2 int
// 			for j := first2; j != last2; j = j.Next() {
// 				if eq(i, j) {
// 					c2++
// 				}
// 			}
// 			if c2 == 0 {
// 				return false
// 			}
// 			c1 := 1
// 			for j := Advance(i).(ForwardIter); j != last1; j = j.Next() {
// 				if eq(i, j) {
// 					c1++
// 				}
// 			}
// 			if c1 != c2 {
// 				return false
// 			}
// 		}
// 	}
// 	return true
// }
