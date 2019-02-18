import React, { Component } from 'react';
import './app.css';
import BarChart from './charts/BarChart';

import {
  DisplayLogEvents,

  DisplayMainEvent,
	DisplayLogEvent,
	DisplayGraphsEvent,
	DisplayStatsEvent
} from './nav/Pages';

import NavBar from './nav/NavBar';
import LogDisplay from './logs/LogDisplay';
import LogFetcher from './logs/LogFetcher';

export default class App extends Component {
  constructor(props) {
    super(props);

    this.state = {
      username : null,

      // metrics
      load : [1,5,45,29,12,32,12,31,13], // some dummy data

      // events
      events : [],

      // displayable pages
      pages : this.initialDisplayedPages(),

      // current log being processed
      logName : "mocklog",

      // current log entries
      logEntries : [],
    };

    this.handleDisplayEvents = this.handleDisplayEvents.bind(this);
    this.fetchLogs = this.fetchLogs.bind(this);

    // fetch logs
    this.fetchLogs(this.state.logName);

    // add listeners for events associated with each page change
    for (let e of DisplayLogEvents) {
        window.addEventListener(e, this.handleDisplayEvents);
    }
  }

  fetchLogs() {
    new LogFetcher(this.state.logName).fetchLog().then((entries) => {
      this.setState(Object.assign(this.state, {
        logEntries : entries,
      }));
    });
  }

  initialDisplayedPages() {
    let pages = {};
    for (let e of DisplayLogEvents) {
      pages[e] = false;
    }
    pages[DisplayMainEvent] = true;
    return pages;
  }


  handleDisplayEvents(event) {
    // change the displayed page
    let pages = this.state.pages;
    for (let k in pages) {
      if (pages.hasOwnProperty(k)) {
        pages[k] = false;
      }
    }
    pages[event.type] = true;

    this.setState(state => ({
      'pages' : pages
    }));
  }

  componentDidMount() {
    fetch('/api/getUsername')
      .then(res => res.json())
      .then(user => this.setState({ username: user.username }));
  }

  render() {
    const username = this.state.username;
    const load     = this.state.load;
    const pages    = this.state.pages;
    const entries  = this.state.logEntries;

    return (<div>
      <NavBar/>

      <div>
        {pages[DisplayLogEvent] ?
          <LogDisplay entries={entries} /> :
          null
        }
      </div>

      <div>
        {pages[DisplayStatsEvent] ?
          <div>
            <BarChart
              id="chart1"
              data={load}
              size={[500,500]}
              margin={60}
              dataSize={40}
            />
            <BarChart
              id="chart2"
              data={load}
              size={[300,300]}
              margin={60}
              dataSize={40}
            />
            <BarChart
              id="chart3"
              data={load}
              size={[500,500]}
              margin={60}
              dataSize={40}
            />
          </div>:
            null
        }
      </div>
    </div>);
  }
}
