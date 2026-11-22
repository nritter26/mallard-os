package memory

type RegionType string

const (
	RegionInternalRAM RegionType = "internal_ram"
	RegionExternalRAM RegionType = "external_ram"
	RegionDMA RegionType = "dma"
	RegionFlash RegionType = "flash"
)

type Region struct {
	Name string
	Type RegionType
	Total uint64
	Available uint64
}
