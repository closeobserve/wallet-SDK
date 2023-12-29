package base

import (
	"errors"
	"math/big"
	"strconv"
	"strings"
	"sync"
)

// This method will traverse the array concurrently and map each object in the array.
// @param list: [TYPE1], a list that all item is TYPE1
// @param limit: maximum number of tasks to execute, 0 means no limit
// @param maper: func(TYPE1) (TYPE2, error), a function that input TYPE1, return TYPE2
//
//	you can throw an error to finish task.
//
// @return : [TYPE2], a list that all item is TYPE2
// @example : ```
//
//	nums := []interface{}{1, 2, 3, 4, 5, 6}
//	res, _ := MapListConcurrent(nums, func(i interface{}) (interface{}, error) {
//	    return strconv.Itoa(i.(int) * 100), nil
//	})
//	println(res) // ["100" "200" "300" "400" "500" "600"]
//
// ```
func MapListConcurrent(list []interface{}, limit int, maper func(interface{}) (interface{}, error)) ([]interface{}, error) {
	thread := 0
	max := limit
	wg := sync.WaitGroup{}

	mapContainer := newSafeMap()
	var firstError error
	for _, item := range list {
		if firstError != nil {
			continue
		}
		if max == 0 {
			wg.Add(1)
			// no limit
		} else {
			if thread == max {
				wg.Wait()
				thread = 0
			}
			if thread < max {
				wg.Add(1)
			}
		}

		go func(w *sync.WaitGroup, item interface{}, mapContainer *safeMap, firstError *error) {
			maped, err := maper(item)
			if *firstError == nil && err != nil {
				*firstError = err
			} else {
				mapContainer.writeMap(item, maped)
			}
			wg.Done()
		}(&wg, item, mapContainer, &firstError)
		thread++
	}
	wg.Wait()
	if firstError != nil {
		return nil, firstError
	}

	result := []interface{}{}
	for _, item := range list {
		result = append(result, mapContainer.Map[item])
	}
	return result, nil
}

// The encapsulation of MapListConcurrent.
func MapListConcurrentStringToString(strList []string, maper func(string) (string, error)) ([]string, error) {
	list := make([]interface{}, len(strList))
	for i, s := range strList {
		list[i] = s
	}
	temp, err := MapListConcurrent(list, 10, func(i interface{}) (interface{}, error) {
		return maper(i.(string))
	})
	if err != nil {
		return nil, err
	}

	result := make([]string, len(temp))
	for i, v := range temp {
		result[i] = v.(string)
	}
	return result, nil
}

// Return the more biger of the two numbers
func MaxBigInt(x, y *big.Int) *big.Int {
	if x.Cmp(y) > 0 {
		return x
	} else {
		return y
	}
}

// @note float64 should use `math.Max()`
func Max[T int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | string](x, y T) T {
	if x >= y {
		return x
	} else {
		return y
	}
}

// @note float64 should use `math.Min()`
func Min[T int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | string](x, y T) T {
	if x <= y {
		return x
	} else {
		return y
	}
}

/* [zh] 该方法会捕捉 panic 抛出的值，并转成一个 error 对象通过参数指针返回
 *      注意: 如果想要返回它抓住的 error, 必须使用命名返回值！！
 * [en] This method will catch the value thrown by panic, and turn it into an error object and return it through the parameter pointer
 *		Note: If you want to return the error it caught, you must use a named return value! !
 *  ```
 *  func actionWillThrowError(parameters...) (namedErr error, other...) {
 *      defer CatchPanicAndMapToBasicError(&namedErr)
 *      // action code ...
 *      return namedErr, other...
 *  }
 *  ```
 */
func CatchPanicAndMapToBasicError(errOfResult *error) {
	// first we have to recover()
	errOfPanic := recover()
	if errOfResult == nil {
		return
	}
	if errOfPanic != nil {
		*errOfResult = MapAnyToBasicError(errOfPanic)
	} else {
		*errOfResult = MapAnyToBasicError(*errOfResult)
	}
}

func MapAnyToBasicError(e any) error {
	if e == nil {
		return nil
	}

	err, ok := e.(error)
	if ok {
		return errors.New(err.Error())
	}

	msg, ok := e.(string)
	if ok {
		return errors.New("panic error: " + msg)
	}

	code, ok := e.(int)
	if ok {
		return errors.New("panic error: code = " + strconv.Itoa(code))
	}

	return errors.New("panic error: unexpected error.")
}

// ParseNumber
// @param num any format number, such as decimal "1237890123", hex "0x123ef0", "123ef0"
func ParseNumber(num string) (*big.Int, error) {
	if strings.HasPrefix(num, "0x") || strings.HasPrefix(num, "0X") {
		num = num[2:]
		if b, ok := big.NewInt(0).SetString(num, 16); ok {
			return b, nil
		}
	}
	if b, ok := big.NewInt(0).SetString(num, 10); ok {
		return b, nil
	}
	if b, ok := big.NewInt(0).SetString(num, 16); ok {
		return b, nil
	}
	return nil, errors.New("invalid number")
}

// ParseNumberToHex
// @param num any format number, such as decimal "1237890123", hex "0x123ef0", "123ef0"
// @return hex number start with 0x, characters include 0-9 a-f
func ParseNumberToHex(num string) string {
	if b, err := ParseNumber(num); err == nil {
		return "0x" + b.Text(16)
	}
	return "0x0"
}

// ParseNumberToDecimal
// @param num any format number, such as decimal "1237890123", hex "0x123ef0", "123ef0"
// @return decimal number, characters include 0-9
func ParseNumberToDecimal(num string) string {
	if b, err := ParseNumber(num); err == nil {
		return b.Text(10)
	}
	return "0"
}

func BigIntMultiply(b *big.Int, ratio float64) *big.Int {
	f1 := new(big.Float).SetInt(b)
	product := f1.Mul(f1, big.NewFloat(ratio))
	res, _ := product.Int(big.NewInt(0))
	return res
}
