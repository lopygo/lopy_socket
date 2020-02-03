package packet

type FilterInterface interface {
	Filter() error
}

type FilterResult struct {
	packageBuffer []byte

}

func (filterResult *FilterResult) GetPackageLength() int {
	return len(filterResult.packageBuffer)
}