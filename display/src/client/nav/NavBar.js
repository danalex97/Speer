import React, { Component } from 'react';
import { Navbar, Nav, NavItem, NavDropdown, MenuItem } from 'react-bootstrap';

import {
	DisplayLogEvents,

	DisplayMainEvent,
	DisplayLogEvent,
	DisplayGraphsEvent,
	DisplayStatsEvent
} from './Pages';

export default class NavBar extends Component {
	render() {
		const Handlers = {
		  "Main" : () => window.dispatchEvent(new Event(DisplayMainEvent)),
		  "Logs" : () => window.dispatchEvent(new Event(DisplayLogEvent)),
		  "Graphs" : () => window.dispatchEvent(new Event(DisplayGraphsEvent)),
		  "Stats" : () => window.dispatchEvent(new Event(DisplayStatsEvent))
		};

		return (<Navbar bg="light" expand="lg">
		  <Navbar.Brand onClick={Handlers["Main"]}>Speer</Navbar.Brand>
		  <Navbar.Toggle aria-controls="basic-navbar-nav" />
		  <Navbar.Collapse id="basic-navbar-nav">
		    <Nav className="mr-auto">
		      <Nav.Link onClick={Handlers["Main"]}>Home</Nav.Link>
		      <Nav.Link onClick={Handlers["Logs"]}>Logs</Nav.Link>
					<Nav.Link onClick={Handlers["Graphs"]}>Graphs</Nav.Link>
		      <Nav.Link onClick={Handlers["Stats"]}>Stats</Nav.Link>
		    </Nav>
		  </Navbar.Collapse>
		</Navbar>);
	}
}
