import React, { Component } from 'react';
import LogExporter from './LogExporter';
import { parseEntry } from './Entries';

const SERVER = "http://localhost:8080";

export default class LogDisplay extends Component {
  constructor(props) {
    super(props);
    this.state = {
      "entries" : []
    };

    this.fetchLog = this.fetchLog.bind(this);
    this.fetchLog();
  }

  fetchLog() {
    fetch(`${SERVER}/api/getLog/${this.props.logName}`)
      .then(res => res.json())
      .then(res => {
        this.setState({
          "entries" : res
        });
      });
  }

  render() {
    const events = this.state.entries.map(parseEntry);
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
