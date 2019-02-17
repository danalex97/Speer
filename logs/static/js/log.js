// TODO: for debugging
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
		return (<div className="join row btn-logs">
			<span className="timestamp col-sm-4"> {this.time} </span>
			<span className="node col-sm-4"> Node {this.node} </span>
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
		return (<div className="packet row btn-logs">
			<span className="timestamp col-sm-4"> {this.time} </span>
			<span className="src col-sm-4"> Src {this.src} </span>
			<span className="dst col-sm-4"> Dest {this.dst} </span>
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
		// TODO: for debugging
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

	// Return log at a point in time; this allows vizualization
	getLogForTime(time) {
		let events = new Array();
		for (let event of this.log) {
			if (event.time <= time) {
				events.push(event);
			}
		}
		return events;
	}

	events() {
		return this.log;
	}
}

class LogExporter extends React.Component {
	constructor(props) {
		super(props);
		this.state = {name: ''};

		this.handleChange = this.handleChange.bind(this);
		this.handleClick = this.handleClick.bind(this);
 	}

	handleClick() {
		const events  = this.props.events;
		const altName = this.props.placeholder;
		const name    = this.state.name;

		let content = JSON.stringify(events.map(e => e.ref));
		let uriContent = "data:application/octet-stream,"
			+ encodeURIComponent(content);

		let link = document.createElement('a');
	    link.setAttribute('href', uriContent);
	    link.setAttribute('download', name ? name : altName);

		let event = document.createEvent('MouseEvents');
        event.initEvent('click', true, true);
        link.dispatchEvent(event);
	}

	handleChange(event) {
		this.setState({name : event.target.value});
	}

	render() {
		return (<div className="row">
			<div className="col-sm-4"/>
			<div className="col-sm-4"/>
			<div className="form-group col-sm-4">
				<input
					type="text"
					placeholder={this.props.placeholder}
					className="formControl"
					value={this.state.name}
					onChange={this.handleChange} />
		  		<button
		  			type="submit"
		  			className="btn btn-outline-dark"
		  			onClick={this.handleClick}>
		  				Download Log
		  		</button>
			</div>
		</div>);
	}
}

class LogDisplay extends React.Component {
	render() {
		const events = this.props.events;
		const entries = events.map((entry, step) => {
			return (
				<div key={step}>
					{entry.render()}
				</div>
			);
		});

		return (<div>
			<LogExporter
				events={events}
				placeholder="file.json"/>
			<div className="container">
				{entries}
			</div>
		</div>);
	}
}
