var ctr = 5;

class LogEntry {
	constructor(json) {
		this.ref = json;
		this.time = json.time;
	}
}

class JoinEntry extends LogEntry {
	constructor(json) {
		super(json);

		this.node = json.node;
	}

	render() {
		return (<div class="join">
			<span class="timestamp"> {this.time} </span>
			<span class="node"> Node {this.node} </span>
		</div>);
	}
}

class PacketEntry extends LogEntry {
	constructor(json) {
		super(json);

		this.src = json.src;
		this.dst = json.dst;
	}

	render() {
		return (<div class="packet">
			<span class="timestamp"> {this.time} </span>
			<span class="src"> Src {this.src} </span>
			<span class="dst"> Dest {this.dst} </span>
		</div>);
	}
}

class LogTracker {
	constructor() {
		this.url = "http://localhost:8000/new_events";
		this.log = new Array();
		this.updateInterval = 1000;

		// auto-update the log
		setInterval(
			() => this.fetchNewEvents(),
			this.updateInterval,
		);
	}

	getEntry(json) {
		if (json["node"] !== undefined) {
			return new JoinEntry(json);
		}
		if (json["src"] !== undefined) {
			return new PacketEntry(json);
		}
		return null;
	}

	fetchNewEvents() {
		if (ctr == 0) {
			return;
		}
		ctr -= 1;
		// assumes each new call gets new events and the events are orderedx
		fetch(this.url)
			.then(res => res.json())
			.then((events) => {
				for (let raw_event of events) {
					let event = this.getEntry(raw_event);
					this.log.push(event);
				}
			});
	}

	// Return log at a point in time
	getLogForTime() {
		// TODO
	}

	events() {
		return this.log;
	}
}

class LogDisplay extends React.Component {
	constructor(props) {
		super(props);

		this.state = {
			logTracker : new LogTracker(),
			log : [],
		};
	}

	updateLog() {
		if (ctr == 0) {
			return;
		}
		this.setState({
			logTracker: this.state.logTracker,
			log: this.state.logTracker.events(),
		});
		console.log(this.state);
	}

	componentDidMount() {
		// Update the UI at the same interval of requesting data
		this.interval = setInterval(
			() => this.updateLog(),
			this.state.logTracker.updateInterval,
		);
	}

	componentWillUnmount() {
	  clearInterval(this.interval);
	}

	render() {
		const log = this.state.log;
		const entries = log.map((entry, step) => {
			return (
				<li key={step}>
					{entry.render()}
				</li>
			);
		});

		return (<ul>
			{entries}
		</ul>);
	}
}
