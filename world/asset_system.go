package world

type AssetSystem struct {
	Estates map[int]int
}

func NewAssetSystem() *AssetSystem {
	return &AssetSystem{}
}
