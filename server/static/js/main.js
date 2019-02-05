class Main extends React.Component {
    constructor(props) {
        super(props);

        this.state = {
            tracker : new LogTracker(),

            // metrics
            load : [],

            // events
            events : [],

            // displayable pages
            pages : this.initialDisplayedPages()
        };

        this.handleDisplayEvents = this.handleDisplayEvents.bind(this);

        for (let e of DisplayLogEvents) {
            window.addEventListener(e, this.handleDisplayEvents);
        }
    }

    onUpdate() {
        const tracker = this.state.tracker;
        const windowSize = 10;

        this.setState({
            load: new Metric(windowSize, tracker, Metrics.load).data(),
            events : tracker.events(),
            tracker : tracker,
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
        const tracker = this.state.tracker;
        this.interval = setInterval(
			() => this.onUpdate(),
			tracker.updateInterval,
		);
    }

    componentWillUnmount() {
      clearInterval(this.interval);
    }

    render() {
        const load = this.state.load;
        const events = this.state.events;
        const pages = this.state.pages;

        return (<div>
            <nav className="navbar navbar-default">
                <div className="container-fluid">
                   <NavBar items={navbarData}/>
                </div>
            </nav>
            <div>
                {pages[DisplayLogEvent] ?
                    <LogDisplay events={events}/> :
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

const e = React.createElement;
const domContainer = document.querySelector('#main');
ReactDOM.render(e(Main), domContainer);
