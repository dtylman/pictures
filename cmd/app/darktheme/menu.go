package darktheme

import (
	"github.com/dtylman/gowd"
)

type Menu struct {
	*gowd.Element
	top  *gowd.Element
	side *gowd.Element
}

func NewMenu() *Menu {
	m := &Menu{}
	m.Element, _ = gowd.ParseElement(
		`<nav class="navbar navbar-inverse navbar-fixed-top" role="navigation">
			<div class="navbar-header">
                		<button type="button" class="navbar-toggle" data-toggle="collapse" data-target=".navbar-ex1-collapse">
	                    		<span class="sr-only">Toggle navigation</span>
	                    		<span class="icon-bar"></span>
        	            		<span class="icon-bar"></span>
	                    		<span class="icon-bar"></span>
        	        	</button>
	                	<a class="navbar-brand" href="#">Back to Admin</a>
        	        </div>
			<div class="collapse navbar-collapse navbar-ex1-collapse">
                		<ul class="nav navbar-nav side-nav">
				</ul>
                		<ul class="nav navbar-nav navbar-right navbar-user">
                		</ul>
            		</div>
        	</nav>`)

	m.side = m.Kids[3].Kids[1]
	m.top = m.Kids[3].Kids[3]
	return m
}

//<li><a href="index.html"><i class="fa fa-bullseye"></i> Dashboard</a></li>
func (m*Menu) addButton(submenu*gowd.Element, caption string, icon string, handler gowd.EventHandler) *gowd.Element {
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
	submenu.AddElement(li)
	return link
}

func (m*Menu) AddTopButton(caption string, icon string, handler gowd.EventHandler) *gowd.Element {
	return m.addButton(m.top, caption, icon, handler)
}

func (m*Menu) AddSideButton(caption string, icon string, handler gowd.EventHandler) *gowd.Element {
	return m.addButton(m.side, caption, icon, handler)
}