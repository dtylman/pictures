package view

import (
	"github.com/dtylman/pictures/view/components"
	"github.com/dtylman/pictures/webkit"
)

//OnConfigChanged called when configuration had changed
func OnConfigChanged() {
	//IndexView.onConfigChanged()
}

//OnConfigChanged called when indexer had stopped
func OnIndexerStopped() {
	//IndexView.onIndexerStopped()
}

//RootElement returns the root "body" container
func RootElement() *webkit.Element {
	return components.Root.Element
}
