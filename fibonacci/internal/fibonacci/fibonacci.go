package fibonacci

import (
	"fmt"
	"math/big"

	"github.com/mrvin/tasks-go/fibonacci/internal/cache"
)

func get(prevNumStr, numStr string, from, to uint64) []string {
	prevNum := big.NewInt(0)
	num := big.NewInt(0)

	prevNum.SetString(prevNumStr, 10)
	num.SetString(numStr, 10)

	slValFib := make([]string, 0, to-from+1)
	for i := from; i <= to; i++ {
		prevNum.Add(prevNum, num)
		prevNum, num = num, prevNum
		slValFib = append(slValFib, prevNum.String())
	}

	return slValFib
}

func GetFibNumbers(cacheFib cache.Cache, from, to uint64) ([]string, error) {
	var slValFib []string
	maxCachedNums := cacheFib.GetMaxCachedNum()
	if maxCachedNums >= to {
		var err error
		slValFib, err = cacheFib.GetFromCache(from, to)
		if err != nil {
			return nil, fmt.Errorf("can't get from cache [%d, %d]: %v", from, to, err)
		}
	} else {
		partLeftSlValFib, err := cacheFib.GetFromCache(0, maxCachedNums)
		if err != nil {
			return nil, fmt.Errorf("can't get from cache [%d, %d]: %v", 0, maxCachedNums, err)
		}

		partRightSlValFib := get(partLeftSlValFib[len(partLeftSlValFib)-2], partLeftSlValFib[len(partLeftSlValFib)-1], maxCachedNums, to)
		if err := cacheFib.SetToCache(partRightSlValFib, maxCachedNums, to); err != nil {
			return nil, fmt.Errorf("can't set to cache [%d, %d]: %v", maxCachedNums, to, err)
		}

		slValFib = make([]string, to-from+1)
		if from < uint64(len(partLeftSlValFib)) {
			copy(slValFib, partLeftSlValFib[from:])
			copy(slValFib[maxCachedNums-from+1:], partRightSlValFib[1:])
		} else {
			copy(slValFib, partRightSlValFib[from-1:])
		}

	}

	return slValFib, nil
}
