import React, { Component } from 'react';

export default class LogExporter extends Component {
	constructor(props) {
		super(props);
		this.state = {name: ''};

		this.handleChange = this.handleChange.bind(this);
		this.handleClick = this.handleClick.bind(this);
 	}

	handleClick() {
    // TODO: replace with backend callback
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
