import React, { Component } from 'react';

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

const parseEntry = (son) => {
  if (son.node) {
    return new JoinEntry(son);
  } else if (son.src && son.dst) {
    return new PacketEntry(son);
  } else {
    return new LogEntry(son);
  }
};

export {
  LogEntry,
  JoinEntry,
  PacketEntry,

  parseEntry
};
