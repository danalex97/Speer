import React, { Component } from 'react';
import LogExporter from './LogExporter';
import LogUploader from './LogUploader';

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
      <div className="row">
        <div className="col-sm-4"/>
        <LogUploader
          className="col-sm-4"
          placeholder="file.json"
          handleLogUpdate={this.props.handleLogUpdate}/>
        <LogExporter
          className="col-sm-4"
          events={this.props.entries}
          placeholder="file.json"/>
    	</div>

			<div className="container">
				{entries}
			</div>
		</div>);
  }
}
