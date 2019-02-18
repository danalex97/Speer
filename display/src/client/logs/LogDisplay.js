import React, { Component } from 'react';
import LogExporter from './LogExporter';

export default class LogDisplay extends Component {
  constructor(props) {
    super(props);
  }

  render() {
    const entries = this.props.entries.map((entry, step) => {
      return (
        <div key={step}>
          {entry.render()}
        </div>
      );
    });

    return (<div>
			<LogExporter
				events={this.props.entries}
				placeholder="file.json"/>
			<div className="container">
				{entries}
			</div>
		</div>);
  }
}
