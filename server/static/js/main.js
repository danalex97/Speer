class Main extends React.Component {
    render() {
        return (<div>
            <h1>Hello World</h1>
            <div>
                <LogDisplay />
            </div>
            <div>
                <BarChart data={[5,10,1,3]} size={[500,500]} />
            </div>
        </div>);
    }
}

const e = React.createElement;
const domContainer = document.querySelector('#main');
ReactDOM.render(e(Main), domContainer);
