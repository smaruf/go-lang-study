type GCD struct {
  first int64
  second int64
}

func (numbers *GCD) calculateGcd() int64 {
    divisor := getMin(numbers)
    divident := getMax(numbers)
    if(divident % divisor == 0) {
        return divisor
    } else {
      modulus = divident % divisor
      return calculateGcd(GCD(first: modulus, second: divisor))
    }
}

private func (numbers *GCD) getMin() int64 {
    return numbers.first < numbers.second ? numbers.first : numbers.second
}

private func (numbers *GCD) getMax() int64 {
    return numbers.first > numbers.second ? numbers.first : numbers.second
}
