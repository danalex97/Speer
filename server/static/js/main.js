class Main extends React.Component {
    constructor(props) {
        super(props);

        this.state = {
            tracker : new LogTracker(),

            // metrics
            load : [],

            // events
            events : [],
        };
    }

    onUpdate() {
        const tracker = this.state.tracker;
        this.setState({
            load: new Metric(1, tracker, Metrics.load).data(),
            events : tracker.events(),
            tracker : tracker,
        });
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

        return (<div>
            <nav className="navbar navbar-default">
                <div className="container-fluid">
                   <NavBar items={navbarData}/>
                </div>
            </nav>
            <div>
                <LogDisplay events={events}/>
            </div>
            <div>
                <BarChart data={load} size={[500,500]} />
            </div>
        </div>);
    }
}

const e = React.createElement;
const domContainer = document.querySelector('#main');
ReactDOM.render(e(Main), domContainer);
