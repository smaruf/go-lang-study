package main

// "GCD struct for gcd generation"
type GCD struct {
  first int64
  second int64
}
// "Calculate calculate gdc"
func (numbers *GCD) Calculate() int64 {
    divisor := numbers.getMin()
    divident := numbers.getMax()
    if divident % divisor == 0 {
        return divisor
	};
	
	modulus := divident % divisor
	newNumbers := GCD{modulus, divisor}
    return newNumbers.Calculate()
}

func (numbers *GCD) getMin() int64 {
	if numbers.first < numbers.second {
		return numbers.first
	};
	return numbers.second
}

func (numbers *GCD) getMax() int64 {
    if numbers.first > numbers.second {
		return numbers.first
	};

	return numbers.second
}
