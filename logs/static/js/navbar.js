const DisplayMainEvent = "DisplayMainEvent";
const DisplayLogEvent = "DisplayLogEvent";
const DisplayGraphsEvent = "DisplayGraphsEvent";
const DisplayStatsEvent = "DisplayStatsEvent";

const DisplayLogEvents = [
	DisplayMainEvent,
	DisplayLogEvent,
	DisplayGraphsEvent,
	DisplayStatsEvent
];

const navbarData = [{
	"text" : "Main",
	"key" : "l1",
	"handler": () => window.dispatchEvent(new Event(DisplayMainEvent))
}, {
	"text" : "Logs",
	"key" : "l2",
	"handler": () => window.dispatchEvent(new Event(DisplayLogEvent))
}, {
	"text" : "Graphs",
	"key" : "l3",
	"handler" : () => window.dispatchEvent(new Event(DisplayGraphsEvent))
}, {
	"text" : "Stats",
	"key" : "l4",
	"handler" : () => window.dispatchEvent(new Event(DisplayStatsEvent))
}];

class NavBarLink extends React.Component {
	render() {
		return (<a onClick={this.props.handler}>
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
			handler={this.props.handler}
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
			handler={item.handler}
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
