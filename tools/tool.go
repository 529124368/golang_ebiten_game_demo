package tools

import "math"

func CaluteDir(x, y, x_tar, y_tar int64) int {
	if x < x_tar && float64(y) == math.Abs(float64(y_tar)) {
		return 2
	}
	if float64(x) == math.Abs(float64(x_tar)) && y > y_tar {
		return 0
	}
	if x < x_tar && y > y_tar {
		return 1
	}
	if x < x_tar && y < y_tar {
		return 3
	}

	if float64(x) == math.Abs(float64(x_tar)) && y < y_tar {
		return 4
	}
	if x > x_tar && y < y_tar {
		return 5
	}
	if x > x_tar && float64(y) == math.Abs(float64(y_tar)) {
		return 6
	}
	if x > x_tar && y > y_tar {
		return 7
	}
	return 0
}
