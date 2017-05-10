package darktheme

import (
	"github.com/dtylman/gowd"
)

type Menu struct {
	*gowd.Element
	TopLeft  *gowd.Element
	TopRight *gowd.Element
	Side     *gowd.Element
}

func NewMenu() *Menu {
	m := &Menu{}
	em:=gowd.NewElementMap()
	m.Element, _ = gowd.ParseElement(
		`<nav class="navbar navbar-inverse navbar-fixed-top" role="navigation">
			<div class="navbar-header">
                		<button type="button" class="navbar-toggle" data-toggle="collapse" data-target=".navbar-ex1-collapse">
	                    		<span class="sr-only">Toggle navigation</span>
	                    		<span class="icon-bar"></span>
        	            		<span class="icon-bar"></span>
	                    		<span class="icon-bar"></span>
        	        	</button>
	                	<a class="navbar-brand" href="#">Bome</a>
        	        </div>
			<div class="collapse navbar-collapse navbar-ex1-collapse">
                		<ul class="nav navbar-nav side-nav" id="pnlSide">
				</ul>
                		<ul class="nav navbar-nav navbar-user" style="margin-left: 148px" id="pnlTopLeft">
                		</ul>
                		<ul class="nav navbar-nav navbar-right navbar-user" id="pnlTopRight">
            		</div>
        	</nav>`,em)

	m.Side = em["pnlSide"]
	m.TopLeft = em["pnlTopLeft"]
	m.TopRight = em["pnlTopRight"]
	return m
}

//<li><a href="index.html"><i class="fa fa-bullseye"></i> Dashboard</a></li>
func (m *Menu) AddButton(panel *gowd.Element, caption string, icon string, handler gowd.EventHandler) *gowd.Element {
	link := gowd.NewElement("a")
	link.SetAttribute("href", "#")
	if icon != "" {
		i := gowd.NewElement("i")
		i.SetAttribute("class", icon)
		link.AddElement(i)
	}
	link.AddElement(gowd.NewText(" " + caption))
	if handler != nil {
		link.OnEvent(gowd.OnClick, handler)
	}
	li := gowd.NewElement("li")
	li.AddElement(link)
	panel.AddElement(li)
	return link
}
