package osutil

import "github.com/shopspring/decimal"

// Round 四舍五入.
func Round(x float64, place int32) float64 {
	// strconv.ParseFloat(fmt.Sprintf("%.2f", x), 64) 结果不一致问题，有的能四舍五入有的不可以，很坑的额
	// return math.Floor(x + 0.5) 只能舍成整数
	// return math.Trunc(x*1e2+0.5) * 1e-2  会有精度问题
	// 目前最佳方案使用 decimal 包
	y, _ := decimal.NewFromFloat(x).Round(place).Float64()
	return y
}
