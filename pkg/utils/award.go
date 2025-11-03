package utils

import "strings"

// 根据批发文件对award进行一次模糊鉴定
// 返回字符类型:特等奖、一等奖、二等奖、三等奖、优秀奖
func AppraisalReward(award string) string {
	var result string
	switch award {
	case "特等奖", "金奖", "金牌", "一等奖特等", "特等金奖":
		result = "特等奖"

	case "一等奖", "壹等奖", "银奖", "银牌":
		result = "一等奖"

	case "二等奖", "贰等奖", "铜奖", "铜牌":
		result = "二等奖"

	case "三等奖", "叁等奖":
		result = "三等奖"

	case "优秀奖", "鼓励奖", "参与奖", "纪念奖":
		result = "优秀奖"

	default:
		// 做一次模糊匹配
		if contains(award, "特等") || contains(award, "金牌") || contains(award, "金奖") {
			result = "特等奖"
		} else if contains(award, "一等") || contains(award, "壹等") || contains(award, "银奖") || contains(award, "银牌") {
			result = "一等奖"
		} else if contains(award, "二等") || contains(award, "贰等") || contains(award, "铜奖") || contains(award, "铜牌") {
			result = "二等奖"
		} else if contains(award, "三等") || contains(award, "叁等") {
			result = "三等奖"
		} else if contains(award, "优秀") || contains(award, "鼓励") || contains(award, "参与") || contains(award, "纪念") {
			result = "优秀奖"
		} else {
			result = "优秀奖"
		}
	}
	return result
}

// 辅助函数：检查字符串是否包含子串
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}
