
package constC

type PlatformSourceType int //平台来源

const (
	Paoshuzw PlatformSourceType = iota+1 //新笔趣阁
)

func platformSourceSlice() []PlatformSourceType {
	return []PlatformSourceType{Paoshuzw}
}

func CheckPlatformSourceInSlice(ps PlatformSourceType) bool {
	for _, v := range platformSourceSlice() {
		if ps == v {
			return true
		}
	}
	return false
}