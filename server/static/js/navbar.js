var navbarData = [{
	"text" : "Main",
	"key" : "l1",
	"url": "#"
}, {
	"text" : "Logs",
	"key" : "l2",
	"url": "#",
}, {
	"text" : "Display",
	"key" : "l3",
	"url" : "#",
}, {
	"text" : "Stats",
	"key" : "l4",
	"url" : "#",
}];

class NavBarLink extends React.Component {
	render() {
		return (<a href={this.props.url}>
			{this.props.text}
		</a>);
	}
}

class NavBarItem extends React.Component {
	constructor(props) {
		super(props);

		let dropdown = !props.submenu ? {} : {
			"className" : "dropdown-toggle",
			"data-toggle" : "dropdown",
			"role" : "buttton",
		};

		this.state = {
			dropdown : dropdown,
		};
	}

	generateLink() {
		return (<NavBarLink
			{...this.state.dropdown}
			url={this.props.url}
			text={this.props.text}
		/>);
	}

	generateSubmenu() {
		return (<NavBar
			items={this.props.submenu}
			dropdown={true}
		/>);
	}

	render() {


		let link = this.generateLink();
		let submenu = <div/>;
		if (this.props.submenu) {
			submenu = this.generateSubmenu();
		}

		return (<li key={this.props.subkey}>
			{link}
			{submenu}
		</li>);
	}
}

class NavBar extends React.Component {
	generateItem(item) {
		return (<NavBarItem
			key={"_"+item.key}
			subkey={item.key}
			text={item.text}
			url={item.url}
			submenu={item.submenu}
		/>);
	}

	render() {
		const navbar   = "nav navbar-nav";
		const dropdown = "dropdown-menu";

		const items = this.props.items.map(this.generateItem);
		const type  = this.props.dropdown ? dropdown : navbar;

		return (<ul className={type}>
			{items}
		</ul>);
	}
}
